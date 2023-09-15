package server

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/agent_registrar/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/agent_tracker"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_agent_tracker"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_modserver"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestRegister(t *testing.T) {
	mockRpcApi, mockAgentTracker, s, req, ctx := setupServer(t)

	mockRpcApi.EXPECT().Log().Return(zaptest.NewLogger(t))
	mockRpcApi.EXPECT().AgentInfo(gomock.Any(), gomock.Any()).Return(&api.AgentInfo{Id: 123, ProjectId: 456}, nil)
	mockAgentTracker.EXPECT().RegisterConnection(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, connectedAgentInfo *agent_tracker.ConnectedAgentInfo) {
		assert.EqualValues(t, 123, connectedAgentInfo.AgentId)
		assert.EqualValues(t, 456, connectedAgentInfo.ProjectId)
		assert.EqualValues(t, 123456789, connectedAgentInfo.ConnectionId)
	})

	resp, err := s.Register(ctx, req)
	assert.NotNil(t, resp)
	assert.NoError(t, err)
}

func TestRegister_AgentInfo_Error(t *testing.T) {
	mockRpcApi, _, s, req, ctx := setupServer(t)

	mockRpcApi.EXPECT().Log().Return(zaptest.NewLogger(t))
	mockRpcApi.EXPECT().AgentInfo(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unavailable, "Failed to register agent"))

	resp, err := s.Register(ctx, req)
	assert.Nil(t, resp)
	assert.Equal(t, codes.Unavailable, status.Code(err))
}

func TestRegister_registerAgent_Error(t *testing.T) {
	mockRpcApi, mockAgentTracker, s, req, ctx := setupServer(t)

	expectedErr := errors.New("expected error")

	mockRpcApi.EXPECT().Log().Return(zaptest.NewLogger(t))
	mockRpcApi.EXPECT().AgentInfo(gomock.Any(), gomock.Any()).Return(&api.AgentInfo{Id: 1, ProjectId: 1}, nil)
	mockAgentTracker.EXPECT().RegisterConnection(gomock.Any(), gomock.Any()).Return(expectedErr)
	mockRpcApi.EXPECT().HandleProcessingError(gomock.Any(), gomock.Any(), gomock.Any(), expectedErr)

	resp, err := s.Register(ctx, req)
	assert.Nil(t, resp)
	assert.Equal(t, codes.Unavailable, status.Code(err))
}

func setupServer(t *testing.T) (*mock_modserver.MockAgentRpcApi,
	*mock_agent_tracker.MockTracker, *server, *rpc.RegisterRequest, context.Context) {
	ctrl := gomock.NewController(t)

	mockRpcApi := mock_modserver.NewMockAgentRpcApi(ctrl)
	mockAgentTracker := mock_agent_tracker.NewMockTracker(ctrl)

	s := &server{
		agentRegisterer: mockAgentTracker,
	}

	req := &rpc.RegisterRequest{
		AgentMeta: mock_modserver.AgentMeta(),
		PodId:     123456789,
	}

	ctx := modserver.InjectAgentRpcApi(context.Background(), mockRpcApi)

	return mockRpcApi, mockAgentTracker, s, req, ctx
}
