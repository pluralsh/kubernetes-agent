package server

import (
	"context"
	"errors"
	"fmt"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitlab"
	gapi "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitlab/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/flux/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/cache"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/retry"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	rpc.UnimplementedGitLabFluxServer
	serverApi           modserver.Api
	pollCfgFactory      retry.PollConfigFactory
	projectAccessClient *projectAccessClient
}

func (s *server) ReconcileProjects(req *rpc.ReconcileProjectsRequest, server rpc.GitLabFlux_ReconcileProjectsServer) error {
	ctx := server.Context()
	rpcApi := modserver.AgentRpcApiFromContext(ctx)
	agentToken := rpcApi.AgentToken()
	projects := req.ToProjectSet()

	_ = retry.PollWithBackoff(ctx, s.pollCfgFactory(), func(ctx context.Context) (error, retry.AttemptResult) {
		log := rpcApi.Log()
		agentInfo, err := rpcApi.AgentInfo(ctx, log)
		if err != nil {
			if status.Code(err) == codes.Unavailable {
				return nil, retry.Backoff
			}
			return err, retry.Done // no wrap
		}

		log = log.With(logz.AgentId(agentInfo.Id))
		log.Debug("Started reconcile projects ...")
		defer log.Debug("Stopped reconcile projects ...")
		s.serverApi.OnGitPushEvent(ctx, func(ctx context.Context, message *modserver.Project) {
			if _, ok := projects[message.FullPath]; !ok {
				// NOTE: it's probably not a good idea to log here as we'd get one for every event,
				// which on GitLab.com is thousands per minute.
				return
			}

			hasAccess, err := s.verifyProjectAccess(ctx, log, rpcApi, agentInfo.Id, agentToken, message.FullPath)
			if err != nil {
				rpcApi.HandleProcessingError(log, agentInfo.Id, fmt.Sprintf("failed to check if project %s is accessible by agent", message.FullPath), err)
				return
			}
			if !hasAccess {
				return
			}

			err = server.Send(&rpc.ReconcileProjectsResponse{
				Project: &rpc.Project{Id: message.FullPath},
			})
			if err != nil {
				_ = rpcApi.HandleIoError(log, fmt.Sprintf("failed to send reconcile message for project %s", message.FullPath), err)
			}
		})

		return nil, retry.Done
	})

	return nil
}

// verifyProjectAccess verifies if the given agent has access to the given project.
// If this is not the case `false` is returned, otherwise `true`.
// If the error has the code Unavailable a caller my retry.
func (s *server) verifyProjectAccess(ctx context.Context, log *zap.Logger, rpcApi modserver.RpcApi, agentId int64,
	agentToken api.AgentToken, projectId string) (bool, error) {
	hasAccess, err := s.projectAccessClient.VerifyProjectAccess(ctx, agentToken, projectId)
	switch {
	case err == nil:
		return hasAccess, nil
	case errors.Is(err, context.Canceled):
		err = status.Error(codes.Canceled, err.Error())
	case errors.Is(err, context.DeadlineExceeded):
		err = status.Error(codes.DeadlineExceeded, err.Error())
	default:
		rpcApi.HandleProcessingError(log, agentId, "VerifyProjectAccess()", err)
		err = status.Error(codes.Unavailable, "unavailable")
	}
	return false, err
}

type projectAccessClient struct {
	gitLabClient       gitlab.ClientInterface
	projectAccessCache *cache.CacheWithErr[projectAccessCacheKey, bool]
}

func (c *projectAccessClient) VerifyProjectAccess(ctx context.Context, agentToken api.AgentToken, projectId string) (bool, error) {
	key := projectAccessCacheKey{agentToken: agentToken, projectId: projectId}
	return c.projectAccessCache.GetItem(ctx, key, func() (bool, error) {
		return gapi.VerifyProjectAccess(ctx, c.gitLabClient, agentToken, projectId, gitlab.WithoutRetries())
	})
}

type projectAccessCacheKey struct {
	agentToken api.AgentToken
	projectId  string
}
