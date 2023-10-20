package server

import (
	"context"
	"math"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/pluralsh/kuberentes-agent/internal/module/agent_tracker"
	"github.com/pluralsh/kuberentes-agent/internal/module/agent_tracker/rpc"
	"github.com/pluralsh/kuberentes-agent/internal/module/modserver"
	"github.com/pluralsh/kuberentes-agent/internal/module/modshared"
	"github.com/pluralsh/kuberentes-agent/internal/tool/metric"
)

type Factory struct {
	AgentQuerier agent_tracker.Querier
}

func (f *Factory) New(config *modserver.Config) (modserver.Module, error) {
	connectedAgentsCountGaugeFunc := f.constructConnectedAgentsCountGaugeFunc()
	err := metric.Register(config.Registerer, connectedAgentsCountGaugeFunc)
	if err != nil {
		return nil, err
	}

	rpc.RegisterAgentTrackerServer(config.ApiServer, &server{
		agentQuerier: f.AgentQuerier,
	})

	return &module{}, nil
}

func (f *Factory) Name() string {
	return agent_tracker.ModuleName
}

func (f *Factory) StartStopPhase() modshared.ModuleStartStopPhase {
	return modshared.ModuleStartBeforeServers
}

func (f *Factory) constructConnectedAgentsCountGaugeFunc() prometheus.GaugeFunc {
	return prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "connected_agents_count",
		Help: "The number of unique connected agents",
	}, func() float64 {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		size, err := f.AgentQuerier.GetConnectedAgentsCount(ctx)
		if err != nil {
			return math.NaN()
		}
		return float64(size)
	},
	)
}
