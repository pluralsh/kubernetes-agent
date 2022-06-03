package grpctool

import (
	"bufio"
	"context"
	"errors"
	"io"
	"net"
	"net/http"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/memz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/prototool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InboundGrpcToOutboundHttpStream interface {
	Send(*HttpResponse) error
	grpc.ServerStream
}

type HandleProcessingErrorFunc func(msg string, err error)
type HandleIoErrorFunc func(msg string, err error) error

type DoResponse struct {
	// Resp is the server's response to a request.
	Resp *http.Response
	// UpgradeConn is the underlying network connection to the server.
	// May be nil if request was not an Upgrade request or if server decided not to switch protocols
	// (non-101 response status code).
	UpgradeConn net.Conn
	// ConnReader is a buffered reader, wrapping UpgradeConn. Is set when UpgradeConn is set.
	// Must be used for reading as it may contain buffered bytes that are no longer available directly via UpgradeConn.
	ConnReader *bufio.Reader
}

// HttpDo makes an HTTP request and returns a response. If an HTTP upgrade was requested, the underlying network
// connection is also returned. Implementations that don't support Upgrade should return an error.
type HttpDo func(ctx context.Context, header *HttpRequest_Header, body io.Reader) (DoResponse, error)

type InboundGrpcToOutboundHttp struct {
	Log                   *zap.Logger
	HandleProcessingError HandleProcessingErrorFunc
	HandleIoError         HandleIoErrorFunc
	HttpDo                HttpDo
}

func (x *InboundGrpcToOutboundHttp) Pipe(inbound InboundGrpcToOutboundHttpStream) (retErr error) {
	var upgradeConn net.Conn
	defer func() {
		if upgradeConn != nil {
			err := upgradeConn.Close()
			if retErr == nil {
				retErr = x.maybeHandleIoError("error closing connection", err)
			}
		}
	}()

	ctx := inbound.Context()

	pr, pw := io.Pipe()
	headerC := make(chan *HttpRequest_Header)
	// buffered to not block the sender as receiver might encounter an error and exit before even trying to receive.
	respC := make(chan DoResponse, 1)
	s := InboundStreamToOutboundStream{
		// Pipe gRPC request -> HTTP request
		PipeInboundToOutbound: func() error {
			// unblock the PipeOutboundToInbound goroutine if we exited before sending the header due to an error.
			defer close(headerC)
			return x.pipeInboundToOutbound(inbound, headerC, respC, pw)
		},
		// Pipe HTTP response -> gRPC response
		PipeOutboundToInbound: func() error {
			// Make sure the writer is unblocked if we exit abruptly
			// The error is ignored because it will always occur if things go normally - the pipe will have been
			// closed already when this code is reached (and that's an error).
			defer pr.Close() // nolint: errcheck
			// unblock the PipeInboundToOutbound goroutine if we exited before sending the response object due to an error.
			defer close(respC)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case header, ok := <-headerC:
				if !ok {
					// Something went wrong in the PipeInboundToOutbound goroutine, exit.
					return nil
				}
				r, err := x.HttpDo(ctx, header, pr)
				if err != nil {
					return err
				}
				respC <- r
				// this store is not synchronized and that's ok because PipeOutboundToInbound is executed
				// on the caller's goroutine.
				upgradeConn = r.UpgradeConn
				return x.pipeOutboundToInbound(inbound, r, header.Request.IsUpgrade())
			}
		},
	}
	err := s.Pipe()
	switch {
	case err == nil:
	case IsStatusError(err):
		// A gRPC status already
	case errors.Is(err, context.Canceled):
		x.Log.Debug("gRPC -> HTTP", logz.Error(err))
		err = status.Errorf(codes.Canceled, "gRPC -> HTTP: %v", err)
	case errors.Is(err, context.DeadlineExceeded):
		x.Log.Debug("gRPC -> HTTP", logz.Error(err))
		err = status.Errorf(codes.DeadlineExceeded, "gRPC -> HTTP: %v", err)
	default:
		x.HandleProcessingError("gRPC -> HTTP", err)
		err = status.Errorf(codes.Unavailable, "gRPC -> HTTP: %v", err)
	}
	return err
}

