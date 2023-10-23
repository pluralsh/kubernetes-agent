package agent

import (
	"context"

	"github.com/pluralsh/kuberentes-agent/pkg/agentcfg"
	"github.com/pluralsh/kuberentes-agent/pkg/module/kubernetes_api"
)

type module struct {
}

func (m *module) Run(ctx context.Context, cfg <-chan *agentcfg.AgentConfiguration) error {
	return nil
}

func (m *module) DefaultAndValidateConfiguration(config *agentcfg.AgentConfiguration) error {
	return nil
}

func (m *module) Name() string {
	return kubernetes_api.ModuleName
}
