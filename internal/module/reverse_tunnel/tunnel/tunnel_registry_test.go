package tunnel

import (
	"context"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/reverse_tunnel/info"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/reverse_tunnel/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/prototool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/matcher"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_reverse_tunnel_rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_tool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/testhelpers"
	"go.uber.org/mock/gomock"
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

func TestStopUnregistersAllConnections(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	rep := mock_tool.NewMockErrReporter(ctrl)
	connectServer := mock_reverse_tunnel_rpc.NewMockReverseTunnel_ConnectServer(ctrl)
	tunnelRegisterer := NewMockRegisterer(ctrl)
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
		rep.EXPECT().
			HandleProcessingError(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
	)
	r, err := NewTunnelRegistry(zaptest.NewLogger(t), rep, tunnelRegisterer)
	require.NoError(t, err)
	var wg wait.Group
	defer wg.Wait()
	wg.Start(func() {
		err = r.HandleTunnel(context.Background(), testhelpers.AgentInfoObj(), connectServer)
		assert.NoError(t, err)
	})
	time.Sleep(100 * time.Millisecond)
	tl, fl := r.stopInternal()
	assert.EqualValues(t, 1, tl)
	assert.Zero(t, fl)
	assert.Empty(t, r.tunsByAgentId)
	assert.Empty(t, r.findRequestsByAgentId)
	tl, fl = r.stopInternal()
	assert.Zero(t, tl)
	assert.Zero(t, fl)
}

func TestTunnelDoneRegistersUnusedTunnel(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	rep := mock_tool.NewMockErrReporter(ctrl)
	connectServer := mock_reverse_tunnel_rpc.NewMockReverseTunnel_ConnectServer(ctrl)
	tunnelRegisterer := NewMockRegisterer(ctrl)
	reg := make(chan struct{})
	gomock.InOrder(
		connectServer.EXPECT().
			Recv().
			Return(&rpc.ConnectRequest{
				Msg: &rpc.ConnectRequest_Descriptor_{
					Descriptor_: descriptor(),
				},
			}, nil),
		tunnelRegisterer.EXPECT(). // HandleTunnel()
						RegisterTunnel(gomock.Any(), gomock.Any()).
						Do(func(ctx context.Context, agentId int64) {
				close(reg)
			}),
		tunnelRegisterer.EXPECT(). // FindTunnel()
						UnregisterTunnel(gomock.Any(), gomock.Any()),
		tunnelRegisterer.EXPECT(). // Done()
						RegisterTunnel(gomock.Any(), gomock.Any()),
		tunnelRegisterer.EXPECT(). // FindTunnel()
						UnregisterTunnel(gomock.Any(), gomock.Any()),
		tunnelRegisterer.EXPECT(). // Done()
						RegisterTunnel(gomock.Any(), gomock.Any()),
		tunnelRegisterer.EXPECT(). // stopInternal()
						UnregisterTunnel(gomock.Any(), gomock.Any()),
		rep.EXPECT().
			HandleProcessingError(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
	)
	agentInfo := testhelpers.AgentInfoObj()
	r, err := NewTunnelRegistry(zaptest.NewLogger(t), rep, tunnelRegisterer)
	require.NoError(t, err)
	var wg wait.Group
	defer wg.Wait()
	wg.Start(func() {
		err = r.HandleTunnel(context.Background(), agentInfo, connectServer)
		assert.NoError(t, err)
	})
	<-reg
	found, th := r.FindTunnel(agentInfo.Id, serviceName, methodName)
	assert.True(t, found)
	tun, err := th.Get(context.Background())
	require.NoError(t, err)
	tun.Done()
	th.Done()
	found, th = r.FindTunnel(agentInfo.Id, serviceName, methodName)
	assert.True(t, found)
	tun, err = th.Get(context.Background())
	require.NoError(t, err)
	tun.Done()
	th.Done()
	tl, fl := r.stopInternal()
	assert.EqualValues(t, 1, tl)
	assert.Zero(t, fl)
}

func TestTunnelDoneDonePanics(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	rep := mock_tool.NewMockErrReporter(ctrl)
	connectServer := mock_reverse_tunnel_rpc.NewMockReverseTunnel_ConnectServer(ctrl)
	tunnelRegisterer := NewMockRegisterer(ctrl)
	reg := make(chan struct{})
	regCnt := 0
	connectServer.EXPECT().
		Recv().
		Return(&rpc.ConnectRequest{
			Msg: &rpc.ConnectRequest_Descriptor_{
				Descriptor_: descriptor(),
			},
		}, nil)
	tunnelRegisterer.EXPECT().
		RegisterTunnel(gomock.Any(), gomock.Any()).
		Do(func(ctx context.Context, agentId int64) {
			regCnt++
			if regCnt == 1 {
				close(reg)
			}
		}).
		Times(2)
	tunnelRegisterer.EXPECT().
		UnregisterTunnel(gomock.Any(), gomock.Any()).
		Times(2)
	rep.EXPECT().
		HandleProcessingError(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
	agentInfo := testhelpers.AgentInfoObj()
	r, err := NewTunnelRegistry(zaptest.NewLogger(t), rep, tunnelRegisterer)
	require.NoError(t, err)
	var wg wait.Group
	defer wg.Wait()
	wg.Start(func() {
		err = r.HandleTunnel(context.Background(), agentInfo, connectServer)
		assert.NoError(t, err)
	})
	<-reg
	found, th := r.FindTunnel(agentInfo.Id, serviceName, methodName)
	assert.True(t, found)
	tun, err := th.Get(context.Background())
	require.NoError(t, err)
	tun.Done()
	assert.PanicsWithError(t, "unreachable: ready -> done should never happen", func() {
		tun.Done()
	})
	th.Done()
	tl, fl := r.stopInternal()
	assert.EqualValues(t, 1, tl)
	assert.Zero(t, fl)
}

func TestHandleTunnelIsUnblockedByContext(t *testing.T) {
	t.Parallel()
	ctxConn, cancelConn := context.WithTimeout(context.Background(), 50*time.Millisecond) // will unblock HandleTunnel()
	defer cancelConn()

	ctrl := gomock.NewController(t)
	rep := mock_tool.NewMockErrReporter(ctrl)
	connectServer := mock_reverse_tunnel_rpc.NewMockReverseTunnel_ConnectServer(ctrl)
	tunnelRegisterer := NewMockRegisterer(ctrl)
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
	r, err := NewTunnelRegistry(zaptest.NewLogger(t), rep, tunnelRegisterer)
	require.NoError(t, err)
	err = r.HandleTunnel(ctxConn, testhelpers.AgentInfoObj(), connectServer)
	assert.NoError(t, err)
	tl, fl := r.stopInternal()
	assert.Zero(t, tl)
	assert.Zero(t, fl)
}

// Two tunnels with the same agent id. Both register. Then one of them is retrieved via FindTunnel()
// and then its context is cancelled. If this test gets stuck, we have a problem.
// Reproducer for https://gitlab.com/gitlab-org/cluster-integration/gitlab-agent/-/issues/183.
func TestHandleTunnelIsUnblockedByContext_WithTwoTunnels(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	rep := mock_tool.NewMockErrReporter(ctrl)
	connectServer1 := mock_reverse_tunnel_rpc.NewMockReverseTunnel_ConnectServer(ctrl)
	connectServer2 := mock_reverse_tunnel_rpc.NewMockReverseTunnel_ConnectServer(ctrl)
	tunnelRegisterer := NewMockRegisterer(ctrl)
	var regWg sync.WaitGroup
	regWg.Add(2)
	d1 := descriptor()
	connectServer1.EXPECT().
		Recv().
		Return(&rpc.ConnectRequest{
			Msg: &rpc.ConnectRequest_Descriptor_{
				Descriptor_: d1,
			},
		}, nil)
	connectServer2.EXPECT().
		Recv().
		Return(&rpc.ConnectRequest{
			Msg: &rpc.ConnectRequest_Descriptor_{
				Descriptor_: descriptor(),
			},
		}, nil)
	tunnelRegisterer.EXPECT().
		RegisterTunnel(gomock.Any(), gomock.Any()).
		Do(func(ctx context.Context, agentId int64) {
			regWg.Done()
		}).Times(2)
	tunnelRegisterer.EXPECT().
		UnregisterTunnel(gomock.Any(), gomock.Any()).
		Times(2)
	rep.EXPECT().
		HandleProcessingError(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
	r, err := NewTunnelRegistry(zaptest.NewLogger(t), rep, tunnelRegisterer)
	require.NoError(t, err)
	var wg wait.Group
	defer wg.Wait()
	defer r.Stop()
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
	// wait for both to register
	regWg.Wait()
	found, th := r.FindTunnel(agentInfo.Id, serviceName, methodName)
	assert.True(t, found)
	tun, err := th.Get(context.Background())
	require.NoError(t, err)
	// cancel context for the found tunnel
	switch tun.(*tunnelImpl).tunnel {
	case connectServer1:
		cancel1()
	case connectServer2:
		cancel2()
	default:
		t.FailNow()
	}
	assert.Eventually(t, func() bool {
		r.mu.Lock()
		defer r.mu.Unlock()
		return tun.(*tunnelImpl).state == stateContextDone
	}, time.Second, 10*time.Millisecond)
	tun.Done()
	th.Done()
}

func TestHandleTunnelReturnErrOnRecvErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	rep := mock_tool.NewMockErrReporter(ctrl)
	connectServer := mock_reverse_tunnel_rpc.NewMockReverseTunnel_ConnectServer(ctrl)
	connectServer.EXPECT().
		Recv().
		Return(nil, status.Error(codes.DataLoss, "expected err"))
	tunnelRegisterer := NewMockRegisterer(ctrl)
	r, err := NewTunnelRegistry(zaptest.NewLogger(t), rep, tunnelRegisterer)
	require.NoError(t, err)
	err = r.HandleTunnel(context.Background(), testhelpers.AgentInfoObj(), connectServer)
	assert.EqualError(t, err, "rpc error: code = DataLoss desc = expected err")
}

func TestHandleTunnelReturnErrOnInvalidMsg(t *testing.T) {
	ctrl := gomock.NewController(t)
	rep := mock_tool.NewMockErrReporter(ctrl)
	connectServer := mock_reverse_tunnel_rpc.NewMockReverseTunnel_ConnectServer(ctrl)
	connectServer.EXPECT().
		Recv().
		Return(&rpc.ConnectRequest{
			Msg: &rpc.ConnectRequest_Header{
				Header: &rpc.Header{},
			},
		}, nil)
	tunnelRegisterer := NewMockRegisterer(ctrl)
	r, err := NewTunnelRegistry(zaptest.NewLogger(t), rep, tunnelRegisterer)
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
	defer func() {
		tl, fl := r.stopInternal()
		assert.Zero(t, tl)
		assert.Zero(t, fl)
	}()
	wg.Start(func() {
		assert.NoError(t, r.HandleTunnel(context.Background(), agentInfo, tunnel))
	})
	time.Sleep(50 * time.Millisecond)
	found, th := r.FindTunnel(agentInfo.Id, serviceName, methodName)
	defer th.Done()
	assert.True(t, found)
	tun, err := th.Get(context.Background())
	require.NoError(t, err)
	defer tun.Done()
	err = tun.ForwardStream(zaptest.NewLogger(t), rpcApi, incomingStream, cb)
	require.NoError(t, err)
}

func TestHandleTunnelIsNotMatchedToIncomingConnectionForMissingMethod(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	rep := mock_tool.NewMockErrReporter(ctrl)
	tunnelRegisterer := NewMockRegisterer(ctrl)
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
		rep.EXPECT().
			HandleProcessingError(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
	)
	r, err := NewTunnelRegistry(zaptest.NewLogger(t), rep, tunnelRegisterer)
	require.NoError(t, err)
	agentInfo := testhelpers.AgentInfoObj()
	var wg wait.Group
	defer wg.Wait()
	defer r.Stop()
	wg.Start(func() {
		assert.NoError(t, r.HandleTunnel(context.Background(), agentInfo, connectServer))
	})
	time.Sleep(50 * time.Millisecond)
	ctx2, cancel2 := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel2()
	found, th := r.FindTunnel(agentInfo.Id, "missing_service", "missing_method")
	defer th.Done()
	assert.False(t, found)
	_, err = th.Get(ctx2)
	assert.EqualError(t, err, "rpc error: code = DeadlineExceeded desc = FindTunnel request aborted: context deadline exceeded")
}

func TestForwardStreamIsMatchedToHandleTunnel(t *testing.T) {
	t.Parallel()
	incomingStream, rpcApi, cb, tunnel, r := setupStreams(t, false)
	agentInfo := testhelpers.AgentInfoObj()
	var wg wait.Group
	defer wg.Wait()
	defer func() {
		tl, fl := r.stopInternal()
		assert.Zero(t, tl)
		assert.Zero(t, fl)
	}()
	wg.Start(func() {
		_, th := r.FindTunnel(agentInfo.Id, serviceName, methodName)
		defer th.Done()
		tun, err := th.Get(context.Background())
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
	rep := mock_tool.NewMockErrReporter(ctrl)
	tunnelRegisterer := NewMockRegisterer(ctrl)
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
		rep.EXPECT().
			HandleProcessingError(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
	)
	r, err := NewTunnelRegistry(zaptest.NewLogger(t), rep, tunnelRegisterer)
	require.NoError(t, err)
	agentInfo := testhelpers.AgentInfoObj()
	var wg wait.Group
	defer wg.Wait()
	defer r.Stop()
	wg.Start(func() {
		found, th := r.FindTunnel(agentInfo.Id, "missing_service", "missing_method")
		defer th.Done()
		assert.False(t, found)
		_, findErr := th.Get(context.Background())
		assert.EqualError(t, findErr, "rpc error: code = Unavailable desc = kas is shutting down")
	})
	time.Sleep(50 * time.Millisecond)
	ctx2, cancel2 := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel2()
	err = r.HandleTunnel(ctx2, agentInfo, connectServer)
	assert.NoError(t, err)
}

func TestFindTunnelIsUnblockedByContext(t *testing.T) {
	t.Parallel()
	ctxConn, cancelConn := context.WithTimeout(context.Background(), 50*time.Millisecond) // will unblock FindTunnel()
	defer cancelConn()

	ctrl := gomock.NewController(t)
	rep := mock_tool.NewMockErrReporter(ctrl)
	tunnelRegisterer := NewMockRegisterer(ctrl)
	r, err := NewTunnelRegistry(zaptest.NewLogger(t), rep, tunnelRegisterer)
	require.NoError(t, err)
	found, th := r.FindTunnel(testhelpers.AgentId, serviceName, methodName)
	defer th.Done()
	assert.False(t, found)
	_, err = th.Get(ctxConn)
	assert.EqualError(t, err, "rpc error: code = DeadlineExceeded desc = FindTunnel request aborted: context deadline exceeded")
}

func setupStreams(t *testing.T, expectRegisterTunnel bool) (*mock_rpc.MockServerStream, *mock_modserver.MockAgentRpcApi, *MockDataCallback, *mock_reverse_tunnel_rpc.MockReverseTunnel_ConnectServer, *TunnelRegistry) {
	const metaKey = "Cba"
	meta := metadata.MD{}
	meta.Set(metaKey, "3", "4")
	ctrl := gomock.NewController(t)
	rep := mock_tool.NewMockErrReporter(ctrl)
	sts := mock_rpc.NewMockServerTransportStream(ctrl)
	cb := NewMockDataCallback(ctrl)

	rpcApi := mock_modserver.NewMockAgentRpcApi(ctrl)
	incomingCtx := grpc.NewContextWithServerTransportStream(context.Background(), sts)
	incomingCtx = metadata.NewIncomingContext(incomingCtx, meta)
	incomingStream := mock_rpc.NewMockServerStream(ctrl)
	incomingStream.EXPECT().
		Context().
		Return(incomingCtx).
		MinTimes(1)

	tunnelRegisterer := NewMockRegisterer(ctrl)
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

	r, err := NewTunnelRegistry(zaptest.NewLogger(t), rep, tunnelRegisterer)
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
