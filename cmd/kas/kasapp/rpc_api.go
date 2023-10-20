package kasapp

import (
	"context"
	"sync"

	"github.com/getsentry/sentry-go"
	"github.com/pluralsh/kuberentes-agent/internal/module/modserver"
	"github.com/pluralsh/kuberentes-agent/internal/module/modshared"
	"github.com/pluralsh/kuberentes-agent/internal/tool/grpctool"
	"github.com/pluralsh/kuberentes-agent/internal/tool/logz"
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
	transaction   string
}

func (a *serverRpcApi) HandleProcessingError(log *zap.Logger, agentId int64, msg string, err error) {
	handleProcessingError(a.StreamCtx, a.hub, log, agentId, msg, err)
}

func (a *serverRpcApi) HandleIoError(log *zap.Logger, msg string, err error) error {
	// The problem is almost certainly with the client's connection.
	// Still log it on Debug.
	log.Debug(msg, logz.Error(err))
	return grpctool.HandleIoError(msg, err)
}

func (a *serverRpcApi) hub() (SentryHub, string) {
	a.sentryHubOnce.Do(a.hubOnce)
	return a.sentryHub, a.transaction
}

func (a *serverRpcApi) hubOnce() {
	hub := a.sentryHubRoot.Clone()
	scope := hub.Scope()
	scope.SetTag(modserver.GrpcServiceSentryField, a.service)
	scope.SetTag(modserver.GrpcMethodSentryField, a.method)
	a.transaction = a.service + "::" + a.method                            // Like in Gitaly
	scope.SetFingerprint([]string{"{{ default }}", "grpc", a.transaction}) // use Sentry's default error hash but also split by gRPC transaction
	if a.traceID.IsValid() {
		scope.SetTag(modserver.SentryFieldTraceId, a.traceID.String())
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
