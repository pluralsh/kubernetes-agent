package agent

import (
	"context"
	"testing"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/agent_registrar/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/mathz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_agent_registrar"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/testhelpers"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/entity"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc"
)

func TestModule_Run(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	client := mock_agent_registrar.NewMockAgentRegistrarClient(ctrl)
	client.EXPECT().
		Register(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, request *rpc.RegisterRequest, opts ...grpc.CallOption) (*rpc.RegisterResponse, error) {
			cancel()
			return &rpc.RegisterResponse{}, nil
		})

	m := &module{
		Log:        zaptest.NewLogger(t),
		AgentMeta:  &entity.AgentMeta{},
		PodId:      mathz.Int63(),
		PollConfig: testhelpers.NewPollConfig(0),
		Client:     client,
	}
	_ = m.Run(ctx, nil)
}
