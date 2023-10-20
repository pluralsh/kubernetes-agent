package agentkapp

import (
	"context"
	"net"

	"github.com/ash2k/stager"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"
	"github.com/pluralsh/kuberentes-agent/internal/module/modagent"
	"github.com/pluralsh/kuberentes-agent/internal/module/modshared"
	"github.com/pluralsh/kuberentes-agent/internal/tool/grpctool"
	"github.com/pluralsh/kuberentes-agent/internal/tool/logz"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type internalServer struct {
	server   *grpc.Server
	conn     *grpc.ClientConn
	listener net.Listener
}

func newInternalServer(log *zap.Logger, tp trace.TracerProvider, mp otelmetric.MeterProvider, p propagation.TextMapPropagator,
	streamProm grpc.StreamServerInterceptor, unaryProm grpc.UnaryServerInterceptor) (*internalServer, error) {
	// Internal gRPC client->listener pipe
	listener := grpctool.NewDialListener()

	// Construct connection to internal gRPC server
	conn, err := grpc.DialContext(context.Background(), "passthrough:pipe", // nolint: contextcheck
		grpc.WithSharedWriteBuffer(true),
		grpc.WithContextDialer(listener.DialContext),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.ForceCodec(grpctool.RawCodec{})),
	)
	if err != nil {
		return nil, err
	}
	factory := func(ctx context.Context, method string) modagent.RpcApi {
		return &agentRpcApi{
			RpcApiStub: modshared.RpcApiStub{
				Logger:    log.With(logz.TraceIdFromContext(ctx)),
				StreamCtx: ctx,
			},
		}
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
				streamProm, // 1. measure all invocations
				modagent.StreamRpcApiInterceptor(factory), // 2. inject RPC API
				grpc_validator.StreamServerInterceptor(),  // x. wrap with validator
			),
			grpc.ChainUnaryInterceptor(
				unaryProm,                                // 1. measure all invocations
				modagent.UnaryRpcApiInterceptor(factory), // 2. inject RPC API
				grpc_validator.UnaryServerInterceptor(),  // x. wrap with validator
			),
		),
		conn:     conn,
		listener: listener,
	}, nil
}

func (s *internalServer) Start(stage stager.Stage) {
	grpctool.StartServer(stage, s.server, func() (net.Listener, error) {
		return s.listener, nil
	}, func() {})
}
