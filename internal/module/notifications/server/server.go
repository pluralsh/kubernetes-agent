package server

import (
	"context"
	"errors"
	"fmt"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modserver/notifications"
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
	if err := s.publisher.Publish(ctx, notifications.GitPushEventsChannel, req.Project.ToNotificationsProject()); err != nil {
		rpcApi := modserver.RpcApiFromContext(ctx)
		if ioHandleErr := rpcApi.HandleIoError(rpcApi.Log(), "failed to publish received git push event", err); ioHandleErr != nil {
			return nil, fmt.Errorf("failed to handle io error %q: %w", err, ioHandleErr)
		}
		return nil, errors.New("failed to handle git push event")
	}
	return &rpc.GitPushEventResponse{}, nil
}
