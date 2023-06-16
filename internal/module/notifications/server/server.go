package server

import (
	"context"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/notifications/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/event"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type gitPushPublisherFunc func(ctx context.Context, e *event.GitPushEvent) error

func newServer(gitPushPublisher gitPushPublisherFunc) *server {
	return &server{
		gitPushPublisher: gitPushPublisher,
	}
}

type server struct {
	rpc.UnimplementedNotificationsServer
	gitPushPublisher gitPushPublisherFunc
}

func (s *server) GitPushEvent(ctx context.Context, req *rpc.GitPushEventRequest) (*rpc.GitPushEventResponse, error) {
	err := s.gitPushPublisher(ctx, req.Event)
	if err != nil {
		rpcApi := modserver.RpcApiFromContext(ctx)
		rpcApi.HandleProcessingError(rpcApi.Log(), modshared.NoAgentId, "failed to publish received git push event", err)
		return nil, status.Errorf(codes.Unavailable, "Failed to publish received git push event: %v", err)
	}
	return &rpc.GitPushEventResponse{}, nil
}
