package server

import (
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/agent_registrar"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/agent_registrar/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/agent_tracker"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modshared"
)

type Factory struct {
	AgentRegisterer agent_tracker.Registerer
}

func (f *Factory) New(config *modserver.Config) (modserver.Module, error) {
	rpc.RegisterAgentRegistrarServer(config.AgentServer, &server{
		agentRegisterer: f.AgentRegisterer,
	})

	return &module{}, nil
}

func (f *Factory) Name() string {
	return agent_registrar.ModuleName
}

func (f *Factory) StartStopPhase() modshared.ModuleStartStopPhase {
	return modshared.ModuleStartBeforeServers
}
