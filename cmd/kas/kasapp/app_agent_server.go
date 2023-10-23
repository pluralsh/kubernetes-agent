package kasapp

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"github.com/ash2k/stager"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"
	"github.com/pluralsh/kuberentes-agent/internal/module/modserver"
	"github.com/pluralsh/kuberentes-agent/internal/module/observability"
	"github.com/pluralsh/kuberentes-agent/internal/module/reverse_tunnel/tunnel"
	"github.com/pluralsh/kuberentes-agent/internal/tool/grpctool"
	"github.com/pluralsh/kuberentes-agent/internal/tool/httpz"
	"github.com/pluralsh/kuberentes-agent/internal/tool/logz"
	"github.com/pluralsh/kuberentes-agent/internal/tool/metric"
	"github.com/pluralsh/kuberentes-agent/internal/tool/redistool"
	"github.com/pluralsh/kuberentes-agent/internal/tool/tlstool"
	"github.com/pluralsh/kuberentes-agent/internal/tool/wstunnel"
	"github.com/pluralsh/kuberentes-agent/pkg/kascfg"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/rueidis"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/stats"
	"nhooyr.io/websocket"
)

const (
	defaultMaxMessageSize                 = 10 * 1024 * 1024
	agentConnectionRateExceededMetricName = "agent_server_rate_exceeded_total"
)

type agentServer struct {
	log            *zap.Logger
	listenCfg      *kascfg.ListenAgentCF
	tlsConfig      *tls.Config
	server         *grpc.Server
	tunnelRegistry *tunnel.Registry
	auxCancel      context.CancelFunc
	ready          func()
}

func newAgentServer(log *zap.Logger, cfg *kascfg.ConfigurationFile, srvApi modserver.Api, dt trace.Tracer, dm otelmetric.Meter,
	tp trace.TracerProvider, mp otelmetric.MeterProvider, redisClient rueidis.Client, ssh stats.Handler, factory modserver.AgentRpcApiFactory,
	ownPrivateApiUrl string, probeRegistry *observability.ProbeRegistry, reg *prometheus.Registry,
	streamProm grpc.StreamServerInterceptor, unaryProm grpc.UnaryServerInterceptor,
	grpcServerErrorReporter grpctool.ServerErrorReporter) (*agentServer, error) {
	listenCfg := cfg.Agent.Listen
	tlsConfig, err := tlstool.MaybeDefaultServerTLSConfig(listenCfg.CertificateFile, listenCfg.KeyFile)
	if err != nil {
		return nil, err
	}
	rateExceededCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: agentConnectionRateExceededMetricName,
		Help: "The total number of times configured rate limit of new agent connections was exceeded",
	})
	err = reg.Register(rateExceededCounter)
	if err != nil {
		return nil, err
	}
	// Tunnel registry
	tunnelRegistry, err := tunnel.NewRegistry(
		log,
		srvApi,
		dt,
		cfg.Agent.RedisConnInfoRefresh.AsDuration(),
		cfg.Agent.RedisConnInfoTtl.AsDuration(),
		tunnel.NewRedisTracker(redisClient, cfg.Redis.KeyPrefix+":tunnel_tracker2", ownPrivateApiUrl),
	)
	if err != nil {
		return nil, err
	}
	var agentConnectionLimiter grpctool.ServerLimiter
	agentConnectionLimiter = redistool.NewTokenLimiter(
		redisClient,
		cfg.Redis.KeyPrefix+":agent_limit",
		uint64(listenCfg.ConnectionsPerTokenPerMinute),
		rateExceededCounter,
		func(ctx context.Context) redistool.RpcApi {
			return &tokenLimiterApi{
				rpcApi: modserver.AgentRpcApiFromContext(ctx),
			}
		},
	)
	agentConnectionLimiter, err = metric.NewAllowLimiterInstrumentation(
		"agent_connection",
		float64(listenCfg.ConnectionsPerTokenPerMinute),
		"{connection/token/m}",
		dt,
		dm,
		agentConnectionLimiter,
	)
	if err != nil {
		return nil, err
	}
	auxCtx, auxCancel := context.WithCancel(context.Background())
	traceContextProp := propagation.TraceContext{} // only want trace id, not baggage from external clients/agents
	keepaliveOpt, sh := grpctool.MaxConnectionAge2GrpcKeepalive(auxCtx, listenCfg.MaxConnectionAge.AsDuration())
	serverOpts := []grpc.ServerOption{
		grpc.StatsHandler(otelgrpc.NewServerHandler(
			otelgrpc.WithTracerProvider(tp),
			otelgrpc.WithMeterProvider(mp),
			otelgrpc.WithPropagators(traceContextProp),
			otelgrpc.WithMessageEvents(otelgrpc.ReceivedEvents, otelgrpc.SentEvents),
		)),
		grpc.StatsHandler(ssh),
		grpc.StatsHandler(sh),
		grpc.SharedWriteBuffer(true),
		grpc.ChainStreamInterceptor(
			streamProm, // 1. measure all invocations
			modserver.StreamAgentRpcApiInterceptor(factory), // 2. inject RPC API
			grpc_validator.StreamServerInterceptor(),        // x. wrap with validator
			grpctool.StreamServerLimitingInterceptor(agentConnectionLimiter),
			grpctool.StreamServerErrorReporterInterceptor(grpcServerErrorReporter),
		),
		grpc.ChainUnaryInterceptor(
			unaryProm, // 1. measure all invocations
			modserver.UnaryAgentRpcApiInterceptor(factory), // 2. inject RPC API
			grpc_validator.UnaryServerInterceptor(),        // x. wrap with validator
			grpctool.UnaryServerLimitingInterceptor(agentConnectionLimiter),
			grpctool.UnaryServerErrorReporterInterceptor(grpcServerErrorReporter),
		),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             20 * time.Second,
			PermitWithoutStream: true,
		}),
		keepaliveOpt,
	}

	if !listenCfg.Websocket && tlsConfig != nil {
		// If we are listening for WebSocket connections, gRPC server doesn't need TLS as it's handled by the
		// HTTP/WebSocket server. Otherwise, we handle it here (if configured).
		serverOpts = append(serverOpts, grpc.Creds(credentials.NewTLS(tlsConfig)))
	}

	return &agentServer{
		log:            log,
		listenCfg:      listenCfg,
		tlsConfig:      tlsConfig,
		server:         grpc.NewServer(serverOpts...),
		tunnelRegistry: tunnelRegistry,
		auxCancel:      auxCancel,
		ready:          probeRegistry.RegisterReadinessToggle("agentServer"),
	}, nil
}

