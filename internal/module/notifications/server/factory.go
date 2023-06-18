package server

import (
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/notifications"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/notifications/rpc"
)

type Factory struct {
	// PublishEvent provides a `Publish` interface to emit events.
	PublishEvent      PublishEvent
	SubscribeToEvents SubscribeToEvents
}

func (f *Factory) New(config *modserver.Config) (modserver.Module, error) {
	rpc.RegisterNotificationsServer(config.ApiServer, newServer(f.PublishEvent))
	return &module{
		subscribeToEvents: f.SubscribeToEvents,
	}, nil
}

func (f *Factory) Name() string {
	return notifications.ModuleName
}

func (f *Factory) StartStopPhase() modshared.ModuleStartStopPhase {
	return modshared.ModuleStartBeforeServers
}
