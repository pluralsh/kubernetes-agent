package kasapp

import (
	"context"
	"net"

	"github.com/ash2k/stager"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/observability"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/grpctool"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type internalServer struct {
	server        *grpc.Server
	inMemConn     *grpc.ClientConn
	inMemListener net.Listener
	ready         func()
}

func newInternalServer(tp trace.TracerProvider, p propagation.TextMapPropagator,
	factory modserver.RpcApiFactory, probeRegistry *observability.ProbeRegistry) (*internalServer, error) {

	// In-memory gRPC client->listener pipe
	listener := grpctool.NewDialListener()

	// Construct connection to internal gRPC server
	conn, err := grpc.DialContext(context.Background(), "pipe", // nolint: contextcheck
		grpc.WithContextDialer(listener.DialContext),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainStreamInterceptor(
			otelgrpc.StreamClientInterceptor(otelgrpc.WithTracerProvider(tp), otelgrpc.WithPropagators(p)),
			grpctool.StreamClientValidatingInterceptor,
		),
		grpc.WithChainUnaryInterceptor(
			otelgrpc.UnaryClientInterceptor(otelgrpc.WithTracerProvider(tp), otelgrpc.WithPropagators(p)),
			grpctool.UnaryClientValidatingInterceptor,
		),
	)
	if err != nil {
		return nil, err
	}
	return &internalServer{
		server: grpc.NewServer(
			grpc.StatsHandler(grpctool.ServerNoopMaxConnAgeStatsHandler{}),
			grpc.ChainStreamInterceptor(
				otelgrpc.StreamServerInterceptor(otelgrpc.WithTracerProvider(tp), otelgrpc.WithPropagators(p)), // 1. trace
				modserver.StreamRpcApiInterceptor(factory),                                                     // 2. inject RPC API
			),
			grpc.ChainUnaryInterceptor(
				otelgrpc.UnaryServerInterceptor(otelgrpc.WithTracerProvider(tp), otelgrpc.WithPropagators(p)), // 1. trace
				modserver.UnaryRpcApiInterceptor(factory),                                                     // 2. inject RPC API
			),
			grpc.ForceServerCodec(grpctool.RawCodec{}),
		),
		inMemConn:     conn,
		inMemListener: listener,
		ready:         probeRegistry.RegisterReadinessToggle("internalServer"),
	}, nil
}

func (s *internalServer) Start(stage stager.Stage) {
	grpctool.StartServer(stage, s.server, func() (net.Listener, error) {
		s.ready()
		return s.inMemListener, nil
	}, func() {})
}
