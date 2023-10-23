package metric

import (
	"fmt"

	"github.com/pluralsh/kuberentes-agent/internal/tool/logz"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

func Register(registerer prometheus.Registerer, toRegister ...prometheus.Collector) error {
	for _, c := range toRegister {
		if err := registerer.Register(c); err != nil {
			return fmt.Errorf("registering %T: %w", c, err)
		}
	}
	return nil
}

type OtelErrorHandler zap.Logger

func (h *OtelErrorHandler) Handle(err error) {
	(*zap.Logger)(h).Warn("OpenTelemetry error", logz.Error(err))
}
