package kasapp

import (
	"context"
	"errors"
	"sync"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/pluralsh/kuberentes-agent/internal/api"
	fake "github.com/pluralsh/kuberentes-agent/internal/fake/api"
	"github.com/pluralsh/kuberentes-agent/internal/gitlab"
	"github.com/pluralsh/kuberentes-agent/internal/module/modserver"
	"github.com/pluralsh/kuberentes-agent/internal/module/modshared"
	"github.com/pluralsh/kuberentes-agent/internal/tool/cache"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAgentRpcApi struct {
	modserver.RpcApi
	Token           api.AgentToken
	GitLabClient    gitlab.ClientInterface
	AgentInfoCache  *cache.CacheWithErr[api.AgentToken, *api.AgentInfo]
	agentIdAttrOnce sync.Once
}

func (a *serverAgentRpcApi) AgentToken() api.AgentToken {
	return a.Token
}

func (a *serverAgentRpcApi) AgentInfo(ctx context.Context, log *zap.Logger) (*api.AgentInfo, error) {
	agentInfo, err := a.getAgentInfoCached(ctx)
	switch {
	case err == nil:
		a.agentIdAttrOnce.Do(func() {
			trace.SpanFromContext(ctx).SetAttributes(api.TraceAgentIdAttr.Int64(agentInfo.Id))
		})
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
	return a.AgentInfoCache.GetItem(ctx, a.Token, func() (*api.AgentInfo, error) {
		return fake.GetAgentInfo(ctx, a.Token, gitlab.WithoutRetries())
	})
}

type serverAgentRpcApiFactory struct {
	rpcApiFactory  modserver.RpcApiFactory
	gitLabClient   gitlab.ClientInterface
	agentInfoCache *cache.CacheWithErr[api.AgentToken, *api.AgentInfo]
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
