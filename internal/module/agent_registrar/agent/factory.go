package agent

import (
	"time"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/agent_registrar"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/agent_registrar/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/retry"
)

const (
	registerAttemptInterval = 5 * time.Minute
	registerInitBackoff     = 10 * time.Second
	registerMaxBackoff      = 5 * time.Minute
	registerResetDuration   = 10 * time.Minute
	registerBackoffFactor   = 2.0
	registerJitter          = 1.0
)

type Factory struct {
	PodId int64
}

func (f *Factory) IsProducingLeaderModules() bool {
	return false
}

func (f *Factory) Name() string {
	return agent_registrar.ModuleName
}

func (f *Factory) New(config *modagent.Config) (modagent.Module, error) {
	m := &module{
		Log:       config.Log,
		AgentMeta: config.AgentMeta,
		PodId:     f.PodId,
		PollConfig: retry.NewPollConfigFactory(registerAttemptInterval, retry.NewExponentialBackoffFactory(
			registerInitBackoff,
			registerMaxBackoff,
			registerResetDuration,
			registerBackoffFactor,
			registerJitter,
		)),
		Client: rpc.NewAgentRegistrarClient(config.KasConn),
	}
	return m, nil
}

func (f *Factory) StartStopPhase() modshared.ModuleStartStopPhase {
	return modshared.ModuleStartBeforeServers
}
