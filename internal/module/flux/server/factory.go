package server

import (
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/api"
	gapi "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitlab/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/flux"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/flux/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/cache"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/metric"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/prototool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/redistool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/retry"
)

const (
	reconcileProjectsInitBackoff   = 10 * time.Second
	reconcileProjectsMaxBackoff    = 5 * time.Minute
	reconcileProjectsResetDuration = 10 * time.Minute
	reconcileProjectsBackoffFactor = 2.0
	reconcileProjectsJitter        = 1.0
	projectAccessCacheTtl          = 5 * time.Minute
	projectAccessCacheErrTtl       = 1 * time.Minute

	fluxNotifiedCounterMetricName             = "flux_git_push_notifications_total"
	fluxDroppedNotificationsCounterMetricName = "flux_dropped_git_push_notifications_total"
	fluxNotifiedProjectsCounterMetricName     = "flux_git_push_notified_unique_projects"
)

type Factory struct {
}

func (f *Factory) New(config *modserver.Config) (modserver.Module, error) {
	promNotifiedCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: fluxNotifiedCounterMetricName,
		Help: "The total number of sent Git Push notifications to agentks in Flux module",
	})
	promDroppedCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: fluxDroppedNotificationsCounterMetricName,
		Help: "The total number of dropped Git push notifications in Flux module",
	})
	if err := metric.Register(config.Registerer, promNotifiedCounter, promDroppedCounter); err != nil {
		return nil, err
	}
	rpc.RegisterGitLabFluxServer(config.AgentServer, &server{
		serverApi:               config.Api,
		notifiedCounter:         metric.NewMultiCounter(promNotifiedCounter, config.UsageTracker.RegisterCounter(fluxNotifiedCounterMetricName)),
		notifiedProjectsCounter: config.UsageTracker.RegisterUniqueCounter(fluxNotifiedProjectsCounterMetricName),
		droppedCounter:          promDroppedCounter,
		pollCfgFactory: retry.NewPollConfigFactory(0, retry.NewExponentialBackoffFactory(
			reconcileProjectsInitBackoff, reconcileProjectsMaxBackoff, reconcileProjectsResetDuration, reconcileProjectsBackoffFactor, reconcileProjectsJitter)),
		projectAccessClient: &projectAccessClient{
			gitLabClient: config.GitLabClient,
			projectAccessCache: cache.NewWithError[projectAccessCacheKey, bool](
				projectAccessCacheTtl,
				projectAccessCacheErrTtl,
				&redistool.ErrCacher[projectAccessCacheKey]{
					Log:          config.Log,
					ErrRep:       modshared.ApiToErrReporter(config.Api),
					Client:       config.RedisClient,
					ErrMarshaler: prototool.ProtoErrMarshaler{},
					KeyToRedisKey: func(cacheKey projectAccessCacheKey) string {
						var result strings.Builder
						result.WriteString(config.Config.Redis.KeyPrefix)
						result.WriteString(":verify_project_access_errs:")
						result.Write(api.AgentToken2key(cacheKey.agentToken))
						result.WriteByte(':')
						result.WriteString(cacheKey.projectId)
						return result.String()
					},
				},
				config.TraceProvider.Tracer(flux.ModuleName),
				gapi.IsCacheableError,
			),
		},
	})
	return &module{}, nil
}

func (f *Factory) Name() string {
	return flux.ModuleName
}

func (f *Factory) StartStopPhase() modshared.ModuleStartStopPhase {
	return modshared.ModuleStartBeforeServers
}
