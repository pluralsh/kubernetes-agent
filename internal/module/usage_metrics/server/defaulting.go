package server

import (
	"time"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/prototool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/pkg/kascfg"
)

const (
	defaultObservabilityUsageReportingPeriod = 1 * time.Minute
)

func ApplyDefaults(config *kascfg.ConfigurationFile) {
	prototool.NotNil(&config.Observability)
	prototool.Duration(&config.Observability.UsageReportingPeriod, defaultObservabilityUsageReportingPeriod)
}
