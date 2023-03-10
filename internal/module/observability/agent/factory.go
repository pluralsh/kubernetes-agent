package agent

import (
	"crypto/tls"
	"fmt"
	"net"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/observability"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/tlstool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/pkg/agentcfg"
	"go.uber.org/zap"
)

type Factory struct {
	LogLevel            zap.AtomicLevel
	GrpcLogLevel        zap.AtomicLevel
	ListenNetwork       string
	ListenAddress       string
	CertFile            string
	KeyFile             string
	DefaultGrpcLogLevel agentcfg.LogLevelEnum
}

func (f *Factory) New(config *modagent.Config) (modagent.Module, error) {
	tlsConfig, err := tlstool.MaybeDefaultServerTLSConfig(f.CertFile, f.KeyFile)
	if err != nil {
		return nil, err
	}
	var listener func() (net.Listener, error)
	if tlsConfig != nil {
		listener = func() (net.Listener, error) {
			return tls.Listen(f.ListenNetwork, f.ListenAddress, tlsConfig) // nolint:gosec
		}
	} else {
		listener = func() (net.Listener, error) {
			return net.Listen(f.ListenNetwork, f.ListenAddress) // nolint:gosec
		}
	}
	return &module{
		log:                 config.Log,
		logLevel:            f.LogLevel,
		grpcLogLevel:        f.GrpcLogLevel,
		defaultGrpcLogLevel: f.DefaultGrpcLogLevel,
		api:                 config.Api,
		listener:            listener,
		serverName:          fmt.Sprintf("%s/%s/%s", config.AgentName, config.AgentMeta.Version, config.AgentMeta.CommitId),
	}, nil
}

func (f *Factory) Name() string {
	return observability.ModuleName
}

func (f *Factory) StartStopPhase() modshared.ModuleStartStopPhase {
	return modshared.ModuleStartBeforeServers
}
