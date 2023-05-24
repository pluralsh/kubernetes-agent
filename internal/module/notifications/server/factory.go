package server

import (
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/notifications"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/notifications/rpc"
)

type Factory struct {
	// GitPushPublisher provides a `Publish` interface to emit notifications about Git push events.
	GitPushPublisher  gitPushPublisherFunc
	GitPushSubscriber gitPushSubscriberFunc
}

func (f *Factory) New(config *modserver.Config) (modserver.Module, error) {
	rpc.RegisterNotificationsServer(config.ApiServer, newServer(f.GitPushPublisher))
	return &module{
		gitPushSubscriber: f.GitPushSubscriber,
	}, nil
}

func (f *Factory) Name() string {
	return notifications.ModuleName
}

func (f *Factory) StartStopPhase() modshared.ModuleStartStopPhase {
	return modshared.ModuleStartBeforeServers
}
