package kasapp

import (
	"context"
	"time"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/logz"
	otelcodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// RouteToKasStreamHandler is a gRPC handler that routes the request to another kas instance.
// Must return a gRPC status-compatible error.
func (r *router) RouteToKasStreamHandler(srv interface{}, stream grpc.ServerStream) error {
	// 0. boilerplate
	startRouting := time.Now()
	ctx := stream.Context()
	md, _ := metadata.FromIncomingContext(ctx)
	agentId, err := agentIdFromMeta(md)
	if err != nil {
		return err
	}
	rpcApi := modserver.RpcApiFromContext(ctx)

	// 1. find a ready, suitable tunnel
	rt, err := r.findReadyTunnel(ctx, rpcApi, md, agentId)
	routingDuration := time.Since(startRouting).Seconds()
	if err != nil {
		r.kasRoutingDurationError.Observe(routingDuration)
		return err
	}
	defer rt.Done()
	r.kasRoutingDurationSuccess.Observe(routingDuration)

	// 2. start streaming via the found tunnel
	f := kasStreamForwarder{
		log:               rpcApi.Log().With(logz.AgentId(agentId), logz.KasUrl(rt.kasUrl)),
		rpcApi:            rpcApi,
		gatewayKasVisitor: r.gatewayKasVisitor,
	}
	return f.ForwardStream(rt.kasStream, stream)
}

func (r *router) findReadyTunnel(ctx context.Context, rpcApi modserver.RpcApi, md metadata.MD, agentId int64) (readyTunnel, error) {
	tr := trace.SpanFromContext(ctx).TracerProvider().Tracer("tunnel-router")
	findCtx, span := tr.Start(ctx, "find tunnel",
		trace.WithSpanKind(trace.SpanKindInternal),
	)
	defer span.End()
	tChan := make(chan readyTunnel)
	tf := tunnelFinder{
		log:              rpcApi.Log().With(logz.AgentId(agentId)),
		kasPool:          r.kasPool,
		tunnelQuerier:    r.tunnelQuerier,
		rpcApi:           rpcApi,
		fullMethod:       grpc.ServerTransportStreamFromContext(ctx).Method(),
		ownPrivateApiUrl: r.ownPrivateApiUrl,
		agentId:          agentId,
		outgoingCtx:      metadata.NewOutgoingContext(ctx, md),
		pollConfig:       r.pollConfig,
		foundTunnel:      tChan,
		connections:      make(map[string]kasConnAttempt),
	}
	t := time.NewTimer(r.tunnelFindTimeout)
	defer t.Stop()
	go tf.Run(findCtx)
	select {
	case <-findCtx.Done():
		span.SetStatus(otelcodes.Error, "Aborted")
		return readyTunnel{}, grpctool.StatusErrorFromContext(findCtx, "RouteToKasStreamHandler request aborted")
	case <-t.C:
		span.SetStatus(otelcodes.Error, "Timed out")
		// No need to cancel ctx explicitly here.
		// ctx will be cancelled when we return from the RPC handler and tf.Run() will stop.
		return readyTunnel{}, status.Error(codes.DeadlineExceeded, "Agent connection not found. Is agent up to date and connected?")
	case rt := <-tChan:
		span.SetStatus(otelcodes.Ok, "")
		return rt, nil
	}
}
