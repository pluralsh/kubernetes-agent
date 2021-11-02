package kasapp

import (
	"context"
	"errors"
	"fmt"
	"io"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/reverse_tunnel/tracker"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/mathz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/retry"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	proxyStreamDesc = grpc.StreamDesc{
		ServerStreams: true,
		ClientStreams: true,
	}
)

// RouteToCorrectKasHandler is a gRPC handler that routes the request to another kas instance.
// Must return a gRPC status-compatible error.
func (r *router) RouteToCorrectKasHandler(srv interface{}, stream grpc.ServerStream) error {
	ctx := stream.Context()
	md, _ := metadata.FromIncomingContext(ctx)
	agentId, err := agentIdFromMeta(md)
	if err != nil {
		return err
	}
	err = retry.PollImmediateUntil(ctx, r.routeAttemptInterval, r.attemptToRoute(agentId, stream))
	if errors.Is(err, retry.ErrWaitTimeout) {
		return status.FromContextError(ctx.Err()).Err()
	}
	return err // nil or some gRPC status value
}

// attemptToRoute
// must return a gRPC status-compatible error or retry.ErrWaitTimeout.
func (r *router) attemptToRoute(agentId int64, stream grpc.ServerStream) retry.ConditionFunc {
	ctx := stream.Context()
	rpcApi := modserver.RpcApiFromContext(ctx)
	log := rpcApi.Log().With(logz.AgentId(agentId))
	pollConfig := r.pollConfig()
	sts := grpc.ServerTransportStreamFromContext(ctx)
	service, method := grpctool.SplitGrpcMethod(sts.Method())
	return func() (bool /* done */, error) {
		var tunnels []*tracker.TunnelInfo
		err := rpcApi.PollWithBackoff(pollConfig, r.attemptToGetTunnels(ctx, log, rpcApi, agentId, service, method, &tunnels))
		if err != nil {
			return false, err
		}
		mathz.Shuffle(len(tunnels), func(i, j int) {
			tunnels[i], tunnels[j] = tunnels[j], tunnels[i]
		})
		for _, tunnel := range tunnels {
			// Redefines log variable to eliminate the chance of using the original one
			log := log.With(logz.ConnectionId(tunnel.ConnectionId), logz.KasUrl(tunnel.KasUrl)) // nolint:govet
			log.Debug("Trying tunnel")
			err, done := r.attemptToRouteViaTunnel(log, rpcApi, tunnel, stream)
			switch {
			case done:
				// Request was routed successfully. The remote may have returned an error, but that's still a
				// successful response as far as we are concerned. Our job is to route the request and return what
				// the remote responded with.
				return true, err
			case err == nil:
				// No error to log, but also not a success. Continue to try the next tunnel.
			case grpctool.RequestCanceled(err):
				return false, status.Error(codes.Canceled, err.Error()) // TODO cleanup message
			case grpctool.RequestTimedOut(err):
				return false, status.Error(codes.DeadlineExceeded, err.Error()) // TODO cleanup message
			default:
				// There was an error routing the request via this tunnel. Log and try another one.
				rpcApi.HandleProcessingError(log, agentId, "Failed to route request", err)
			}
		}
		return false, nil
	}
}

// attemptToGetTunnels
// must return a gRPC status-compatible error or retry.ErrWaitTimeout.
func (r *router) attemptToGetTunnels(ctx context.Context, log *zap.Logger, rpcApi modserver.RpcApi, agentId int64,
	service, method string, infosTarget *[]*tracker.TunnelInfo) retry.PollWithBackoffFunc {
	return func() (error, retry.AttemptResult) {
		var infos tunnelInfoCollector
		err := r.tunnelQuerier.GetTunnelsByAgentId(ctx, agentId, infos.Collect(service, method))
		if err != nil {
			rpcApi.HandleProcessingError(log, agentId, "GetTunnelsByAgentId()", err)
			return nil, retry.Backoff
		}
		*infosTarget = infos
		return nil, retry.Done
	}
}

// attemptToRouteViaTunnel attempts to route the stream via the tunnel.
// Unusual signature to signal that the done bool should be checked to determine what the error value means.
func (r *router) attemptToRouteViaTunnel(log *zap.Logger, rpcApi modserver.RpcApi, tunnel *tracker.TunnelInfo, stream grpc.ServerStream) (error, bool) {
	ctx := stream.Context()
	kasClient, err := r.kasPool.Dial(ctx, tunnel.KasUrl)
	if err != nil {
		return err, false
	}
	defer kasClient.Done()
	md, _ := metadata.FromIncomingContext(ctx)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // ensure outbound stream is cleaned up
	kasStream, err := kasClient.NewStream(
		metadata.NewOutgoingContext(ctx, md),
		&proxyStreamDesc,
		grpc.ServerTransportStreamFromContext(ctx).Method(),
		grpc.ForceCodec(grpctool.RawCodecWithProtoFallback{}),
	)
	if err != nil {
		return fmt.Errorf("NewStream(): %w", err), false
	}

	// The gateway kas will block until it has a matching tunnel if it does not have one already. Yes, we found the
	// "correct" kas by looking at Redis, but that tunnel may be no longer available (e.g. disconnected or
	// has been used by another request) and there might not be any other matching tunnels.
	// So in the future we may establish connections to one or more other suitable kas instances concurrently (after
	// a small delay or immediately) to ensure the stream is routed to an agent ASAP.
	// To ensure that the stream is only routed to a single agent, we need the gateway kas to tell us (the routing kas)
	// that it has a matching tunnel. That way we know to which connection to forward the stream to.
	// We then need to tell the gateway kas that we are starting to route the stream to it. If we don't and just
	// close the connection, it does not have to use the tunnel it found and can put it back into it's
	// "ready to be used" list of tunnels.

	var kasResponse GatewayKasResponse
	err = kasStream.RecvMsg(&kasResponse) // Wait for the tunnel to be found
	if err != nil {
		if errors.Is(err, io.EOF) {
			// Gateway kas closed the connection cleanly, perhaps it's been open for too long
			return nil, false
		}
		return fmt.Errorf("kas RecvMsg(): %w", err), false
	}
	tunnelReady := kasResponse.GetTunnelReady()
	if tunnelReady == nil {
		return fmt.Errorf("invalid oneof value type: %T", kasResponse.Msg), false
	}
	err = kasStream.SendMsg(&StartStreaming{})
	if err != nil {
		if errors.Is(err, io.EOF) {
			var frame grpctool.RawFrame
			err = kasStream.RecvMsg(&frame) // get the real error
		}
		return rpcApi.HandleSendError(log, "SendMsg(StartStreaming) failed", err), false
	}
	f := kasStreamForwarder{
		log:               log,
		rpcApi:            rpcApi,
		gatewayKasVisitor: r.gatewayKasVisitor,
	}
	return f.ForwardStream(kasStream, stream), true
}

type tunnelInfoCollector []*tracker.TunnelInfo

func (c *tunnelInfoCollector) Collect(service, method string) tracker.GetTunnelsByAgentIdCallback {
	return func(info *tracker.TunnelInfo) (bool /* done */, error) {
		if info.KasUrl == "" {
			// kas without a private API endpoint. Ignore it.
			// TODO this can be made mandatory if/when the env var the address is coming from is mandatory
			return false, nil
		}
		if !info.SupportsServiceAndMethod(service, method) {
			// This tunnel doesn't support required API. Ignore it.
			return false, nil
		}
		*c = append(*c, info)
		return false, nil
	}
}
