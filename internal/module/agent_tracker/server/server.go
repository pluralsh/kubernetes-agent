package server

import (
	"context"

	"github.com/pluralsh/kuberentes-agent/internal/module/agent_tracker"
	"github.com/pluralsh/kuberentes-agent/internal/module/agent_tracker/rpc"
	"github.com/pluralsh/kuberentes-agent/internal/module/modserver"
	"github.com/pluralsh/kuberentes-agent/internal/module/modshared"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	rpc.UnimplementedAgentTrackerServer
	agentQuerier agent_tracker.Querier
}

func (s *server) GetConnectedAgents(ctx context.Context, req *rpc.GetConnectedAgentsRequest) (*rpc.GetConnectedAgentsResponse, error) {
	rpcApi := modserver.RpcApiFromContext(ctx)
	log := rpcApi.Log()
	switch v := req.GetRequest().(type) {
	case *rpc.GetConnectedAgentsRequest_AgentId:
		var infos agent_tracker.ConnectedAgentInfoCollector
		err := s.agentQuerier.GetConnectionsByAgentId(ctx, v.AgentId, infos.Collect)
		if err != nil {
			rpcApi.HandleProcessingError(log, modshared.NoAgentId, "GetConnectionsByAgentId() failed", err)
			return nil, status.Error(codes.Unavailable, "GetConnectionsByAgentId() failed")
		}
		return &rpc.GetConnectedAgentsResponse{
			Agents: infos,
		}, nil
	case *rpc.GetConnectedAgentsRequest_ProjectId:
		var infos agent_tracker.ConnectedAgentInfoCollector
		err := s.agentQuerier.GetConnectionsByProjectId(ctx, v.ProjectId, infos.Collect)
		if err != nil {
			rpcApi.HandleProcessingError(log, modshared.NoAgentId, "GetConnectionsByProjectId() failed", err)
			return nil, status.Error(codes.Unavailable, "GetConnectionsByProjectId() failed")
		}
		return &rpc.GetConnectedAgentsResponse{
			Agents: infos,
		}, nil
	default:
		// Should never happen
		return nil, status.Errorf(codes.InvalidArgument, "Unexpected field type: %T", req.GetRequest())
	}
}
