package server

import (
	"context"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitlab"
	gapi "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitlab/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/cache"
)

type projectInfoClient struct {
	GitLabClient     gitlab.ClientInterface
	ProjectInfoCache *cache.CacheWithErr[projectInfoCacheKey, *api.ProjectInfo]
}

func (c *projectInfoClient) GetProjectInfo(ctx context.Context, agentToken api.AgentToken, projectId string) (*api.ProjectInfo, error) {
	key := projectInfoCacheKey{agentToken: agentToken, projectId: projectId}
	return c.ProjectInfoCache.GetItem(ctx, key, func() (*api.ProjectInfo, error) {
		return gapi.GetProjectInfo(ctx, c.GitLabClient, agentToken, projectId, gitlab.WithoutRetries())
	})
}

type projectInfoCacheKey struct {
	agentToken api.AgentToken
	projectId  string
}
