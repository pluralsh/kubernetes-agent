package server

import (
	"github.com/pluralsh/kuberentes-agent/internal/module/modserver"
	"github.com/pluralsh/kuberentes-agent/internal/module/modshared"
	"github.com/pluralsh/kuberentes-agent/internal/module/usage_metrics"
)

type Factory struct {
	UsageTracker usage_metrics.UsageTrackerCollector
}

func (f *Factory) New(config *modserver.Config) (modserver.Module, error) {
	return &module{
		log:                  config.Log,
		api:                  config.Api,
		usageTracker:         f.UsageTracker,
		usageReportingPeriod: config.Config.GetObservability().GetUsageReportingPeriod().AsDuration(),
	}, nil
}

func (f *Factory) Name() string {
	return usage_metrics.ModuleName
}

func (f *Factory) StartStopPhase() modshared.ModuleStartStopPhase {
	return modshared.ModuleStartBeforeServers
}