func (s *agentServer) Start(stage stager.Stage) {
	registryCtx, registryCancel := context.WithCancel(context.Background())
	stage.Go(func(ctx context.Context) error {
		return s.tunnelRegistry.Run(registryCtx) // use a separate ctx to stop when the server starts stopping
	})
	grpctool.StartServer(stage, s.server, func() (retLis net.Listener, retErr error) {
		defer func() {
			if retErr != nil { // something went wrong here, stop the registry
				registryCancel()
			}
		}()
		var lis net.Listener
		var err error
		if s.listenCfg.Websocket { // Explicitly handle TLS for a WebSocket server
			if s.tlsConfig != nil {
				s.tlsConfig.NextProtos = []string{httpz.TLSNextProtoH2, httpz.TLSNextProtoH1} // h2 for gRPC, http/1.1 for WebSocket
				lis, err = tls.Listen(*s.listenCfg.Network, s.listenCfg.Address, s.tlsConfig)
			} else {
				lis, err = net.Listen(*s.listenCfg.Network, s.listenCfg.Address)
			}
			if err != nil {
				return nil, err
			}
			wsWrapper := wstunnel.ListenerWrapper{
				AcceptOptions: websocket.AcceptOptions{
					CompressionMode: websocket.CompressionDisabled,
				},
				// TODO set timeouts
				ReadLimit:  defaultMaxMessageSize,
				ServerName: kasServerName(),
			}
			lis = wsWrapper.Wrap(lis, s.tlsConfig != nil)
		} else {
			lis, err = net.Listen(*s.listenCfg.Network, s.listenCfg.Address)
			if err != nil {
				return nil, err
			}
		}
		addr := lis.Addr()
		s.log.Info("Agentk API endpoint is up",
			logz.NetNetworkFromAddr(addr),
			logz.NetAddressFromAddr(addr),
			logz.IsWebSocket(s.listenCfg.Websocket),
		)

		s.ready()

		return lis, nil
	}, func() {
		time.Sleep(s.listenCfg.ListenGracePeriod.AsDuration())
		s.auxCancel()
		registryCancel()
	})
}
