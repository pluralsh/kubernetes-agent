package server

import (
	"github.com/pluralsh/kuberentes-agent/internal/tool/prototool"
	"github.com/pluralsh/kuberentes-agent/pkg/kascfg"
)

const (
	defaultObservabilityListenNetwork         = "tcp"
	defaultObservabilityListenAddress         = "127.0.0.1:8151"
	defaultObservabilityPrometheusUrlPath     = "/metrics"
	defaultObservabilityLivenessProbeUrlPath  = "/liveness"
	defaultObservabilityReadinessProbeUrlPath = "/readiness"

	defaultGrpcLogLevel = kascfg.LogLevelEnum_error
)

func ApplyDefaults(config *kascfg.ConfigurationFile) {
	prototool.NotNil(&config.Observability)
	o := config.GetObservability()

	prototool.NotNil(&o.Listen)
	prototool.StringPtr(&o.Listen.Network, defaultObservabilityListenNetwork)
	prototool.String(&o.Listen.Address, defaultObservabilityListenAddress)

	prototool.NotNil(&o.Prometheus)
	prototool.String(&o.Prometheus.UrlPath, defaultObservabilityPrometheusUrlPath)

	prototool.NotNil(&o.Sentry)

	prototool.NotNil(&o.Logging)
	if o.GetLogging().GetGrpcLevel() == nil {
		x := defaultGrpcLogLevel
		o.Logging.GrpcLevel = &x
	}

	prototool.NotNil(&o.LivenessProbe)
	prototool.String(&o.LivenessProbe.UrlPath, defaultObservabilityLivenessProbeUrlPath)

	prototool.NotNil(&o.ReadinessProbe)
	prototool.String(&o.ReadinessProbe.UrlPath, defaultObservabilityReadinessProbeUrlPath)
}
