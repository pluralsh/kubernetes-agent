package kasapp

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/ash2k/stager"
	"github.com/getsentry/sentry-go"
	"github.com/go-redis/redis/v8"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/cmd"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/gitaly"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/gitlab"
	gapi "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/gitlab/api"
	agent_configuration_server "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/agent_configuration/server"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/agent_tracker"
	agent_tracker_server "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/agent_tracker/server"
	configuration_project_server "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/configuration_project/server"
	gitlab_access_server "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/gitlab_access/server"
	gitops_server "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/gitops/server"
	google_profiler_server "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/google_profiler/server"
	kubernetes_api_server "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/kubernetes_api/server"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/observability"
	observability_server "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/observability/server"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/reverse_tunnel"
	reverse_tunnel_server "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/reverse_tunnel/server"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/reverse_tunnel/tracker"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/usage_metrics"
	usage_metrics_server "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/usage_metrics/server"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/cache"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/errz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/filez"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/metric"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/redistool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/retry"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/tlstool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/tracing"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/wstunnel"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/pkg/kascfg"
	"gitlab.com/gitlab-org/gitaly/v14/client"
	grpccorrelation "gitlab.com/gitlab-org/labkit/correlation/grpc"
	"gitlab.com/gitlab-org/labkit/monitoring"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	_ "google.golang.org/grpc/encoding/gzip" // Install the gzip compressor
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/stats"
)

const (
	routingAttemptInterval = 50 * time.Millisecond
	routingInitBackoff     = 100 * time.Millisecond
	routingMaxBackoff      = 1 * time.Second
	routingResetDuration   = 10 * time.Second
	routingBackoffFactor   = 2.0
	routingJitter          = 1.0

	authSecretLength      = 32
	defaultMaxMessageSize = 10 * 1024 * 1024

	kasName = "gitlab-kas"

	kasRoutingMetricName              = "k8s_api_proxy_routing_duration_seconds"
	kasRoutingStatusLabelName         = "status"
	kasRoutingStatusSuccessLabelValue = "success"
	kasRoutingStatusErrorLabelValue   = "error"
)

type ConfiguredApp struct {
	Log              *zap.Logger
	Configuration    *kascfg.ConfigurationFile
	OwnPrivateApiUrl string
}

