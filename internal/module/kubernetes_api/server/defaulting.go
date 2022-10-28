package server

import (
	"strings"
	"time"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/prototool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/pkg/kascfg"
)

const (
	defaultKubernetesApiListenNetwork    = "tcp"
	defaultKubernetesApiListenAddress    = "0.0.0.0:8154"
	defaultListenGracePeriod             = 5 * time.Second
	defaultAllowedAgentInfoCacheTTL      = 1 * time.Minute
	defaultAllowedAgentInfoCacheErrorTTL = 10 * time.Second
)

func ApplyDefaults(config *kascfg.ConfigurationFile) {
	prototool.NotNil(&config.Agent)
	o := config.Agent.KubernetesApi

	if o == nil {
		return
	}
	prototool.NotNil(&o.Listen)
	prototool.StringPtr(&o.Listen.Network, defaultKubernetesApiListenNetwork)
	prototool.String(&o.Listen.Address, defaultKubernetesApiListenAddress)
	prototool.Duration(&o.Listen.ListenGracePeriod, defaultListenGracePeriod)
	if !strings.HasSuffix(o.UrlPathPrefix, "/") {
		o.UrlPathPrefix = o.UrlPathPrefix + "/"
	}
	prototool.Duration(&o.AllowedAgentCacheTtl, defaultAllowedAgentInfoCacheTTL)
	prototool.Duration(&o.AllowedAgentCacheErrorTtl, defaultAllowedAgentInfoCacheErrorTTL)
}
