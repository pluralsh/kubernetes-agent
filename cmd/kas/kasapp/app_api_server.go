package kasapp

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/ash2k/stager"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/observability"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/ioz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/kascfg"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
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

func newApiServer(log *zap.Logger, cfg *kascfg.ConfigurationFile, tp trace.TracerProvider,
	p propagation.TextMapPropagator, ssh stats.Handler, factory modserver.RpcApiFactory,
	probeRegistry *observability.ProbeRegistry) (*apiServer, error) {
	listenCfg := cfg.Api.Listen
	jwtSecret, err := ioz.LoadBase64Secret(listenCfg.AuthenticationSecretFile)
	if err != nil {
		return nil, fmt.Errorf("auth secret file: %w", err)
	}
	credsOpt, err := maybeTLSCreds(listenCfg.CertificateFile, listenCfg.KeyFile)
	if err != nil {
		return nil, err
	}

	jwtAuther := grpctool.NewJWTAuther(jwtSecret, "", kasName, func(ctx context.Context) *zap.Logger {
		return modserver.RpcApiFromContext(ctx).Log()
	})

	auxCtx, auxCancel := context.WithCancel(context.Background())
	keepaliveOpt, sh := grpctool.MaxConnectionAge2GrpcKeepalive(auxCtx, listenCfg.MaxConnectionAge.AsDuration())
	serverOpts := []grpc.ServerOption{
		grpc.StatsHandler(ssh),
		grpc.StatsHandler(sh),
		grpc.ChainStreamInterceptor(
			grpc_prometheus.StreamServerInterceptor,                                                        // 1. measure all invocations
			otelgrpc.StreamServerInterceptor(otelgrpc.WithTracerProvider(tp), otelgrpc.WithPropagators(p)), // 2. trace
			modserver.StreamRpcApiInterceptor(factory),                                                     // 3. inject RPC API
			jwtAuther.StreamServerInterceptor,                                                              // 4. auth and maybe log
			grpc_validator.StreamServerInterceptor(),                                                       // x. wrap with validator
		),
		grpc.ChainUnaryInterceptor(
			grpc_prometheus.UnaryServerInterceptor,                                                        // 1. measure all invocations
			otelgrpc.UnaryServerInterceptor(otelgrpc.WithTracerProvider(tp), otelgrpc.WithPropagators(p)), // 2. trace
			modserver.UnaryRpcApiInterceptor(factory),                                                     // 3. inject RPC API
			jwtAuther.UnaryServerInterceptor,                                                              // 4. auth and maybe log
			grpc_validator.UnaryServerInterceptor(),                                                       // x. wrap with validator
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
		lis, err := net.Listen(*s.listenCfg.Network, s.listenCfg.Address)
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
		time.Sleep(s.listenCfg.ListenGracePeriod.AsDuration())
		s.auxCancel()
	})
}
