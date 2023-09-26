package kasapp

import (
	"time"

	agent_configuration_server "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/agent_configuration/server"
	gitops_server "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/gitops/server"
	google_profiler_server "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/google_profiler/server"
	kubernetes_api_server "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/kubernetes_api/server"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modserver"
	observability_server "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/observability/server"
	usage_metrics_server "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/usage_metrics/server"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/prototool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/kascfg"
)

const (
	defaultGitLabApiRateLimitRefillRate = 50.0
	defaultGitLabApiRateLimitBucketSize = 250

	defaultListenGracePeriod = 5 * time.Second

	defaultAgentInfoCacheTTL         = 5 * time.Minute
	defaultAgentInfoCacheErrorTTL    = 1 * time.Minute
	defaultAgentRedisConnInfoTTL     = 5 * time.Minute
	defaultAgentRedisConnInfoRefresh = 4 * time.Minute
	defaultAgentRedisConnInfoGC      = 10 * time.Minute

	defaultAgentListenNetwork                      = "tcp"
	defaultAgentListenAddress                      = "127.0.0.1:8150"
	defaultAgentListenConnectionsPerTokenPerMinute = 40000
	defaultAgentListenMaxConnectionAge             = 2 * time.Hour

	defaultGitalyGlobalApiRefillRate    = 30.0
	defaultGitalyGlobalApiBucketSize    = 70
	defaultGitalyPerServerApiRate       = 15.0
	defaultGitalyPerServerApiBucketSize = 40

	defaultRedisDialTimeout  = 5 * time.Second
	defaultRedisWriteTimeout = 3 * time.Second
	defaultRedisKeyPrefix    = "gitlab-kas"
	defaultRedisNetwork      = "tcp"

	defaultApiListenNetwork          = "tcp"
	defaultApiListenAddress          = "127.0.0.1:8153"
	defaultApiListenMaxConnectionAge = 2 * time.Hour

	defaultPrivateApiListenNetwork = "tcp"
	defaultPrivateApiListenAddress = "127.0.0.1:8155"
	// Should be equal to the defaultAgentListenMaxConnectionAge as agent's tunnel requests go via private API server.
	defaultPrivateApiListenMaxConnectionAge = defaultAgentListenMaxConnectionAge
)

var (
	defaulters = []modserver.ApplyDefaults{
		observability_server.ApplyDefaults,
		google_profiler_server.ApplyDefaults,
		agent_configuration_server.ApplyDefaults,
		usage_metrics_server.ApplyDefaults,
		gitops_server.ApplyDefaults,
		kubernetes_api_server.ApplyDefaults,
	}
)

func ApplyDefaultsToKasConfigurationFile(cfg *kascfg.ConfigurationFile) {
	prototool.NotNil(&cfg.Gitlab)
	defaultGitLab(cfg.Gitlab)

	prototool.NotNil(&cfg.Agent)
	defaultAgent(cfg.Agent)

	prototool.NotNil(&cfg.Gitaly)
	defaultGitaly(cfg.Gitaly)

	prototool.NotNil(&cfg.Redis)
	defaultRedis(cfg.Redis)

	prototool.NotNil(&cfg.Api)
	defaultApi(cfg.Api)

	prototool.NotNil(&cfg.PrivateApi)
	defaultPrivateApi(cfg.PrivateApi)

	for _, defaulter := range defaulters {
		defaulter(cfg)
	}
}

func defaultApi(api *kascfg.ApiCF) {
	prototool.NotNil(&api.Listen)
	prototool.StringPtr(&api.Listen.Network, defaultApiListenNetwork)
	prototool.String(&api.Listen.Address, defaultApiListenAddress)
	prototool.Duration(&api.Listen.MaxConnectionAge, defaultApiListenMaxConnectionAge)
	prototool.Duration(&api.Listen.ListenGracePeriod, defaultListenGracePeriod)
}

func defaultPrivateApi(api *kascfg.PrivateApiCF) {
	prototool.NotNil(&api.Listen)
	prototool.StringPtr(&api.Listen.Network, defaultPrivateApiListenNetwork)
	prototool.String(&api.Listen.Address, defaultPrivateApiListenAddress)
	prototool.Duration(&api.Listen.MaxConnectionAge, defaultPrivateApiListenMaxConnectionAge)
	prototool.Duration(&api.Listen.ListenGracePeriod, defaultListenGracePeriod)
}

func defaultGitLab(g *kascfg.GitLabCF) {
	prototool.NotNil(&g.ApiRateLimit)
	prototool.Float64(&g.ApiRateLimit.RefillRatePerSecond, defaultGitLabApiRateLimitRefillRate)
	prototool.Uint32(&g.ApiRateLimit.BucketSize, defaultGitLabApiRateLimitBucketSize)
}

func defaultAgent(a *kascfg.AgentCF) {
	prototool.NotNil(&a.Listen)
	prototool.StringPtr(&a.Listen.Network, defaultAgentListenNetwork)
	prototool.String(&a.Listen.Address, defaultAgentListenAddress)
	prototool.Uint32(&a.Listen.ConnectionsPerTokenPerMinute, defaultAgentListenConnectionsPerTokenPerMinute)
	prototool.Duration(&a.Listen.MaxConnectionAge, defaultAgentListenMaxConnectionAge)
	prototool.Duration(&a.Listen.ListenGracePeriod, defaultListenGracePeriod)

	prototool.Duration(&a.InfoCacheTtl, defaultAgentInfoCacheTTL)
	prototool.Duration(&a.InfoCacheErrorTtl, defaultAgentInfoCacheErrorTTL)
	prototool.Duration(&a.RedisConnInfoTtl, defaultAgentRedisConnInfoTTL)
	prototool.Duration(&a.RedisConnInfoRefresh, defaultAgentRedisConnInfoRefresh)
	prototool.Duration(&a.RedisConnInfoGc, defaultAgentRedisConnInfoGC)
}

func defaultGitaly(g *kascfg.GitalyCF) {
	prototool.NotNil(&g.GlobalApiRateLimit)
	prototool.Float64(&g.GlobalApiRateLimit.RefillRatePerSecond, defaultGitalyGlobalApiRefillRate)
	prototool.Uint32(&g.GlobalApiRateLimit.BucketSize, defaultGitalyGlobalApiBucketSize)

	prototool.NotNil(&g.PerServerApiRateLimit)
	prototool.Float64(&g.PerServerApiRateLimit.RefillRatePerSecond, defaultGitalyPerServerApiRate)
	prototool.Uint32(&g.PerServerApiRateLimit.BucketSize, defaultGitalyPerServerApiBucketSize)
}

func defaultRedis(r *kascfg.RedisCF) {
	prototool.Duration(&r.DialTimeout, defaultRedisDialTimeout)
	prototool.Duration(&r.WriteTimeout, defaultRedisWriteTimeout)
	prototool.String(&r.KeyPrefix, defaultRedisKeyPrefix)
	prototool.String(&r.Network, defaultRedisNetwork)
}
