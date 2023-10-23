package server

import (
	"context"

	"github.com/pluralsh/kuberentes-agent/pkg/module/reverse_tunnel"
)

type module struct {
}

func (m *module) Run(ctx context.Context) error {
	return nil
}

func (m *module) Name() string {
	return reverse_tunnel.ModuleName
}
