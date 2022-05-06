package grpctool

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/errz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/httpz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/memz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/prototool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

var (
	// See https://httpwg.org/http-core/draft-ietf-httpbis-semantics-latest.html#field.connection
	// See https://datatracker.ietf.org/doc/html/rfc2616#section-13.5.1
	// See https://github.com/golang/go/blob/81ea89adf38b90c3c3a8c4eed9e6c093a8634d59/src/net/http/httputil/reverseproxy.go#L169-L184
	hopHeaders = []string{
		httpz.ConnectionHeader,
		httpz.ProxyConnectionHeader,
		httpz.KeepAliveHeader,
		httpz.ProxyAuthenticateHeader,
		httpz.ProxyAuthorizationHeader,
		httpz.TeHeader,
		httpz.TrailerHeader,
		httpz.TransferEncodingHeader,
		httpz.UpgradeHeader,
	}
)

type HttpRequestClient interface {
	Send(*HttpRequest) error
	Recv() (*HttpResponse, error)
	grpc.ClientStream
}

type MergeHeadersFunc func(outboundResponse, inboundResponse http.Header)

type InboundHttpToOutboundGrpc struct {
	Log                   *zap.Logger
	HandleProcessingError HandleProcessingErrorFunc
	MergeHeaders          MergeHeadersFunc
}

func (x *InboundHttpToOutboundGrpc) Pipe(outboundClient HttpRequestClient, w http.ResponseWriter, r *http.Request, headerExtra proto.Message) {
	headerWritten, errF := x.pipe(outboundClient, w, r, headerExtra)
	if errF != nil {
		if headerWritten {
			// HTTP status has been written already as part of the normal response flow.
			// But then something went wrong and an error happened. To let the client know that something isn't right
			// we have only one thing we can do - abruptly close the connection. To do that we panic with a special
			// error value that the "http" package provides. See its description.
			// If we try to write the status again here, http package would log a warning, which is not nice.
			panic(http.ErrAbortHandler)
		} else {
			errF(w)
		}
	}
}

func (x *InboundHttpToOutboundGrpc) pipe(outboundClient HttpRequestClient, w http.ResponseWriter, r *http.Request, headerExtra proto.Message) (bool, errFunc) {
	// http.ResponseWriter does not support concurrent request body reads and response writes so
	// consume the request body first and then write the response from remote.
	// See https://github.com/golang/go/issues/15527
	// See https://github.com/golang/go/blob/go1.17.2/src/net/http/server.go#L118-L139

	// Pipe client -> remote
	ef := x.pipeInboundToOutbound(outboundClient, r, headerExtra)
	if ef != nil {
		return false, ef
	}
	// Pipe remote -> client
	return x.pipeOutboundToInbound(outboundClient, w)
}

func (x *InboundHttpToOutboundGrpc) pipeOutboundToInbound(outboundClient HttpRequestClient, w http.ResponseWriter) (bool, errFunc) {
	writeFailed := false
	headerWritten := false
	// ResponseWriter buffers headers and response body writes and that may break use cases like long polling or streaming.
	// Flusher is used so that when HTTP headers and response body chunks are received from the outbound connection,
	// they are flushed to the inbound stream ASAP.
	flusher, ok := w.(http.Flusher)
	if !ok {
		x.Log.Sugar().Warnf("HTTP->gRPC: %T does not implement http.Flusher, cannot flush data to client", w)
	}
	err := HttpResponseStreamVisitor().Visit(outboundClient,
		WithCallback(HttpResponseHeaderFieldNumber, func(header *HttpResponse_Header) error {
			outboundResponse := header.Response.HttpHeader()
			httpz.RemoveConnectionHeaders(outboundResponse)
			x.MergeHeaders(outboundResponse, w.Header())
			w.WriteHeader(int(header.Response.StatusCode))
			if flusher != nil {
				flusher.Flush()
			}
			headerWritten = true
			return nil
		}),
		WithCallback(HttpResponseDataFieldNumber, func(data *HttpResponse_Data) error {
			_, err := w.Write(data.Data)
			if err != nil {
				writeFailed = true
			} else {
				if flusher != nil {
					flusher.Flush()
				}
			}
			return err
		}),
		WithCallback(HttpResponseTrailerFieldNumber, func(trailer *HttpResponse_Trailer) error {
			return nil
		}),
	)
	if err != nil {
		if writeFailed {
			// there is likely a connection problem so the client will likely not receive this
			err = errz.NewUserErrorWithCause(err, "")
			return headerWritten, x.handleProcessingError("HTTP->gRPC: failed to write HTTP response", err)
		}
		return headerWritten, x.handleProcessingError("HTTP->gRPC: failed to read gRPC response", err)
	}
	return headerWritten, nil
}

