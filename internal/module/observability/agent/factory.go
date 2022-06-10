package agent

import (
	"fmt"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/observability"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/pkg/agentcfg"
	"go.uber.org/zap"
)

type Factory struct {
	LogLevel            zap.AtomicLevel
	GrpcLogLevel        zap.AtomicLevel
	DefaultGrpcLogLevel agentcfg.LogLevelEnum
}

func (f *Factory) New(config *modagent.Config) (modagent.Module, error) {
	return &module{
		log:                 config.Log,
		logLevel:            f.LogLevel,
		grpcLogLevel:        f.GrpcLogLevel,
		defaultGrpcLogLevel: f.DefaultGrpcLogLevel,
		api:                 config.Api,
		serverName:          fmt.Sprintf("%s/%s/%s", config.AgentName, config.AgentMeta.Version, config.AgentMeta.CommitId),
	}, nil
}

func (f *Factory) Name() string {
	return observability.ModuleName
}

func (f *Factory) UsesInternalServer() bool {
	return false
}
