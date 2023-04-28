package reverse_tunnel

import (
	"io"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/reverse_tunnel/info"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/reverse_tunnel/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/prototool"
	"go.uber.org/zap"
	statuspb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type stateType int

const (
	// zero value is invalid to catch initialization bugs.
	_ stateType = iota
	// stateReady - tunnel is owned by the registry and is ready to be found and used for forwarding.
	stateReady
	// stateFound - tunnel is not owned by registry, was found and about to be used for forwarding.
	stateFound
	// stateForwarding - tunnel is not owned by registry, is being used for forwarding.
	stateForwarding
	// stateDone - tunnel is not owned by anyone, it has been used for forwarding, Done() has been called.
	stateDone
	// stateContextDone - tunnel is not owned by anyone, reverse tunnel's context signalled done in HandleTunnel().
	stateContextDone
)

const (
	agentDescriptorNumber protoreflect.FieldNumber = 1
	headerNumber          protoreflect.FieldNumber = 2
	messageNumber         protoreflect.FieldNumber = 3
	trailerNumber         protoreflect.FieldNumber = 4
	errorNumber           protoreflect.FieldNumber = 5
)

type TunnelDataCallback interface {
	Header(map[string]*prototool.Values) error
	Message([]byte) error
	Trailer(map[string]*prototool.Values) error
	Error(*statuspb.Status) error
}

type RpcApi interface {
	HandleIoError(log *zap.Logger, msg string, err error) error
}

type Tunnel interface {
	// ForwardStream performs bi-directional message forwarding between incomingStream and the tunnel.
	// cb is called with header, messages and trailer coming from the tunnel. It's the callers
	// responsibility to forward them into the incomingStream.
	ForwardStream(log *zap.Logger, rpcApi RpcApi, incomingStream grpc.ServerStream, cb TunnelDataCallback) error
	// Done must be called when the caller is done with the Tunnel.
	Done()
}

type tunnel struct {
	tunnel              rpc.ReverseTunnel_ConnectServer
	tunnelStreamVisitor *grpctool.StreamVisitor
	tunnelRetErr        chan<- error
	agentId             int64
	agentDescriptor     *info.AgentDescriptor
	state               stateType

	onForward func(*tunnel) error
	onDone    func(*tunnel)
}

func (t *tunnel) ForwardStream(log *zap.Logger, rpcApi RpcApi, incomingStream grpc.ServerStream, cb TunnelDataCallback) error {
	if err := t.onForward(t); err != nil {
		return err
	}
	pair := t.forwardStream(log, rpcApi, incomingStream, cb)
	t.tunnelRetErr <- pair.forTunnel
	return pair.forIncomingStream
}

func (t *tunnel) forwardStream(log *zap.Logger, rpcApi RpcApi, incomingStream grpc.ServerStream, cb TunnelDataCallback) errPair {
	// Here we have a situation where we need to pipe one server stream into another server stream.
	// One stream is incoming request stream and the other one is incoming tunnel stream.
	// We need to use at least one extra goroutine in addition to the current one (or two separate ones) to
	// implement full duplex bidirectional stream piping. One goroutine reads and writes in one direction and the other
	// one in the opposite direction.
	// What if one of them returns an error? We need to unblock the other one, ideally ASAP, to release resources. If
	// it's not unblocked, it'll sit there until it hits a timeout or is aborted by peer. Ok-ish, but far from ideal.
	// To abort request processing on the server side, gRPC stream handler should just return from the call.
	// See https://github.com/grpc/grpc-go/issues/465#issuecomment-179414474
	// To implement this, we read and write in both directions in separate goroutines and return from both
	// handlers whenever there is an error, aborting both connections:
	// - Returning from this function means returning from the incoming request handler.
	// - Sending to c.tunnelRetErr leads to returning that value from the tunnel handler.

	// Channel of size 1 to ensure that if we return early, the second goroutine has space for the value.
	// We don't care about the second value if the first one has at least one non-nil error.
	res := make(chan errPair, 1)
	incomingCtx := incomingStream.Context()
	// Pipe incoming stream (i.e. data a client is sending us) into the tunnel stream
	goErrPair(res, func() (error /* forTunnel */, error /* forIncomingStream */) {
		md, _ := metadata.FromIncomingContext(incomingCtx)
		err := t.tunnel.Send(&rpc.ConnectResponse{
			Msg: &rpc.ConnectResponse_RequestInfo{
				RequestInfo: &rpc.RequestInfo{
					MethodName: grpc.ServerTransportStreamFromContext(incomingCtx).Method(),
					Meta:       grpctool.MetaToValuesMap(md),
				},
			},
		})
		if err != nil {
			err = rpcApi.HandleIoError(log, "Send(ConnectResponse_RequestInfo)", err)
			return err, err
		}
		// Outside the loop to allocate once vs on each message
		var frame grpctool.RawFrame
		var msg rpc.Message
		response := &rpc.ConnectResponse{
			Msg: &rpc.ConnectResponse_Message{
				Message: &msg,
			},
		}
		for {
			err = incomingStream.RecvMsg(&frame)
			if err != nil {
				if err == io.EOF { // nolint:errorlint
					break
				}
				return status.Error(codes.Canceled, "read from incoming stream"), err
			}
			msg.Data = frame.Data
			err = t.tunnel.Send(response)
			if err != nil {
				err = rpcApi.HandleIoError(log, "Send(ConnectResponse_Message)", err)
				return err, err
			}
		}
		err = t.tunnel.Send(&rpc.ConnectResponse{
			Msg: &rpc.ConnectResponse_CloseSend{
				CloseSend: &rpc.CloseSend{},
			},
		})
		if err != nil {
			err = rpcApi.HandleIoError(log, "Send(ConnectResponse_CloseSend)", err)
			return err, err
		}
		return nil, nil
	})
	// Pipe tunnel stream (i.e. data agentk is sending us) into the incoming stream
	goErrPair(res, func() (error /* forTunnel */, error /* forIncomingStream */) {
		var forTunnel, forIncomingStream error
		fromVisitor := t.tunnelStreamVisitor.Visit(t.tunnel,
			grpctool.WithStartState(agentDescriptorNumber),
			grpctool.WithCallback(headerNumber, func(header *rpc.Header) error {
				return cb.Header(header.Meta)
			}),
			grpctool.WithCallback(messageNumber, func(message *rpc.Message) error {
				return cb.Message(message.Data)
			}),
			grpctool.WithCallback(trailerNumber, func(trailer *rpc.Trailer) error {
				return cb.Trailer(trailer.Meta)
			}),
			grpctool.WithCallback(errorNumber, func(rpcError *rpc.Error) error {
				forIncomingStream = cb.Error(rpcError.Status)
				// Not returning an error since we must be reading from the tunnel stream till io.EOF
				// to properly consume it. There is no need to abort it in this scenario.
				// The server is expected to close the stream (i.e. we'll get io.EOF) right after we got this message.
				return nil
			}),
		)
		if fromVisitor != nil {
			forIncomingStream = fromVisitor
			forTunnel = fromVisitor
		}
		return forTunnel, forIncomingStream
	})
	pair := <-res
	if !pair.isNil() {
		return pair
	}
	select {
	case <-incomingCtx.Done():
		// incoming stream finished sending all data (i.e. io.EOF was read from it) but
		// now it signals that it's closing. We need to abort the potentially stuck t.tunnel.RecvMsg().
		err := grpctool.StatusErrorFromContext(incomingCtx, "Incoming stream closed")
		pair = errPair{
			forTunnel:         err,
			forIncomingStream: err,
		}
	case pair = <-res:
	}
	return pair
}

func (t *tunnel) Done() {
	t.onDone(t)
}

type errPair struct {
	forTunnel         error
	forIncomingStream error
}

func (p errPair) isNil() bool {
	return p.forTunnel == nil && p.forIncomingStream == nil
}

func goErrPair(c chan<- errPair, f func() (error /* forTunnel */, error /* forIncomingStream */)) {
	go func() {
		var pair errPair
		pair.forTunnel, pair.forIncomingStream = f()
		c <- pair
	}()
}