func (x *InboundHttpToOutboundGrpc) pipeInboundToOutbound(outboundClient HttpRequestClient, r *http.Request, headerExtra proto.Message) errFunc {
	extra, err := anypb.New(headerExtra)
	if err != nil {
		return x.handleProcessingError("HTTP->gRPC: failed to marshal header extra proto", err)
	}
	errF := x.send(outboundClient, "HTTP->gRPC: failed to send request header", &HttpRequest{
		Message: &HttpRequest_Header_{
			Header: &HttpRequest_Header{
				Request: &prototool.HttpRequest{
					Method:  r.Method,
					Header:  headerFromHttpRequestHeader(r.Header),
					UrlPath: r.URL.Path,
					Query:   prototool.UrlValuesToValuesMap(r.URL.Query()),
				},
				Extra: extra,
			},
		},
	})
	if errF != nil {
		return errF
	}

	errF = x.sendRequestBody(outboundClient, r.Body)
	if errF != nil {
		return errF
	}
	errF = x.send(outboundClient, "HTTP->gRPC: failed to send trailer", &HttpRequest{
		Message: &HttpRequest_Trailer_{
			Trailer: &HttpRequest_Trailer{},
		},
	})
	if errF != nil {
		return errF
	}
	err = outboundClient.CloseSend()
	if err != nil {
		return x.handleSendError("HTTP->gRPC: failed to send close frame", err)
	}
	return nil
}

func (x *InboundHttpToOutboundGrpc) sendRequestBody(outboundClient HttpRequestClient, body io.Reader) errFunc {
	buffer := memz.Get32k()
	defer memz.Put32k(buffer)
	for {
		n, err := body.Read(buffer)
		if err != nil && !errors.Is(err, io.EOF) {
			// There is likely a connection problem so the client will likely not receive this
			err = errz.NewUserErrorWithCause(err, "")
			return x.handleProcessingError("HTTP->gRPC: failed to read request body", err)
		}
		if n > 0 { // handle n=0, err=io.EOF case
			errF := x.send(outboundClient, "HTTP->gRPC: failed to send request body", &HttpRequest{
				Message: &HttpRequest_Data_{
					Data: &HttpRequest_Data{
						Data: buffer[:n],
					},
				},
			})
			if errF != nil {
				return errF
			}
		}
		if errors.Is(err, io.EOF) {
			break
		}
	}
	return nil
}

func (x *InboundHttpToOutboundGrpc) send(client HttpRequestClient, errMsg string, msg *HttpRequest) errFunc {
	err := client.Send(msg)
	if err != nil {
		if errors.Is(err, io.EOF) {
			_, err = client.Recv()
		}
		return x.handleSendError(errMsg, err)
	}
	return nil
}

func (x *InboundHttpToOutboundGrpc) handleSendError(msg string, err error) errFunc {
	x.Log.Debug(msg, logz.Error(err))
	return writeError(msg, err)
}

func (x *InboundHttpToOutboundGrpc) handleProcessingError(msg string, err error) errFunc {
	x.HandleProcessingError(msg, err)
	return writeError(msg, err)
}

func headerFromHttpRequestHeader(header http.Header) map[string]*prototool.Values {
	header = header.Clone()
	delete(header, httpz.HostHeader) // Use the destination host name

	// Remove hop-by-hop headers
	// 1. Remove headers listed in the Connection header
	httpz.RemoveConnectionHeaders(header)
	// 2. Remove well-known headers
	for _, name := range hopHeaders {
		delete(header, name)
	}

	return prototool.HttpHeaderToValuesMap(header)
}

func writeError(msg string, err error) errFunc {
	return func(w http.ResponseWriter) {
		// See https://datatracker.ietf.org/doc/html/rfc7231#section-6.6.3
		http.Error(w, fmt.Sprintf("%s: %v", msg, err), http.StatusBadGateway)
	}
}

// errFunc enhances type safety.
type errFunc func(http.ResponseWriter)
