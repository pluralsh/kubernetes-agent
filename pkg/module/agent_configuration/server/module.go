//go:build exclude

package server

import (
	"context"

	"github.com/pluralsh/kuberentes-agent/pkg/module/agent_configuration"
)

type module struct {
}

func (m *module) Run(ctx context.Context) error {
	return nil
}

func (m *module) Name() string {
	return agent_configuration.ModuleName
}