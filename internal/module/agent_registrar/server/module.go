package server

import (
	"context"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/agent_registrar"
)

type module struct{}

func (m *module) Run(ctx context.Context) error {
	return nil
}

func (m *module) Name() string {
	return agent_registrar.ModuleName
}