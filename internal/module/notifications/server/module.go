package server

import (
	"context"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/notifications"
)

type gitPushSubscriberFunc func(ctx context.Context)

type module struct {
	gitPushSubscriber gitPushSubscriberFunc
}

func (m *module) Run(ctx context.Context) error {
	m.gitPushSubscriber(ctx)
	return nil
}

func (m *module) Name() string {
	return notifications.ModuleName
}
