package server

import (
	"time"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/agent_configuration"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/agent_configuration/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/agent_tracker"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/retry"
)

const (
	getConfigurationInitBackoff   = 10 * time.Second
	getConfigurationMaxBackoff    = time.Minute
	getConfigurationResetDuration = 2 * time.Minute
	getConfigurationBackoffFactor = 2.0
	getConfigurationJitter        = 1.0
)

type Factory struct {
	AgentRegisterer agent_tracker.Registerer
}

func (f *Factory) New(config *modserver.Config) (modserver.Module, error) {
	agentCfg := config.Config.Agent.Configuration
	rpc.RegisterAgentConfigurationServer(config.AgentServer, &fakeServer{
		serverApi:                config.Api,
		gitaly:                   config.Gitaly,
		gitLabClient:             config.GitLabClient,
		agentRegisterer:          f.AgentRegisterer,
		maxConfigurationFileSize: int64(agentCfg.MaxConfigurationFileSize),
		getConfigurationPollConfig: retry.NewPollConfigFactory(agentCfg.PollPeriod.AsDuration(), retry.NewExponentialBackoffFactory(
			getConfigurationInitBackoff,
			getConfigurationMaxBackoff,
			getConfigurationResetDuration,
			getConfigurationBackoffFactor,
			getConfigurationJitter,
		)),
		gitLabExternalUrl: config.Config.Gitlab.GetExternalUrl(),
	})
	return &module{}, nil
}

func (f *Factory) Name() string {
	return agent_configuration.ModuleName
}

func (f *Factory) StartStopPhase() modshared.ModuleStartStopPhase {
	return modshared.ModuleStartBeforeServers
}