func (a *ConfiguredApp) Run(ctx context.Context) (retErr error) {
	// TODO This should become required later?
	if a.OwnPrivateApiUrl == "" {
		a.Log.Warn(envVarOwnPrivateApiUrl + " environment variable is not set, this kas instance will not be accessible to other kas instances and certain features will not work")
	}
	// Metrics
	// TODO use an independent registry with https://gitlab.com/gitlab-org/cluster-integration/gitlab-agent/-/issues/32
	// reg := prometheus.NewPedanticRegistry()
	registerer := prometheus.DefaultRegisterer
	gatherer := prometheus.DefaultGatherer
	ssh := metric.ServerStatsHandler()
	csh := metric.ClientStatsHandler()
	//goCollector := prometheus.NewGoCollector()
	err := metric.Register(registerer, ssh, csh, gitlabBuildInfoGauge())
	if err != nil {
		return err
	}

	cfg := a.Configuration

	// Tracing
	tracer, closer, err := tracing.ConstructTracer(kasName, cfg.Observability.Tracing.ConnectionString)
	if err != nil {
		return fmt.Errorf("tracing: %w", err)
	}
	defer errz.SafeClose(closer, &retErr)

	// GitLab REST client
	gitLabClient, err := a.constructGitLabClient(tracer)
	if err != nil {
		return err
	}

	// Sentry
	sentryHub, err := a.constructSentryHub()
	if err != nil {
		return fmt.Errorf("error tracker: %w", err)
	}

	// Redis
	redisClient, err := a.constructRedisClient()
	if err != nil {
		return err
	}

	// RPC API factory
	rpcApiFactory, agentRpcApiFactory := a.constructRpcApiFactory(sentryHub, gitLabClient)

	// Server for handling agentk requests
	agentServer, err := a.constructAgentServer(ctx, tracer, redisClient, ssh, agentRpcApiFactory)
	if err != nil {
		return fmt.Errorf("agent server: %w", err)
	}

	// Server for handling external requests e.g. from GitLab
	apiServer, err := a.constructApiServer(ctx, tracer, ssh, rpcApiFactory)
	if err != nil {
		return fmt.Errorf("API server: %w", err)
	}

	// Server for handling API requests from other kas instances
	privateApiServer, err := a.constructPrivateApiServer(ctx, tracer, ssh, rpcApiFactory)
	if err != nil {
		return fmt.Errorf("private API server: %w", err)
	}

	// Internal gRPC client->listener pipe
	internalListener := grpctool.NewDialListener()

	// Construct connection to internal gRPC server
	internalServerConn, err := a.constructInternalServerConn(ctx, tracer, internalListener.DialContext)
	if err != nil {
		return err
	}
	defer errz.SafeClose(internalServerConn, &retErr)

	// Reverse gRPC tunnel tracker
	tunnelTracker := a.constructTunnelTracker(redisClient)

	// Tunnel registry
	tunnelRegistry, err := reverse_tunnel.NewTunnelRegistry(a.Log, tunnelTracker, a.OwnPrivateApiUrl)
	if err != nil {
		return err
	}

	// Construct internal gRPC server
	internalServer := a.constructInternalServer(ctx, tracer, rpcApiFactory)

	// Kas to agentk router
	kasToAgentRouter, err := a.constructKasToAgentRouter(
		tracer,
		csh,
		tunnelTracker,
		tunnelRegistry,
		internalServer,
		privateApiServer,
		registerer)
	if err != nil {
		return err
	}

	// Agent tracker
	agentTracker := a.constructAgentTracker(redisClient)

	// Usage tracker
	usageTracker := usage_metrics.NewUsageTracker()

	// Gitaly client
	gitalyClientPool := a.constructGitalyPool(csh, tracer)
	defer errz.SafeClose(gitalyClientPool, &retErr)

	// Module factories
	factories := []modserver.Factory{
		&observability_server.Factory{
			Gatherer:       gatherer,
			LivenessProbe:  observability.NoopProbe,
			ReadinessProbe: constructReadinessProbe(redisClient),
		},
		&google_profiler_server.Factory{},
		&agent_configuration_server.Factory{
			AgentRegisterer: agentTracker,
		},
		&configuration_project_server.Factory{},
		&gitops_server.Factory{},
		&usage_metrics_server.Factory{
			UsageTracker: usageTracker,
		},
		&gitlab_access_server.Factory{},
		&agent_tracker_server.Factory{
			AgentQuerier: agentTracker,
		},
		&reverse_tunnel_server.Factory{
			TunnelHandler: tunnelRegistry,
		},
		&kubernetes_api_server.Factory{},
	}

	// Construct modules
	serverApi := &serverApi{
		Hub: sentryHub,
	}
	poolWrapper := &gitaly.Pool{
		ClientPool: gitalyClientPool,
	}
	modules := make([]modserver.Module, 0, len(factories))
	for _, factory := range factories {
		// factory.New() must be called from the main goroutine because it may mutate a gRPC server (register an API)
		// and that can only be done before Serve() is called on the server.
		module, err := factory.New(&modserver.Config{
			Log:              a.Log.With(logz.ModuleName(factory.Name())),
			Api:              serverApi,
			Config:           cfg,
			GitLabClient:     gitLabClient,
			Registerer:       registerer,
			UsageTracker:     usageTracker,
			AgentServer:      agentServer,
			ApiServer:        apiServer,
			RegisterAgentApi: kasToAgentRouter.RegisterAgentApi,
			AgentConn:        internalServerConn,
			Gitaly:           poolWrapper,
			KasName:          kasName,
			Version:          cmd.Version,
			CommitId:         cmd.Commit,
		})
		if err != nil {
			return fmt.Errorf("%s: %w", factory.Name(), err)
		}
		modules = append(modules, module)
	}

	// Start things up. Stages are shut down in reverse order.
	return stager.RunStages(ctx,
		// connRegistry depends on tunnelTracker so it must be stopped last
		func(stage stager.Stage) {
			stage.Go(tunnelTracker.Run)
		},
		// Start things that modules use.
		func(stage stager.Stage) {
			stage.Go(agentTracker.Run)
			stage.Go(tunnelRegistry.Run)
		},
		// Start modules.
		func(stage stager.Stage) {
			for _, module := range modules {
				module := module // closure captures the right variable
				stage.Go(func(ctx context.Context) error {
					err := module.Run(ctx)
					if err != nil {
						return fmt.Errorf("%s: %w", module.Name(), err)
					}
					return nil
				})
			}
		},
		// Start gRPC servers.
		func(stage stager.Stage) {
			a.startAgentServer(stage, agentServer)
			a.startApiServer(stage, apiServer)
			a.startPrivateApiServer(stage, privateApiServer)
			a.startInternalServer(stage, internalServer, internalListener)
		},
	)
}

func (a *ConfiguredApp) constructRpcApiFactory(sentryHub *sentry.Hub, gitLabClient gitlab.ClientInterface) (modserver.RpcApiFactory, modserver.AgentRpcApiFactory) {
	aCfg := a.Configuration.Agent
	f := serverRpcApiFactory{
		log:       a.Log,
		sentryHub: sentryHub,
	}
	fAgent := serverAgentRpcApiFactory{
		rpcApiFactory:  f.New,
		gitLabClient:   gitLabClient,
		agentInfoCache: cache.NewWithError(aCfg.InfoCacheTtl.AsDuration(), aCfg.InfoCacheErrorTtl.AsDuration(), gapi.IsCacheableError),
	}
	return f.New, fAgent.New
}

