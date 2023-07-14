package agent

import (
	"crypto/tls"
	"fmt"
	"net"

	"github.com/prometheus/client_golang/prometheus"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/observability"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/tlstool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/agentcfg"
	"go.uber.org/zap"
)

type Factory struct {
	LogLevel            zap.AtomicLevel
	GrpcLogLevel        zap.AtomicLevel
	DefaultGrpcLogLevel agentcfg.LogLevelEnum
	Gatherer            prometheus.Gatherer
	Registerer          prometheus.Registerer
	ListenNetwork       string
	ListenAddress       string
	CertFile            string
	KeyFile             string
}

func (f *Factory) IsProducingLeaderModules() bool {
	return false
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
		gatherer:            f.Gatherer,
		registerer:          f.Registerer,
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
