package server

import (
	"context"
	"fmt"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/flux/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/retry"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	rpc.UnimplementedGitLabFluxServer
	serverApi         modserver.Api
	pollCfgFactory    retry.PollConfigFactory
	projectInfoClient *projectInfoClient
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

			// FIXME: actually check if that project is allowed to be accessed by this agent ...
			// I think we cannot really check access when the receiver is created (or we receive a new project)
			// in the server, because the user could grant access after the fact.
			// So, here is actually a not too bad place.
			// We could use a redis cache to cache access allowance for some time, so that we don't spam Rails.
			// See comments in `project_info_client.go` and https://gitlab.com/gitlab-org/cluster-integration/gitlab-agent/-/issues/401
			_, err = s.getProjectInfo(ctx, log, rpcApi, agentInfo.Id, agentToken, message.FullPath)
			if err != nil {
				rpcApi.HandleProcessingError(log, agentInfo.Id, fmt.Sprintf("failed to check if project %s is accessible by agent", message.FullPath), err)
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
