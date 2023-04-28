package server

import (
	"context"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/notifications/rpc"
)

func newServer(publisher Publisher) *server {
	return &server{
		publisher: publisher,
	}
}

type server struct {
	rpc.UnimplementedNotificationsServer
	publisher Publisher
}

func (s *server) GitPushEvent(ctx context.Context, req *rpc.GitPushEventRequest) (*rpc.GitPushEventResponse, error) {
	if err := s.publisher.Publish(ctx, modserver.GitPushEventsChannel, req.Project.ToNotificationsProject()); err != nil {
		rpcApi := modserver.RpcApiFromContext(ctx)
		return nil, rpcApi.HandleIoError(rpcApi.Log(), "failed to publish received git push event", err)
	}
	return &rpc.GitPushEventResponse{}, nil
}
