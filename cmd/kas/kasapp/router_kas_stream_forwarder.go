package kasapp

import (
	"errors"
	"io"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/grpctool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	tunnelReadyFieldNumber protoreflect.FieldNumber = 1
	headerFieldNumber      protoreflect.FieldNumber = 2
	messageFieldNumber     protoreflect.FieldNumber = 3
	trailerFieldNumber     protoreflect.FieldNumber = 4
	errorFieldNumber       protoreflect.FieldNumber = 5
)

type kasStreamForwarder struct {
	log               *zap.Logger
	rpcApi            modserver.RpcApi
	gatewayKasVisitor *grpctool.StreamVisitor
}

// ForwardStream does bi-directional stream forwarding.
// Returns a gRPC status-compatible error.
func (f *kasStreamForwarder) ForwardStream(kasStream grpc.ClientStream, stream grpc.ServerStream) error {
	// Cancellation
	//
	// kasStream is an outbound client stream (this/routing kas -> gateway kas)
	// stream is an inbound server stream (internal/external gRPC client -> this/routing kas)
	//
	// If one of the streams breaks, the other one needs to be aborted too ASAP. Waiting for a timeout
	// is a waste of resources and a bad API with unpredictable latency.
	//
	// kasStream is automatically aborted if there is a problem with stream because kasStream uses stream's context.
	// Unlike the above, if there is a problem with kasStream, stream.RecvMsg()/stream.SendMsg() are unaffected
	// so can stay blocked for an arbitrary amount of time.
	// To make gRPC abort those method calls, gRPC stream handler (i.e. this method) should just return from the call.
	// See https://github.com/grpc/grpc-go/issues/465#issuecomment-179414474
	// To implement this, we read and write in both directions in separate goroutines and return from this
	// handler whenever there is an error, aborting both connections.

	// Channel of size 1 to ensure that if we return early, the other goroutine has space for the value.
	// We don't care about that value if we already got a non-nil error.
	res := make(chan error, 1)
	go func() {
		res <- f.pipeFromStreamToKas(kasStream, stream)
	}()
	err := f.pipeFromKasToStream(kasStream, stream)
	if err != nil {
		return err // unblocks stream.RecvMsg() in the other goroutine if it is stuck
	}
	// Wait for the other goroutine to return to cleanly finish reading from stream
	return <-res
}

func (f *kasStreamForwarder) pipeFromKasToStream(kasStream grpc.ClientStream, stream grpc.ServerStream) error {
	var statusFromKasStream error
	err := f.gatewayKasVisitor.Visit(kasStream,
		grpctool.WithStartState(tunnelReadyFieldNumber),
		grpctool.WithCallback(headerFieldNumber, func(header *GatewayKasResponse_Header) error {
			err := stream.SetHeader(header.Metadata())
			if err != nil {
				return f.rpcApi.HandleSendError(f.log, "router kas->stream SetHeader() failed", err)
			}
			return nil
		}),
		grpctool.WithCallback(messageFieldNumber, func(message *GatewayKasResponse_Message) error {
			err := stream.SendMsg(&grpctool.RawFrame{
				Data: message.Data,
			})
			if err != nil {
				return f.rpcApi.HandleSendError(f.log, "router kas->stream SendMsg() failed", err)
			}
			return nil
		}),
		grpctool.WithCallback(trailerFieldNumber, func(trailer *GatewayKasResponse_Trailer) error {
			stream.SetTrailer(trailer.Metadata())
			return nil
		}),
		grpctool.WithCallback(errorFieldNumber, func(err *GatewayKasResponse_Error) error {
			statusFromKasStream = status.ErrorProto(err.Status)
			return nil
		}),
	)
	if err != nil {
		return err
	}
	return statusFromKasStream
}

// pipeFromStreamToKas pipes data kasStream -> stream
// must return gRPC status compatible error or nil.
func (f *kasStreamForwarder) pipeFromStreamToKas(kasStream grpc.ClientStream, stream grpc.ServerStream) error {
	var frame grpctool.RawFrame
	for {
		err := stream.RecvMsg(&frame)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		err = kasStream.SendMsg(&frame)
		if err != nil {
			if errors.Is(err, io.EOF) { // the other goroutine will receive the error in RecvMsg()
				return nil
			}
			return f.rpcApi.HandleSendError(f.log, "stream->router kas SendMsg() failed", err)
		}
	}
	err := kasStream.CloseSend()
	if err != nil {
		return f.rpcApi.HandleSendError(f.log, "stream->router kas CloseSend() failed", err)
	}
	return nil
}
