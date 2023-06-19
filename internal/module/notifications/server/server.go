package server

import (
	"context"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/notifications/rpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type PublishEvent func(ctx context.Context, e proto.Message) error

func newServer(publishEvent PublishEvent) *server {
	return &server{
		publishEvent: publishEvent,
	}
}

type server struct {
	rpc.UnimplementedNotificationsServer
	publishEvent PublishEvent
}

func (s *server) GitPushEvent(ctx context.Context, req *rpc.GitPushEventRequest) (*rpc.GitPushEventResponse, error) {
	err := s.publishEvent(ctx, req.Event)
	if err != nil {
		rpcApi := modserver.RpcApiFromContext(ctx)
		rpcApi.HandleProcessingError(rpcApi.Log(), modshared.NoAgentId, "Failed to publish received git push event", err)
		return nil, status.Errorf(codes.Unavailable, "Failed to publish received git push event: %v", err)
	}
	return &rpc.GitPushEventResponse{}, nil
}
