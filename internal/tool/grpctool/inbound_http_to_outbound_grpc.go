package grpctool

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/errz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/httpz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/prototool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	hostHeader = "Host"

	httpResponseHeaderFieldNumber  protoreflect.FieldNumber = 1
	httpResponseDataFieldNumber    protoreflect.FieldNumber = 2
	httpResponseTrailerFieldNumber protoreflect.FieldNumber = 3
)

var (
	// See https://httpwg.org/http-core/draft-ietf-httpbis-semantics-latest.html#field.connection
	// See https://tools.ietf.org/html/rfc2616#section-13.5.1
	// See https://github.com/golang/go/blob/81ea89adf38b90c3c3a8c4eed9e6c093a8634d59/src/net/http/httputil/reverseproxy.go#L169-L184
	hopHeaders = []string{
		"Connection",
		"Proxy-Connection",
		"Keep-Alive",
		"Proxy-Authenticate",
		"Proxy-Authorization",
		"Te",      // canonicalized version of "TE"
		"Trailer", // not Trailers as per rfc2616; See errata https://www.rfc-editor.org/errata_search.php?eid=4522
		"Transfer-Encoding",
		"Upgrade",
	}
)

type HttpRequestClient interface {
	Send(*HttpRequest) error
	Recv() (*HttpResponse, error)
	grpc.ClientStream
}

type MergeHeadersFunc func(fromOutbound, toInbound http.Header)

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
	var (
		wg            wait.Group
		headerWritten = false
		errFuncRet2   errFunc
	)
	// Pipe remote -> client
	wg.Start(func() {
		headerWritten, errFuncRet2 = x.pipeOutboundToInbound(outboundClient, w)
	})
	// Pipe client -> remote
	errFuncRet1 := x.pipeInboundToOutbound(outboundClient, r, headerExtra)
	wg.Wait()
	if errFuncRet1 != nil {
		return headerWritten, errFuncRet1
	}
	return headerWritten, errFuncRet2
}

func (x *InboundHttpToOutboundGrpc) pipeOutboundToInbound(outboundClient HttpRequestClient, w http.ResponseWriter) (bool, errFunc) {
	writeFailed := false
	headerWritten := false
	err := HttpResponseStreamVisitor().Visit(outboundClient,
		WithCallback(httpResponseHeaderFieldNumber, func(header *HttpResponse_Header) error {
			fromOutbound := header.Response.HttpHeader()
			httpz.RemoveConnectionHeaders(fromOutbound)
			x.MergeHeaders(fromOutbound, w.Header())
			w.WriteHeader(int(header.Response.StatusCode))
			headerWritten = true
			return nil
		}),
		WithCallback(httpResponseDataFieldNumber, func(data *HttpResponse_Data) error {
			_, err := w.Write(data.Data)
			if err != nil {
				writeFailed = true
			}
			return err
		}),
		WithCallback(httpResponseTrailerFieldNumber, func(trailer *HttpResponse_Trailer) error {
			return nil
		}),
	)
	if err != nil {
		if writeFailed {
			// there is likely a connection problem so the client will likely not receive this
			err = errz.NewUserErrorWithCause(err, "")
			return headerWritten, x.handleProcessingError("Proxy failed to write response to client", err)
		}
		return headerWritten, x.handleProcessingError("Proxy failed to read response from agent", err)
	}
	return headerWritten, nil
}

func (x *InboundHttpToOutboundGrpc) pipeInboundToOutbound(outboundClient HttpRequestClient, r *http.Request, headerExtra proto.Message) errFunc {
	extra, err := anypb.New(headerExtra)
	if err != nil {
		return x.handleProcessingError("Proxy failed to marshal header extra proto", err)
	}
	err = outboundClient.Send(&HttpRequest{
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
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil // the other goroutine will receive the error in RecvMsg()
		}
		return x.handleSendError("Proxy failed to send request header", err)
	}

	buffer := make([]byte, maxDataChunkSize)
	for {
		var n int
		n, err = r.Body.Read(buffer)
		if err != nil && !errors.Is(err, io.EOF) {
			// There is likely a connection problem so the client will likely not receive this
			err = errz.NewUserErrorWithCause(err, "")
			return x.handleProcessingError("Proxy failed to read request body from client", err)
		}
		if n > 0 { // handle n=0, err=io.EOF case
			sendErr := outboundClient.Send(&HttpRequest{
				Message: &HttpRequest_Data_{
					Data: &HttpRequest_Data{
						Data: buffer[:n],
					},
				},
			})
			if sendErr != nil {
				if errors.Is(sendErr, io.EOF) {
					return nil // the other goroutine will receive the error in RecvMsg()
				}
				return x.handleSendError("Proxy failed to send request body", sendErr)
			}
		}
		if errors.Is(err, io.EOF) {
			break
		}
	}
	err = outboundClient.Send(&HttpRequest{
		Message: &HttpRequest_Trailer_{
			Trailer: &HttpRequest_Trailer{},
		},
	})
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil // the other goroutine will receive the error in RecvMsg()
		}
		return x.handleSendError("Proxy failed to send trailer", err)
	}
	err = outboundClient.CloseSend()
	if err != nil {
		return x.handleSendError("Proxy failed to send close frame to agent", err)
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
	header.Del(hostHeader) // Use the destination host name

	// Remove hop-by-hop headers
	// 1. Remove headers listed in the Connection header
	httpz.RemoveConnectionHeaders(header)
	// 2. Remove well-known headers
	for _, name := range hopHeaders {
		header.Del(name)
	}

	return prototool.HttpHeaderToValuesMap(header)
}

func writeError(msg string, err error) errFunc {
	return func(w http.ResponseWriter) {
		// See https://tools.ietf.org/html/rfc7231#section-6.6.3
		http.Error(w, fmt.Sprintf("%s: %v", msg, err), http.StatusBadGateway)
	}
}

// errFunc enhances type safety.
type errFunc func(http.ResponseWriter)
