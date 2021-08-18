package server

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/reverse_tunnel/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/testing/mock_modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/testing/mock_reverse_tunnel"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/testing/mock_reverse_tunnel_rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/testing/testhelpers"
	"go.uber.org/zap/zaptest"
)

var (
	_ rpc.ReverseTunnelServer = &server{}
)

func TestConnectAllowsValidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	h := mock_reverse_tunnel.NewMockTunnelHandler(ctrl)
	mockRpcApi := mock_modserver.NewMockAgentRpcApi(ctrl)
	mockRpcApi.EXPECT().
		Log().
		Return(zaptest.NewLogger(t)).
		AnyTimes()
	s := &server{
		tunnelHandler: h,
	}
	agentInfo := testhelpers.AgentInfoObj()
	ctx := grpctool.AddMaxConnectionAgeContext(context.Background(), context.Background())
	ctx = modserver.InjectAgentRpcApi(ctx, mockRpcApi)
	connectServer := mock_reverse_tunnel_rpc.NewMockReverseTunnel_ConnectServer(ctrl)
	connectServer.EXPECT().
		Context().
		Return(ctx).
		MinTimes(1)
	gomock.InOrder(
		mockRpcApi.EXPECT().
			AgentInfo(gomock.Any(), gomock.Any()).
			Return(agentInfo, nil),
		h.EXPECT().
			HandleTunnel(gomock.Any(), agentInfo, connectServer),
	)
	err := s.Connect(connectServer)
	require.NoError(t, err)
}

func TestConnectRejectsInvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	h := mock_reverse_tunnel.NewMockTunnelHandler(ctrl)
	mockRpcApi := mock_modserver.NewMockAgentRpcApi(ctrl)
	mockRpcApi.EXPECT().
		Log().
		Return(zaptest.NewLogger(t)).
		AnyTimes()
	s := &server{
		tunnelHandler: h,
	}
	ctx := grpctool.AddMaxConnectionAgeContext(context.Background(), context.Background())
	ctx = modserver.InjectAgentRpcApi(ctx, mockRpcApi)
	connectServer := mock_reverse_tunnel_rpc.NewMockReverseTunnel_ConnectServer(ctrl)
	connectServer.EXPECT().
		Context().
		Return(ctx).
		MinTimes(1)
	mockRpcApi.EXPECT().
		AgentInfo(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("expected err"))
	err := s.Connect(connectServer)
	assert.EqualError(t, err, "expected err")
}