func (a *ConfiguredApp) constructKasToAgentRouter(tracer opentracing.Tracer, csh stats.Handler, tunnelQuerier tracker.Querier,
	tunnelFinder reverse_tunnel.TunnelFinder, internalServer, privateApiServer grpc.ServiceRegistrar,
	registerer prometheus.Registerer) (kasRouter, error) {
	// TODO this should become required
	if a.Configuration.PrivateApi == nil {
		return nopKasRouter{}, nil
	}
	listenCfg := a.Configuration.PrivateApi.Listen
	jwtSecret, err := filez.LoadBase64Secret(listenCfg.AuthenticationSecretFile)
	if err != nil {
		return nil, fmt.Errorf("auth secret file: %w", err)
	}
	gatewayKasVisitor, err := grpctool.NewStreamVisitor(&GatewayKasResponse{})
	if err != nil {
		return nil, err
	}
	kasRoutingDuration := constructKasRoutingDurationHistogram()
	err = metric.Register(registerer, kasRoutingDuration)
	if err != nil {
		return nil, err
	}
	return &router{
		kasPool: grpctool.NewPool(a.Log,
			grpc.WithTransportCredentials(insecure.NewCredentials()), // TODO support TLS
			grpc.WithUserAgent(kasServerName()),
			grpc.WithStatsHandler(csh),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                55 * time.Second,
				PermitWithoutStream: true,
			}),
			grpc.WithPerRPCCredentials(&grpctool.JwtCredentials{
				Secret:   jwtSecret,
				Audience: kasName,
				Issuer:   kasName,
				Insecure: true, // TODO support TLS
			}),
			grpc.WithChainStreamInterceptor(
				grpc_prometheus.StreamClientInterceptor,
				grpccorrelation.StreamClientCorrelationInterceptor(grpccorrelation.WithClientName(kasName)),
				grpc_opentracing.StreamClientInterceptor(grpc_opentracing.WithTracer(tracer)),
				grpctool.StreamClientValidatingInterceptor,
			),
			grpc.WithChainUnaryInterceptor(
				grpc_prometheus.UnaryClientInterceptor,
				grpccorrelation.UnaryClientCorrelationInterceptor(grpccorrelation.WithClientName(kasName)),
				grpc_opentracing.UnaryClientInterceptor(grpc_opentracing.WithTracer(tracer)),
				grpctool.UnaryClientValidatingInterceptor,
			),
		),
		tunnelQuerier: tunnelQuerier,
		tunnelFinder:  tunnelFinder,
		pollConfig: retry.NewPollConfigFactory(routingAttemptInterval, retry.NewExponentialBackoffFactory(
			routingInitBackoff,
			routingMaxBackoff,
			routingResetDuration,
			routingBackoffFactor,
			routingJitter,
		)),
		internalServer:            internalServer,
		privateApiServer:          privateApiServer,
		gatewayKasVisitor:         gatewayKasVisitor,
		kasRoutingDurationSuccess: kasRoutingDuration.WithLabelValues(kasRoutingStatusSuccessLabelValue),
		kasRoutingDurationError:   kasRoutingDuration.WithLabelValues(kasRoutingStatusErrorLabelValue),
	}, nil
}

func (a *ConfiguredApp) startAgentServer(stage stager.Stage, agentServer *grpc.Server) {
	grpctool.StartServer(stage, agentServer, func() (net.Listener, error) {
		listenCfg := a.Configuration.Agent.Listen
		var lis net.Listener
		if listenCfg.Websocket { // Explicitly handle TLS for a WebSocket server
			tlsConfig, err := tlstool.MaybeDefaultServerTLSConfig(listenCfg.CertificateFile, listenCfg.KeyFile)
			if err != nil {
				return nil, err
			}
			if tlsConfig != nil {
				tlsConfig.NextProtos = []string{"http/1.1"} // like http.Server.ServeTLS(), but only http/1.1
				lis, err = tls.Listen(listenCfg.Network.String(), listenCfg.Address, tlsConfig)
			} else {
				lis, err = net.Listen(listenCfg.Network.String(), listenCfg.Address)
			}
			if err != nil {
				return nil, err
			}
			wsWrapper := wstunnel.ListenerWrapper{
				// TODO set timeouts
				ReadLimit:  defaultMaxMessageSize,
				ServerName: kasServerName(),
			}
			lis = wsWrapper.Wrap(lis)
		} else {
			var err error
			lis, err = net.Listen(listenCfg.Network.String(), listenCfg.Address)
			if err != nil {
				return nil, err
			}
		}
		addr := lis.Addr()
		a.Log.Info("Agentk API endpoint is up",
			logz.NetNetworkFromAddr(addr),
			logz.NetAddressFromAddr(addr),
			logz.IsWebSocket(listenCfg.Websocket),
		)
		return lis, nil
	})
}

func (a *ConfiguredApp) startApiServer(stage stager.Stage, apiServer *grpc.Server) {
	// TODO this should become required
	if a.Configuration.Api == nil {
		return
	}
	grpctool.StartServer(stage, apiServer, func() (net.Listener, error) {
		listenCfg := a.Configuration.Api.Listen
		lis, err := net.Listen(listenCfg.Network.String(), listenCfg.Address)
		if err != nil {
			return nil, err
		}
		addr := lis.Addr()
		a.Log.Info("API endpoint is up",
			logz.NetNetworkFromAddr(addr),
			logz.NetAddressFromAddr(addr),
		)
		return lis, nil
	})
}

