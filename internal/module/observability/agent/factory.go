package agent

import (
	"fmt"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/observability"
	"go.uber.org/zap"
)

type Factory struct {
	LogLevel zap.AtomicLevel
}

func (f *Factory) New(config *modagent.Config) (modagent.Module, error) {
	return &module{
		log:        config.Log,
		logLevel:   f.LogLevel,
		api:        config.Api,
		serverName: fmt.Sprintf("%s/%s/%s", config.AgentName, config.AgentMeta.Version, config.AgentMeta.CommitId),
	}, nil
}

func (f *Factory) Name() string {
	return observability.ModuleName
}
