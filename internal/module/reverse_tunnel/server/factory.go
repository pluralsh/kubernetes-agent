package server

import (
	"time"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/reverse_tunnel"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/reverse_tunnel/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/retry"
)

const (
	getAgentInfoInitBackoff   = 10 * time.Second
	getAgentInfoMaxBackoff    = 5 * time.Minute
	getAgentInfoResetDuration = 10 * time.Minute
	getAgentInfoBackoffFactor = 2.0
	getAgentInfoJitter        = 1.0
)

type Factory struct {
	TunnelHandler reverse_tunnel.TunnelHandler
}

func (f *Factory) New(config *modserver.Config) (modserver.Module, error) {
	rpc.RegisterReverseTunnelServer(config.AgentServer, &server{
		tunnelHandler: f.TunnelHandler,
		getAgentInfoPollConfig: retry.NewPollConfigFactory(0, retry.NewExponentialBackoffFactory(
			getAgentInfoInitBackoff,
			getAgentInfoMaxBackoff,
			getAgentInfoResetDuration,
			getAgentInfoBackoffFactor,
			getAgentInfoJitter,
		)),
	})
	return &module{}, nil
}

func (f *Factory) Name() string {
	return reverse_tunnel.ModuleName
}

func (f *Factory) StartStopPhase() modshared.ModuleStartStopPhase {
	return modshared.ModuleStartBeforeServers
}
