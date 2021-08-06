package server

import (
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/reverse_tunnel"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/reverse_tunnel/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/grpctool"
)

type server struct {
	rpc.UnimplementedReverseTunnelServer
	tunnelHandler reverse_tunnel.TunnelHandler
}

func (s *server) Connect(server rpc.ReverseTunnel_ConnectServer) error {
	ctx := server.Context()
	log := grpctool.LoggerFromContext(ctx)
	ageCtx := grpctool.MaxConnectionAgeContextFromStreamContext(ctx)
	rpcApi := modserver.RpcApiFromContext(ctx)
	agentInfo, err := rpcApi.AgentInfo(ageCtx, log)
	if err != nil {
		return err // no wrap
	}
	return s.tunnelHandler.HandleTunnel(ageCtx, agentInfo, server)
}
