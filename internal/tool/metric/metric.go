package metric

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

func Register(registerer prometheus.Registerer, toRegister ...prometheus.Collector) error {
	for _, c := range toRegister {
		if err := registerer.Register(c); err != nil {
			return fmt.Errorf("registering %T: %w", c, err)
		}
	}
	return nil
}
