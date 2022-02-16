package server

import (
	"crypto/tls"
	"fmt"
	"net"

	gapi "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/gitlab/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/kubernetes_api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/kubernetes_api/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/cache"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/tlstool"
	"gitlab.com/gitlab-org/labkit/metrics"
)

const (
	k8sApiRequestCountKnownMetric = "k8s_api_proxy_request_count"

	httpMetricsNamespace = "kas_k8s_proxy"
)

type Factory struct {
}

func (f *Factory) New(config *modserver.Config) (modserver.Module, error) {
	k8sApi := config.Config.Agent.KubernetesApi
	if k8sApi == nil {
		return nopModule{}, nil
	}
	listenCfg := k8sApi.Listen
	certFile := listenCfg.CertificateFile
	keyFile := listenCfg.KeyFile
	var listener func() (net.Listener, error)

	tlsConfig, err := tlstool.MaybeDefaultServerTLSConfig(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	if tlsConfig != nil {
		listener = func() (net.Listener, error) {
			return tls.Listen(listenCfg.Network.String(), listenCfg.Address, tlsConfig)
		}
	} else {
		listener = func() (net.Listener, error) {
			return net.Listen(listenCfg.Network.String(), listenCfg.Address)
		}
	}
	m := &module{
		log: config.Log,
		proxy: kubernetesApiProxy{
			log:                       config.Log,
			api:                       config.Api,
			kubernetesApiClient:       rpc.NewKubernetesApiClient(config.AgentConn),
			gitLabClient:              config.GitLabClient,
			allowedAgentsCache:        cache.NewWithError(k8sApi.AllowedAgentCacheTtl.AsDuration(), k8sApi.AllowedAgentCacheErrorTtl.AsDuration(), gapi.IsCacheableError),
			requestCount:              config.UsageTracker.RegisterCounter(k8sApiRequestCountKnownMetric),
			metricsHttpHandlerFactory: metrics.NewHandlerFactory(metrics.WithNamespace(httpMetricsNamespace)),
			serverName:                fmt.Sprintf("%s/%s/%s", config.KasName, config.Version, config.CommitId),
			urlPathPrefix:             k8sApi.UrlPathPrefix,
		},
		listener: listener,
	}
	config.RegisterAgentApi(&rpc.KubernetesApi_ServiceDesc)
	return m, nil
}

func (f *Factory) Name() string {
	return kubernetes_api.ModuleName
}
