package agent

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/pluralsh/kuberentes-agent/internal/module/reverse_tunnel/info"
	"github.com/pluralsh/kuberentes-agent/internal/module/reverse_tunnel/rpc"
	"github.com/pluralsh/kuberentes-agent/internal/tool/grpctool"
	"github.com/pluralsh/kuberentes-agent/internal/tool/prototool"
	"github.com/pluralsh/kuberentes-agent/internal/tool/testing/matcher"
	"github.com/pluralsh/kuberentes-agent/internal/tool/testing/mock_reverse_tunnel_rpc"
	"github.com/pluralsh/kuberentes-agent/internal/tool/testing/mock_rpc"
	"github.com/pluralsh/kuberentes-agent/internal/tool/testing/testhelpers"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	_ connectionInterface = (*connection)(nil)
)

func TestPropagateUntilStop(t *testing.T) {
	ctxParent, cancelParent := context.WithCancel(context.Background())
	ctx, cancel, stop := propagateUntil(ctxParent)
	stop()
	// Let the Go runtime schedule the other goroutine.
	// It should exit so that this test doesn't flake.
	time.Sleep(10 * time.Millisecond)
	cancelParent()
	select {
	case <-ctx.Done():
		require.FailNow(t, "Unexpected context cancellation")
	default:
	}
	cancel()
	<-ctx.Done()
}

func TestPropagateUntilNoStop(t *testing.T) {
	ctxParent, cancelParent := context.WithCancel(context.Background())
	cancelParent()
	ctx, cancel, _ := propagateUntil(ctxParent)
	defer cancel()
	<-ctx.Done()
}

func TestConnectUnblocksIfNotStartedStreaming(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, _, _, c := setupConnection(t)

	client.EXPECT().
		Connect(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, opts ...grpc.CallOption) (rpc.ReverseTunnel_ConnectClient, error) {
			cancel()
			<-ctx.Done()
			return nil, ctx.Err()
		})

	err := c.attempt(ctx)
	require.EqualError(t, err, "Connect(): context canceled")
}

// Visitor can get io.EOF after getting rpc.RequestInfo if client sent an rpc.Error, which was forwarded to the tunnel
// and then the tunnel closed the stream.
func TestNoErrorOnEofAfterRequestInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, conn, tunnel, c := setupConnection(t)
	clientStream := mock_rpc.NewMockClientStream(ctrl)

	gomock.InOrder(
		clientStream.EXPECT().
			Header().
			Return(nil, errors.New("header err")),
		tunnel.EXPECT().Send(gomock.Any()),
		tunnel.EXPECT().CloseSend(),
	)

	gomock.InOrder(
		client.EXPECT().
			Connect(gomock.Any(), gomock.Any()).
			Return(tunnel, nil),
		tunnel.EXPECT().
			Send(gomock.Any()),
		tunnel.EXPECT().
			RecvMsg(gomock.Any()).
			Do(testhelpers.RecvMsg(&rpc.ConnectResponse{
				Msg: &rpc.ConnectResponse_RequestInfo{RequestInfo: &rpc.RequestInfo{}},
			})),
		conn.EXPECT().
			NewStream(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(clientStream, nil),
		tunnel.EXPECT().
			RecvMsg(gomock.Any()).
			Return(io.EOF),
	)

	err := c.attempt(ctx)
	require.NoError(t, err)
}

// Visitor can get io.EOF after getting rpc.Message if client sent an rpc.Error, which was forwarded to the tunnel
// and then the tunnel closed the stream.
func TestNoErrorOnEofAfterMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, conn, tunnel, c := setupConnection(t)
	clientStream := mock_rpc.NewMockClientStream(ctrl)

	gomock.InOrder(
		clientStream.EXPECT().
			Header().
			Return(nil, errors.New("header err")),
		tunnel.EXPECT().Send(gomock.Any()),
		tunnel.EXPECT().CloseSend(),
	)

	gomock.InOrder(
		client.EXPECT().
			Connect(gomock.Any(), gomock.Any()).
			Return(tunnel, nil),
		tunnel.EXPECT().
			Send(gomock.Any()),
		tunnel.EXPECT().
			RecvMsg(gomock.Any()).
			Do(testhelpers.RecvMsg(&rpc.ConnectResponse{
				Msg: &rpc.ConnectResponse_RequestInfo{RequestInfo: &rpc.RequestInfo{}},
			})),
		conn.EXPECT().
			NewStream(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(clientStream, nil),
		tunnel.EXPECT().
			RecvMsg(gomock.Any()).
			Do(testhelpers.RecvMsg(&rpc.ConnectResponse{
				Msg: &rpc.ConnectResponse_Message{
					Message: &rpc.Message{Data: []byte{1, 2, 3}},
				},
			})),
		clientStream.EXPECT().
			SendMsg(gomock.Any()),
		tunnel.EXPECT().
			RecvMsg(gomock.Any()).
			Return(io.EOF),
	)

	err := c.attempt(ctx)
	require.NoError(t, err)
}

