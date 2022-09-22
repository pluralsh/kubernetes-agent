package server

import (
	"crypto/tls"
	"fmt"
	"net"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/api"
	gapi "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/gitlab/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/kubernetes_api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/kubernetes_api/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/cache"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/prototool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/redistool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/tlstool"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

const (
	k8sApiRequestCountKnownMetric        = "k8s_api_proxy_request"
	usersCiTunnelInteractionsCountMetric = "agent_users_using_ci_tunnel"
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
			return tls.Listen(*listenCfg.Network, listenCfg.Address, tlsConfig)
		}
	} else {
		listener = func() (net.Listener, error) {
			return net.Listen(*listenCfg.Network, listenCfg.Address)
		}
	}
	serverName := fmt.Sprintf("%s/%s/%s", config.KasName, config.Version, config.CommitId)
	m := &module{
		log: config.Log,
		proxy: kubernetesApiProxy{
			log:                 config.Log,
			api:                 config.Api,
			kubernetesApiClient: rpc.NewKubernetesApiClient(config.AgentConn),
			gitLabClient:        config.GitLabClient,
			allowedAgentsCache: cache.NewWithError(
				k8sApi.AllowedAgentCacheTtl.AsDuration(),
				k8sApi.AllowedAgentCacheErrorTtl.AsDuration(),
				&redistool.ErrCacher{
					Log:          config.Log,
					Client:       config.RedisClient,
					ErrMarshaler: prototool.ProtoErrMarshaler{},
					KeyToRedisKey: func(agentToken interface{}) string {
						key := api.AgentToken2key(agentToken.(api.AgentToken))
						return config.Config.Redis.KeyPrefix + ":allowed_agents_errs:" + string(key)
					},
				},
				gapi.IsCacheableError,
			),
			requestCounter:       config.UsageTracker.RegisterCounter(k8sApiRequestCountKnownMetric),
			ciTunnelUsersCounter: config.UsageTracker.RegisterUniqueCounter(usersCiTunnelInteractionsCountMetric),
			responseSerializer:   serializer.NewCodecFactory(runtime.NewScheme()),
			traceProvider:        config.TraceProvider,
			tracePropagator:      config.TracePropagator,
			meterProvider:        config.MeterProvider,
			serverName:           serverName,
			serverVia:            "gRPC/1.0 " + serverName,
			urlPathPrefix:        k8sApi.UrlPathPrefix,
		},
		listener: listener,
	}
	config.RegisterAgentApi(&rpc.KubernetesApi_ServiceDesc)
	return m, nil
}

func (f *Factory) Name() string {
	return kubernetes_api.ModuleName
}
