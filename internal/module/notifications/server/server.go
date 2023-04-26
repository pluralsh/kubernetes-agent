package server

import (
	"context"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/notifications/rpc"
)

type server struct {
	rpc.UnimplementedNotificationsServer
}

func (s *server) GitPushEvent(ctx context.Context, req *rpc.GitPushEventRequest) (*rpc.GitPushEventResponse, error) {
	rpcApi := modserver.RpcApiFromContext(ctx)
	log := rpcApi.Log()
	log.Sugar().Debugf("received git push notifications event for project %s", req.Project)
	return &rpc.GitPushEventResponse{}, nil
}
