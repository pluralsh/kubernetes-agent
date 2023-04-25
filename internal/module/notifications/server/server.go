package server

import (
	"context"
	"errors"
	"fmt"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/notifications/rpc"
)

const (
	gitPushEventsChannel = "git_push_events"
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
	rpcApi := modserver.RpcApiFromContext(ctx)
	log := rpcApi.Log()
	log.Sugar().Debugf("received git push notifications event for project %s", req.Project)
	err := s.publisher.Publish(ctx, gitPushEventsChannel, req.Project.ToNotificationsProject())
	if err != nil {
		if ioHandleErr := rpcApi.HandleIoError(log, "failed to publish received git push event", err); ioHandleErr != nil {
			return nil, fmt.Errorf("failed to handle io error %q: %w", err, ioHandleErr)
		}
		return nil, errors.New("failed to handle git push event")
	}
	return &rpc.GitPushEventResponse{}, nil
}