func (x *InboundGrpcToOutboundHttp) pipeInboundToOutbound(inbound InboundGrpcToOutboundHttpStream,
	headerC chan<- *HttpRequest_Header, respC <-chan DoResponse, pw *io.PipeWriter) error {
	var isUpgrade bool
	var upgradeConn net.Conn
	return HttpRequestStreamVisitor().Visit(inbound,
		WithCallback(HttpRequestHeaderFieldNumber, func(header *HttpRequest_Header) error {
			x.logRequest(header)
			isUpgrade = header.Request.IsUpgrade()
			ctx := inbound.Context()
			select {
			case <-ctx.Done():
				return ctx.Err()
			case headerC <- header:
				return nil
			}
		}),
		WithCallback(HttpRequestDataFieldNumber, func(data *HttpRequest_Data) error {
			_, err := pw.Write(data.Data)
			return x.maybeHandleIoError("request body write", err)
		}),
		WithCallback(HttpRequestTrailerFieldNumber, func(trailer *HttpRequest_Trailer) error {
			if isUpgrade {
				// Nothing more to send, close the write end of the pipe
				err := pw.Close()
				return x.maybeHandleIoError("request body close", err)
			}
			// Nothing to do
			return nil
		}),
		WithCallback(HttpRequestUpgradeDataFieldNumber, func(data *HttpRequest_UpgradeData) error {
			if !isUpgrade {
				// Inbound client didn't request a connection upgrade but sent an upgrade data frame.
				return status.Error(codes.Internal, "unexpected HttpRequest_UpgradeData message for non-upgrade request")
			}
			if upgradeConn == nil {
				r, ok := <-respC
				if !ok {
					// error in the other goroutine, abort.
					return context.Canceled
				}
				if r.Resp.StatusCode != http.StatusSwitchingProtocols {
					// Outbound server doesn't want to switch protocols but inbound client sent an upgrade data frame.
					return status.Errorf(codes.Internal, "unexpected HttpRequest_UpgradeData message for HTTP status code %d", r.Resp.StatusCode)
				}
				upgradeConn = r.UpgradeConn
			}
			_, err := upgradeConn.Write(data.Data)
			return x.maybeHandleIoError("upgrade request write", err)
		}),
		WithEOFCallback(func() error {
			if !isUpgrade {
				// Nothing more to send, close the write end of the pipe
				err := pw.Close()
				return x.maybeHandleIoError("request body close", err)
			}
			return nil
		}),
	)
}

func (x *InboundGrpcToOutboundHttp) logRequest(header *HttpRequest_Header) {
	if !x.Log.Core().Enabled(zap.DebugLevel) {
		return
	}
	req := header.Request
	sugar := x.Log.Sugar()
	if len(req.Query) > 0 {
		sugar.Debugf("Handling %s %s?%s", req.Method, req.UrlPath, req.UrlQuery().Encode())
	} else {
		sugar.Debugf("Handling %s %s", req.Method, req.UrlPath)
	}
}

func (x *InboundGrpcToOutboundHttp) pipeOutboundToInbound(inbound InboundGrpcToOutboundHttpStream, r DoResponse, isUpgrade bool) error {
	err := x.sendResponseHeaderAndBody(inbound, r.Resp)
	if err != nil {
		return err
	}

	err = inbound.Send(&HttpResponse{
		Message: &HttpResponse_Trailer_{
			Trailer: &HttpResponse_Trailer{},
		},
	})
	if err != nil {
		return x.handleIoError("SendMsg(HttpResponse_Trailer) failed", err)
	}
	if isUpgrade && r.Resp.StatusCode == http.StatusSwitchingProtocols {
		// Only stream if upgrade was requested AND outbound server is switching protocols.
		return x.sendUpgradeResponseStream(inbound, r.ConnReader)
	}
	return nil
}

func (x *InboundGrpcToOutboundHttp) sendResponseHeaderAndBody(inbound InboundGrpcToOutboundHttpStream, resp *http.Response) (retErr error) {
	defer func() {
		err := resp.Body.Close()
		if retErr == nil {
			retErr = x.maybeHandleIoError("response body close", err)
		}
	}()
	err := inbound.Send(&HttpResponse{
		Message: &HttpResponse_Header_{
			Header: &HttpResponse_Header{
				Response: &prototool.HttpResponse{
					StatusCode: int32(resp.StatusCode),
					Status:     resp.Status,
					Header:     prototool.HttpHeaderToValuesMap(resp.Header),
				},
			},
		},
	})
	if err != nil {
		return x.handleIoError("SendMsg(HttpResponse_Header) failed", err)
	}

	buffer := memz.Get32k()
	defer memz.Put32k(buffer)
	for {
		n, readErr := resp.Body.Read(buffer)
		if n > 0 { // handle n>0 before readErr != nil to ensure any consumed data gets forwarded
			sendErr := inbound.Send(&HttpResponse{
				Message: &HttpResponse_Data_{
					Data: &HttpResponse_Data{
						Data: buffer[:n],
					},
				},
			})
			if sendErr != nil {
				return x.handleIoError("SendMsg(HttpResponse_Data) failed", sendErr)
			}
		}
		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				break
			}
			return x.handleIoError("read HTTP response body", readErr)
		}
	}
	return nil
}

func (x *InboundGrpcToOutboundHttp) sendUpgradeResponseStream(inbound InboundGrpcToOutboundHttpStream, upgradeConn *bufio.Reader) error {
	buffer := memz.Get32k()
	defer memz.Put32k(buffer)
	for {
		n, readErr := upgradeConn.Read(buffer)
		if n > 0 { // handle n>0 before readErr != nil to ensure any consumed data gets forwarded
			sendErr := inbound.Send(&HttpResponse{
				Message: &HttpResponse_UpgradeData_{
					UpgradeData: &HttpResponse_UpgradeData{
						Data: buffer[:n],
					},
				},
			})
			if sendErr != nil {
				return x.handleIoError("SendMsg(HttpResponse_UpgradeData) failed", sendErr)
			}
		}
		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				break
			}
			return x.handleIoError("read upgrade response body", readErr)
		}
	}
	return nil
}

func (x *InboundGrpcToOutboundHttp) maybeHandleIoError(msg string, err error) error {
	if err != nil {
		return x.handleIoError(msg, err)
	}
	return nil
}

func (x *InboundGrpcToOutboundHttp) handleIoError(msg string, err error) error {
	return x.HandleIoError("gRPC -> HTTP: "+msg, err)
}
