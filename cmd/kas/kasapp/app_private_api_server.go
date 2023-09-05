package kasapp

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/ash2k/stager"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/observability"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/errz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/ioz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/tlstool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/kascfg"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/stats"
)

var (
	_ grpc.ServiceRegistrar = (*privateApiServer)(nil)
)

type privateApiServer struct {
	log           *zap.Logger
	listenCfg     *kascfg.ListenPrivateApiCF
	server        *grpc.Server
	inMemServer   *grpc.Server
	inMemListener net.Listener
	kasPool       grpctool.PoolInterface
	auxCancel     context.CancelFunc
	ready         func()
}

func newPrivateApiServer(log *zap.Logger, errRep errz.ErrReporter, cfg *kascfg.ConfigurationFile, tp trace.TracerProvider,
	p propagation.TextMapPropagator, csh, ssh stats.Handler, factory modserver.RpcApiFactory,
	ownPrivateApiUrl, ownPrivateApiHost string, probeRegistry *observability.ProbeRegistry,
	streamProm grpc.StreamServerInterceptor, unaryProm grpc.UnaryServerInterceptor,
	streamClientProm grpc.StreamClientInterceptor, unaryClientProm grpc.UnaryClientInterceptor,
	grpcServerErrorReporter grpctool.ServerErrorReporter) (*privateApiServer, error) {
	listenCfg := cfg.PrivateApi.Listen
	jwtSecret, err := ioz.LoadBase64Secret(listenCfg.AuthenticationSecretFile)
	if err != nil {
		return nil, fmt.Errorf("auth secret file: %w", err)
	}

	// In-memory gRPC client->listener pipe
	listener := grpctool.NewDialListener()

	// Client pool
	kasPool, err := newKasPool(log, errRep, tp, p, csh, jwtSecret, ownPrivateApiUrl, ownPrivateApiHost,
		listenCfg.CaCertificateFile, listener.DialContext, streamClientProm, unaryClientProm)
	if err != nil {
		return nil, fmt.Errorf("kas pool: %w", err)
	}

	// Server
	auxCtx, auxCancel := context.WithCancel(context.Background()) // nolint: govet
	server, inMemServer, err := newPrivateApiServerImpl(auxCtx, cfg, tp, p, ssh, jwtSecret, factory, ownPrivateApiHost, streamProm, unaryProm, grpcServerErrorReporter)
	if err != nil {
		return nil, fmt.Errorf("new server: %w", err) // nolint: govet
	}
	return &privateApiServer{
		log:           log,
		listenCfg:     listenCfg,
		server:        server,
		inMemServer:   inMemServer,
		inMemListener: listener,
		kasPool:       kasPool,
		auxCancel:     auxCancel,
		ready:         probeRegistry.RegisterReadinessToggle("privateApiServer"),
	}, nil
}

func (s *privateApiServer) Start(stage stager.Stage) {
	stopInMem := make(chan struct{})
	grpctool.StartServer(stage, s.inMemServer, func() (net.Listener, error) {
		return s.inMemListener, nil
	}, func() {
		<-stopInMem
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
		close(stopInMem)
		s.auxCancel()
	})
}

// RegisterService should be used rather than directly registering on the field servers.
func (s *privateApiServer) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
	s.server.RegisterService(desc, impl)
	s.inMemServer.RegisterService(desc, impl)
}

