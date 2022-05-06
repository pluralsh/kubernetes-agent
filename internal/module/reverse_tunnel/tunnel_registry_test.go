package reverse_tunnel

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/reverse_tunnel/info"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/reverse_tunnel/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/prototool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/matcher"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_reverse_tunnel_rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_reverse_tunnel_tracker"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/testhelpers"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"k8s.io/apimachinery/pkg/util/wait"
)

// "slow" tests in this file are marked for concurrent execution with t.Parallel()

const (
	serviceName    = "gitlab.service1"
	methodName     = "DoSomething"
	fullMethodName = "/" + serviceName + "/" + methodName
)

var (
	_ TunnelHandler = &TunnelRegistry{}
	_ TunnelFinder  = &TunnelRegistry{}
)

func TestRunUnregistersAllConnections(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	connectServer := mock_reverse_tunnel_rpc.NewMockReverseTunnel_ConnectServer(ctrl)
	tunnelRegisterer := mock_reverse_tunnel_tracker.NewMockRegisterer(ctrl)
	gomock.InOrder(
		connectServer.EXPECT().
			Recv().
			Return(&rpc.ConnectRequest{
				Msg: &rpc.ConnectRequest_Descriptor_{
					Descriptor_: descriptor(),
				},
			}, nil),
		tunnelRegisterer.EXPECT().
			RegisterTunnel(gomock.Any(), gomock.Any()),
		tunnelRegisterer.EXPECT().
			UnregisterTunnel(gomock.Any(), gomock.Any()),
	)
	r, err := NewTunnelRegistry(zaptest.NewLogger(t), tunnelRegisterer, "grpc://127.0.0.1:123")
	require.NoError(t, err)
	var wg wait.Group
	defer wg.Wait()
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond) // will unblock Run()
	defer cancel()
	wg.Start(func() {
		assert.NoError(t, r.Run(ctx))
	})
	err = r.HandleTunnel(context.Background(), testhelpers.AgentInfoObj(), connectServer)
	assert.NoError(t, err)
}

func TestHandleTunnelIsUnblockedByContext(t *testing.T) {
	t.Parallel()
	ctxConn, cancelConn := context.WithTimeout(context.Background(), 50*time.Millisecond) // will unblock HandleTunnel()
	defer cancelConn()

	ctrl := gomock.NewController(t)
	connectServer := mock_reverse_tunnel_rpc.NewMockReverseTunnel_ConnectServer(ctrl)
	tunnelRegisterer := mock_reverse_tunnel_tracker.NewMockRegisterer(ctrl)
	gomock.InOrder(
		connectServer.EXPECT().
			Recv().
			Return(&rpc.ConnectRequest{
				Msg: &rpc.ConnectRequest_Descriptor_{
					Descriptor_: descriptor(),
				},
			}, nil),
		tunnelRegisterer.EXPECT().
			RegisterTunnel(gomock.Any(), gomock.Any()),
		tunnelRegisterer.EXPECT().
			UnregisterTunnel(gomock.Any(), gomock.Any()),
	)
	r, err := NewTunnelRegistry(zaptest.NewLogger(t), tunnelRegisterer, "grpc://127.0.0.1:123")
	require.NoError(t, err)
	var wg wait.Group
	defer wg.Wait()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg.Start(func() {
		assert.NoError(t, r.Run(ctx))
	})
	err = r.HandleTunnel(ctxConn, testhelpers.AgentInfoObj(), connectServer)
	assert.NoError(t, err)
}

