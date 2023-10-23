package server

import (
	"crypto/tls"
	"fmt"
	"net"

	"github.com/pluralsh/kuberentes-agent/internal/module/modserver"
	"github.com/pluralsh/kuberentes-agent/internal/module/modshared"
	"github.com/pluralsh/kuberentes-agent/internal/module/observability"
	"github.com/pluralsh/kuberentes-agent/internal/tool/tlstool"
	"github.com/prometheus/client_golang/prometheus"
)

type Factory struct {
	Gatherer prometheus.Gatherer
}

func (f *Factory) New(config *modserver.Config) (modserver.Module, error) {
	listenCfg := config.Config.GetObservability().GetListen()
	var listener func() (net.Listener, error)

	tlsConfig, err := tlstool.MaybeDefaultServerTLSConfig(listenCfg.GetCertificateFile(), listenCfg.GetKeyFile())
	if err != nil {
		return nil, err
	}
	if tlsConfig != nil {
		listener = func() (net.Listener, error) {
			return tls.Listen(*listenCfg.GetNetwork(), listenCfg.GetAddress(), tlsConfig)
		}
	} else {
		listener = func() (net.Listener, error) {
			return net.Listen(*listenCfg.GetNetwork(), listenCfg.GetAddress())
		}
	}
	return &module{
		log:           config.Log,
		api:           config.Api,
		cfg:           config.Config.GetObservability(),
		listener:      listener,
		gatherer:      f.Gatherer,
		registerer:    config.Registerer,
		serverName:    fmt.Sprintf("%s/%s/%s", config.KasName, config.Version, config.CommitId),
		probeRegistry: config.ProbeRegistry,
	}, nil
}

func (f *Factory) Name() string {
	return observability.ModuleName
}

func (f *Factory) StartStopPhase() modshared.ModuleStartStopPhase {
	return modshared.ModuleStartBeforeServers
}