func TestNoTrailerAfterHeaderError(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, conn, tunnel, c := setupConnection(t)
	clientStream := mock_rpc.NewMockClientStream(ctrl)

	done := make(chan struct{})

	headerErr := status.Error(codes.InvalidArgument, "expected header err")
	gomock.InOrder(
		clientStream.EXPECT().
			Header().
			Return(nil, headerErr),
		tunnel.EXPECT().
			Send(matcher.ProtoEq(nil, &rpc.ConnectRequest{
				Msg: &rpc.ConnectRequest_Error{
					Error: &rpc.Error{
						Status: status.Convert(headerErr).Proto(),
					},
				},
			})),
		tunnel.EXPECT().
			CloseSend().
			Do(func() error {
				close(done)
				return nil
			}),
	)

	gomock.InOrder(
		client.EXPECT().
			Connect(gomock.Any(), gomock.Any()).
			Return(tunnel, nil),
		tunnel.EXPECT().
			Send(gomock.Any()),
		tunnel.EXPECT().
			RecvMsg(gomock.Any()).
			Do(testhelpers.RecvMsg(&rpc.ConnectResponse{
				Msg: &rpc.ConnectResponse_RequestInfo{RequestInfo: &rpc.RequestInfo{}},
			})),
		conn.EXPECT().
			NewStream(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(clientStream, nil),
		tunnel.EXPECT().
			RecvMsg(gomock.Any()).
			DoAndReturn(func(m interface{}) error {
				<-done
				return io.EOF
			}),
	)

	err := c.attempt(ctx)
	require.NoError(t, err)
}

func TestTrailerAfterRecvMsgEof(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, conn, tunnel, c := setupConnection(t)
	clientStream := mock_rpc.NewMockClientStream(ctrl)

	done := make(chan struct{})

	gomock.InOrder(
		clientStream.EXPECT().
			Header(),
		tunnel.EXPECT().
			Send(matcher.ProtoEq(nil, &rpc.ConnectRequest{
				Msg: &rpc.ConnectRequest_Header{Header: &rpc.Header{}},
			})),
		clientStream.EXPECT().
			RecvMsg(gomock.Any()).
			Return(io.EOF),
		clientStream.EXPECT().
			Trailer().
			Return(metadata.MD{"abc": []string{"a", "b"}}),
		tunnel.EXPECT().
			Send(matcher.ProtoEq(nil, &rpc.ConnectRequest{
				Msg: &rpc.ConnectRequest_Trailer{
					Trailer: &rpc.Trailer{
						Meta: map[string]*prototool.Values{
							"abc": {Value: []string{"a", "b"}},
						},
					},
				},
			})),
		tunnel.EXPECT().
			CloseSend().
			Do(func() error {
				close(done)
				return nil
			}),
	)

	gomock.InOrder(
		client.EXPECT().
			Connect(gomock.Any(), gomock.Any()).
			Return(tunnel, nil),
		tunnel.EXPECT().
			Send(gomock.Any()),
		tunnel.EXPECT().
			RecvMsg(gomock.Any()).
			Do(testhelpers.RecvMsg(&rpc.ConnectResponse{
				Msg: &rpc.ConnectResponse_RequestInfo{RequestInfo: &rpc.RequestInfo{}},
			})),
		conn.EXPECT().
			NewStream(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(clientStream, nil),
		tunnel.EXPECT().
			RecvMsg(gomock.Any()).
			DoAndReturn(func(m interface{}) error {
				<-done
				return io.EOF
			}),
	)

	err := c.attempt(ctx)
	require.NoError(t, err)
}

func TestTrailerAndErrorAfterRecvMsgError(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, conn, tunnel, c := setupConnection(t)
	clientStream := mock_rpc.NewMockClientStream(ctrl)

	done := make(chan struct{})

	recvErr := status.Error(codes.InvalidArgument, "expected RecvMsg err")
	gomock.InOrder(
		clientStream.EXPECT().
			Header(),
		tunnel.EXPECT().
			Send(matcher.ProtoEq(nil, &rpc.ConnectRequest{
				Msg: &rpc.ConnectRequest_Header{Header: &rpc.Header{}},
			})),
		clientStream.EXPECT().
			RecvMsg(gomock.Any()).
			Return(recvErr),
		clientStream.EXPECT().
			Trailer().
			Return(metadata.MD{"abc": []string{"a", "b"}}),
		tunnel.EXPECT().
			Send(matcher.ProtoEq(nil, &rpc.ConnectRequest{
				Msg: &rpc.ConnectRequest_Trailer{
					Trailer: &rpc.Trailer{
						Meta: map[string]*prototool.Values{
							"abc": {Value: []string{"a", "b"}},
						},
					},
				},
			})),
		tunnel.EXPECT().
			Send(matcher.ProtoEq(nil, &rpc.ConnectRequest{
				Msg: &rpc.ConnectRequest_Error{
					Error: &rpc.Error{
						Status: status.Convert(recvErr).Proto(),
					},
				},
			})),
		tunnel.EXPECT().
			CloseSend().
			Do(func() error {
				close(done)
				return nil
			}),
	)

	gomock.InOrder(
		client.EXPECT().
			Connect(gomock.Any(), gomock.Any()).
			Return(tunnel, nil),
		tunnel.EXPECT().
			Send(gomock.Any()),
		tunnel.EXPECT().
			RecvMsg(gomock.Any()).
			Do(testhelpers.RecvMsg(&rpc.ConnectResponse{
				Msg: &rpc.ConnectResponse_RequestInfo{RequestInfo: &rpc.RequestInfo{}},
			})),
		conn.EXPECT().
			NewStream(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(clientStream, nil),
		tunnel.EXPECT().
			RecvMsg(gomock.Any()).
			DoAndReturn(func(m interface{}) error {
				<-done
				return io.EOF
			}),
	)

	err := c.attempt(ctx)
	require.NoError(t, err)
}

func TestRecvMsgUnblocksIfNotStartedStreaming(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, _, tunnel, c := setupConnection(t)

	var connectCtx context.Context

	gomock.InOrder(
		client.EXPECT().
			Connect(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, opts ...grpc.CallOption) (rpc.ReverseTunnel_ConnectClient, error) {
				connectCtx = ctx
				return tunnel, nil
			}),
		tunnel.EXPECT().
			Send(gomock.Any()),
		tunnel.EXPECT().
			RecvMsg(gomock.Any()).
			DoAndReturn(func(m interface{}) error {
				cancel()
				<-connectCtx.Done()
				return connectCtx.Err()
			}),
	)

	err := c.attempt(ctx)
	require.EqualError(t, err, "context canceled")
}

func TestContextIgnoredIfStartedStreaming(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, conn, tunnel, c := setupConnection(t)
	clientStream := mock_rpc.NewMockClientStream(ctrl)

	gomock.InOrder(
		clientStream.EXPECT().
			Header().
			Return(nil, errors.New("header err")),
		tunnel.EXPECT().Send(gomock.Any()),
		tunnel.EXPECT().CloseSend(),
	)
	var connectCtx context.Context

	gomock.InOrder(
		client.EXPECT().
			Connect(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, opts ...grpc.CallOption) (rpc.ReverseTunnel_ConnectClient, error) {
				connectCtx = ctx
				return tunnel, nil
			}),
		tunnel.EXPECT().
			Send(gomock.Any()),
		tunnel.EXPECT().
			RecvMsg(gomock.Any()).
			Do(testhelpers.RecvMsg(&rpc.ConnectResponse{
				Msg: &rpc.ConnectResponse_RequestInfo{RequestInfo: &rpc.RequestInfo{}},
			})),
		conn.EXPECT().
			NewStream(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(clientStream, nil),
		tunnel.EXPECT().
			RecvMsg(gomock.Any()).
			DoAndReturn(func(m interface{}) error {
				cancel()
				select {
				case <-connectCtx.Done():
					require.FailNow(t, "Unexpected context cancellation")
				default:
				}
				return errors.New("expected err")
			}),
	)

	err := c.attempt(ctx)
	require.EqualError(t, err, "expected err")
}

func TestAgentDescriptorIsSent(t *testing.T) {
	client, _, tunnel, c := setupConnection(t)
	gomock.InOrder(
		client.EXPECT().
			Connect(gomock.Any(), gomock.Any()).
			Return(tunnel, nil),
		tunnel.EXPECT().
			Send(matcher.ProtoEq(t, &rpc.ConnectRequest{
				Msg: &rpc.ConnectRequest_Descriptor_{
					Descriptor_: &rpc.Descriptor{
						AgentDescriptor: descriptor(),
					},
				},
			})).
			Return(errors.New("expected err")),
	)
	err := c.attempt(context.Background())
	require.EqualError(t, err, "Send(descriptor): expected err")
}

func TestAttemptIsUnblockedOnTunnelRecvMessageError(t *testing.T) {
	client, conn, tunnel, c := setupConnection(t)
	ctrl := gomock.NewController(t)
	clientStream := mock_rpc.NewMockClientStream(ctrl)
	var newStreamCtx context.Context
	gomock.InOrder(
		client.EXPECT().
			Connect(gomock.Any(), gomock.Any()).
			Return(tunnel, nil),
		tunnel.EXPECT().
			Send(gomock.Any()), // ConnectRequest_Descriptor_
		tunnel.EXPECT().
			RecvMsg(gomock.Any()).
			Do(testhelpers.RecvMsg(&rpc.ConnectResponse{
				Msg: &rpc.ConnectResponse_RequestInfo{RequestInfo: &rpc.RequestInfo{}},
			})),
		conn.EXPECT().
			NewStream(gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
				newStreamCtx = ctx
				return clientStream, nil
			}),
		tunnel.EXPECT().
			RecvMsg(gomock.Any()).
			Return(errors.New("expected recv error")),
	)
	gomock.InOrder(
		clientStream.EXPECT().
			Header(),
		tunnel.EXPECT().
			Send(matcher.ProtoEq(nil, &rpc.ConnectRequest{
				Msg: &rpc.ConnectRequest_Header{Header: &rpc.Header{}},
			})),
		clientStream.EXPECT().
			RecvMsg(gomock.Any()).
			DoAndReturn(func(m interface{}) error {
				<-newStreamCtx.Done() // block until context is cancelled
				return newStreamCtx.Err()
			}),
		clientStream.EXPECT().
			Trailer(),
		tunnel.EXPECT().
			Send(gomock.Any()).
			Return(errors.New("expected send error")),
	)

	err := c.attempt(context.Background())
	require.EqualError(t, err, "expected recv error")
}

func TestAttemptIsUnblockedOnTunnelHeaderSendError(t *testing.T) {
	client, conn, tunnel, c := setupConnection(t)
	ctrl := gomock.NewController(t)
	clientStream := mock_rpc.NewMockClientStream(ctrl)
	gomock.InOrder(
		client.EXPECT().
			Connect(gomock.Any(), gomock.Any()).
			Return(tunnel, nil),
		tunnel.EXPECT().
			Send(gomock.Any()), // ConnectRequest_Descriptor_
		tunnel.EXPECT().
			RecvMsg(gomock.Any()).
			Do(testhelpers.RecvMsg(&rpc.ConnectResponse{
				Msg: &rpc.ConnectResponse_RequestInfo{
					RequestInfo: &rpc.RequestInfo{},
				},
			})),
		conn.EXPECT().
			NewStream(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(clientStream, nil),
		tunnel.EXPECT().
			RecvMsg(gomock.Any()).
			Return(errors.New("expected recv error")),
	)
	gomock.InOrder(
		clientStream.EXPECT().
			Header(),
		tunnel.EXPECT().
			Send(matcher.ProtoEq(nil, &rpc.ConnectRequest{
				Msg: &rpc.ConnectRequest_Header{Header: &rpc.Header{}},
			})).
			Return(io.EOF),
	)

	err := c.attempt(context.Background())
	require.EqualError(t, err, "expected recv error")
}

func setupConnection(t *testing.T) (*mock_reverse_tunnel_rpc.MockReverseTunnelClient, *mock_rpc.MockClientConnInterface, *mock_reverse_tunnel_rpc.MockReverseTunnel_ConnectClient, *connection) {
	ctrl := gomock.NewController(t)
	client := mock_reverse_tunnel_rpc.NewMockReverseTunnelClient(ctrl)
	conn := mock_rpc.NewMockClientConnInterface(ctrl)
	tunnel := mock_reverse_tunnel_rpc.NewMockReverseTunnel_ConnectClient(ctrl)
	sv, err := grpctool.NewStreamVisitor(&rpc.ConnectResponse{})
	require.NoError(t, err)
	c := &connection{
		log:                zaptest.NewLogger(t),
		descriptor:         descriptor(),
		client:             client,
		internalServerConn: conn,
		streamVisitor:      sv,
		onIdle:             func(c connectionInterface) {},
		onActive:           func(c connectionInterface) {},
	}
	return client, conn, tunnel, c
}

func descriptor() *info.AgentDescriptor {
	return &info.AgentDescriptor{
		Services: []*info.Service{
			{
				Name: "bla",
				Methods: []*info.Method{
					{
						Name: "bab",
					},
				},
			},
		},
	}
}