// Two tunnels with the same agent id. Both register. Then one of them is retrieved via FindTunnel()
// and then its context is cancelled. If this test gets stuck, we have a problem.
// Reproducer for https://gitlab.com/gitlab-org/cluster-integration/gitlab-agent/-/issues/183.
func TestHandleTunnelIsUnblockedByContext_WithTwoTunnels(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	connectServer1 := mock_reverse_tunnel_rpc.NewMockReverseTunnel_ConnectServer(ctrl)
	connectServer2 := mock_reverse_tunnel_rpc.NewMockReverseTunnel_ConnectServer(ctrl)
	tunnelRegisterer := mock_reverse_tunnel_tracker.NewMockRegisterer(ctrl)
	gomock.InOrder(
		connectServer1.EXPECT().
			Recv().
			Return(&rpc.ConnectRequest{
				Msg: &rpc.ConnectRequest_Descriptor_{
					Descriptor_: descriptor(),
				},
			}, nil),
		tunnelRegisterer.EXPECT().
			RegisterTunnel(gomock.Any(), gomock.Any()),
		tunnelRegisterer.EXPECT().
			UnregisterTunnel(gomock.Any(), gomock.Any()),
	)
	gomock.InOrder(
		connectServer2.EXPECT().
			Recv().
			Return(&rpc.ConnectRequest{
				Msg: &rpc.ConnectRequest_Descriptor_{
					Descriptor_: descriptor(),
				},
			}, nil),
		tunnelRegisterer.EXPECT().
			RegisterTunnel(gomock.Any(), gomock.Any()),
		tunnelRegisterer.EXPECT().
			UnregisterTunnel(gomock.Any(), gomock.Any()),
	)
	r, err := NewTunnelRegistry(zaptest.NewLogger(t), tunnelRegisterer, "grpc://127.0.0.1:123")
	require.NoError(t, err)
	var wg wait.Group
	defer wg.Wait()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg.Start(func() {
		assert.NoError(t, r.Run(ctx))
	})
	agentInfo := testhelpers.AgentInfoObj()
	ctx1, cancel1 := context.WithCancel(context.Background())
	defer cancel1()
	wg.Start(func() {
		assert.NoError(t, r.HandleTunnel(ctx1, agentInfo, connectServer1))
	})
	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()
	wg.Start(func() {
		assert.NoError(t, r.HandleTunnel(ctx2, agentInfo, connectServer2))
	})
	time.Sleep(50 * time.Millisecond) // wait for both to register
	tun, err := r.FindTunnel(context.Background(), agentInfo.Id, serviceName, methodName)
	require.NoError(t, err)
	// cancel context for the found tunnel
	switch tun.(*tunnel).tunnel {
	case connectServer1:
		cancel1()
	case connectServer2:
		cancel2()
	default:
		t.FailNow()
	}
	tun.Done()
}

func TestHandleTunnelReturnErrOnRecvErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	connectServer := mock_reverse_tunnel_rpc.NewMockReverseTunnel_ConnectServer(ctrl)
	connectServer.EXPECT().
		Recv().
		Return(nil, status.Error(codes.DataLoss, "expected err"))
	tunnelRegisterer := mock_reverse_tunnel_tracker.NewMockRegisterer(ctrl)
	r, err := NewTunnelRegistry(zaptest.NewLogger(t), tunnelRegisterer, "grpc://127.0.0.1:123")
	require.NoError(t, err)
	err = r.HandleTunnel(context.Background(), testhelpers.AgentInfoObj(), connectServer)
	assert.EqualError(t, err, "rpc error: code = DataLoss desc = expected err")
}

func TestHandleTunnelReturnErrOnInvalidMsg(t *testing.T) {
	ctrl := gomock.NewController(t)
	connectServer := mock_reverse_tunnel_rpc.NewMockReverseTunnel_ConnectServer(ctrl)
	connectServer.EXPECT().
		Recv().
		Return(&rpc.ConnectRequest{
			Msg: &rpc.ConnectRequest_Header{
				Header: &rpc.Header{},
			},
		}, nil)
	tunnelRegisterer := mock_reverse_tunnel_tracker.NewMockRegisterer(ctrl)
	r, err := NewTunnelRegistry(zaptest.NewLogger(t), tunnelRegisterer, "grpc://127.0.0.1:123")
	require.NoError(t, err)
	err = r.HandleTunnel(context.Background(), testhelpers.AgentInfoObj(), connectServer)
	assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = Invalid oneof value type: *rpc.ConnectRequest_Header")
}

func TestHandleTunnelIsMatchedToIncomingConnection(t *testing.T) {
	t.Parallel()
	incomingStream, rpcApi, cb, tunnel, r := setupStreams(t, true)
	agentInfo := testhelpers.AgentInfoObj()
	var wg wait.Group
	defer wg.Wait()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg.Start(func() {
		assert.NoError(t, r.Run(ctx))
	})
	wg.Start(func() {
		assert.NoError(t, r.HandleTunnel(context.Background(), agentInfo, tunnel))
	})
	time.Sleep(50 * time.Millisecond)
	tun, err := r.FindTunnel(context.Background(), agentInfo.Id, serviceName, methodName)
	require.NoError(t, err)
	defer tun.Done()
	err = tun.ForwardStream(zaptest.NewLogger(t), rpcApi, incomingStream, cb)
	require.NoError(t, err)
}