func (a *ConfiguredApp) startPrivateApiServer(stage stager.Stage, apiServer *grpc.Server) {
	// TODO this should become required
	if a.Configuration.PrivateApi == nil {
		return
	}
	grpctool.StartServer(stage, apiServer, func() (net.Listener, error) {
		listenCfg := a.Configuration.PrivateApi.Listen
		lis, err := net.Listen(listenCfg.Network.String(), listenCfg.Address)
		if err != nil {
			return nil, err
		}
		addr := lis.Addr()
		a.Log.Info("Private API endpoint is up",
			logz.NetNetworkFromAddr(addr),
			logz.NetAddressFromAddr(addr),
		)
		return lis, nil
	})
}

func (a *ConfiguredApp) startInternalServer(stage stager.Stage, internalServer *grpc.Server, internalListener net.Listener) {
	grpctool.StartServer(stage, internalServer, func() (net.Listener, error) {
		return internalListener, nil
	})
}

func (a *ConfiguredApp) constructAgentServer(ctx context.Context, tracer opentracing.Tracer,
	redisClient redis.UniversalClient, ssh stats.Handler, factory modserver.AgentRpcApiFactory) (*grpc.Server, error) {
	listenCfg := a.Configuration.Agent.Listen
	// TODO construct independent metrics interceptors with https://gitlab.com/gitlab-org/cluster-integration/gitlab-agent/-/issues/32
	grpcStreamServerInterceptors := []grpc.StreamServerInterceptor{
		grpc_prometheus.StreamServerInterceptor, // 1. measure all invocations
		grpccorrelation.StreamServerCorrelationInterceptor( // 2. add correlation id
			grpccorrelation.WithoutPropagation(),
			grpccorrelation.WithReversePropagation()),
		modserver.StreamAgentRpcApiInterceptor(factory),                               // 3. inject RPC API
		grpc_validator.StreamServerInterceptor(),                                      // x. wrap with validator
		grpc_opentracing.StreamServerInterceptor(grpc_opentracing.WithTracer(tracer)), // x. trace
	}
	grpcUnaryServerInterceptors := []grpc.UnaryServerInterceptor{
		grpc_prometheus.UnaryServerInterceptor, // 1. measure all invocations
		grpccorrelation.UnaryServerCorrelationInterceptor( // 2. add correlation id
			grpccorrelation.WithoutPropagation(),
			grpccorrelation.WithReversePropagation()),
		modserver.UnaryAgentRpcApiInterceptor(factory),                               // 3. inject RPC API
		grpc_validator.UnaryServerInterceptor(),                                      // x. wrap with validator
		grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(tracer)), // x. trace
	}

	if redisClient != nil {
		agentConnectionLimiter := redistool.NewTokenLimiter(
			redisClient,
			a.Configuration.Redis.KeyPrefix+":agent_limit",
			uint64(listenCfg.ConnectionsPerTokenPerMinute),
			func(ctx context.Context) redistool.RpcApi {
				return &tokenLimiterApi{
					rpcApi: modserver.AgentRpcApiFromContext(ctx),
				}
			},
		)
		grpcStreamServerInterceptors = append(grpcStreamServerInterceptors, grpctool.StreamServerLimitingInterceptor(agentConnectionLimiter))
		grpcUnaryServerInterceptors = append(grpcUnaryServerInterceptors, grpctool.UnaryServerLimitingInterceptor(agentConnectionLimiter))
	}
	keepaliveOpt, sh := grpctool.MaxConnectionAge2GrpcKeepalive(ctx, listenCfg.MaxConnectionAge.AsDuration())
	serverOpts := []grpc.ServerOption{
		grpc.StatsHandler(grpctool.NewJoinStatHandlers(ssh, sh)),
		grpc.ChainStreamInterceptor(grpcStreamServerInterceptors...),
		grpc.ChainUnaryInterceptor(grpcUnaryServerInterceptors...),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             20 * time.Second,
			PermitWithoutStream: true,
		}),
		keepaliveOpt,
	}

	if !listenCfg.Websocket {
		// If we are listening for WebSocket connections, gRPC server doesn't need TLS.
		// TLS is handled by the HTTP/WebSocket server in this case.
		credsOpt, err := maybeTLSCreds(listenCfg.CertificateFile, listenCfg.KeyFile)
		if err != nil {
			return nil, err
		}
		serverOpts = append(serverOpts, credsOpt...)
	}

	return grpc.NewServer(serverOpts...), nil
}

