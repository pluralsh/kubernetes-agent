package server

import (
	"crypto/sha256"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	gapi "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/gitlab/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/kubernetes_api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/kubernetes_api/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/cache"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/prototool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/redistool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/tlstool"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

const (
	shutdownTimeout = 15 * time.Second

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
			allowedAgentsCache: cache.NewWithError[string, *gapi.AllowedAgentsForJob](
				k8sApi.AllowedAgentCacheTtl.AsDuration(),
				k8sApi.AllowedAgentCacheErrorTtl.AsDuration(),
				&redistool.ErrCacher[string]{
					Log:          config.Log,
					ErrRep:       modshared.ApiToErrReporter(config.Api),
					Client:       config.RedisClient,
					ErrMarshaler: prototool.ProtoErrMarshaler{},
					KeyToRedisKey: func(jobToken string) string {
						// Hash half of the token. Even if that hash leaks, it's not a big deal.
						// We do the same in api.AgentToken2key().
						n := len(jobToken) / 2
						tokenHash := sha256.Sum256([]byte(jobToken[:n]))
						return config.Config.Redis.KeyPrefix + ":allowed_agents_errs:" + string(tokenHash[:])
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
			listenerGracePeriod:  listenCfg.ListenGracePeriod.AsDuration(),
			shutdownTimeout:      shutdownTimeout,
		},
		listener: listener,
	}
	config.RegisterAgentApi(&rpc.KubernetesApi_ServiceDesc)
	return m, nil
}

func (f *Factory) Name() string {
	return kubernetes_api.ModuleName
}

func (f *Factory) StartStopPhase() modshared.ModuleStartStopPhase {
	// Start after servers because proxy uses agent connection (config.AgentConn), which works by accessing
	// in-memory private API server. So proxy needs to start after and stop before that server.
	return modshared.ModuleStartAfterServers
}