func TestHandleTunnelIsNotMatchedToIncomingConnectionForMissingMethod(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	tunnelRegisterer := mock_reverse_tunnel_tracker.NewMockRegisterer(ctrl)
	connectServer := mock_reverse_tunnel_rpc.NewMockReverseTunnel_ConnectServer(ctrl)
	connectServer.EXPECT().
		Recv().
		Return(&rpc.ConnectRequest{
			Msg: &rpc.ConnectRequest_Descriptor_{
				Descriptor_: descriptor(),
			},
		}, nil)
	gomock.InOrder(
		tunnelRegisterer.EXPECT().
			RegisterTunnel(gomock.Any(), gomock.Any()),
		tunnelRegisterer.EXPECT().
			UnregisterTunnel(gomock.Any(), gomock.Any()),
	)
	r, err := NewTunnelRegistry(zaptest.NewLogger(t), tunnelRegisterer, "grpc://127.0.0.1:123")
	require.NoError(t, err)
	agentInfo := testhelpers.AgentInfoObj()
	var wg wait.Group
	defer wg.Wait()
	ctx1, cancel1 := context.WithCancel(context.Background())
	defer cancel1()
	wg.Start(func() {
		assert.NoError(t, r.Run(ctx1))
	})
	wg.Start(func() {
		assert.NoError(t, r.HandleTunnel(context.Background(), agentInfo, connectServer))
	})
	time.Sleep(50 * time.Millisecond)
	ctx2, cancel2 := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel2()
	_, err = r.FindTunnel(ctx2, agentInfo.Id, "missing_service", "missing_method")
	assert.Equal(t, context.DeadlineExceeded, err)
}

func TestForwardStreamIsMatchedToHandleTunnel(t *testing.T) {
	t.Parallel()
	incomingStream, rpcApi, cb, tunnel, r := setupStreams(t, false)
	agentInfo := testhelpers.AgentInfoObj()
	var wg wait.Group
	defer wg.Wait()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg.Start(func() {
		assert.NoError(t, r.Run(ctx))
	})
	wg.Start(func() {
		tun, err := r.FindTunnel(context.Background(), agentInfo.Id, serviceName, methodName)
		if !assert.NoError(t, err) {
			return
		}
		defer tun.Done()
		err = tun.ForwardStream(zaptest.NewLogger(t), rpcApi, incomingStream, cb)
		assert.NoError(t, err)
	})
	time.Sleep(50 * time.Millisecond)
	err := r.HandleTunnel(context.Background(), agentInfo, tunnel)
	require.NoError(t, err)
}

func TestForwardStreamIsNotMatchedToHandleTunnelForMissingMethod(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	tunnelRegisterer := mock_reverse_tunnel_tracker.NewMockRegisterer(ctrl)
	connectServer := mock_reverse_tunnel_rpc.NewMockReverseTunnel_ConnectServer(ctrl)
	connectServer.EXPECT().
		Recv().
		Return(&rpc.ConnectRequest{
			Msg: &rpc.ConnectRequest_Descriptor_{
				Descriptor_: descriptor(),
			},
		}, nil)
	gomock.InOrder(
		tunnelRegisterer.EXPECT().
			RegisterTunnel(gomock.Any(), gomock.Any()),
		tunnelRegisterer.EXPECT().
			UnregisterTunnel(gomock.Any(), gomock.Any()),
	)
	r, err := NewTunnelRegistry(zaptest.NewLogger(t), tunnelRegisterer, "grpc://127.0.0.1:123")
	require.NoError(t, err)
	agentInfo := testhelpers.AgentInfoObj()
	var wg wait.Group
	defer wg.Wait()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg.Start(func() {
		assert.NoError(t, r.Run(ctx))
	})
	wg.Start(func() {
		_, findErr := r.FindTunnel(context.Background(), agentInfo.Id, "missing_service", "missing_method")
		assert.Equal(t, context.Canceled, findErr)
	})
	time.Sleep(50 * time.Millisecond)
	ctx2, cancel2 := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel2()
	err = r.HandleTunnel(ctx2, agentInfo, connectServer)
	assert.NoError(t, err)
}