func (a *ConfiguredApp) constructApiServer(ctx context.Context, tracer opentracing.Tracer, ssh stats.Handler,
	factory modserver.RpcApiFactory) (*grpc.Server, error) {
	// TODO this should become required
	if a.Configuration.Api == nil {
		return grpc.NewServer(), nil
	}
	listenCfg := a.Configuration.Api.Listen
	jwtSecret, err := filez.LoadBase64Secret(listenCfg.AuthenticationSecretFile)
	if err != nil {
		return nil, fmt.Errorf("auth secret file: %w", err)
	}

	jwtAuther := grpctool.NewJWTAuther(jwtSecret, "", kasName, func(ctx context.Context) *zap.Logger {
		return modserver.RpcApiFromContext(ctx).Log()
	})

	// TODO construct independent metrics interceptors with https://gitlab.com/gitlab-org/cluster-integration/gitlab-agent/-/issues/32
	grpcStreamServerInterceptors := []grpc.StreamServerInterceptor{
		grpc_prometheus.StreamServerInterceptor,                                       // 1. measure all invocations
		grpccorrelation.StreamServerCorrelationInterceptor(),                          // 2. add correlation id
		modserver.StreamRpcApiInterceptor(factory),                                    // 3. inject RPC API
		jwtAuther.StreamServerInterceptor,                                             // 4. auth and maybe log
		grpc_validator.StreamServerInterceptor(),                                      // x. wrap with validator
		grpc_opentracing.StreamServerInterceptor(grpc_opentracing.WithTracer(tracer)), // x. trace
	}
	grpcUnaryServerInterceptors := []grpc.UnaryServerInterceptor{
		grpc_prometheus.UnaryServerInterceptor,                                       // 1. measure all invocations
		grpccorrelation.UnaryServerCorrelationInterceptor(),                          // 2. add correlation id
		modserver.UnaryRpcApiInterceptor(factory),                                    // 3. inject RPC API
		jwtAuther.UnaryServerInterceptor,                                             // 4. auth and maybe log
		grpc_validator.UnaryServerInterceptor(),                                      // x. wrap with validator
		grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(tracer)), // x. trace
	}
	keepaliveOpt, sh := grpctool.MaxConnectionAge2GrpcKeepalive(ctx, listenCfg.MaxConnectionAge.AsDuration())
	serverOpts := []grpc.ServerOption{
		grpc.StatsHandler(grpctool.NewJoinStatHandlers(ssh, sh)),
		grpc.ChainStreamInterceptor(grpcStreamServerInterceptors...),
		grpc.ChainUnaryInterceptor(grpcUnaryServerInterceptors...),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             20 * time.Second,
			PermitWithoutStream: true,
		}),
		keepaliveOpt,
	}

	credsOpt, err := maybeTLSCreds(listenCfg.CertificateFile, listenCfg.KeyFile)
	if err != nil {
		return nil, err
	}
	serverOpts = append(serverOpts, credsOpt...)

	return grpc.NewServer(serverOpts...), nil
}

func (a *ConfiguredApp) constructPrivateApiServer(ctx context.Context, tracer opentracing.Tracer, ssh stats.Handler,
	factory modserver.RpcApiFactory) (*grpc.Server, error) {
	// TODO this should become required
	if a.Configuration.PrivateApi == nil {
		return grpc.NewServer(), nil
	}
	listenCfg := a.Configuration.PrivateApi.Listen
	jwtSecret, err := filez.LoadBase64Secret(listenCfg.AuthenticationSecretFile)
	if err != nil {
		return nil, fmt.Errorf("auth secret file: %w", err)
	}

	jwtAuther := grpctool.NewJWTAuther(jwtSecret, kasName, kasName, func(ctx context.Context) *zap.Logger {
		return modserver.RpcApiFromContext(ctx).Log()
	})

	// TODO construct independent metrics interceptors with https://gitlab.com/gitlab-org/cluster-integration/gitlab-agent/-/issues/32
	grpcStreamServerInterceptors := []grpc.StreamServerInterceptor{
		grpc_prometheus.StreamServerInterceptor,                                       // 1. measure all invocations
		grpccorrelation.StreamServerCorrelationInterceptor(),                          // 2. add correlation id
		modserver.StreamRpcApiInterceptor(factory),                                    // 3. inject RPC API
		jwtAuther.StreamServerInterceptor,                                             // 4. auth and maybe log
		grpc_validator.StreamServerInterceptor(),                                      // x. wrap with validator
		grpc_opentracing.StreamServerInterceptor(grpc_opentracing.WithTracer(tracer)), // x. trace
	}
	grpcUnaryServerInterceptors := []grpc.UnaryServerInterceptor{
		grpc_prometheus.UnaryServerInterceptor,                                       // 1. measure all invocations
		grpccorrelation.UnaryServerCorrelationInterceptor(),                          // 2. add correlation id
		modserver.UnaryRpcApiInterceptor(factory),                                    // 3. inject RPC API
		jwtAuther.UnaryServerInterceptor,                                             // 4. auth and maybe log
		grpc_validator.UnaryServerInterceptor(),                                      // x. wrap with validator
		grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(tracer)), // x. trace
	}
	keepaliveOpt, sh := grpctool.MaxConnectionAge2GrpcKeepalive(ctx, listenCfg.MaxConnectionAge.AsDuration())
	serverOpts := []grpc.ServerOption{
		grpc.StatsHandler(grpctool.NewJoinStatHandlers(ssh, sh)),
		grpc.ChainStreamInterceptor(grpcStreamServerInterceptors...),
		grpc.ChainUnaryInterceptor(grpcUnaryServerInterceptors...),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             20 * time.Second,
			PermitWithoutStream: true,
		}),
		keepaliveOpt,
		grpc.ForceServerCodec(grpctool.RawCodecWithProtoFallback{}),
	}
	credsOpt, err := maybeTLSCreds(listenCfg.CertificateFile, listenCfg.KeyFile)
	if err != nil {
		return nil, err
	}
	serverOpts = append(serverOpts, credsOpt...)

	return grpc.NewServer(serverOpts...), nil
}

