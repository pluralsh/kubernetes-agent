package grpctool

type InboundGrpcToOutboundStream struct {
	PipeInboundToOutbound func() error
	PipeOutboundToInbound func() error
}

func (x *InboundGrpcToOutboundStream) Pipe() error {
	// Cancellation
	//
	// If one of the streams breaks, the other one needs to be aborted too ASAP. Waiting for a timeout
	// is a waste of resources and a bad API with unpredictable latency.
	//
	// The outbound stream is automatically aborted if there is a problem with inbound stream because
	// it uses the inbound stream's context.
	// Unlike the above, if there is a problem with the outbound stream, inbound.RecvMsg() in
	// PipeInboundToOutbound() are unaffected so can stay blocked for an arbitrary amount of time.
	// To make gRPC abort those method calls, gRPC stream handler (i.e. this method) should just return from the call.
	// See https://github.com/grpc/grpc-go/issues/465#issuecomment-179414474
	// To implement this, we read from the inbound stream in a separate goroutine and return from this
	// handler whenever there is an error, aborting reads from the incoming connection.

	// Channel of size 1 to ensure that if we return early, the other goroutine has space for the value.
	// We don't care about that value if we already got a non-nil error.
	res := make(chan error, 1)
	go func() {
		res <- x.PipeInboundToOutbound()
	}()
	err := x.PipeOutboundToInbound()
	if err != nil {
		return err // unblocks inbound.RecvMsg() in the other goroutine if it is stuck
	}
	// Wait for the other goroutine to return to cleanly finish reading from the inbound stream
	return <-res
}
