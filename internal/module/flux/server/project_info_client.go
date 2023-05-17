package server

import (
	"context"
	"errors"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitlab"
	gapi "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitlab/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/cache"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// FIXME: the code in this file is a temporary solution to check if a given agent has access to a given project.
// The code is copy & pasted from the /internal/modules/gitops module.
// We should refactor this once https://gitlab.com/gitlab-org/cluster-integration/gitlab-agent/-/issues/401
// is implemented.

// getProjectInfo returns an error with code Unavailable if there was a retriable error.
func (s *server) getProjectInfo(ctx context.Context, log *zap.Logger, rpcApi modserver.RpcApi, agentId int64,
	agentToken api.AgentToken, projectId string) (*api.ProjectInfo, error) {
	projectInfo, err := s.projectInfoClient.GetProjectInfo(ctx, agentToken, projectId)
	switch {
	case err == nil:
		return projectInfo, nil
	case errors.Is(err, context.Canceled):
		err = status.Error(codes.Canceled, err.Error())
	case errors.Is(err, context.DeadlineExceeded):
		err = status.Error(codes.DeadlineExceeded, err.Error())
	case gitlab.IsForbidden(err):
		err = status.Error(codes.PermissionDenied, "forbidden")
	case gitlab.IsUnauthorized(err):
		err = status.Error(codes.Unauthenticated, "unauthenticated")
	case gitlab.IsNotFound(err):
		err = status.Error(codes.NotFound, "project not found")
	default:
		rpcApi.HandleProcessingError(log, agentId, "GetProjectInfo()", err)
		err = status.Error(codes.Unavailable, "unavailable")
	}
	return nil, err
}

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
