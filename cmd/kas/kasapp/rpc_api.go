package kasapp

import (
	"context"
	"errors"

	"github.com/getsentry/sentry-go"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/gitlab"
	gapi "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/gitlab/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/cache"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/logz"
	"gitlab.com/gitlab-org/labkit/correlation"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	_ modserver.RpcApi = (*serverRpcApi)(nil)
)

type serverRpcApi struct {
	modshared.RpcApiStub
	Hub            SentryHub
	GitLabClient   gitlab.ClientInterface
	AgentInfoCache *cache.CacheWithErr
}

func (a *serverRpcApi) AgentToken() api.AgentToken {
	return api.AgentTokenFromContext(a.StreamCtx)
}

func (a *serverRpcApi) AgentInfo(ctx context.Context, log *zap.Logger) (*api.AgentInfo, error) {
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
		logAndCapture(ctx, a.Hub, log, modshared.NoAgentId, "AgentInfo()", err)
		err = status.Error(codes.Unavailable, "unavailable")
	}
	return nil, err
}

func (a *serverRpcApi) HandleProcessingError(log *zap.Logger, agentId int64, msg string, err error) {
	handleProcessingError(a.StreamCtx, a.Hub, log, agentId, msg, err)
}

func (a *serverRpcApi) HandleSendError(log *zap.Logger, msg string, err error) error {
	// The problem is almost certainly with the client's connection.
	// Still log it on Debug.
	if !grpctool.RequestCanceled(err) {
		log.Debug(msg, logz.Error(err))
	}
	return status.Error(codes.Unavailable, "gRPC send failed")
}

func (a *serverRpcApi) getAgentInfoCached(ctx context.Context) (*api.AgentInfo, error) {
	agentToken := a.AgentToken()
	agentInfo, err := a.AgentInfoCache.GetItem(ctx, agentToken, func() (interface{}, error) {
		return gapi.GetAgentInfo(ctx, a.GitLabClient, agentToken)
	})
	if err != nil {
		return nil, err
	}
	return agentInfo.(*api.AgentInfo), nil
}

type serverRpcApiFactory struct {
	log            *zap.Logger
	sentryHub      *sentry.Hub
	gitLabClient   gitlab.ClientInterface
	agentInfoCache *cache.CacheWithErr
}

func (f *serverRpcApiFactory) New(ctx context.Context, fullMethodName string) modserver.RpcApi {
	service, method := grpctool.SplitGrpcMethod(fullMethodName)
	correlationId := correlation.ExtractFromContext(ctx)
	return &serverRpcApi{
		RpcApiStub: modshared.RpcApiStub{
			Logger: f.log.With(
				logz.CorrelationId(correlationId),
				logz.GrpcService(service),
				logz.GrpcMethod(method),
			),
			StreamCtx: ctx,
		},
		Hub:            f.hub(service, method, correlationId),
		GitLabClient:   f.gitLabClient,
		AgentInfoCache: f.agentInfoCache,
	}
}

func (f *serverRpcApiFactory) hub(service, method, correlationId string) SentryHub {
	hub := f.sentryHub.Clone()
	scope := hub.Scope()
	scope.SetTag(modserver.GrpcServiceSentryField, service)
	scope.SetTag(modserver.GrpcMethodSentryField, method)
	transaction := service + "::" + method              // Like in Gitaly
	scope.SetTransaction(transaction)                   // Like in Gitaly
	scope.SetFingerprint([]string{"grpc", transaction}) // Like in Gitaly
	if correlationId != "" {
		scope.SetTag(modserver.CorrelationIdSentryField, correlationId)
	}
	return hub
}
