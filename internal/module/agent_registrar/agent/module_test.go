package agent

import (
	"context"
	"testing"

	"github.com/pluralsh/kuberentes-agent/internal/module/agent_registrar/rpc"
	"github.com/pluralsh/kuberentes-agent/internal/tool/mathz"
	"github.com/pluralsh/kuberentes-agent/internal/tool/testing/mock_agent_registrar"
	"github.com/pluralsh/kuberentes-agent/internal/tool/testing/testhelpers"
	"github.com/pluralsh/kuberentes-agent/pkg/entity"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes/fake"
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
		Log:         zaptest.NewLogger(t),
		AgentMeta:   &entity.AgentMeta{KubernetesVersion: &entity.KubernetesVersion{}},
		PodId:       mathz.Int63(),
		PollConfig:  testhelpers.NewPollConfig(0),
		Client:      client,
		KubeVersion: fake.NewSimpleClientset().Discovery(),
	}
	_ = m.Run(ctx, nil)
}
