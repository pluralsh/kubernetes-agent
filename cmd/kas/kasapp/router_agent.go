package kasapp

import (
	"errors"
	"io"
	"strings"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/reverse_tunnel"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/prototool"
	"go.uber.org/zap"
	statuspb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (r *router) RouteToCorrectAgentHandler(srv interface{}, stream grpc.ServerStream) error {
	ctx := stream.Context()
	md, _ := metadata.FromIncomingContext(ctx)
	agentId, err := agentIdFromMeta(md)
	if err != nil {
		return err
	}
	sts := grpc.ServerTransportStreamFromContext(ctx)
	service, method := grpctool.SplitGrpcMethod(sts.Method())
	wrappedStream := grpc_middleware.WrapServerStream(stream)
	// Overwrite incoming MD with sanitized MD
	wrappedStream.WrappedContext = metadata.NewIncomingContext(
		wrappedStream.WrappedContext,
		removeHopMeta(md),
	)
	stream = wrappedStream
	tunnel, err := r.tunnelFinder.FindTunnel(wrappedStream.WrappedContext, agentId, service, method)
	if err != nil {
		return status.FromContextError(err).Err()
	}
	defer tunnel.Done()
	rpcApi := modserver.RpcApiFromContext(ctx)
	log := rpcApi.Log().With(logz.AgentId(agentId))
	err = stream.SendMsg(&GatewayKasResponse{
		Msg: &GatewayKasResponse_TunnelReady_{
			TunnelReady: &GatewayKasResponse_TunnelReady{},
		},
	})
	if err != nil {
		return rpcApi.HandleIoError(log, "SendMsg(GatewayKasResponse_TunnelReady) failed", err)
	}
	var start StartStreaming
	err = stream.RecvMsg(&start)
	if err != nil {
		if errors.Is(err, io.EOF) {
			// Routing kas decided not to proceed
			return nil
		}
		return err
	}
	return tunnel.ForwardStream(log, rpcApi, stream, wrappingCallback{
		log:    log,
		rpcApi: rpcApi,
		stream: stream,
	})
}

func removeHopMeta(md metadata.MD) metadata.MD {
	md = md.Copy()
	for k := range md {
		if strings.HasPrefix(k, modserver.RoutingHopPrefix) {
			delete(md, k)
		}
	}
	return md
}

var (
	_ reverse_tunnel.TunnelDataCallback = wrappingCallback{}
)

type wrappingCallback struct {
	log    *zap.Logger
	rpcApi modserver.RpcApi
	stream grpc.ServerStream
}

func (c wrappingCallback) Header(md map[string]*prototool.Values) error {
	return c.sendMsg("SendMsg(GatewayKasResponse_Header) failed", &GatewayKasResponse{
		Msg: &GatewayKasResponse_Header_{
			Header: &GatewayKasResponse_Header{
				Meta: md,
			},
		},
	})
}

func (c wrappingCallback) Message(data []byte) error {
	return c.sendMsg("SendMsg(GatewayKasResponse_Message) failed", &GatewayKasResponse{
		Msg: &GatewayKasResponse_Message_{
			Message: &GatewayKasResponse_Message{
				Data: data,
			},
		},
	})
}

func (c wrappingCallback) Trailer(md map[string]*prototool.Values) error {
	return c.sendMsg("SendMsg(GatewayKasResponse_Trailer) failed", &GatewayKasResponse{
		Msg: &GatewayKasResponse_Trailer_{
			Trailer: &GatewayKasResponse_Trailer{
				Meta: md,
			},
		},
	})
}

func (c wrappingCallback) Error(stat *statuspb.Status) error {
	return c.sendMsg("SendMsg(GatewayKasResponse_Error) failed", &GatewayKasResponse{
		Msg: &GatewayKasResponse_Error_{
			Error: &GatewayKasResponse_Error{
				Status: stat,
			},
		},
	})
}

func (c wrappingCallback) sendMsg(errMsg string, msg *GatewayKasResponse) error {
	err := c.stream.SendMsg(msg)
	if err != nil {
		return c.rpcApi.HandleIoError(c.log, errMsg, err)
	}
	return nil
}
