package kasapp

import (
	"context"
	"errors"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/gitlab"
	gapi "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/gitlab/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/cache"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAgentRpcApi struct {
	modserver.RpcApi
	Token          api.AgentToken
	GitLabClient   gitlab.ClientInterface
	AgentInfoCache *cache.CacheWithErr
}

func (a *serverAgentRpcApi) AgentToken() api.AgentToken {
	return a.Token
}

func (a *serverAgentRpcApi) AgentInfo(ctx context.Context, log *zap.Logger) (*api.AgentInfo, error) {
	agentInfo, err := a.getAgentInfoCached(ctx)
	switch {
	case err == nil:
		return agentInfo, nil
	case errors.Is(err, context.Canceled):
		err = status.Error(codes.Canceled, err.Error())
	case errors.Is(err, context.DeadlineExceeded):
		err = status.Error(codes.DeadlineExceeded, err.Error())
	case gitlab.IsForbidden(err):
		err = status.Error(codes.PermissionDenied, "forbidden")
	case gitlab.IsUnauthorized(err):
		err = status.Error(codes.Unauthenticated, "unauthenticated")
	case gitlab.IsNotFound(err):
		err = status.Error(codes.NotFound, "agent not found")
	default:
		a.HandleProcessingError(log, modshared.NoAgentId, "AgentInfo()", err)
		err = status.Error(codes.Unavailable, "unavailable")
	}
	return nil, err
}

func (a *serverAgentRpcApi) getAgentInfoCached(ctx context.Context) (*api.AgentInfo, error) {
	agentInfo, err := a.AgentInfoCache.GetItem(ctx, a.Token, func() (interface{}, error) {
		return gapi.GetAgentInfo(ctx, a.GitLabClient, a.Token, gitlab.WithoutRetries())
	})
	if err != nil {
		return nil, err
	}
	return agentInfo.(*api.AgentInfo), nil
}

type serverAgentRpcApiFactory struct {
	rpcApiFactory  modserver.RpcApiFactory
	gitLabClient   gitlab.ClientInterface
	agentInfoCache *cache.CacheWithErr
}

func (f *serverAgentRpcApiFactory) New(ctx context.Context, fullMethodName string) (modserver.AgentRpcApi, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}
	return &serverAgentRpcApi{
		RpcApi:         f.rpcApiFactory(ctx, fullMethodName),
		Token:          api.AgentToken(token),
		GitLabClient:   f.gitLabClient,
		AgentInfoCache: f.agentInfoCache,
	}, nil
}
