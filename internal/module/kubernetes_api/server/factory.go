package server

import (
	"crypto/sha256"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"net"

	gapi "github.com/pluralsh/kuberentes-agent/internal/gitlab/api"
	"github.com/pluralsh/kuberentes-agent/internal/module/kubernetes_api"
	"github.com/pluralsh/kuberentes-agent/internal/module/kubernetes_api/rpc"
	"github.com/pluralsh/kuberentes-agent/internal/module/modserver"
	"github.com/pluralsh/kuberentes-agent/internal/module/modshared"
	"github.com/pluralsh/kuberentes-agent/internal/tool/cache"
	"github.com/pluralsh/kuberentes-agent/internal/tool/prototool"
	"github.com/pluralsh/kuberentes-agent/internal/tool/redistool"
	"github.com/pluralsh/kuberentes-agent/internal/tool/tlstool"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
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
	var allowedOriginUrls []string
	if u := config.Config.Gitlab.GetExternalUrl(); u != "" {
		allowedOriginUrls = append(allowedOriginUrls, u)
	}
	allowedAgentCacheTtl := k8sApi.AllowedAgentCacheTtl.AsDuration()
	allowedAgentCacheErrorTtl := k8sApi.AllowedAgentCacheErrorTtl.AsDuration()
	tracer := config.TraceProvider.Tracer(kubernetes_api.ModuleName)
	m := &module{
		log: config.Log,
		proxy: kubernetesApiProxy{
			log:                 config.Log,
			api:                 config.Api,
			kubernetesApiClient: rpc.NewKubernetesApiClient(config.AgentConn),
			allowedOriginUrls:   allowedOriginUrls,
			allowedAgentsCache: cache.NewWithError[string, *gapi.AllowedAgentsForJob](
				allowedAgentCacheTtl,
				allowedAgentCacheErrorTtl,
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
				tracer,
				gapi.IsCacheableError,
			),
			authorizeProxyUserCache: cache.NewWithError[proxyUserCacheKey, *gapi.AuthorizeProxyUserResponse](
				allowedAgentCacheTtl,
				allowedAgentCacheErrorTtl,
				&redistool.ErrCacher[proxyUserCacheKey]{
					Log:           config.Log,
					ErrRep:        modshared.ApiToErrReporter(config.Api),
					Client:        config.RedisClient,
					ErrMarshaler:  prototool.ProtoErrMarshaler{},
					KeyToRedisKey: getAuthorizedProxyUserCacheKey(config.Config.Redis.KeyPrefix),
				},
				tracer,
				gapi.IsCacheableError,
			),
			responseSerializer:  serializer.NewCodecFactory(runtime.NewScheme()),
			traceProvider:       config.TraceProvider,
			tracePropagator:     config.TracePropagator,
			meterProvider:       config.MeterProvider,
			serverName:          serverName,
			serverVia:           "gRPC/1.0 " + serverName,
			urlPathPrefix:       k8sApi.UrlPathPrefix,
			listenerGracePeriod: listenCfg.ListenGracePeriod.AsDuration(),
			shutdownGracePeriod: listenCfg.ShutdownGracePeriod.AsDuration(),
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

func getAuthorizedProxyUserCacheKey(redisKeyPrefix string) redistool.KeyToRedisKey[proxyUserCacheKey] {
	return func(key proxyUserCacheKey) string {
		// Hash half of the token. Even if that hash leaks, it's not a big deal.
		// We do the same in api.AgentToken2key().
		n := len(key.accessKey) / 2

		// Use delimiters between fields to ensure hash of "ab" + "c" is different from "a" + "bc".
		h := sha256.New()
		id := make([]byte, 8)
		binary.LittleEndian.PutUint64(id, uint64(key.agentId))
		h.Write(id)
		// Don't need a delimiter here because id is fixed size in bytes
		h.Write([]byte(key.accessType))
		h.Write([]byte{11}) // delimiter
		h.Write([]byte(key.accessKey[:n]))
		h.Write([]byte{11}) // delimiter
		h.Write([]byte(key.csrfToken))
		tokenHash := h.Sum(nil)
		return redisKeyPrefix + ":auth_proxy_user_errs:" + string(tokenHash)
	}
}
