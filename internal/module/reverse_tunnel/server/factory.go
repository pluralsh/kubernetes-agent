package server

import (
	"time"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/reverse_tunnel"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/reverse_tunnel/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/reverse_tunnel/tunnel"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/retry"
)

const (
	getAgentInfoInitBackoff   = 20 * time.Millisecond
	getAgentInfoMaxBackoff    = 100 * time.Millisecond
	getAgentInfoResetDuration = 10 * time.Second
	getAgentInfoBackoffFactor = 2.0
	getAgentInfoJitter        = 1.0
)

type Factory struct {
	TunnelHandler tunnel.Handler
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
