package server

import (
	"github.com/pluralsh/kuberentes-agent/internal/module/modserver"
	"github.com/pluralsh/kuberentes-agent/internal/module/reverse_tunnel/rpc"
	"github.com/pluralsh/kuberentes-agent/internal/module/reverse_tunnel/tunnel"
	"github.com/pluralsh/kuberentes-agent/internal/tool/grpctool"
	"github.com/pluralsh/kuberentes-agent/internal/tool/retry"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	rpc.UnimplementedReverseTunnelServer
	tunnelHandler          tunnel.Handler
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
