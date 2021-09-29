package grpctool

import (
	"context"
	"errors"
	"io"
	"net/http"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/errz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/prototool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	headerFieldNumber  protoreflect.FieldNumber = 1
	dataFieldNumber    protoreflect.FieldNumber = 2
	trailerFieldNumber protoreflect.FieldNumber = 3

	maxDataChunkSize = 32 * 1024
)

type InboundGrpcToOutboundHttpStream interface {
	Send(*HttpResponse) error
	grpc.ServerStream
}

// RpcApi is a reduced version on modshared.RpcApi.
// It's here to avoid the dependency.
type RpcApi interface {
	Log() *zap.Logger
	HandleProcessingError(log *zap.Logger, agentId int64, msg string, err error)
	HandleSendError(log *zap.Logger, msg string, err error) error
}

type HttpDo func(ctx context.Context, header *HttpRequest_Header, body io.Reader) (*http.Response, error)

type InboundGrpcToOutboundHttp struct {
	streamVisitor *StreamVisitor
	httpDo        HttpDo
}

func NewInboundGrpcToOutboundHttp(httpDo HttpDo) *InboundGrpcToOutboundHttp {
	sv, err := NewStreamVisitor(&HttpRequest{})
	if err != nil {
		panic(err) // this will never panic as long as the proto file is correct
	}
	return &InboundGrpcToOutboundHttp{
		streamVisitor: sv,
		httpDo:        httpDo,
	}
}

func (x *InboundGrpcToOutboundHttp) Pipe(rpcApi RpcApi, server InboundGrpcToOutboundHttpStream, agentId int64) error {
	ctx := server.Context()

	pr, pw := io.Pipe()
	headerMsg := make(chan *HttpRequest_Header)
	s := InboundGrpcToOutboundStream{
		// Pipe gRPC request -> HTTP request
		PipeInboundToOutbound: func() error {
			return x.pipeGrpcIntoHttp(ctx, server, headerMsg, pw)
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
				resp, err := x.httpDo(ctx, header, pr)
				if err != nil {
					return err
				}
				return x.pipeHttpIntoGrpc(rpcApi, server, resp)
			}
		},
	}
	err := s.Pipe()
	switch {
	case err == nil:
	case IsStatusError(err):
		// A gRPC status already
	case errors.Is(err, context.Canceled):
		rpcApi.Log().Debug("gRPC -> HTTP", logz.Error(err))
		err = status.Error(codes.Canceled, err.Error())
	case errors.Is(err, context.DeadlineExceeded):
		rpcApi.Log().Debug("gRPC -> HTTP", logz.Error(err))
		err = status.Error(codes.DeadlineExceeded, err.Error())
	default:
		rpcApi.HandleProcessingError(rpcApi.Log(), agentId, "gRPC -> HTTP", err)
		err = status.Errorf(codes.Unavailable, "gRPC -> HTTP: %v", err)
	}
	return err
}

func (x *InboundGrpcToOutboundHttp) pipeGrpcIntoHttp(ctx context.Context, server grpc.ServerStream, headerMsg chan *HttpRequest_Header, pw *io.PipeWriter) error {
	return x.streamVisitor.Visit(server,
		WithCallback(headerFieldNumber, func(header *HttpRequest_Header) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case headerMsg <- header:
				return nil
			}
		}),
		WithCallback(dataFieldNumber, func(data *HttpRequest_Data) error {
			_, err := pw.Write(data.Data)
			return err
		}),
		WithCallback(trailerFieldNumber, func(trailer *HttpRequest_Trailer) error {
			// Nothing to do
			return nil
		}),
		WithEOFCallback(pw.Close), // Nothing more to send, close the write end of the pipe
	)
}

func (x *InboundGrpcToOutboundHttp) pipeHttpIntoGrpc(rpcApi RpcApi, server grpc.ServerStream, resp *http.Response) error {
	err := func() (retErr error) { // closure to close resp.Body ASAP
		defer errz.SafeClose(resp.Body, &retErr)
		err := server.SendMsg(&HttpResponse{
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
			return rpcApi.HandleSendError(rpcApi.Log(), "SendMsg(HttpResponse_Header) failed", err)
		}

		buffer := make([]byte, maxDataChunkSize)
		for err == nil { // loop while not EOF
			var n int
			n, err = resp.Body.Read(buffer)
			if err != nil && !errors.Is(err, io.EOF) {
				return status.Errorf(codes.Canceled, "read HTTP response body: %v", err)
			}
			if n > 0 { // handle n=0, err=io.EOF case
				sendErr := server.SendMsg(&HttpResponse{
					Message: &HttpResponse_Data_{
						Data: &HttpResponse_Data{
							Data: buffer[:n],
						},
					},
				})
				if sendErr != nil {
					return rpcApi.HandleSendError(rpcApi.Log(), "SendMsg(HttpResponse_Data) failed", sendErr)
				}
			}
		}
		return nil
	}()
	if err != nil {
		return err
	}

	err = server.SendMsg(&HttpResponse{
		Message: &HttpResponse_Trailer_{
			Trailer: &HttpResponse_Trailer{},
		},
	})
	if err != nil {
		return rpcApi.HandleSendError(rpcApi.Log(), "SendMsg(HttpResponse_Trailer) failed", err)
	}
	return nil
}
