package server

import (
	"context"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/notifications"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/notifications/rpc"
)

type Factory struct {
	// GitPushPublisher provides a `Publish` interface to emit notifications about Git push events.
	GitPushPublisher func(ctx context.Context, e *modserver.Project) error
}

func (f *Factory) New(config *modserver.Config) (modserver.Module, error) {
	rpc.RegisterNotificationsServer(config.ApiServer, newServer(f.GitPushPublisher))
	return &module{}, nil
}

func (f *Factory) Name() string {
	return notifications.ModuleName
}

func (f *Factory) StartStopPhase() modshared.ModuleStartStopPhase {
	return modshared.ModuleStartBeforeServers
}
