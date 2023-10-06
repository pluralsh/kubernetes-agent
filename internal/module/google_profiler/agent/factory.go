package agent

import (
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/google_profiler"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modshared"
)

type Factory struct {
}

func (f *Factory) IsProducingLeaderModules() bool {
	return false
}

func (f *Factory) New(config *modagent.Config) (modagent.Module, error) {
	return &module{
		log: config.Log,
		runner: &googleProfilerRunner{
			service: config.AgentName,
			version: config.AgentMeta.Version,
		},
	}, nil
}

func (f *Factory) Name() string {
	return google_profiler.ModuleName
}

func (f *Factory) StartStopPhase() modshared.ModuleStartStopPhase {
	return modshared.ModuleStartBeforeServers
}