func (a *ConfiguredApp) constructInternalServer(auxCtx context.Context, tracer opentracing.Tracer,
	factory modserver.RpcApiFactory) *grpc.Server {
	// TODO construct independent metrics interceptors with https://gitlab.com/gitlab-org/cluster-integration/gitlab-agent/-/issues/32
	return grpc.NewServer(
		grpc.StatsHandler(grpctool.NewServerMaxConnAgeStatsHandler(auxCtx, 0)),
		grpc.ChainStreamInterceptor(
			grpccorrelation.StreamServerCorrelationInterceptor(),                          // 1. add correlation id
			modserver.StreamRpcApiInterceptor(factory),                                    // 2. inject RPC API
			grpc_opentracing.StreamServerInterceptor(grpc_opentracing.WithTracer(tracer)), // x. trace
		),
		grpc.ChainUnaryInterceptor(
			grpccorrelation.UnaryServerCorrelationInterceptor(),                          // 1. add correlation id
			modserver.UnaryRpcApiInterceptor(factory),                                    // 2. inject RPC API
			grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(tracer)), // x. trace
		),
		grpc.ForceServerCodec(grpctool.RawCodec{}),
	)
}

func (a *ConfiguredApp) constructInternalServerConn(ctx context.Context, tracer opentracing.Tracer, dialContext func(ctx context.Context, addr string) (net.Conn, error)) (*grpc.ClientConn, error) {
	return grpc.DialContext(ctx, "pipe",
		grpc.WithContextDialer(dialContext),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainStreamInterceptor(
			grpccorrelation.StreamClientCorrelationInterceptor(grpccorrelation.WithClientName(kasName)),
			grpc_opentracing.StreamClientInterceptor(grpc_opentracing.WithTracer(tracer)),
			grpctool.StreamClientValidatingInterceptor,
		),
		grpc.WithChainUnaryInterceptor(
			grpccorrelation.UnaryClientCorrelationInterceptor(grpccorrelation.WithClientName(kasName)),
			grpc_opentracing.UnaryClientInterceptor(grpc_opentracing.WithTracer(tracer)),
			grpctool.UnaryClientValidatingInterceptor,
		),
	)
}

func (a *ConfiguredApp) constructAgentTracker(redisClient redis.UniversalClient) agent_tracker.Tracker {
	if redisClient == nil {
		return nopAgentTracker{}
	}
	cfg := a.Configuration
	return agent_tracker.NewRedisTracker(
		a.Log,
		redisClient,
		cfg.Redis.KeyPrefix+":agent_tracker",
		cfg.Agent.RedisConnInfoTtl.AsDuration(),
		cfg.Agent.RedisConnInfoRefresh.AsDuration(),
		cfg.Agent.RedisConnInfoGc.AsDuration(),
	)
}

func (a *ConfiguredApp) constructTunnelTracker(redisClient redis.UniversalClient) tracker.Tracker {
	if redisClient == nil {
		return nopTunnelTracker{}
	}
	cfg := a.Configuration
	return tracker.NewRedisTracker(
		a.Log,
		redisClient,
		cfg.Redis.KeyPrefix+":tunnel_tracker",
		cfg.Agent.RedisConnInfoTtl.AsDuration(),
		cfg.Agent.RedisConnInfoRefresh.AsDuration(),
		cfg.Agent.RedisConnInfoGc.AsDuration(),
	)
}

func (a *ConfiguredApp) constructSentryHub() (*sentry.Hub, error) {
	s := a.Configuration.Observability.Sentry
	a.Log.Debug("Initializing Sentry error tracking", logz.SentryDSN(s.Dsn), logz.SentryEnv(s.Environment))
	dialer := net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	sentryClient, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:         s.Dsn, // empty DSN disables Sentry transport
		SampleRate:  1,     // no sampling
		Release:     fmt.Sprintf("%s/%s", cmd.Version, cmd.Commit),
		Environment: s.Environment,
		HTTPTransport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			DialContext:           dialer.DialContext,
			TLSClientConfig:       tlstool.DefaultClientTLSConfig(),
			MaxIdleConns:          10,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 20 * time.Second,
			ExpectContinueTimeout: 20 * time.Second,
			ForceAttemptHTTP2:     true,
		},
	})
	if err != nil {
		return nil, err
	}
	return sentry.NewHub(sentryClient, sentry.NewScope()), nil
}

func (a *ConfiguredApp) loadGitLabClientAuthSecret() ([]byte, error) {
	decodedAuthSecret, err := filez.LoadBase64Secret(a.Configuration.Gitlab.AuthenticationSecretFile)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	if len(decodedAuthSecret) != authSecretLength {
		return nil, fmt.Errorf("decoding: expecting %d bytes, was %d", authSecretLength, len(decodedAuthSecret))
	}
	return decodedAuthSecret, nil
}