func setupStreams(t *testing.T, expectRegisterTunnel bool) (*mock_rpc.MockServerStream, *mock_modserver.MockAgentRpcApi, *MockTunnelDataCallback, *mock_reverse_tunnel_rpc.MockReverseTunnel_ConnectServer, *TunnelRegistry) {
	const metaKey = "Cba"
	meta := metadata.MD{}
	meta.Set(metaKey, "3", "4")
	ctrl := gomock.NewController(t)
	sts := mock_rpc.NewMockServerTransportStream(ctrl)
	cb := NewMockTunnelDataCallback(ctrl)

	rpcApi := mock_modserver.NewMockAgentRpcApi(ctrl)
	incomingCtx := grpc.NewContextWithServerTransportStream(context.Background(), sts)
	incomingCtx = metadata.NewIncomingContext(incomingCtx, meta)
	incomingStream := mock_rpc.NewMockServerStream(ctrl)
	incomingStream.EXPECT().
		Context().
		Return(incomingCtx).
		MinTimes(1)

	tunnelRegisterer := mock_reverse_tunnel_tracker.NewMockRegisterer(ctrl)
	connectServer := mock_reverse_tunnel_rpc.NewMockReverseTunnel_ConnectServer(ctrl)
	connectServer.EXPECT().
		Recv().
		Return(&rpc.ConnectRequest{
			Msg: &rpc.ConnectRequest_Descriptor_{
				Descriptor_: descriptor(),
			},
		}, nil)
	if expectRegisterTunnel {
		gomock.InOrder(
			tunnelRegisterer.EXPECT().
				RegisterTunnel(gomock.Any(), gomock.Any()),
			tunnelRegisterer.EXPECT().
				UnregisterTunnel(gomock.Any(), gomock.Any()),
		)
	}
	frame := grpctool.RawFrame{
		Data: []byte{1, 2, 3},
	}
	gomock.InOrder(
		sts.EXPECT().
			Method().
			Return(fullMethodName).
			MinTimes(1),
		connectServer.EXPECT().
			Send(matcher.ProtoEq(t, &rpc.ConnectResponse{
				Msg: &rpc.ConnectResponse_RequestInfo{
					RequestInfo: &rpc.RequestInfo{
						MethodName: fullMethodName,
						Meta: map[string]*prototool.Values{
							"cba": {Value: []string{"3", "4"}},
						},
					},
				},
			})),
		incomingStream.EXPECT().
			RecvMsg(gomock.Any()).
			Do(testhelpers.RecvMsg(&frame)),
		connectServer.EXPECT().
			Send(matcher.ProtoEq(t, &rpc.ConnectResponse{
				Msg: &rpc.ConnectResponse_Message{
					Message: &rpc.Message{
						Data: frame.Data,
					},
				},
			})),
		incomingStream.EXPECT().
			RecvMsg(gomock.Any()).
			Return(io.EOF),
		connectServer.EXPECT().
			Send(matcher.ProtoEq(t, &rpc.ConnectResponse{
				Msg: &rpc.ConnectResponse_CloseSend{
					CloseSend: &rpc.CloseSend{},
				},
			})),
	)
	gomock.InOrder(
		connectServer.EXPECT().
			RecvMsg(gomock.Any()).
			Do(testhelpers.RecvMsg(&rpc.ConnectRequest{
				Msg: &rpc.ConnectRequest_Header{
					Header: &rpc.Header{
						Meta: map[string]*prototool.Values{
							"resp": {Value: []string{"1", "2"}},
						},
					},
				},
			})),
		cb.EXPECT().
			Header(map[string]*prototool.Values{
				"resp": {Value: []string{"1", "2"}},
			}),
		connectServer.EXPECT().
			RecvMsg(gomock.Any()).
			Do(testhelpers.RecvMsg(&rpc.ConnectRequest{
				Msg: &rpc.ConnectRequest_Message{
					Message: &rpc.Message{
						Data: []byte{5, 6, 7},
					},
				},
			})),
		cb.EXPECT().
			Message([]byte{5, 6, 7}),
		connectServer.EXPECT().
			RecvMsg(gomock.Any()).
			Do(testhelpers.RecvMsg(&rpc.ConnectRequest{
				Msg: &rpc.ConnectRequest_Trailer{
					Trailer: &rpc.Trailer{
						Meta: map[string]*prototool.Values{
							"trailer": {Value: []string{"8", "9"}},
						},
					},
				},
			})),
		cb.EXPECT().
			Trailer(map[string]*prototool.Values{
				"trailer": {Value: []string{"8", "9"}},
			}),
		connectServer.EXPECT().
			RecvMsg(gomock.Any()).
			Return(io.EOF),
	)

	r, err := NewTunnelRegistry(zaptest.NewLogger(t), tunnelRegisterer, "grpc://127.0.0.1:123")
	require.NoError(t, err)
	return incomingStream, rpcApi, cb, connectServer, r
}

func descriptor() *rpc.Descriptor {
	return &rpc.Descriptor{
		AgentDescriptor: &info.AgentDescriptor{
			Services: []*info.Service{
				{
					Name: serviceName,
					Methods: []*info.Method{
						{
							Name: methodName,
						},
					},
				},
			},
		},
	}
}
