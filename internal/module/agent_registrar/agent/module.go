package agent

import (
	"context"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/agent_registrar"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/agent_registrar/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/retry"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/agentcfg"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/entity"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"k8s.io/client-go/discovery"
)

type module struct {
	Log         *zap.Logger
	AgentMeta   *entity.AgentMeta
	PodId       int64
	PollConfig  retry.PollConfigFactory
	Client      rpc.AgentRegistrarClient
	KubeVersion discovery.ServerVersionInterface
}

func (m *module) Run(ctx context.Context, cfg <-chan *agentcfg.AgentConfiguration) error {
	// Create a deep copy of agentMeta to prevent unexpected mutations
	agentMeta := proto.Clone(m.AgentMeta).(*entity.AgentMeta)

	_ = retry.PollWithBackoff(ctx, m.PollConfig(), func(ctx context.Context) (error, retry.AttemptResult) {
		// Retrieve and set the Kubernetes version
		version, err := m.KubeVersion.ServerVersion()
		if err == nil {
			agentMeta.KubernetesVersion.Major = version.Major
			agentMeta.KubernetesVersion.Minor = version.Minor
			agentMeta.KubernetesVersion.GitVersion = version.GitVersion
			agentMeta.KubernetesVersion.Platform = version.Platform
		} else {
			m.Log.Warn("Failed to fetch Kubernetes version", logz.Error(err))
		}

		_, err = m.Client.Register(ctx, &rpc.RegisterRequest{
			AgentMeta: agentMeta,
			PodId:     m.PodId,
		})
		if err != nil {
			if !grpctool.RequestCanceledOrTimedOut(err) {
				m.Log.Error("Failed to register agent pod. Please make sure the agent version matches the server version", logz.Error(err))
			}
			return nil, retry.Backoff
		}

		return nil, retry.Continue
	})
	return nil
}

func (m *module) DefaultAndValidateConfiguration(config *agentcfg.AgentConfiguration) error {
	return nil
}

func (m *module) Name() string {
	return agent_registrar.ModuleName
}
