package server

import (
	"strings"
	"time"

	"github.com/redis/rueidis"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/api"
	gapi "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitlab/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/gitops"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/gitops/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/cache"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/prototool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/redistool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/retry"
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
	s := newServerFromConfig(config, config.RedisClient, config.Api)
	rpc.RegisterGitopsServer(config.AgentServer, s)
	return &module{}, nil
}

func (f *Factory) Name() string {
	return gitops.ModuleName
}

func (f *Factory) StartStopPhase() modshared.ModuleStartStopPhase {
	return modshared.ModuleStartBeforeServers
}

func newServerFromConfig(config *modserver.Config, redisClient rueidis.Client, serverApi modserver.Api) *server {
	gitopsCfg := config.Config.Agent.Gitops
	return &server{
		serverApi:  serverApi,
		gitalyPool: config.Gitaly,
		projectInfoClient: &projectInfoClient{
			GitLabClient: config.GitLabClient,
			ProjectInfoCache: cache.NewWithError[projectInfoCacheKey, *api.ProjectInfo](
				gitopsCfg.ProjectInfoCacheTtl.AsDuration(),
				gitopsCfg.ProjectInfoCacheErrorTtl.AsDuration(),
				&redistool.ErrCacher[projectInfoCacheKey]{
					Log:          config.Log,
					ErrRep:       modshared.ApiToErrReporter(config.Api),
					Client:       redisClient,
					ErrMarshaler: prototool.ProtoErrMarshaler{},
					KeyToRedisKey: func(cacheKey projectInfoCacheKey) string {
						var result strings.Builder
						result.WriteString(config.Config.Redis.KeyPrefix)
						result.WriteString(":project_info_errs:")
						result.Write(api.AgentToken2key(cacheKey.agentToken))
						result.WriteByte(':')
						result.WriteString(cacheKey.projectId)
						return result.String()
					},
				},
				config.TraceProvider.Tracer(gitops.ModuleName),
				gapi.IsCacheableError,
			),
		},
		syncCount: config.UsageTracker.RegisterCounter(gitopsSyncCountKnownMetric),
		getObjectsPollConfig: retry.NewPollConfigFactory(gitopsCfg.PollPeriod.AsDuration(), retry.NewExponentialBackoffFactory(
			getObjectsToSynchronizeInitBackoff,
			getObjectsToSynchronizeMaxBackoff,
			getObjectsToSynchronizeResetDuration,
			getObjectsToSynchronizeBackoffFactor,
			getObjectsToSynchronizeJitter,
		)),
		maxManifestFileSize:      int64(gitopsCfg.MaxManifestFileSize),
		maxTotalManifestFileSize: int64(gitopsCfg.MaxTotalManifestFileSize),
		maxNumberOfPaths:         gitopsCfg.MaxNumberOfPaths,
		maxNumberOfFiles:         gitopsCfg.MaxNumberOfFiles,
	}
}
