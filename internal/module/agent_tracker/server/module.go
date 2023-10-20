package server

import (
	"context"

	"github.com/pluralsh/kuberentes-agent/internal/module/agent_tracker"
)

type module struct{}

func (m *module) Run(ctx context.Context) error {
	return nil
}

func (m *module) Name() string {
	return agent_tracker.ModuleName
}
