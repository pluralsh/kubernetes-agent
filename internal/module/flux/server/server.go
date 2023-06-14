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
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/metric"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/retry"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	maxBufferedNotifications = 10
)

type server struct {
	rpc.UnimplementedGitLabFluxServer
	serverApi           modserver.Api
	notifiedCounter     metric.Counter
	droppedCounter      metric.Counter
	pollCfgFactory      retry.PollConfigFactory
	projectAccessClient *projectAccessClient
}

func (s *server) ReconcileProjects(req *rpc.ReconcileProjectsRequest, server rpc.GitLabFlux_ReconcileProjectsServer) error {
	ctx := server.Context()
	rpcApi := modserver.AgentRpcApiFromContext(ctx)
	log := rpcApi.Log()
	var agentInfo *api.AgentInfo
	var err error

	err = rpcApi.PollWithBackoff(s.pollCfgFactory(), func() (error, retry.AttemptResult) {
		agentInfo, err = rpcApi.AgentInfo(ctx, log)
		if err != nil {
			if status.Code(err) == codes.Unavailable {
				return nil, retry.Backoff
			}
			return err, retry.Done // no wrap
		}
		return nil, retry.Done
	})
	if agentInfo == nil {
		return err // ctx done, err may be nil but must return
	}

	log = log.With(logz.AgentId(agentInfo.Id))

	pipe := make(chan *modserver.Project, maxBufferedNotifications)
	var wg wait.Group
	defer wg.Wait()
	defer close(pipe)
	wg.Start(func() {
		for project := range pipe {
			hasAccess, err := s.verifyProjectAccess(ctx, log, rpcApi, agentInfo.Id, project.FullPath)
			if err != nil {
				rpcApi.HandleProcessingError(log, agentInfo.Id, fmt.Sprintf("Failed to check if project %s is accessible by agent", project.FullPath), err)
				continue
			}
			if !hasAccess {
				continue
			}

			// increase Flux Git push event notification counter
			s.notifiedCounter.Inc()

			err = server.Send(&rpc.ReconcileProjectsResponse{
				Project: &rpc.Project{Id: project.FullPath},
			})
			if err != nil {
				_ = rpcApi.HandleIoError(log, fmt.Sprintf("Failed to send reconcile message for project %s", project.FullPath), err)
			}
		}
	})

	log.Debug("Started reconcile projects ...")
	defer log.Debug("Stopped reconcile projects ...")
	projects := req.ToProjectSet()
	s.serverApi.OnGitPushEvent(ctx, func(ctx context.Context, project *modserver.Project) {
		if _, ok := projects[project.FullPath]; !ok {
			// NOTE: it's probably not a good idea to log here as we'd get one for every event,
			// which on GitLab.com is thousands per minute.
			return
		}

		select {
		case pipe <- project:
		default:
			s.droppedCounter.Inc()
			// NOTE: if for whatever reason the other goroutine isn't able to keep up with the events,
			// we just drop them for now.
			log.Debug("Dropping Git push event", logz.ProjectId(project.FullPath))
		}
	})

	return nil
}

// verifyProjectAccess verifies if the given agent has access to the given project.
// If this is not the case `false` is returned, otherwise `true`.
// If the error has the code Unavailable a caller my retry.
func (s *server) verifyProjectAccess(ctx context.Context, log *zap.Logger, rpcApi modserver.AgentRpcApi, agentId int64,
	projectId string) (bool, error) {
	hasAccess, err := s.projectAccessClient.VerifyProjectAccess(ctx, rpcApi.AgentToken(), projectId)
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
