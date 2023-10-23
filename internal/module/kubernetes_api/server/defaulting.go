package server

import (
	"strings"
	"time"

	"github.com/pluralsh/kuberentes-agent/internal/tool/prototool"
	"github.com/pluralsh/kuberentes-agent/pkg/kascfg"
)

const (
	defaultKubernetesApiListenNetwork    = "tcp"
	defaultKubernetesApiListenAddress    = "0.0.0.0:8154"
	defaultListenGracePeriod             = 5 * time.Second
	defaultAllowedAgentInfoCacheTTL      = 1 * time.Minute
	defaultAllowedAgentInfoCacheErrorTTL = 10 * time.Second
	defaultShutdownGracePeriod           = 1 * time.Hour
)

func ApplyDefaults(config *kascfg.ConfigurationFile) {
	prototool.NotNil(&config.Agent)
	o := config.GetAgent().GetKubernetesApi()

	if o == nil {
		return
	}
	prototool.NotNil(&o.Listen)
	prototool.StringPtr(&o.Listen.Network, defaultKubernetesApiListenNetwork)
	prototool.String(&o.Listen.Address, defaultKubernetesApiListenAddress)
	prototool.Duration(&o.Listen.ListenGracePeriod, defaultListenGracePeriod)
	prototool.Duration(&o.Listen.ShutdownGracePeriod, defaultShutdownGracePeriod)
	if !strings.HasSuffix(o.GetUrlPathPrefix(), "/") {
		o.UrlPathPrefix = o.GetUrlPathPrefix() + "/"
	}
	prototool.Duration(&o.AllowedAgentCacheTtl, defaultAllowedAgentInfoCacheTTL)
	prototool.Duration(&o.AllowedAgentCacheErrorTtl, defaultAllowedAgentInfoCacheErrorTTL)
}
