package kasapp

import (
	"context"
	"time"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/logz"
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
	ctx := stream.Context()
	md, _ := metadata.FromIncomingContext(ctx)
	agentId, err := agentIdFromMeta(md)
	if err != nil {
		return err
	}
	rpcApi := modserver.RpcApiFromContext(ctx)

	// 1. find a ready, suitable tunnel
	rt, err := r.findReadyTunnel(ctx, rpcApi, md, agentId)
	if err != nil {
		return err
	}
	defer rt.Done()

	// 2. start streaming via the found tunnel
	f := kasStreamForwarder{
		log:               rpcApi.Log().With(logz.AgentId(agentId), logz.KasUrl(rt.kasUrl)),
		rpcApi:            rpcApi,
		gatewayKasVisitor: r.gatewayKasVisitor,
	}
	return f.ForwardStream(rt.kasStream, stream)
}

func (r *router) findReadyTunnel(ctx context.Context, rpcApi modserver.RpcApi, md metadata.MD, agentId int64) (readyTunnel, error) {
	startRouting := time.Now()
	tr := trace.SpanFromContext(ctx).TracerProvider().Tracer("tunnel-router")
	findCtx, span := tr.Start(ctx, "find tunnel",
		trace.WithSpanKind(trace.SpanKindInternal),
	)
	defer span.End()
	tf := newTunnelFinder(
		rpcApi.Log().With(logz.AgentId(agentId)),
		r.kasPool,
		r.tunnelQuerier,
		rpcApi,
		grpc.ServerTransportStreamFromContext(ctx).Method(),
		r.ownPrivateApiUrl,
		agentId,
		metadata.NewOutgoingContext(ctx, md),
		r.pollConfig,
	)
	findCtx, findCancel := context.WithTimeout(findCtx, r.tunnelFindTimeout)
	defer findCancel()

	rt, err := tf.Find(findCtx)
	if err != nil {
		switch { // Order is important here.
		case ctx.Err() != nil: // Incoming stream cancelled.
			r.kasRoutingDurationAborted.Observe(time.Since(startRouting).Seconds())
			span.SetStatus(otelcodes.Error, "Aborted")
			return readyTunnel{}, grpctool.StatusErrorFromContext(ctx, "RouteToKasStreamHandler request aborted")
		case findCtx.Err() != nil: // Find tunnel timed out.
			r.kasRoutingDurationTimeout.Inc()
			span.SetStatus(otelcodes.Error, "Timed out")
			return readyTunnel{}, status.Error(codes.DeadlineExceeded, "Agent connection not found. Is agent up to date and connected?")
		default: // This should never happen, but let's handle a non-ctx error for completeness and future-proofing.
			span.SetStatus(otelcodes.Error, "Failed")
			span.RecordError(err)
			return readyTunnel{}, status.Errorf(codes.Unavailable, "Find tunnel failed: %v", err)
		}
	}
	r.kasRoutingDurationSuccess.Observe(time.Since(startRouting).Seconds())
	span.SetStatus(otelcodes.Ok, "")
	return rt, nil
}