func newPrivateApiServerImpl(auxCtx context.Context, cfg *kascfg.ConfigurationFile, tp trace.TracerProvider,
	p propagation.TextMapPropagator, ssh stats.Handler, jwtSecret []byte, factory modserver.RpcApiFactory,
	ownPrivateApiHost string, streamProm grpc.StreamServerInterceptor, unaryProm grpc.UnaryServerInterceptor,
	grpcServerErrorReporter grpctool.ServerErrorReporter) (*grpc.Server, *grpc.Server, error) {
	listenCfg := cfg.PrivateApi.Listen
	credsOpt, err := maybeTLSCreds(listenCfg.CertificateFile, listenCfg.KeyFile)
	if err != nil {
		return nil, nil, err
	}
	if ownPrivateApiHost == "" && len(credsOpt) > 0 {
		return nil, nil, fmt.Errorf("%s environment variable is not set. Set it to the kas' host name if you want to use TLS for kas->kas communication", envVarOwnPrivateApiHost)
	}

	jwtAuther := grpctool.NewJWTAuther(jwtSecret, kasName, kasName, func(ctx context.Context) *zap.Logger {
		return modserver.RpcApiFromContext(ctx).Log()
	})

	keepaliveOpt, sh := grpctool.MaxConnectionAge2GrpcKeepalive(auxCtx, listenCfg.MaxConnectionAge.AsDuration())
	sharedOpts := []grpc.ServerOption{
		keepaliveOpt,
		grpc.StatsHandler(ssh),
		grpc.StatsHandler(sh),
		grpc.ChainStreamInterceptor(
			streamProm, // 1. measure all invocations
			otelgrpc.StreamServerInterceptor(otelgrpc.WithTracerProvider(tp), otelgrpc.WithPropagators(p),
				otelgrpc.WithMessageEvents(otelgrpc.ReceivedEvents, otelgrpc.SentEvents)), // 2. trace
			modserver.StreamRpcApiInterceptor(factory),                             // 3. inject RPC API
			jwtAuther.StreamServerInterceptor,                                      // 4. auth and maybe log
			grpc_validator.StreamServerInterceptor(),                               // x. wrap with validator
			grpctool.StreamServerErrorReporterInterceptor(grpcServerErrorReporter), // nolint:contextcheck
		),
		grpc.ChainUnaryInterceptor(
			unaryProm, // 1. measure all invocations
			otelgrpc.UnaryServerInterceptor(otelgrpc.WithTracerProvider(tp), otelgrpc.WithPropagators(p),
				otelgrpc.WithMessageEvents(otelgrpc.ReceivedEvents, otelgrpc.SentEvents)), // 2. trace
			modserver.UnaryRpcApiInterceptor(factory), // 3. inject RPC API
			jwtAuther.UnaryServerInterceptor,          // 4. auth and maybe log
			grpc_validator.UnaryServerInterceptor(),   // x. wrap with validator
			grpctool.UnaryServerErrorReporterInterceptor(grpcServerErrorReporter),
		),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             20 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.ForceServerCodec(grpctool.RawCodecWithProtoFallback{}),
	}
	server := grpc.NewServer(append(credsOpt, sharedOpts...)...)
	inMemServer := grpc.NewServer(sharedOpts...)
	return server, inMemServer, nil
}

func newKasPool(log *zap.Logger, errRep errz.ErrReporter, tp trace.TracerProvider, p propagation.TextMapPropagator,
	csh stats.Handler, jwtSecret []byte, ownPrivateApiUrl, ownPrivateApiHost, caCertificateFile string,
	dialer func(context.Context, string) (net.Conn, error),
	streamClientProm grpc.StreamClientInterceptor, unaryClientProm grpc.UnaryClientInterceptor) (grpctool.PoolInterface, error) {

	sharedPoolOpts := []grpc.DialOption{
		// Default gRPC parameters are good, no need to change them at the moment.
		// Specify them explicitly for discoverability.
		// See https://github.com/grpc/grpc/blob/master/doc/connection-backoff.md.
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff:           backoff.DefaultConfig,
			MinConnectTimeout: 20 * time.Second, // matches the default gRPC value.
		}),
		grpc.WithStatsHandler(csh),
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
			streamClientProm,
			otelgrpc.StreamClientInterceptor(otelgrpc.WithTracerProvider(tp), otelgrpc.WithPropagators(p),
				otelgrpc.WithMessageEvents(otelgrpc.ReceivedEvents, otelgrpc.SentEvents)),
			grpctool.StreamClientValidatingInterceptor,
		),
		grpc.WithChainUnaryInterceptor(
			unaryClientProm,
			otelgrpc.UnaryClientInterceptor(otelgrpc.WithTracerProvider(tp), otelgrpc.WithPropagators(p),
				otelgrpc.WithMessageEvents(otelgrpc.ReceivedEvents, otelgrpc.SentEvents)),
			grpctool.UnaryClientValidatingInterceptor,
		),
	}

	// Construct in-memory connection to private API gRPC server
	inMemConn, err := grpc.DialContext(context.Background(), "passthrough:pipe", // nolint: contextcheck
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
	kasPool := grpctool.NewPool(log, errRep, credentials.NewTLS(tlsCreds), sharedPoolOpts...)
	return grpctool.NewPoolSelf(kasPool, ownPrivateApiUrl, inMemConn), nil
}
