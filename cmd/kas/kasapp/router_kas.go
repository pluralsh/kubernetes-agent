package kasapp

import (
	"context"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/logz"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// RouteToCorrectKasHandler is a gRPC handler that routes the request to another kas instance.
// Must return a gRPC status-compatible error.
func (r *router) RouteToCorrectKasHandler(srv interface{}, stream grpc.ServerStream) error {
	// 0. boilerplate
	ctx := stream.Context()
	md, _ := metadata.FromIncomingContext(ctx)
	agentId, err := agentIdFromMeta(md)
	if err != nil {
		return err
	}
	rpcApi := modserver.RpcApiFromContext(ctx)

	// 1. find a ready, suitable tunnel
	kasStream, kasUrl, done, err := r.findReadyTunnel(ctx, rpcApi, md, agentId)
	if err != nil {
		return err
	}
	defer done()

	// 2. start streaming via the found tunnel
	f := kasStreamForwarder{
		log:               rpcApi.Log().With(logz.AgentId(agentId), logz.KasUrl(kasUrl)),
		rpcApi:            rpcApi,
		gatewayKasVisitor: r.gatewayKasVisitor,
	}
	return f.ForwardStream(kasStream, stream)
}

func (r *router) findReadyTunnel(ctx context.Context, rpcApi modserver.RpcApi, md metadata.MD, agentId int64) (grpc.ClientStream, string, func(), error) {
	tChan := make(chan readyTunnel)
	tf := tunnelFinder{
		log:           rpcApi.Log().With(logz.AgentId(agentId)),
		kasPool:       r.kasPool,
		tunnelQuerier: r.tunnelQuerier,
		rpcApi:        rpcApi,
		fullMethod:    grpc.ServerTransportStreamFromContext(ctx).Method(),
		agentId:       agentId,
		outgoingCtx:   metadata.NewOutgoingContext(ctx, md),
		foundTunnel:   tChan,
		connections:   make(map[string]kasConnAttempt),
	}
	go tf.poll(ctx, r.pollConfig())
	select {
	case <-ctx.Done():
		return nil, "", nil, status.FromContextError(ctx.Err()).Err()
	case rt := <-tChan:
		return rt.kasStream, rt.kasUrl, func() {
			rt.kasStreamCancel()
			rt.kasConn.Done()
		}, nil
	}
}
