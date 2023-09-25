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
)

type module struct {
	Log        *zap.Logger
	AgentMeta  *entity.AgentMeta
	PodId      int64
	PollConfig retry.PollConfigFactory
	Client     rpc.AgentRegistrarClient
}

func (m *module) Run(ctx context.Context, cfg <-chan *agentcfg.AgentConfiguration) error {
	_ = retry.PollWithBackoff(ctx, m.PollConfig(), func(ctx context.Context) (error, retry.AttemptResult) {
		_, err := m.Client.Register(ctx, &rpc.RegisterRequest{
			AgentMeta: m.AgentMeta,
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
