package kasapp

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/ash2k/stager"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/observability"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/ioz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/tlstool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/pkg/kascfg"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type privateApiServer struct {
	log           *zap.Logger
	listenCfg     *kascfg.ListenPrivateApiCF
	server        *grpc.Server
	inMemListener net.Listener
	kasPool       grpctool.PoolInterface
	auxCancel     context.CancelFunc
	ready         func()
}

func newPrivateApiServer(log *zap.Logger, cfg *kascfg.ConfigurationFile, tp trace.TracerProvider,
	p propagation.TextMapPropagator, factory modserver.RpcApiFactory,
	ownPrivateApiUrl, ownPrivateApiHost string, probeRegistry *observability.ProbeRegistry) (*privateApiServer, error) {
	listenCfg := cfg.PrivateApi.Listen
	jwtSecret, err := ioz.LoadBase64Secret(listenCfg.AuthenticationSecretFile)
	if err != nil {
		return nil, fmt.Errorf("auth secret file: %w", err)
	}

	// In-memory gRPC client->listener pipe
	listener := grpctool.NewDialListener()

	// Client pool
	kasPool, err := newKasPool(log, tp, p, jwtSecret, ownPrivateApiUrl, ownPrivateApiHost, listenCfg.CaCertificateFile, listener.DialContext)
	if err != nil {
		return nil, fmt.Errorf("kas pool: %w", err)
	}

	// Server
	auxCtx, auxCancel := context.WithCancel(context.Background()) // nolint: govet
	server, err := newPrivateApiServerImpl(auxCtx, cfg, tp, p, jwtSecret, factory, ownPrivateApiHost)
	if err != nil {
		return nil, fmt.Errorf("new server: %w", err) // nolint: govet
	}
	return &privateApiServer{
		log:           log,
		listenCfg:     listenCfg,
		server:        server,
		inMemListener: listener,
		kasPool:       kasPool,
		auxCancel:     auxCancel,
		ready:         probeRegistry.RegisterReadinessToggle("privateApiServer"),
	}, nil
}

func (s *privateApiServer) Start(stage stager.Stage) {
	stage.Go(func(ctx context.Context) error {
		// Server is terminated by grpctool.StartServer() below.
		return s.server.Serve(s.inMemListener)
	})
	grpctool.StartServer(stage, s.server, func() (net.Listener, error) {
		lis, err := net.Listen(*s.listenCfg.Network, s.listenCfg.Address)
		if err != nil {
			return nil, err
		}
		addr := lis.Addr()
		s.log.Info("Private API endpoint is up",
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

func newPrivateApiServerImpl(auxCtx context.Context, cfg *kascfg.ConfigurationFile, tp trace.TracerProvider,
	p propagation.TextMapPropagator, jwtSecret []byte, factory modserver.RpcApiFactory, ownPrivateApiHost string) (*grpc.Server, error) {
	listenCfg := cfg.PrivateApi.Listen
	credsOpt, err := maybeTLSCreds(listenCfg.CertificateFile, listenCfg.KeyFile)
	if err != nil {
		return nil, err
	}
	if ownPrivateApiHost == "" && len(credsOpt) > 0 {
		return nil, fmt.Errorf("%s environment variable is not set. Set it to the kas' host name if you want to use TLS for kas->kas communication", envVarOwnPrivateApiHost)
	}

	jwtAuther := grpctool.NewJWTAuther(jwtSecret, kasName, kasName, func(ctx context.Context) *zap.Logger {
		return modserver.RpcApiFromContext(ctx).Log()
	})

	keepaliveOpt, sh := grpctool.MaxConnectionAge2GrpcKeepalive(auxCtx, listenCfg.MaxConnectionAge.AsDuration())
	return grpc.NewServer(append(
		credsOpt,
		keepaliveOpt,
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
		grpc.ForceServerCodec(grpctool.RawCodecWithProtoFallback{}),
	)...), nil
}

func newKasPool(log *zap.Logger, tp trace.TracerProvider, p propagation.TextMapPropagator, jwtSecret []byte,
	ownPrivateApiUrl, ownPrivateApiHost, caCertificateFile string, dialer func(context.Context, string) (net.Conn, error)) (grpctool.PoolInterface, error) {

	sharedPoolOpts := []grpc.DialOption{
		grpc.WithUserAgent(kasServerName()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                55 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.WithPerRPCCredentials(&grpctool.JwtCredentials{
			Secret:   jwtSecret,
			Audience: kasName,
			Issuer:   kasName,
			Insecure: true, // We may or may not have TLS setup, so always say creds don't need TLS.
		}),
		grpc.WithChainStreamInterceptor(
			grpc_prometheus.StreamClientInterceptor,
			otelgrpc.StreamClientInterceptor(otelgrpc.WithTracerProvider(tp), otelgrpc.WithPropagators(p)),
			grpctool.StreamClientValidatingInterceptor,
		),
		grpc.WithChainUnaryInterceptor(
			grpc_prometheus.UnaryClientInterceptor,
			otelgrpc.UnaryClientInterceptor(otelgrpc.WithTracerProvider(tp), otelgrpc.WithPropagators(p)),
			grpctool.UnaryClientValidatingInterceptor,
		),
	}

	// Construct in-memory connection to private API gRPC server
	inMemConn, err := grpc.DialContext(context.Background(), "pipe", // nolint: contextcheck
		append([]grpc.DialOption{
			grpc.WithContextDialer(dialer),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}, sharedPoolOpts...)...,
	)
	if err != nil {
		return nil, err
	}
	tlsCreds, err := tlstool.DefaultClientTLSConfigWithCACert(caCertificateFile)
	if err != nil {
		return nil, err
	}
	tlsCreds.ServerName = ownPrivateApiHost
	kasPool := grpctool.NewPool(log, credentials.NewTLS(tlsCreds), sharedPoolOpts...)
	return grpctool.NewPoolSelf(kasPool, ownPrivateApiUrl, inMemConn), nil
}
