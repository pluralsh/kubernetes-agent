package kasapp

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"github.com/ash2k/stager"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/redis/go-redis/v9"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/observability"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/redistool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/tlstool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/wstunnel"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/pkg/kascfg"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

const (
	defaultMaxMessageSize = 10 * 1024 * 1024
)

type agentServer struct {
	log       *zap.Logger
	listenCfg *kascfg.ListenAgentCF
	tlsConfig *tls.Config
	server    *grpc.Server
	auxCancel context.CancelFunc
	ready     func()
}

func newAgentServer(log *zap.Logger, cfg *kascfg.ConfigurationFile, tp trace.TracerProvider,
	redisClient redis.UniversalClient, factory modserver.AgentRpcApiFactory, probeRegistry *observability.ProbeRegistry) (*agentServer, error) {
	listenCfg := cfg.Agent.Listen
	tlsConfig, err := tlstool.MaybeDefaultServerTLSConfig(listenCfg.CertificateFile, listenCfg.KeyFile)
	if err != nil {
		return nil, err
	}
	agentConnectionLimiter := redistool.NewTokenLimiter(
		redisClient,
		cfg.Redis.KeyPrefix+":agent_limit",
		uint64(listenCfg.ConnectionsPerTokenPerMinute),
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
		grpc.StatsHandler(sh),
		grpc.ChainStreamInterceptor(
			grpc_prometheus.StreamServerInterceptor, // 1. measure all invocations
			otelgrpc.StreamServerInterceptor(otelgrpc.WithTracerProvider(tp), otelgrpc.WithPropagators(traceContextProp)), // 2. trace
			modserver.StreamAgentRpcApiInterceptor(factory),                                                               // 3. inject RPC API
			grpc_validator.StreamServerInterceptor(),                                                                      // x. wrap with validator
			grpctool.StreamServerLimitingInterceptor(agentConnectionLimiter),
		),
		grpc.ChainUnaryInterceptor(
			grpc_prometheus.UnaryServerInterceptor, // 1. measure all invocations
			otelgrpc.UnaryServerInterceptor(otelgrpc.WithTracerProvider(tp), otelgrpc.WithPropagators(traceContextProp)), // 2. trace
			modserver.UnaryAgentRpcApiInterceptor(factory),                                                               // 3. inject RPC API
			grpc_validator.UnaryServerInterceptor(),                                                                      // x. wrap with validator
			grpctool.UnaryServerLimitingInterceptor(agentConnectionLimiter),
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
		log:       log,
		listenCfg: listenCfg,
		tlsConfig: tlsConfig,
		server:    grpc.NewServer(serverOpts...),
		auxCancel: auxCancel,
		ready:     probeRegistry.RegisterReadinessToggle("agentServer"),
	}, nil
}

func (s *agentServer) Start(stage stager.Stage) {
	grpctool.StartServer(stage, s.server, func() (net.Listener, error) {
		var lis net.Listener
		var err error
		if s.listenCfg.Websocket { // Explicitly handle TLS for a WebSocket server
			if s.tlsConfig != nil {
				s.tlsConfig.NextProtos = []string{http2.NextProtoTLS, "http/1.1"} // h2 for gRPC, http/1.1 for WebSocket
				lis, err = tls.Listen(*s.listenCfg.Network, s.listenCfg.Address, s.tlsConfig)
			} else {
				lis, err = net.Listen(*s.listenCfg.Network, s.listenCfg.Address)
			}
			if err != nil {
				return nil, err
			}
			wsWrapper := wstunnel.ListenerWrapper{
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
	})
}
