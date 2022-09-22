package server

import (
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/api"
	gapi "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/gitlab/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/gitops"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/gitops/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/cache"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/metric"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/prototool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/redistool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/retry"
)

const (
	getObjectsToSynchronizeInitBackoff   = 10 * time.Second
	getObjectsToSynchronizeMaxBackoff    = 5 * time.Minute
	getObjectsToSynchronizeResetDuration = 10 * time.Minute
	getObjectsToSynchronizeBackoffFactor = 2.0
	getObjectsToSynchronizeJitter        = 1.0

	gitopsSyncCountKnownMetric = "gitops_sync"
)

type Factory struct {
}

func (f *Factory) New(config *modserver.Config) (modserver.Module, error) {
	s, err := newServerFromConfig(config, config.RedisClient)
	if err != nil {
		return nil, err
	}
	rpc.RegisterGitopsServer(config.AgentServer, s)
	return &module{}, nil
}

func (f *Factory) Name() string {
	return gitops.ModuleName
}

func newServerFromConfig(config *modserver.Config, redisClient redis.UniversalClient) (*server, error) {
	gitOpsPollIntervalHistogram := constructGitOpsPollIntervalHistogram()
	err := metric.Register(config.Registerer, gitOpsPollIntervalHistogram)
	if err != nil {
		return nil, err
	}
	gitops := config.Config.Agent.Gitops
	return &server{
		gitalyPool: config.Gitaly,
		projectInfoClient: &projectInfoClient{
			GitLabClient: config.GitLabClient,
			ProjectInfoCache: cache.NewWithError(
				gitops.ProjectInfoCacheTtl.AsDuration(),
				gitops.ProjectInfoCacheErrorTtl.AsDuration(),
				&redistool.ErrCacher{
					Log:          config.Log,
					Client:       redisClient,
					ErrMarshaler: prototool.ProtoErrMarshaler{},
					KeyToRedisKey: func(agentToken interface{}) string {
						key := api.AgentToken2key(agentToken.(api.AgentToken))
						return config.Config.Redis.KeyPrefix + ":project_info_errs:" + string(key)
					},
				},
				gapi.IsCacheableError,
			),
		},
		syncCount:                   config.UsageTracker.RegisterCounter(gitopsSyncCountKnownMetric),
		gitOpsPollIntervalHistogram: gitOpsPollIntervalHistogram,
		getObjectsPollConfig: retry.NewPollConfigFactory(gitops.PollPeriod.AsDuration(), retry.NewExponentialBackoffFactory(
			getObjectsToSynchronizeInitBackoff,
			getObjectsToSynchronizeMaxBackoff,
			getObjectsToSynchronizeResetDuration,
			getObjectsToSynchronizeBackoffFactor,
			getObjectsToSynchronizeJitter,
		)),
		maxManifestFileSize:      int64(gitops.MaxManifestFileSize),
		maxTotalManifestFileSize: int64(gitops.MaxTotalManifestFileSize),
		maxNumberOfPaths:         gitops.MaxNumberOfPaths,
		maxNumberOfFiles:         gitops.MaxNumberOfFiles,
	}, nil
}

func constructGitOpsPollIntervalHistogram() prometheus.Histogram {
	return prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "gitops_poll_interval",
		Help:    "The time between kas calls to Gitaly to look for GitOps updates",
		Buckets: prometheus.LinearBuckets(20, 20, 5), // 5 buckets (20, 40, 60, 80, 100)
	})
}
