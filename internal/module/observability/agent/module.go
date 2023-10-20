package agent

import (
	"context"
	"fmt"
	"net"

	"github.com/ash2k/stager"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/pluralsh/kuberentes-agent/internal/module/modshared"
	"github.com/pluralsh/kuberentes-agent/internal/module/observability"
	"github.com/pluralsh/kuberentes-agent/internal/tool/logz"
	"github.com/pluralsh/kuberentes-agent/internal/tool/prototool"
	"github.com/pluralsh/kuberentes-agent/pkg/agentcfg"
	"go.uber.org/zap"
)

type module struct {
	log                 *zap.Logger
	logLevel            zap.AtomicLevel
	grpcLogLevel        zap.AtomicLevel
	defaultGrpcLogLevel agentcfg.LogLevelEnum
	api                 modshared.Api
	gatherer            prometheus.Gatherer
	registerer          prometheus.Registerer
	listener            func() (net.Listener, error)
	serverName          string
}

const (
	prometheusUrlPath     = "/metrics"
	livenessProbeUrlPath  = "/liveness"
	readinessProbeUrlPath = "/readiness"
)

func (m *module) Run(ctx context.Context, cfg <-chan *agentcfg.AgentConfiguration) error {
	return stager.RunStages(ctx,
		func(stage stager.Stage) {
			// Listen for config changes and apply to logger
			stage.Go(func(ctx context.Context) error {
				done := ctx.Done()
				for {
					select {
					case <-done:
						return nil
					case config, ok := <-cfg:
						if !ok {
							return nil
						}
						err := m.setConfigurationLogging(config.Observability.Logging)
						if err != nil {
							m.log.Error("Failed to apply logging configuration", logz.Error(err))
							continue
						}
					}
				}
			})
			// Start metrics server
			stage.Go(func(ctx context.Context) error {
				lis, err := m.listener()
				if err != nil {
					return fmt.Errorf("observability listener failed to start: %w", err)
				}
				// Error is ignored because metricSrv.Run() closes the listener and
				// a second close always produces an error.
				defer lis.Close() // nolint:errcheck,gosec

				m.log.Info("Observability endpoint is up",
					logz.NetNetworkFromAddr(lis.Addr()),
					logz.NetAddressFromAddr(lis.Addr()),
				)

				metricSrv := observability.MetricServer{
					Log:                   m.log,
					Api:                   m.api,
					Name:                  m.serverName,
					Listener:              lis,
					PrometheusUrlPath:     prometheusUrlPath,
					LivenessProbeUrlPath:  livenessProbeUrlPath,
					ReadinessProbeUrlPath: readinessProbeUrlPath,
					Gatherer:              m.gatherer,
					Registerer:            m.registerer,
					ProbeRegistry:         observability.NewProbeRegistry(),
				}

				return metricSrv.Run(ctx)
			})
		},
	)
}

func (m *module) DefaultAndValidateConfiguration(config *agentcfg.AgentConfiguration) error {
	prototool.NotNil(&config.Observability)
	prototool.NotNil(&config.Observability.Logging)
	err := m.defaultAndValidateLogging(config.Observability.Logging)
	if err != nil {
		return fmt.Errorf("logging: %w", err)
	}
	return nil
}

func (m *module) Name() string {
	return observability.ModuleName
}

func (m *module) defaultAndValidateLogging(logging *agentcfg.LoggingCF) error {
	if logging.GrpcLevel == nil {
		logging.GrpcLevel = &m.defaultGrpcLogLevel
	}
	_, err := logz.LevelFromString(logging.Level.String())
	if err != nil {
		return err
	}
	_, err = logz.LevelFromString(logging.GrpcLevel.String())
	if err != nil {
		return err
	}
	return nil
}

func (m *module) setConfigurationLogging(logging *agentcfg.LoggingCF) error {
	err := setLogLevel(m.logLevel, logging.Level)
	if err != nil {
		return err
	}

	return setLogLevel(m.grpcLogLevel, *logging.GrpcLevel) // not nil after defaulting
}

func setLogLevel(logLevel zap.AtomicLevel, val agentcfg.LogLevelEnum) error {
	level, err := logz.LevelFromString(val.String())
	if err != nil {
		return err
	}
	logLevel.SetLevel(level)
	return nil
}
