package agent

import (
	"context"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/kubernetes_api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/pkg/agentcfg"
)

var (
	_ modagent.Module  = &module{}
	_ modagent.Factory = &Factory{}
)

type module struct {
	api modagent.API
}

func (m *module) Run(ctx context.Context, cfg <-chan *agentcfg.AgentConfiguration) error {
	// The tunnel feature is always required because CI for the agent's configuration project
	// can always access the agent.
	m.api.ToggleFeature(modagent.Tunnel, true)
	return nil
}

func (m *module) DefaultAndValidateConfiguration(config *agentcfg.AgentConfiguration) error {
	return nil
}

func (m *module) Name() string {
	return kubernetes_api.ModuleName
}
