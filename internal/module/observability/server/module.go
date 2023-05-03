package server

import (
	"context"
	"net"

	"github.com/prometheus/client_golang/prometheus"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/observability"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/kascfg"
	"go.uber.org/zap"
)

type module struct {
	log           *zap.Logger
	api           modshared.Api
	cfg           *kascfg.ObservabilityCF
	listener      func() (net.Listener, error)
	gatherer      prometheus.Gatherer
	registerer    prometheus.Registerer
	serverName    string
	probeRegistry *observability.ProbeRegistry
}

func (m *module) Run(ctx context.Context) (retErr error) {
	lis, err := m.listener()
	if err != nil {
		return err
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
		PrometheusUrlPath:     m.cfg.Prometheus.UrlPath,
		LivenessProbeUrlPath:  m.cfg.LivenessProbe.UrlPath,
		ReadinessProbeUrlPath: m.cfg.ReadinessProbe.UrlPath,
		Gatherer:              m.gatherer,
		Registerer:            m.registerer,
		ProbeRegistry:         m.probeRegistry,
	}
	return metricSrv.Run(ctx)
}

func (m *module) Name() string {
	return observability.ModuleName
}
