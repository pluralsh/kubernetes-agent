package server

import (
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/notifications"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/notifications/rpc"
)

type Factory struct {
	Publisher Publisher
}

func (f *Factory) New(config *modserver.Config) (modserver.Module, error) {
	rpc.RegisterNotificationsServer(config.ApiServer, newServer(
		f.Publisher,
	))
	return &module{}, nil
}

func (f *Factory) Name() string {
	return notifications.ModuleName
}

func (f *Factory) StartStopPhase() modshared.ModuleStartStopPhase {
	return modshared.ModuleStartBeforeServers
}
