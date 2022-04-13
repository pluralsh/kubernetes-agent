package kasapp

import (
	"context"
	"io"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/reverse_tunnel"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/reverse_tunnel/tracker"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type kasRouter interface {
	RegisterAgentApi(desc *grpc.ServiceDesc)
}

type KasPool interface {
	Dial(ctx context.Context, targetUrl string) (grpctool.PoolConn, error)
	io.Closer
}

// router routes traffic from kas to another kas to agentk.
// routing kas -> gateway kas -> agentk
type router struct {
	kasPool       KasPool
	tunnelQuerier tracker.Querier
	tunnelFinder  reverse_tunnel.TunnelFinder
	pollConfig    retry.PollConfigFactory
	// internalServer is the internal gRPC server for use inside of kas.
	// Request handlers can obtain the per-request logger using grpctool.LoggerFromContext(requestContext).
	internalServer grpc.ServiceRegistrar
	// privateApiServer is the gRPC server that other kas instances can talk to.
	// Request handlers can obtain the per-request logger using grpctool.LoggerFromContext(requestContext).
	privateApiServer          grpc.ServiceRegistrar
	gatewayKasVisitor         *grpctool.StreamVisitor
	kasRoutingDurationSuccess prometheus.Observer
	kasRoutingDurationError   prometheus.Observer
}

func (r *router) RegisterAgentApi(desc *grpc.ServiceDesc) {
	// 1. Munge the descriptor into the right shape:
	//    - turn all unary calls into streaming calls
	//    - all streaming calls, including the ones from above, are handled by routing handlers
	internalServerDesc := mungeDescriptor(desc, r.RouteToCorrectKasHandler)
	privateApiServerDesc := mungeDescriptor(desc, r.RouteToCorrectAgentHandler)

	// 2. Register on InternalServer gRPC server so that ReverseTunnelClient can be used in kas to send data to
	//    this API within this kas instance. This kas instance then routes the stream to the gateway kas instance.
	r.internalServer.RegisterService(internalServerDesc, nil)

	// 3. Register on PrivateApiServer gRPC server so that this kas instance can act as the gateway kas instance
	//    from above and then route to one of the matching connected agentk instances.
	r.privateApiServer.RegisterService(privateApiServerDesc, nil)
}

func mungeDescriptor(in *grpc.ServiceDesc, handler grpc.StreamHandler) *grpc.ServiceDesc {
	streams := make([]grpc.StreamDesc, 0, len(in.Streams)+len(in.Methods))
	for _, stream := range in.Streams {
		streams = append(streams, grpc.StreamDesc{
			StreamName:    stream.StreamName,
			Handler:       handler,
			ServerStreams: true,
			ClientStreams: true,
		})
	}
	// Turn all methods into streams
	for _, method := range in.Methods {
		streams = append(streams, grpc.StreamDesc{
			StreamName:    method.MethodName,
			Handler:       handler,
			ServerStreams: true,
			ClientStreams: true,
		})
	}
	return &grpc.ServiceDesc{
		ServiceName: in.ServiceName,
		Streams:     streams,
		Metadata:    in.Metadata,
	}
}

func agentIdFromMeta(md metadata.MD) (int64, error) {
	val := md.Get(modserver.RoutingAgentIdMetadataKey)
	if len(val) != 1 {
		return 0, status.Errorf(codes.InvalidArgument, "Expecting a single %s, got %d", modserver.RoutingAgentIdMetadataKey, len(val))
	}
	agentId, err := strconv.ParseInt(val[0], 10, 64)
	if err != nil {
		return 0, status.Errorf(codes.InvalidArgument, "Invalid %s", modserver.RoutingAgentIdMetadataKey)
	}
	return agentId, nil
}
