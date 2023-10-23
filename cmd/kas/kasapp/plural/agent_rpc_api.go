package plural

import (
	"context"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"go.uber.org/zap"

	"github.com/pluralsh/kuberentes-agent/internal/api"
	fake "github.com/pluralsh/kuberentes-agent/internal/fake/api"
	"github.com/pluralsh/kuberentes-agent/internal/gitlab"
	"github.com/pluralsh/kuberentes-agent/internal/module/modserver"
	"github.com/pluralsh/kuberentes-agent/internal/tool/cache"
)

type ServerAgentRpcApi struct {
	modserver.RpcApi
	Token          api.AgentToken
	AgentInfoCache *cache.CacheWithErr[api.AgentToken, *api.AgentInfo]
}

func (a *ServerAgentRpcApi) AgentToken() api.AgentToken {
	return a.Token
}

func (a *ServerAgentRpcApi) AgentInfo(ctx context.Context, log *zap.Logger) (*api.AgentInfo, error) {
	return a.getAgentInfoCached(ctx)
}

func (a *ServerAgentRpcApi) getAgentInfoCached(ctx context.Context) (*api.AgentInfo, error) {
	return a.AgentInfoCache.GetItem(ctx, a.Token, func() (*api.AgentInfo, error) {
		return fake.GetAgentInfo(ctx, a.Token, gitlab.WithoutRetries())
	})
}

type ServerAgentRpcApiFactory struct {
	RPCApiFactory  modserver.RpcApiFactory
	AgentInfoCache *cache.CacheWithErr[api.AgentToken, *api.AgentInfo]
}

func (f *ServerAgentRpcApiFactory) New(ctx context.Context, fullMethodName string) (modserver.AgentRpcApi, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}
	return &ServerAgentRpcApi{
		RpcApi:         f.RPCApiFactory(ctx, fullMethodName),
		Token:          api.AgentToken(token),
		AgentInfoCache: f.AgentInfoCache,
	}, nil
}
