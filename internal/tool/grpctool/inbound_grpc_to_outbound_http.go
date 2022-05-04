package grpctool

import (
	"context"
	"errors"
	"io"
	"net/http"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/errz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/memz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/prototool"
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
type HandleSendErrorFunc func(msg string, err error) error
type HttpDo func(ctx context.Context, header *HttpRequest_Header, body io.Reader) (*http.Response, error)

type InboundGrpcToOutboundHttp struct {
	Log                   *zap.Logger
	HandleProcessingError HandleProcessingErrorFunc
	HandleSendError       HandleSendErrorFunc
	HttpDo                HttpDo
}

func (x *InboundGrpcToOutboundHttp) Pipe(inbound InboundGrpcToOutboundHttpStream) error {
	ctx := inbound.Context()

	pr, pw := io.Pipe()
	headerMsg := make(chan *HttpRequest_Header)
	s := InboundGrpcToOutboundStream{
		// Pipe gRPC request -> HTTP request
		PipeInboundToOutbound: func() error {
			return x.pipeInboundToOutbound(inbound, headerMsg, pw)
		},
		// Pipe HTTP response -> gRPC response
		PipeOutboundToInbound: func() error {
			// Make sure the writer is unblocked if we exit abruptly
			// The error is ignored because it will always occur if things go normally - the pipe will have been
			// closed already when this code is reached (and that's an error).
			defer pr.Close() // nolint: errcheck
			select {
			case <-ctx.Done():
				return ctx.Err()
			case header := <-headerMsg:
				resp, err := x.HttpDo(ctx, header, pr)
				if err != nil {
					return err
				}
				return x.pipeOutboundToInbound(inbound, resp)
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

func (x *InboundGrpcToOutboundHttp) pipeInboundToOutbound(inbound InboundGrpcToOutboundHttpStream, headerMsg chan<- *HttpRequest_Header, pw *io.PipeWriter) error {
	return HttpRequestStreamVisitor().Visit(inbound,
		WithCallback(HttpRequestHeaderFieldNumber, func(header *HttpRequest_Header) error {
			ctx := inbound.Context()
			select {
			case <-ctx.Done():
				return ctx.Err()
			case headerMsg <- header:
				return nil
			}
		}),
		WithCallback(HttpRequestDataFieldNumber, func(data *HttpRequest_Data) error {
			_, err := pw.Write(data.Data)
			return err
		}),
		WithCallback(HttpRequestTrailerFieldNumber, func(trailer *HttpRequest_Trailer) error {
			// Nothing to do
			return nil
		}),
		WithEOFCallback(pw.Close), // Nothing more to send, close the write end of the pipe
	)
}

func (x *InboundGrpcToOutboundHttp) pipeOutboundToInbound(inbound InboundGrpcToOutboundHttpStream, resp *http.Response) error {
	err := x.sendResponseHeaderAndBody(inbound, resp)
	if err != nil {
		return err
	}

	err = inbound.Send(&HttpResponse{
		Message: &HttpResponse_Trailer_{
			Trailer: &HttpResponse_Trailer{},
		},
	})
	if err != nil {
		return x.HandleSendError("SendMsg(HttpResponse_Trailer) failed", err)
	}
	return nil
}

func (x *InboundGrpcToOutboundHttp) sendResponseHeaderAndBody(inbound InboundGrpcToOutboundHttpStream, resp *http.Response) (retErr error) {
	defer errz.SafeClose(resp.Body, &retErr)
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
		return x.HandleSendError("SendMsg(HttpResponse_Header) failed", err)
	}

	buffer := memz.Get32k()
	defer memz.Put32k(buffer)
	for {
		n, err := resp.Body.Read(buffer)
		if err != nil && !errors.Is(err, io.EOF) {
			return status.Errorf(codes.Canceled, "read HTTP response body: %v", err)
		}
		if n > 0 { // handle n=0, err=io.EOF case
			sendErr := inbound.Send(&HttpResponse{
				Message: &HttpResponse_Data_{
					Data: &HttpResponse_Data{
						Data: buffer[:n],
					},
				},
			})
			if sendErr != nil {
				return x.HandleSendError("SendMsg(HttpResponse_Data) failed", sendErr)
			}
		}
		if errors.Is(err, io.EOF) {
			break
		}
	}
	return nil
}
