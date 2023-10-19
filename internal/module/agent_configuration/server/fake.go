package server

import (
	"context"
	"errors"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitaly"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitlab"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/agent_configuration/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/agent_tracker"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/errz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/mathz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/retry"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/syncz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/agentcfg"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type fakeServer struct {
	rpc.UnimplementedAgentConfigurationServer
	serverApi                  modserver.Api
	gitaly                     gitaly.PoolInterface
	gitLabClient               gitlab.ClientInterface
	agentRegisterer            agent_tracker.Registerer
	maxConfigurationFileSize   int64
	getConfigurationPollConfig retry.PollConfigFactory
	gitLabExternalUrl          string
}

func (s *fakeServer) GetConfiguration(req *rpc.ConfigurationRequest, server rpc.AgentConfiguration_GetConfigurationServer) error {
	connectedAgentInfo := &agent_tracker.ConnectedAgentInfo{
		AgentMeta:    req.AgentMeta,
		ConnectedAt:  timestamppb.Now(),
		ConnectionId: mathz.Int63(),
	}
	ctx := server.Context()
	rpcApi := modserver.AgentRpcApiFromContext(ctx)
	log := rpcApi.Log()
	defer s.maybeUnregisterAgent(log, rpcApi, connectedAgentInfo, req.SkipRegister)

	pollCfg := s.getConfigurationPollConfig()

	wh := syncz.NewComparableWorkerHolder[string](
		func(projectId string) syncz.Worker {
			return syncz.WorkerFunc(func(ctx context.Context) {

			})
		},
	)
	defer wh.StopAndWait()

	return rpcApi.PollWithBackoff(pollCfg, func() (error, retry.AttemptResult) {
		// This call is made on each poll because:
		// - it checks that the agent's token is still valid
		// - repository location in Gitaly might have changed
		agentInfo, err := rpcApi.AgentInfo(ctx, log)
		if err != nil {
			if status.Code(err) == codes.Unavailable {
				return nil, retry.Backoff
			}
			return err, retry.Done
		}

		// re-define log to avoid accidentally using the old one
		log := log.With(logz.AgentId(agentInfo.Id)) // nolint:govet
		s.maybeRegisterAgent(ctx, log, rpcApi, connectedAgentInfo, agentInfo, req.SkipRegister)

		_, err = s.fetchConfiguration(ctx, agentInfo)
		if err != nil {
			rpcApi.HandleProcessingError(log, agentInfo.Id, "Config: failed to fetch", err)
			var ue errz.UserError
			if errors.As(err, &ue) {
				// return the error to the client because it's a user error
				return status.Errorf(codes.FailedPrecondition, "Config: %v", err), retry.Done
			}
			return nil, retry.Backoff
		}

		return nil, retry.Continue
	})
}

func (s *fakeServer) poll(ctx context.Context, agentInfo *api.AgentInfo, lastProcessedCommitId string) (*gitaly.PollInfo, error) {
	p, err := s.gitaly.Poller(ctx, agentInfo.GitalyInfo)
	if err != nil {
		return nil, err
	}
	return p.Poll(ctx, agentInfo.Repository, lastProcessedCommitId, "refs/heads/"+agentInfo.DefaultBranch)
}

func (s *fakeServer) sendConfigResponse(server rpc.AgentConfiguration_GetConfigurationServer,
	agentInfo *api.AgentInfo, configFile *agentcfg.ConfigurationFile, commitId string) error {
	return server.Send(&rpc.ConfigurationResponse{
		Configuration: &agentcfg.AgentConfiguration{
			Gitops:            configFile.Gitops,
			Observability:     configFile.Observability,
			AgentId:           agentInfo.Id,
			ProjectId:         agentInfo.ProjectId,
			ProjectPath:       agentInfo.Repository.GlProjectPath,
			CiAccess:          configFile.CiAccess,
			ContainerScanning: configFile.ContainerScanning,
			RemoteDevelopment: configFile.RemoteDevelopment,
			Flux:              configFile.Flux,
			GitlabExternalUrl: s.gitLabExternalUrl,
		},
		CommitId: commitId,
	})
}

// fetchConfiguration fetches agent's configuration from a corresponding repository.
// Assumes configuration is stored in ".gitlab/agents/<agent id>/config.yaml" file.
// fetchConfiguration returns a wrapped context.Canceled, context.DeadlineExceeded or gRPC error if ctx signals done and interrupts a running gRPC call.
func (s *fakeServer) fetchConfiguration(ctx context.Context, agentInfo *api.AgentInfo) (*agentcfg.ConfigurationFile, error) {

	return &agentcfg.ConfigurationFile{
		Gitops:            nil,
		Observability:     nil,
		CiAccess:          nil,
		ContainerScanning: nil,
		UserAccess: &agentcfg.UserAccessCF{
			AccessAs: nil,
			Projects: nil,
			Groups:   nil,
		},
		RemoteDevelopment: nil,
		Flux:              nil,
	}, nil
}

func (s *fakeServer) maybeRegisterAgent(ctx context.Context, log *zap.Logger, rpcApi modserver.AgentRpcApi,
	connectedAgentInfo *agent_tracker.ConnectedAgentInfo, agentInfo *api.AgentInfo, skipRegister bool) {
	// Skip registering agent if skipRegister is true. The agent will call "Register" gRPC method instead.
	if skipRegister {
		return
	}

	if connectedAgentInfo.AgentId != 0 {
		return
	}
	connectedAgentInfo.AgentId = agentInfo.Id
	connectedAgentInfo.ProjectId = agentInfo.ProjectId
	err := s.agentRegisterer.RegisterConnection(ctx, connectedAgentInfo)
	if err != nil {
		rpcApi.HandleProcessingError(log, agentInfo.Id, "Failed to register agent", err)
	}
}

func (s *fakeServer) maybeUnregisterAgent(log *zap.Logger, rpcApi modserver.AgentRpcApi,
	connectedAgentInfo *agent_tracker.ConnectedAgentInfo, skipRegister bool) {
	// Skip unregistering agent if skipRegister is true. GC will clean up the agent from the storage.
	if skipRegister {
		return
	}

	if connectedAgentInfo.AgentId == 0 {
		return
	}
	err := s.agentRegisterer.UnregisterConnection(context.Background(), connectedAgentInfo)
	if err != nil {
		rpcApi.HandleProcessingError(log, connectedAgentInfo.AgentId, "Failed to unregister agent", err)
	}
}
