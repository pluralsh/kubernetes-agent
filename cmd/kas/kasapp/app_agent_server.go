package kasapp

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"github.com/ash2k/stager"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/rueidis"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/observability"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/reverse_tunnel/tunnel"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/httpz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/redistool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/tlstool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/wstunnel"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/kascfg"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
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

func newAgentServer(log *zap.Logger, cfg *kascfg.ConfigurationFile, srvApi modserver.Api, dt trace.Tracer,
	tp trace.TracerProvider, redisClient rueidis.Client, ssh stats.Handler, factory modserver.AgentRpcApiFactory,
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
	prefix := cfg.Redis.KeyPrefix + ":tunnel_tracker2"
	ttl := cfg.Agent.RedisConnInfoTtl.AsDuration()
	tunnelRegistry, err := tunnel.NewRegistry(log, srvApi, dt, cfg.Agent.RedisConnInfoRefresh.AsDuration(), func() tunnel.Tracker {
		return tunnel.NewRedisTracker(redisClient, prefix, ttl, ownPrivateApiUrl)
	})
	if err != nil {
		return nil, err
	}
	agentConnectionLimiter := redistool.NewTokenLimiter(
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
	auxCtx, auxCancel := context.WithCancel(context.Background())
	traceContextProp := propagation.TraceContext{} // only want trace id, not baggage from external clients/agents
	keepaliveOpt, sh := grpctool.MaxConnectionAge2GrpcKeepalive(auxCtx, listenCfg.MaxConnectionAge.AsDuration())
	serverOpts := []grpc.ServerOption{
		grpc.StatsHandler(ssh),
		grpc.StatsHandler(sh),
		grpc.SharedWriteBuffer(true),
		grpc.ChainStreamInterceptor(
			streamProm, // 1. measure all invocations
			otelgrpc.StreamServerInterceptor(otelgrpc.WithTracerProvider(tp), otelgrpc.WithPropagators(traceContextProp),
				otelgrpc.WithMessageEvents(otelgrpc.ReceivedEvents, otelgrpc.SentEvents)), // 2. trace
			modserver.StreamAgentRpcApiInterceptor(factory), // 3. inject RPC API
			grpc_validator.StreamServerInterceptor(),        // x. wrap with validator
			grpctool.StreamServerLimitingInterceptor(agentConnectionLimiter),
			grpctool.StreamServerErrorReporterInterceptor(grpcServerErrorReporter),
		),
		grpc.ChainUnaryInterceptor(
			unaryProm, // 1. measure all invocations
			otelgrpc.UnaryServerInterceptor(otelgrpc.WithTracerProvider(tp), otelgrpc.WithPropagators(traceContextProp),
				otelgrpc.WithMessageEvents(otelgrpc.ReceivedEvents, otelgrpc.SentEvents)), // 2. trace
			modserver.UnaryAgentRpcApiInterceptor(factory), // 3. inject RPC API
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
