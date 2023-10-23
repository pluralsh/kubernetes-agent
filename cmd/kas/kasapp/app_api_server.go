package kasapp

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/ash2k/stager"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"
	"github.com/pluralsh/kuberentes-agent/internal/module/modserver"
	"github.com/pluralsh/kuberentes-agent/internal/module/observability"
	"github.com/pluralsh/kuberentes-agent/internal/tool/grpctool"
	"github.com/pluralsh/kuberentes-agent/internal/tool/ioz"
	"github.com/pluralsh/kuberentes-agent/internal/tool/logz"
	"github.com/pluralsh/kuberentes-agent/pkg/kascfg"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/stats"
)

type apiServer struct {
	log       *zap.Logger
	listenCfg *kascfg.ListenApiCF
	server    *grpc.Server
	auxCancel context.CancelFunc
	ready     func()
}

func newApiServer(log *zap.Logger, cfg *kascfg.ConfigurationFile, tp trace.TracerProvider, mp otelmetric.MeterProvider,
	p propagation.TextMapPropagator, ssh stats.Handler, factory modserver.RpcApiFactory,
	probeRegistry *observability.ProbeRegistry, streamProm grpc.StreamServerInterceptor, unaryProm grpc.UnaryServerInterceptor,
	grpcServerErrorReporter grpctool.ServerErrorReporter) (*apiServer, error) {
	listenCfg := cfg.GetApi().GetListen()
	jwtSecret, err := ioz.LoadBase64Secret(listenCfg.GetAuthenticationSecretFile())
	if err != nil {
		return nil, fmt.Errorf("auth secret file: %w", err)
	}
	credsOpt, err := maybeTLSCreds(listenCfg.GetCertificateFile(), listenCfg.GetKeyFile())
	if err != nil {
		return nil, err
	}

	jwtAuther := grpctool.NewJWTAuther(jwtSecret, "", kasName, func(ctx context.Context) *zap.Logger {
		return modserver.RpcApiFromContext(ctx).Log()
	})

	auxCtx, auxCancel := context.WithCancel(context.Background())
	keepaliveOpt, sh := grpctool.MaxConnectionAge2GrpcKeepalive(auxCtx, listenCfg.GetMaxConnectionAge().AsDuration())
	serverOpts := []grpc.ServerOption{
		grpc.StatsHandler(otelgrpc.NewServerHandler(
			otelgrpc.WithTracerProvider(tp),
			otelgrpc.WithMeterProvider(mp),
			otelgrpc.WithPropagators(p),
			otelgrpc.WithMessageEvents(otelgrpc.ReceivedEvents, otelgrpc.SentEvents),
		)),
		grpc.StatsHandler(ssh),
		grpc.StatsHandler(sh),
		grpc.SharedWriteBuffer(true),
		grpc.ChainStreamInterceptor(
			streamProm, // 1. measure all invocations
			modserver.StreamRpcApiInterceptor(factory), // 2. inject RPC API
			jwtAuther.StreamServerInterceptor,          // 3. auth and maybe log
			grpc_validator.StreamServerInterceptor(),   // x. wrap with validator
			grpctool.StreamServerErrorReporterInterceptor(grpcServerErrorReporter),
		),
		grpc.ChainUnaryInterceptor(
			unaryProm, // 1. measure all invocations
			modserver.UnaryRpcApiInterceptor(factory), // 2. inject RPC API
			jwtAuther.UnaryServerInterceptor,          // 3. auth and maybe log
			grpc_validator.UnaryServerInterceptor(),   // x. wrap with validator
			grpctool.UnaryServerErrorReporterInterceptor(grpcServerErrorReporter),
		),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             20 * time.Second,
			PermitWithoutStream: true,
		}),
		keepaliveOpt,
	}

	serverOpts = append(serverOpts, credsOpt...)

	return &apiServer{
		log:       log,
		listenCfg: listenCfg,
		server:    grpc.NewServer(serverOpts...),
		auxCancel: auxCancel,
		ready:     probeRegistry.RegisterReadinessToggle("apiServer"),
	}, nil
}

func (s *apiServer) Start(stage stager.Stage) {
	grpctool.StartServer(stage, s.server, func() (net.Listener, error) {
		lis, err := net.Listen(*s.listenCfg.GetNetwork(), s.listenCfg.GetAddress())
		if err != nil {
			return nil, err
		}
		addr := lis.Addr()
		s.log.Info("API endpoint is up",
			logz.NetNetworkFromAddr(addr),
			logz.NetAddressFromAddr(addr),
		)
		s.ready()
		return lis, nil
	}, func() {
		time.Sleep(s.listenCfg.GetListenGracePeriod().AsDuration())
		s.auxCancel()
	})
}
