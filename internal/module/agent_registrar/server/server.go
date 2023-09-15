package server

import (
	"context"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/agent_registrar/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/agent_tracker"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modserver"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	rpc.UnimplementedAgentRegistrarServer
	agentRegisterer agent_tracker.Registerer
}

func (s *server) Register(ctx context.Context, req *rpc.RegisterRequest) (*rpc.RegisterResponse, error) {
	rpcApi := modserver.AgentRpcApiFromContext(ctx)
	log := rpcApi.Log()

	// Get agent info
	agentInfo, err := rpcApi.AgentInfo(ctx, log)
	if err != nil {
		return nil, err
	}

	connectedAgentInfo := &agent_tracker.ConnectedAgentInfo{
		AgentMeta:    req.AgentMeta,
		ConnectedAt:  timestamppb.Now(),
		ConnectionId: req.PodId,
		AgentId:      agentInfo.Id,
		ProjectId:    agentInfo.ProjectId,
	}

	// Register agent
	err = s.agentRegisterer.RegisterConnection(ctx, connectedAgentInfo)
	if err != nil {
		rpcApi.HandleProcessingError(log, agentInfo.Id, "Failed to register agent", err)
		return nil, status.Error(codes.Unavailable, "Failed to register agent")
	}

	return &rpc.RegisterResponse{}, nil
}