func (a *ConfiguredApp) constructGitLabClient(tracer opentracing.Tracer) (*gitlab.Client, error) {
	cfg := a.Configuration

	gitLabUrl, err := url.Parse(cfg.Gitlab.Address)
	if err != nil {
		return nil, err
	}
	// TLS cert for talking to GitLab/Workhorse.
	clientTLSConfig, err := tlstool.DefaultClientTLSConfigWithCACert(cfg.Gitlab.CaCertificateFile)
	if err != nil {
		return nil, err
	}
	// Secret for JWT signing
	decodedAuthSecret, err := a.loadGitLabClientAuthSecret()
	if err != nil {
		return nil, fmt.Errorf("authentication secret: %w", err)
	}
	return gitlab.NewClient(
		gitLabUrl,
		decodedAuthSecret,
		gitlab.WithCorrelationClientName(kasName),
		gitlab.WithUserAgent(kasServerName()),
		gitlab.WithTracer(tracer),
		gitlab.WithLogger(a.Log),
		gitlab.WithTLSConfig(clientTLSConfig),
		gitlab.WithRateLimiter(rate.NewLimiter(
			rate.Limit(cfg.Gitlab.ApiRateLimit.RefillRatePerSecond),
			int(cfg.Gitlab.ApiRateLimit.BucketSize),
		)),
	), nil
}

func (a *ConfiguredApp) constructGitalyPool(csh stats.Handler, tracer opentracing.Tracer) *client.Pool {
	g := a.Configuration.Gitaly
	globalGitalyRpcLimiter := rate.NewLimiter(
		rate.Limit(g.GlobalApiRateLimit.RefillRatePerSecond),
		int(g.GlobalApiRateLimit.BucketSize),
	)
	return client.NewPoolWithOptions(
		client.WithDialOptions(
			grpc.WithUserAgent(kasServerName()),
			grpc.WithStatsHandler(csh),
			// Don't put interceptors here as order is important. Put them below.
		),
		client.WithDialer(func(ctx context.Context, address string, dialOptions []grpc.DialOption) (*grpc.ClientConn, error) {
			perServerGitalyRpcLimiter := rate.NewLimiter(
				rate.Limit(g.PerServerApiRateLimit.RefillRatePerSecond),
				int(g.PerServerApiRateLimit.BucketSize))
			// TODO construct independent interceptors with https://gitlab.com/gitlab-org/cluster-integration/gitlab-agent/-/issues/32
			opts := []grpc.DialOption{
				grpc.WithChainStreamInterceptor(
					grpc_prometheus.StreamClientInterceptor,
					grpccorrelation.StreamClientCorrelationInterceptor(grpccorrelation.WithClientName(kasName)),
					grpc_opentracing.StreamClientInterceptor(grpc_opentracing.WithTracer(tracer)),
					grpctool.StreamClientLimitingInterceptor(globalGitalyRpcLimiter),
					grpctool.StreamClientLimitingInterceptor(perServerGitalyRpcLimiter),
				),
				grpc.WithChainUnaryInterceptor(
					grpc_prometheus.UnaryClientInterceptor,
					grpccorrelation.UnaryClientCorrelationInterceptor(grpccorrelation.WithClientName(kasName)),
					grpc_opentracing.UnaryClientInterceptor(grpc_opentracing.WithTracer(tracer)),
					grpctool.UnaryClientLimitingInterceptor(globalGitalyRpcLimiter),
					grpctool.UnaryClientLimitingInterceptor(perServerGitalyRpcLimiter),
				),
			}
			opts = append(opts, dialOptions...)
			return client.DialContext(ctx, address, opts)
		}),
	)
}

func (a *ConfiguredApp) constructRedisClient() (redis.UniversalClient, error) {
	cfg := a.Configuration.Redis
	if cfg == nil {
		return nil, nil
	}
	poolSize := int(cfg.PoolSize)
	dialTimeout := cfg.DialTimeout.AsDuration()
	readTimeout := cfg.ReadTimeout.AsDuration()
	writeTimeout := cfg.WriteTimeout.AsDuration()
	idleTimeout := cfg.IdleTimeout.AsDuration()
	var err error
	var tlsConfig *tls.Config
	if cfg.Tls != nil && cfg.Tls.Enabled {
		tlsConfig, err = tlstool.DefaultClientTLSConfigWithCACertKeyPair(cfg.Tls.CaCertificateFile, cfg.Tls.CertificateFile, cfg.Tls.KeyFile)
		if err != nil {
			return nil, err
		}
	}
	var password string
	if cfg.PasswordFile != "" {
		passwordBytes, err := os.ReadFile(cfg.PasswordFile)
		if err != nil {
			return nil, err
		}
		password = string(passwordBytes)
	}
	switch v := cfg.RedisConfig.(type) {
	case *kascfg.RedisCF_Server:
		if tlsConfig != nil {
			tlsConfig.ServerName = strings.Split(v.Server.Address, ":")[0]
		}
		return redis.NewClient(&redis.Options{
			Addr:         v.Server.Address,
			PoolSize:     poolSize,
			DialTimeout:  dialTimeout,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			IdleTimeout:  idleTimeout,
			Username:     cfg.Username,
			Password:     password,
			Network:      cfg.Network,
			TLSConfig:    tlsConfig,
		}), nil
	case *kascfg.RedisCF_Sentinel:
		var sentinelPassword string
		if v.Sentinel.SentinelPasswordFile != "" {
			sentinelPasswordBytes, err := os.ReadFile(v.Sentinel.SentinelPasswordFile)
			if err != nil {
				return nil, err
			}
			sentinelPassword = string(sentinelPasswordBytes)
		}
		return redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:       v.Sentinel.MasterName,
			SentinelAddrs:    v.Sentinel.Addresses,
			DialTimeout:      dialTimeout,
			ReadTimeout:      readTimeout,
			WriteTimeout:     writeTimeout,
			PoolSize:         poolSize,
			IdleTimeout:      idleTimeout,
			Username:         cfg.Username,
			Password:         password,
			SentinelPassword: sentinelPassword,
			TLSConfig:        tlsConfig,
		}), nil
	default:
		// This should never happen
		return nil, fmt.Errorf("unexpected Redis config type: %T", cfg.RedisConfig)
	}
}

