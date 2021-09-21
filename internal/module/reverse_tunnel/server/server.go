package server

import (
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/reverse_tunnel"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/reverse_tunnel/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/retry"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	rpc.UnimplementedReverseTunnelServer
	tunnelHandler          reverse_tunnel.TunnelHandler
	getAgentInfoPollConfig retry.PollConfigFactory
}

func (s *server) Connect(server rpc.ReverseTunnel_ConnectServer) error {
	ctx := server.Context()
	ageCtx := grpctool.MaxConnectionAgeContextFromStreamContext(ctx)
	rpcApi := modserver.AgentRpcApiFromContext(ctx)
	log := rpcApi.Log()
	return rpcApi.PollWithBackoff(s.getAgentInfoPollConfig(), func() (error, retry.AttemptResult) {
		agentInfo, err := rpcApi.AgentInfo(ctx, log)
		if err != nil {
			if status.Code(err) == codes.Unavailable {
				return nil, retry.Backoff
			}
			return err, retry.Done // no wrap
		}
		return s.tunnelHandler.HandleTunnel(ageCtx, agentInfo, server), retry.Done
	})
}
