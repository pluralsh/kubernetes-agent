package server

import (
	"context"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/notifications"
)

type SubscribeToEvents func(ctx context.Context)

type module struct {
	subscribeToEvents SubscribeToEvents
}

func (m *module) Run(ctx context.Context) error {
	m.subscribeToEvents(ctx)
	return nil
}

func (m *module) Name() string {
	return notifications.ModuleName
}