func constructKasRoutingDurationHistogram() *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    kasRoutingMetricName,
		Help:    "The time it takes the routing kas to find a suitable tunnel in seconds",
		Buckets: prometheus.ExponentialBuckets(time.Millisecond.Seconds(), 4, 9), // 9 buckets of milliseconds as seconds (1,4,16,64,256,1k,4k,16k,64k)
	}, []string{kasRoutingStatusLabelName})
}

func constructReadinessProbe(redisClient redis.UniversalClient) observability.Probe {
	if redisClient == nil {
		return observability.NoopProbe
	}
	return func(ctx context.Context) error {
		status := redisClient.Ping(ctx)
		err := status.Err()
		if err != nil {
			return fmt.Errorf("redis: %w", err)
		}
		return nil
	}
}

func gitlabBuildInfoGauge() prometheus.Gauge {
	buildInfoGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: monitoring.GitlabBuildInfoGaugeMetricName,
		Help: "Current build info for this GitLab Service",
		ConstLabels: prometheus.Labels{
			"version": cmd.Version,
			"built":   cmd.BuildTime,
		},
	})
	buildInfoGauge.Set(1)
	return buildInfoGauge
}

func maybeTLSCreds(certFile, keyFile string) ([]grpc.ServerOption, error) {
	config, err := tlstool.MaybeDefaultServerTLSConfig(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	if config != nil {
		return []grpc.ServerOption{grpc.Creds(credentials.NewTLS(config))}, nil
	}
	return nil, nil
}

func kasServerName() string {
	return fmt.Sprintf("%s/%s/%s", kasName, cmd.Version, cmd.Commit)
}

var (
	_ agent_tracker.Tracker = nopAgentTracker{}
	_ tracker.Tracker       = nopTunnelTracker{}
	_ kasRouter             = nopKasRouter{}
	_ redistool.RpcApi      = (*tokenLimiterApi)(nil)
)

type nopAgentTracker struct {
}

func (n nopAgentTracker) Run(ctx context.Context) error {
	return nil
}

func (n nopAgentTracker) RegisterConnection(ctx context.Context, info *agent_tracker.ConnectedAgentInfo) bool {
	return true
}

func (n nopAgentTracker) UnregisterConnection(ctx context.Context, info *agent_tracker.ConnectedAgentInfo) bool {
	return true
}

func (n nopAgentTracker) GetConnectionsByAgentId(ctx context.Context, agentId int64, cb agent_tracker.ConnectedAgentInfoCallback) error {
	return nil
}

func (n nopAgentTracker) GetConnectionsByProjectId(ctx context.Context, projectId int64, cb agent_tracker.ConnectedAgentInfoCallback) error {
	return nil
}

type nopTunnelTracker struct {
}

func (n nopTunnelTracker) RegisterTunnel(ctx context.Context, info *tracker.TunnelInfo) bool {
	return true
}

func (n nopTunnelTracker) UnregisterTunnel(ctx context.Context, info *tracker.TunnelInfo) bool {
	return true
}

func (n nopTunnelTracker) GetTunnelsByAgentId(ctx context.Context, agentId int64, cb tracker.GetTunnelsByAgentIdCallback) error {
	return nil
}

func (n nopTunnelTracker) Run(ctx context.Context) error {
	return nil
}

type nopKasRouter struct {
}

func (r nopKasRouter) RegisterAgentApi(desc *grpc.ServiceDesc) {
}

type tokenLimiterApi struct {
	rpcApi modserver.AgentRpcApi
}

func (a *tokenLimiterApi) Log() *zap.Logger {
	return a.rpcApi.Log()
}

func (a *tokenLimiterApi) HandleProcessingError(msg string, err error) {
	a.rpcApi.HandleProcessingError(a.rpcApi.Log(), modshared.NoAgentId, msg, err)
}

func (a *tokenLimiterApi) AgentToken() string {
	return string(a.rpcApi.AgentToken())
}
