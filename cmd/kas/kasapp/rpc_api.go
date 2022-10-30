package kasapp

import (
	"context"
	"sync"

	"github.com/getsentry/sentry-go"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/logz"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type serverRpcApi struct {
	modshared.RpcApiStub
	sentryHubRoot *sentry.Hub

	service string
	method  string
	traceID trace.TraceID

	sentryHubOnce sync.Once
	sentryHub     SentryHub
}

func (a *serverRpcApi) HandleProcessingError(log *zap.Logger, agentId int64, msg string, err error) {
	handleProcessingError(a.StreamCtx, a.hub(), log, agentId, msg, err)
}

func (a *serverRpcApi) HandleIoError(log *zap.Logger, msg string, err error) error {
	// The problem is almost certainly with the client's connection.
	// Still log it on Debug.
	log.Debug(msg, logz.Error(err))
	return grpctool.HandleIoError(msg, err)
}

func (a *serverRpcApi) hub() SentryHub {
	a.sentryHubOnce.Do(a.hubOnce)
	return a.sentryHub
}

func (a *serverRpcApi) hubOnce() {
	hub := a.sentryHubRoot.Clone()
	scope := hub.Scope()
	scope.SetTag(modserver.GrpcServiceSentryField, a.service)
	scope.SetTag(modserver.GrpcMethodSentryField, a.method)
	transaction := a.service + "::" + a.method                           // Like in Gitaly
	scope.SetTransaction(transaction)                                    // Like in Gitaly
	scope.SetFingerprint([]string{"{{ default }}", "grpc", transaction}) // use Sentry's default error hash but also split by gRPC transaction
	if a.traceID.IsValid() {
		scope.SetTag(modserver.TraceIdSentryField, a.traceID.String())
	}
	a.sentryHub = hub
}

type serverRpcApiFactory struct {
	log       *zap.Logger
	sentryHub *sentry.Hub
}

func (f *serverRpcApiFactory) New(ctx context.Context, fullMethodName string) modserver.RpcApi {
	service, method := grpctool.SplitGrpcMethod(fullMethodName)
	traceID := trace.SpanContextFromContext(ctx).TraceID()
	return &serverRpcApi{
		RpcApiStub: modshared.RpcApiStub{
			Logger: f.log.With(
				logz.TraceId(traceID),
				logz.GrpcService(service),
				logz.GrpcMethod(method),
			),
			StreamCtx: ctx,
		},
		sentryHubRoot: f.sentryHub,
		service:       service,
		method:        method,
		traceID:       traceID,
	}
}
