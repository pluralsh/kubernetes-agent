package grpctool

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/httpz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/memz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/prototool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

var (
	// See https://httpwg.org/http-core/draft-ietf-httpbis-semantics-latest.html#field.connection
	// See https://datatracker.ietf.org/doc/html/rfc2616#section-13.5.1
	// See https://github.com/golang/go/blob/81ea89adf38b90c3c3a8c4eed9e6c093a8634d59/src/net/http/httputil/reverseproxy.go#L169-L184
	// Must be in canonical form.
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

	// earlyExitError is a sentinel error value to make stream visitor exit early.
	earlyExitError = errors.New("")
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

func (x *InboundHttpToOutboundGrpc) pipe(outboundClient HttpRequestClient, w http.ResponseWriter, r *http.Request,
	headerExtra proto.Message) (bool /* headerWritten */, errFunc) {
	// 0. Check if connection upgrade is requested and if connection can be hijacked.
	var hijacker http.Hijacker
	isUpgrade := len(r.Header[httpz.UpgradeHeader]) > 0
	if isUpgrade {
		// Connection upgrade requested. For that ResponseWriter must support hijacking.
		var ok bool
		hijacker, ok = w.(http.Hijacker)
		if !ok {
			return false, x.handleInternalError("unable to upgrade connection", fmt.Errorf("unable to hijack response: %T does not implement http.Hijacker", w))
		}
	}
	// http.ResponseWriter does not support concurrent request body reads and response writes so
	// consume the request body first and then write the response from remote.
	// See https://github.com/golang/go/issues/15527
	// See https://github.com/golang/go/blob/go1.17.2/src/net/http/server.go#L118-L139

	// 1. Pipe client -> remote
	errF := x.pipeInboundToOutbound(outboundClient, r, headerExtra)
	if errF != nil {
		return false, errF
	}
	if !isUpgrade { // Close outbound connection for writes if it's not an upgraded connection
		errF = x.sendCloseSend(outboundClient)
		if errF != nil {
			return false, errF
		}
	}
	// 2. Pipe remote -> client
	headerWritten, responseStatusCode, errF := x.pipeOutboundToInbound(outboundClient, w, isUpgrade)
	if errF != nil {
		return headerWritten, errF
	}
	// 3. Pipe client <-> remote if connection upgrade is requested
	if !isUpgrade { // nothing to do
		return true, nil
	}
	if responseStatusCode != http.StatusSwitchingProtocols {
		// Remote doesn't want to upgrade the connection
		return true, x.sendCloseSend(outboundClient)
	}
	return true, x.pipeUpgradedConnection(outboundClient, hijacker)
}

func (x *InboundHttpToOutboundGrpc) pipeOutboundToInbound(outboundClient HttpRequestClient, w http.ResponseWriter, isUpgrade bool) (bool, int32, errFunc) {
	writeFailed := false
	headerWritten := false
	var responseStatusCode int32
	flush := x.flush(w)
	err := HttpResponseStreamVisitor().Visit(outboundClient,
		WithCallback(HttpResponseHeaderFieldNumber, func(header *HttpResponse_Header) error {
			responseStatusCode = header.Response.StatusCode
			outboundResponse := header.Response.HttpHeader()
			cleanHeader(outboundResponse)
			x.MergeHeaders(outboundResponse, w.Header())
			w.WriteHeader(int(header.Response.StatusCode))
			flush()
			headerWritten = true
			return nil
		}),
		WithCallback(HttpResponseDataFieldNumber, func(data *HttpResponse_Data) error {
			_, err := w.Write(data.Data)
			if err != nil {
				writeFailed = true
				return err
			}
			flush()
			return nil
		}),
		WithCallback(HttpResponseTrailerFieldNumber, func(trailer *HttpResponse_Trailer) error {
			if isUpgrade && responseStatusCode == http.StatusSwitchingProtocols {
				// Successful upgrade.
				return earlyExitError
			}
			return nil
		}),
		// if it's a successful upgrade, then this field is unreachable because of the early exit above.
		// otherwise, (unsuccessful upgrade or not an upgrade) the remote must not send this field.
		WithNotExpectingToGet(codes.Internal, HttpResponseUpgradeDataFieldNumber),
	)
	if err != nil && err != earlyExitError { // nolint: errorlint
		if writeFailed {
			// there is likely a connection problem so the client will likely not receive this
			return headerWritten, responseStatusCode, x.handleIoError("failed to write HTTP response", err)
		}
		return headerWritten, responseStatusCode, x.handleIoError("failed to read gRPC response", err)
	}
	return headerWritten, responseStatusCode, nil
}

func (x *InboundHttpToOutboundGrpc) flush(w http.ResponseWriter) func() {
	// ResponseWriter buffers headers and response body writes and that may break use cases like long polling or streaming.
	// Flusher is used so that when HTTP headers and response body chunks are received from the outbound connection,
	// they are flushed to the inbound stream ASAP.
	flusher, ok := w.(http.Flusher)
	if !ok {
		x.Log.Sugar().Warnf("HTTP->gRPC: %T does not implement http.Flusher, cannot flush data to client", w)
		return func() {}
	}
	return flusher.Flush
}

func (x *InboundHttpToOutboundGrpc) pipeInboundToOutbound(outboundClient HttpRequestClient, r *http.Request, headerExtra proto.Message) errFunc {
	extra, err := anypb.New(headerExtra)
	if err != nil {
		return x.handleProcessingError("failed to marshal header extra proto", err)
	}
	errF := x.send(outboundClient, "failed to send request header", &HttpRequest{
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
	return x.send(outboundClient, "failed to send trailer", &HttpRequest{
		Message: &HttpRequest_Trailer_{
			Trailer: &HttpRequest_Trailer{},
		},
	})
}

func (x *InboundHttpToOutboundGrpc) sendRequestBody(outboundClient HttpRequestClient, body io.Reader) errFunc {
	buffer := memz.Get32k()
	defer memz.Put32k(buffer)
	for {
		n, err := body.Read(buffer)
		if err != nil && !errors.Is(err, io.EOF) {
			// There is likely a connection problem so the client will likely not receive this
			return x.handleIoError("failed to read request body", err)
		}
		if n > 0 { // handle n=0, err=io.EOF case
			errF := x.send(outboundClient, "failed to send request body", &HttpRequest{
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

func (x *InboundHttpToOutboundGrpc) sendCloseSend(outboundClient HttpRequestClient) errFunc {
	err := outboundClient.CloseSend()
	if err != nil {
		return x.handleIoError("failed to send close frame", err)
	}
	return nil
}

func (x *InboundHttpToOutboundGrpc) send(client HttpRequestClient, errMsg string, msg *HttpRequest) errFunc {
	err := client.Send(msg)
	if err != nil {
		if errors.Is(err, io.EOF) {
			_, err = client.Recv()
		}
		return x.handleIoError(errMsg, err)
	}
	return nil
}

func (x *InboundHttpToOutboundGrpc) handleIoError(msg string, err error) errFunc {
	msg = "HTTP->gRPC: " + msg
	x.Log.Debug(msg, logz.Error(err))
	return writeError(msg, err)
}

func (x *InboundHttpToOutboundGrpc) handleProcessingError(msg string, err error) errFunc {
	x.HandleProcessingError(msg, err)
	return writeError(msg, err)
}

func (x *InboundHttpToOutboundGrpc) handleInternalError(msg string, err error) errFunc {
	msg = "HTTP->gRPC: " + msg
	x.HandleProcessingError(msg, err)
	return func(w http.ResponseWriter) {
		// See https://datatracker.ietf.org/doc/html/rfc7231#section-6.6.1
		http.Error(w, fmt.Sprintf("%s: %v", msg, err), http.StatusInternalServerError)
	}
}

func (x *InboundHttpToOutboundGrpc) pipeUpgradedConnection(outboundClient HttpRequestClient, hijacker http.Hijacker) (errRet errFunc) {
	conn, bufrw, err := hijacker.Hijack()
	if err != nil {
		return x.handleInternalError("unable to upgrade connection: error hijacking response", err)
	}
	defer func() {
		err = conn.Close()
		if err != nil && errRet == nil {
			errRet = x.handleIoError("failed to close upgraded connection", err)
		}
	}()
	// Hijack() docs say we are responsible for managing connection deadlines and a deadline may be set already.
	// We clear the read deadline here because we don't know if the client will be sending any data to us soon.
	err = conn.SetReadDeadline(time.Time{})
	if err != nil {
		return x.handleIoError("failed to clear connection read deadline", err)
	}
	// We don't care if a write deadline is set already, we just wrap the connection in a wrapper that
	// will each time set a new deadline before performing an actual write.
	conn = &httpz.WriteTimeoutConn{
		Conn:    conn,
		Timeout: 20 * time.Second,
	}
	p := InboundStreamToOutboundStream{
		PipeInboundToOutbound: func() error {
			return x.pipeInboundToOutboundUpgraded(outboundClient, bufrw.Reader)
		},
		PipeOutboundToInbound: func() error {
			return x.pipeOutboundToInboundUpgraded(outboundClient, conn)
		},
	}
	err = p.Pipe()
	if err != nil {
		return x.handleIoError("failed to pipe upgraded connection streams", err)
	}
	return nil
}

func (x *InboundHttpToOutboundGrpc) pipeInboundToOutboundUpgraded(outboundClient HttpRequestClient, inboundStream *bufio.Reader) error {
	buffer := memz.Get32k()
	defer memz.Put32k(buffer)
	for {
		n, err := inboundStream.Read(buffer)
		if err != nil && !errors.Is(err, io.EOF) {
			// There is likely a connection problem so the client will likely not receive this
			return fmt.Errorf("read failed: %w", err)
		}
		if n > 0 { // handle n=0, err=io.EOF case
			sendErr := outboundClient.Send(&HttpRequest{
				Message: &HttpRequest_UpgradeData_{
					UpgradeData: &HttpRequest_UpgradeData{
						Data: buffer[:n],
					},
				},
			})
			if sendErr != nil {
				if errors.Is(sendErr, io.EOF) {
					return nil // the other goroutine will receive the error in RecvMsg()
				}
				return fmt.Errorf("Send(HttpRequest_UpgradeData): %w", sendErr)
			}
		}
		if errors.Is(err, io.EOF) {
			break
		}
	}
	err := outboundClient.CloseSend()
	if err != nil {
		return fmt.Errorf("failed to send close frame: %w", err)
	}
	return nil
}

func (x *InboundHttpToOutboundGrpc) pipeOutboundToInboundUpgraded(outboundClient HttpRequestClient, inboundStream io.Writer) error {
	var writeFailed bool
	err := HttpResponseStreamVisitor().Visit(outboundClient,
		WithStartState(HttpResponseTrailerFieldNumber),
		WithCallback(HttpResponseUpgradeDataFieldNumber, func(data *HttpResponse_UpgradeData) error {
			_, err := inboundStream.Write(data.Data)
			if err != nil {
				writeFailed = true
			}
			return err
		}),
	)
	if err != nil {
		if writeFailed {
			// there is likely a connection problem so the client will likely not receive this
			return fmt.Errorf("failed to write upgraded HTTP response: %w", err)
		}
		return fmt.Errorf("failed to read gRPC response: %w", err)
	}
	return nil
}

func headerFromHttpRequestHeader(header http.Header) map[string]*prototool.Values {
	header = header.Clone()
	delete(header, httpz.HostHeader) // Use the destination host name
	cleanHeader(header)
	return prototool.HttpHeaderToValuesMap(header)
}

func cleanHeader(header http.Header) {
	upgrade := header[httpz.UpgradeHeader]

	// 1. Remove hop-by-hop headers listed in the Connection header. See https://datatracker.ietf.org/doc/html/rfc7230#section-6.1
	httpz.RemoveConnectionHeaders(header)
	// 2. Remove well-known hop-by-hop headers
	for _, name := range hopHeaders {
		delete(header, name)
	}
	// 3. Fix up Connection and Upgrade headers if upgrade is requested/confirmed
	if len(upgrade) > 0 {
		header[httpz.UpgradeHeader] = upgrade                // put it back
		header[httpz.ConnectionHeader] = []string{"upgrade"} // this discards any other connection options if they were there
	}
}

func writeError(msg string, err error) errFunc {
	return func(w http.ResponseWriter) {
		// See https://datatracker.ietf.org/doc/html/rfc7231#section-6.6.3
		http.Error(w, fmt.Sprintf("%s: %v", msg, err), http.StatusBadGateway)
	}
}

// errFunc enhances type safety.
type errFunc func(http.ResponseWriter)
