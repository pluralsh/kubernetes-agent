package kasapp

import (
	"context"
	"net"

	"github.com/ash2k/stager"
	"github.com/pluralsh/kuberentes-agent/internal/module/modserver"
	"github.com/pluralsh/kuberentes-agent/internal/module/observability"
	"github.com/pluralsh/kuberentes-agent/internal/tool/grpctool"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	otelmetric "go.opentelemetry.io/otel/metric"
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

func newInternalServer(tp trace.TracerProvider, mp otelmetric.MeterProvider, p propagation.TextMapPropagator,
	factory modserver.RpcApiFactory, probeRegistry *observability.ProbeRegistry,
	grpcServerErrorReporter grpctool.ServerErrorReporter) (*internalServer, error) {

	// In-memory gRPC client->listener pipe
	listener := grpctool.NewDialListener()

	// Construct connection to internal gRPC server
	conn, err := grpc.DialContext(context.Background(), "passthrough:pipe", // nolint: contextcheck
		grpc.WithStatsHandler(otelgrpc.NewServerHandler(
			otelgrpc.WithTracerProvider(tp),
			otelgrpc.WithMeterProvider(mp),
			otelgrpc.WithPropagators(p),
			otelgrpc.WithMessageEvents(otelgrpc.ReceivedEvents, otelgrpc.SentEvents),
		)),
		grpc.WithSharedWriteBuffer(true),
		grpc.WithContextDialer(listener.DialContext),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainStreamInterceptor(
			grpctool.StreamClientValidatingInterceptor,
		),
		grpc.WithChainUnaryInterceptor(
			grpctool.UnaryClientValidatingInterceptor,
		),
	)
	if err != nil {
		return nil, err
	}
	return &internalServer{
		server: grpc.NewServer(
			grpc.StatsHandler(otelgrpc.NewServerHandler(
				otelgrpc.WithTracerProvider(tp),
				otelgrpc.WithMeterProvider(mp),
				otelgrpc.WithPropagators(p),
				otelgrpc.WithMessageEvents(otelgrpc.ReceivedEvents, otelgrpc.SentEvents),
			)),
			grpc.StatsHandler(grpctool.ServerNoopMaxConnAgeStatsHandler{}),
			grpc.SharedWriteBuffer(true),
			grpc.ChainStreamInterceptor(
				modserver.StreamRpcApiInterceptor(factory), // 2. inject RPC API
				grpctool.StreamServerErrorReporterInterceptor(grpcServerErrorReporter),
			),
			grpc.ChainUnaryInterceptor(
				modserver.UnaryRpcApiInterceptor(factory), // 2. inject RPC API
				grpctool.UnaryServerErrorReporterInterceptor(grpcServerErrorReporter),
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
