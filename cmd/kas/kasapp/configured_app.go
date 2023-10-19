package kasapp

import (
	"bytes"
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
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/rueidisotel"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	promexp "go.opentelemetry.io/otel/exporters/prometheus"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	_ "google.golang.org/grpc/encoding/gzip" // Install the gzip compressor
	"google.golang.org/grpc/stats"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/cmd"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/cmd/kas/kasapp/fake"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitaly"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitaly/vendored/client"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitlab"
	gapi "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitlab/api"
	agent_configuration_server "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/agent_configuration/server"
	agent_registrar_server "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/agent_registrar/server"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/agent_tracker"
	agent_tracker_server "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/agent_tracker/server"
	configuration_project_server "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/configuration_project/server"
	flux_server "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/flux/server"
	gitlab_access_server "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/gitlab_access/server"
	gitops_server "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/gitops/server"
	google_profiler_server "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/google_profiler/server"
	kubernetes_api_server "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/kubernetes_api/server"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modshared"
	notifications_server "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/notifications/server"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/observability"
	observability_server "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/observability/server"
	reverse_tunnel_server "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/reverse_tunnel/server"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/reverse_tunnel/tunnel"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/usage_metrics"
	usage_metrics_server "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/usage_metrics/server"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/cache"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/errz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/httpz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/ioz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/metric"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/prototool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/redistool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/retry"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/tlstool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/kascfg"
)

const (
	routingAttemptInterval   = 50 * time.Millisecond
	routingInitBackoff       = 100 * time.Millisecond
	routingMaxBackoff        = 1 * time.Second
	routingResetDuration     = 10 * time.Second
	routingBackoffFactor     = 2.0
	routingJitter            = 1.0
	routingTunnelFindTimeout = 20 * time.Second
	routingCachePeriod       = 5 * time.Minute
	routingTryNewKasInterval = 10 * time.Millisecond

	authSecretLength = 32

	kasName = "gitlab-kas"

	kasTracerName = "kas"
	kasMeterName  = "kas"

	gitlabBuildInfoGaugeMetricName               = "gitlab_build_info"
	kasVersionAttr                 attribute.Key = "version"
	kasBuiltAttr                   attribute.Key = "built"
)

type ConfiguredApp struct {
	Log           *zap.Logger
	Configuration *kascfg.ConfigurationFile
}

func (a *ConfiguredApp) Run(ctx context.Context) (retErr error) {
	// Metrics
	reg := prometheus.NewPedanticRegistry()
	ssh := grpctool.NewServerRequestsInFlightStatsHandler()
	csh := grpctool.NewClientRequestsInFlightStatsHandler()
	goCollector := collectors.NewGoCollector()
	procCollector := collectors.NewProcessCollector(collectors.ProcessCollectorOpts{})
	srvProm := grpc_prometheus.NewServerMetrics()
	clientProm := grpc_prometheus.NewClientMetrics()
	err := metric.Register(reg, ssh, csh, goCollector, procCollector, srvProm, clientProm)
	if err != nil {
		return err
	}
	streamProm := srvProm.StreamServerInterceptor()
	unaryProm := srvProm.UnaryServerInterceptor()
	streamClientProm := clientProm.StreamClientInterceptor()
	unaryClientProm := clientProm.UnaryClientInterceptor()

	// Probe Registry
	probeRegistry := observability.NewProbeRegistry()

	// OTEL resource
	r, err := constructOTELResource()
	if err != nil {
		return err
	}

	// OTEL metrics
	mp, mpStop, err := a.constructOTELMeterProvider(r, reg) // nolint: contextcheck
	if err != nil {
		return err
	}
	defer errz.SafeCall(mpStop, &retErr)
	dm := mp.Meter(kasMeterName)
	err = gitlabBuildInfoGauge(dm)
	if err != nil {
		return err
	}

	// OTEL Tracing
	tp, p, tpStop, err := a.constructOTELTracingTools(ctx, r)
	if err != nil {
		return err
	}
	defer errz.SafeCall(tpStop, &retErr)
	dt := tp.Tracer(kasTracerName) // defaultTracer

	// GitLab REST client
	gitLabClient, err := a.constructGitLabClient(dt, dm, tp, mp, p)
	if err != nil {
		return err
	}

	// Sentry
	sentryHub, err := a.constructSentryHub(tp, mp, p)
	if err != nil {
		return fmt.Errorf("error tracker: %w", err)
	}

	// Redis
	redisClient, err := a.constructRedisClient(tp, mp)
	if err != nil {
		return err
	}
	defer redisClient.Close()
	probeRegistry.RegisterReadinessProbe("redis", constructRedisReadinessProbe(redisClient))

	srvApi := newServerApi(a.Log, sentryHub, redisClient)
	errRep := modshared.ApiToErrReporter(srvApi)
	grpcServerErrorReporter := &serverErrorReporter{log: a.Log, errReporter: errRep}

	// RPC API factory
	// Plural: Use fake factory
	rpcApiFactory, agentRpcApiFactory := a.constructFakeRpcApiFactory(errRep, sentryHub, redisClient, dt)

	// Server for handling API requests from other kas instances
	privateApiSrv, err := newPrivateApiServer(a.Log, errRep, a.Configuration, tp, mp, p, csh, ssh, rpcApiFactory, // nolint: contextcheck
		probeRegistry, streamProm, unaryProm, streamClientProm, unaryClientProm, grpcServerErrorReporter)
	if err != nil {
		return fmt.Errorf("private API server: %w", err)
	}

	// Server for handling agentk requests
	agentSrv, err := newAgentServer(a.Log, a.Configuration, srvApi, dt, dm, tp, mp, redisClient, ssh, agentRpcApiFactory, // nolint: contextcheck
		privateApiSrv.ownUrl, probeRegistry, reg, streamProm, unaryProm, grpcServerErrorReporter)
	if err != nil {
		return fmt.Errorf("agent server: %w", err)
	}

	// Server for handling external requests e.g. from GitLab
	apiSrv, err := newApiServer(a.Log, a.Configuration, tp, mp, p, ssh, rpcApiFactory, probeRegistry, // nolint: contextcheck
		streamProm, unaryProm, grpcServerErrorReporter)
	if err != nil {
		return fmt.Errorf("API server: %w", err)
	}

	// Construct internal gRPC server
	internalSrv, err := newInternalServer(tp, mp, p, rpcApiFactory, probeRegistry, grpcServerErrorReporter) // nolint: contextcheck
	if err != nil {
		return err
	}
	defer errz.SafeClose(internalSrv.inMemConn, &retErr)

	// Kas to agentk router
	pollConfig := retry.NewPollConfigFactory(routingAttemptInterval, retry.NewExponentialBackoffFactory(
		routingInitBackoff,
		routingMaxBackoff,
		routingResetDuration,
		routingBackoffFactor,
		routingJitter,
	))
	tunnelQuerier := tunnel.NewAggregatingQuerier(a.Log, agentSrv.tunnelRegistry, srvApi, pollConfig, routingCachePeriod)
	kasToAgentRouter, err := newRouter(
		privateApiSrv.kasPool,
		tunnelQuerier,
		agentSrv.tunnelRegistry,
		privateApiSrv.ownUrl,
		internalSrv.server,
		privateApiSrv,
		pollConfig,
		tp,
		reg)
	if err != nil {
		return err
	}

	// Agent tracker
	agentTracker := a.constructAgentTracker(errRep, redisClient)

	// Usage tracker
	usageTracker := usage_metrics.NewUsageTracker()

	// Gitaly client
	gitalyClientPool, err := a.constructGitalyPool(csh, dt, dm, tp, mp, p, streamClientProm, unaryClientProm)
	if err != nil {
		return err
	}
	defer errz.SafeClose(gitalyClientPool, &retErr)

	// Module factories
	factories := []modserver.Factory{
		&observability_server.Factory{
			Gatherer: reg,
		},
		&google_profiler_server.Factory{},
		&agent_configuration_server.Factory{
			AgentRegisterer: agentTracker,
		},
		&configuration_project_server.Factory{},
		&notifications_server.Factory{
			PublishEvent:      srvApi.publishEvent,
			SubscribeToEvents: srvApi.subscribeToEvents,
		},
		&flux_server.Factory{},
		&gitops_server.Factory{},
		&usage_metrics_server.Factory{
			UsageTracker: usageTracker,
		},
		&gitlab_access_server.Factory{},
		&agent_registrar_server.Factory{
			AgentRegisterer: agentTracker,
		},
		&agent_tracker_server.Factory{
			AgentQuerier: agentTracker,
		},
		&reverse_tunnel_server.Factory{
			TunnelHandler: agentSrv.tunnelRegistry,
		},
		&kubernetes_api_server.Factory{},
	}

	// Construct modules
	poolWrapper := &gitaly.Pool{
		ClientPool: gitalyClientPool,
	}
	var beforeServersModules, afterServersModules []modserver.Module
	for _, factory := range factories {
		// factory.New() must be called from the main goroutine because it may mutate a gRPC server (register an API)
		// and that can only be done before Serve() is called on the server.
		moduleName := factory.Name()
		module, err := factory.New(&modserver.Config{
			Log:              a.Log.With(logz.ModuleName(moduleName)),
			Api:              srvApi,
			Config:           a.Configuration,
			GitLabClient:     gitLabClient,
			Registerer:       reg,
			UsageTracker:     usageTracker,
			AgentServer:      agentSrv.server,
			ApiServer:        apiSrv.server,
			RegisterAgentApi: kasToAgentRouter.RegisterAgentApi,
			AgentConn:        internalSrv.inMemConn,
			Gitaly:           poolWrapper,
			TraceProvider:    tp,
			TracePropagator:  p,
			MeterProvider:    mp,
			RedisClient:      redisClient,
			KasName:          kasName,
			Version:          cmd.Version,
			CommitId:         cmd.Commit,
			ProbeRegistry:    probeRegistry,
		})
		if err != nil {
			return fmt.Errorf("%s: %w", moduleName, err)
		}
		phase := factory.StartStopPhase()
		switch phase {
		case modshared.ModuleStartBeforeServers:
			beforeServersModules = append(beforeServersModules, module)
		case modshared.ModuleStartAfterServers:
			afterServersModules = append(afterServersModules, module)
		default:
			return fmt.Errorf("invalid StartStopPhase from factory %s: %d", moduleName, phase)
		}
	}

	// Start things up. Stages are shut down in reverse order.
	return stager.RunStages(ctx,
		// Start things that modules use.
		func(stage stager.Stage) {
			stage.Go(agentTracker.Run)
			stage.Go(tunnelQuerier.Run)
		},
		// Start modules.
		func(stage stager.Stage) {
			startModules(stage, beforeServersModules)
		},
		// Start internal gRPC server. This one must be shut down after all other servers have stopped to ensure
		// it's impossible for them to make a request to the internal server and get a failure because
		// it has stopped already.
		func(stage stager.Stage) {
			internalSrv.Start(stage)
		},
		// Start other gRPC servers.
		func(stage stager.Stage) {
			agentSrv.Start(stage)
			apiSrv.Start(stage)
			privateApiSrv.Start(stage)
		},
		// Start modules.
		func(stage stager.Stage) {
			startModules(stage, afterServersModules)
		},
	)
}

func (a *ConfiguredApp) constructRpcApiFactory(errRep errz.ErrReporter, sentryHub *sentry.Hub, gitLabClient gitlab.ClientInterface, redisClient rueidis.Client, dt trace.Tracer) (modserver.RpcApiFactory, modserver.AgentRpcApiFactory) {
	aCfg := a.Configuration.Agent
	f := serverRpcApiFactory{
		log:       a.Log,
		sentryHub: sentryHub,
	}
	fAgent := serverAgentRpcApiFactory{
		rpcApiFactory: f.New,
		gitLabClient:  gitLabClient,
		agentInfoCache: cache.NewWithError[api.AgentToken, *api.AgentInfo](
			aCfg.InfoCacheTtl.AsDuration(),
			aCfg.InfoCacheErrorTtl.AsDuration(),
			&redistool.ErrCacher[api.AgentToken]{
				Log:          a.Log,
				ErrRep:       errRep,
				Client:       redisClient,
				ErrMarshaler: prototool.ProtoErrMarshaler{},
				KeyToRedisKey: func(key api.AgentToken) string {
					return a.Configuration.Redis.KeyPrefix + ":agent_info_errs:" + string(api.AgentToken2key(key))
				},
			},
			dt,
			gapi.IsCacheableError,
		),
	}
	return f.New, fAgent.New
}

func (a *ConfiguredApp) constructFakeRpcApiFactory(errRep errz.ErrReporter, sentryHub *sentry.Hub, redisClient rueidis.Client, dt trace.Tracer) (modserver.RpcApiFactory, modserver.AgentRpcApiFactory) {
	aCfg := a.Configuration.Agent
	f := serverRpcApiFactory{
		log:       a.Log,
		sentryHub: sentryHub,
	}
	fAgent := fake.ServerAgentRpcApiFactory{
		RPCApiFactory: f.New,
		AgentInfoCache: cache.NewWithError[api.AgentToken, *api.AgentInfo](
			aCfg.InfoCacheTtl.AsDuration(),
			aCfg.InfoCacheErrorTtl.AsDuration(),
			&redistool.ErrCacher[api.AgentToken]{
				Log:          a.Log,
				ErrRep:       errRep,
				Client:       redisClient,
				ErrMarshaler: prototool.ProtoErrMarshaler{},
				KeyToRedisKey: func(key api.AgentToken) string {
					return a.Configuration.Redis.KeyPrefix + ":agent_info_errs:" + string(api.AgentToken2key(key))
				},
			},
			dt,
			gapi.IsCacheableError,
		),
	}
	return f.New, fAgent.New
}

func (a *ConfiguredApp) constructAgentTracker(errRep errz.ErrReporter, redisClient rueidis.Client) agent_tracker.Tracker {
	cfg := a.Configuration
	return agent_tracker.NewRedisTracker(
		a.Log,
		errRep,
		redisClient,
		cfg.Redis.KeyPrefix+":agent_tracker2",
		cfg.Agent.RedisConnInfoTtl.AsDuration(),
		cfg.Agent.RedisConnInfoRefresh.AsDuration(),
		cfg.Agent.RedisConnInfoGc.AsDuration(),
	)
}

func (a *ConfiguredApp) constructSentryHub(tp trace.TracerProvider, mp otelmetric.MeterProvider, p propagation.TextMapPropagator) (*sentry.Hub, error) {
	s := a.Configuration.Observability.Sentry
	dialer := net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	sentryClient, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:         s.Dsn, // empty DSN disables Sentry transport
		SampleRate:  1,     // no sampling
		Release:     cmd.Version,
		Environment: s.Environment,
		HTTPTransport: otelhttp.NewTransport(
			&http.Transport{
				Proxy:                 http.ProxyFromEnvironment,
				DialContext:           dialer.DialContext,
				TLSClientConfig:       tlstool.DefaultClientTLSConfig(),
				MaxIdleConns:          10,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ResponseHeaderTimeout: 20 * time.Second,
				ForceAttemptHTTP2:     true,
			},
			otelhttp.WithPropagators(p),
			otelhttp.WithTracerProvider(tp),
			otelhttp.WithMeterProvider(mp),
		),
	})
	if err != nil {
		return nil, err
	}
	return sentry.NewHub(sentryClient, sentry.NewScope()), nil
}

func (a *ConfiguredApp) loadGitLabClientAuthSecret() ([]byte, error) {
	decodedAuthSecret, err := ioz.LoadBase64Secret(a.Configuration.Gitlab.AuthenticationSecretFile)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	if len(decodedAuthSecret) != authSecretLength {
		return nil, fmt.Errorf("decoding: expecting %d bytes, was %d", authSecretLength, len(decodedAuthSecret))
	}
	return decodedAuthSecret, nil
}

func (a *ConfiguredApp) constructGitLabClient(dt trace.Tracer, dm otelmetric.Meter,
	tp trace.TracerProvider, mp otelmetric.MeterProvider, p propagation.TextMapPropagator) (*gitlab.Client, error) {
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
	var limiter httpz.Limiter
	limiter = rate.NewLimiter(
		rate.Limit(cfg.Gitlab.ApiRateLimit.RefillRatePerSecond),
		int(cfg.Gitlab.ApiRateLimit.BucketSize),
	)
	limiter, err = metric.NewWaitLimiterInstrumentation(
		"gitlab_client",
		cfg.Gitlab.ApiRateLimit.RefillRatePerSecond,
		"{refill/s}",
		dt,
		dm,
		limiter,
	)
	if err != nil {
		return nil, err
	}
	return gitlab.NewClient(
		gitLabUrl,
		decodedAuthSecret,
		gitlab.WithTextMapPropagator(p),
		gitlab.WithTracerProvider(tp),
		gitlab.WithMeterProvider(mp),
		gitlab.WithUserAgent(kasServerName()),
		gitlab.WithTLSConfig(clientTLSConfig),
		gitlab.WithRateLimiter(limiter),
	), nil
}

func (a *ConfiguredApp) constructGitalyPool(csh stats.Handler, dt trace.Tracer, dm otelmetric.Meter,
	tp trace.TracerProvider, mp otelmetric.MeterProvider,
	p propagation.TextMapPropagator, streamClientProm grpc.StreamClientInterceptor, unaryClientProm grpc.UnaryClientInterceptor) (*client.Pool, error) {
	g := a.Configuration.Gitaly
	var globalGitalyRpcLimiter grpctool.ClientLimiter
	globalGitalyRpcLimiter = rate.NewLimiter(
		rate.Limit(g.GlobalApiRateLimit.RefillRatePerSecond),
		int(g.GlobalApiRateLimit.BucketSize),
	)
	globalGitalyRpcLimiter, err := metric.NewWaitLimiterInstrumentation(
		"gitaly_client_global",
		g.GlobalApiRateLimit.RefillRatePerSecond,
		"{refill/s}",
		dt,
		dm,
		globalGitalyRpcLimiter,
	)
	if err != nil {
		return nil, err
	}
	return client.NewPoolWithOptions(
		client.WithDialOptions(
			grpc.WithUserAgent(kasServerName()),
			grpc.WithStatsHandler(csh),
			grpc.WithStatsHandler(otelgrpc.NewServerHandler(
				otelgrpc.WithTracerProvider(tp),
				otelgrpc.WithMeterProvider(mp),
				otelgrpc.WithPropagators(p),
				otelgrpc.WithMessageEvents(otelgrpc.ReceivedEvents, otelgrpc.SentEvents),
			)),
			grpc.WithSharedWriteBuffer(true),
			// In https://gitlab.com/groups/gitlab-org/-/epics/8971, we added DNS discovery support to Praefect. This was
			// done by making two changes:
			// - Configure client-side round-robin load-balancing in client dial options. We added that as a default option
			// inside gitaly client in gitaly client since v15.9.0
			// - Configure DNS resolving. Due to some technical limitations, we don't use gRPC's built-in DNS resolver.
			// Instead, we implement our own DNS resolver. This resolver is exposed via the following configuration.
			// Afterward, workhorse can detect and handle DNS discovery automatically. The user needs to setup and set
			// Gitaly address to something like "dns:gitaly.service.dc1.consul"
			client.WithGitalyDNSResolver(client.DefaultDNSResolverBuilderConfig()),
			// Don't put interceptors here as order is important. Put them below.
		),
		client.WithDialer(func(ctx context.Context, address string, dialOptions []grpc.DialOption) (*grpc.ClientConn, error) {
			var perServerGitalyRpcLimiter grpctool.ClientLimiter
			perServerGitalyRpcLimiter = rate.NewLimiter(
				rate.Limit(g.PerServerApiRateLimit.RefillRatePerSecond),
				int(g.PerServerApiRateLimit.BucketSize))
			perServerGitalyRpcLimiter, err := metric.NewWaitLimiterInstrumentation(
				"gitaly_client_"+address,
				g.GlobalApiRateLimit.RefillRatePerSecond,
				"{refill/s}",
				dt,
				dm,
				perServerGitalyRpcLimiter,
			)
			if err != nil {
				return nil, err
			}
			opts := []grpc.DialOption{
				grpc.WithChainStreamInterceptor(
					streamClientProm,
					grpctool.StreamClientLimitingInterceptor(globalGitalyRpcLimiter),
					grpctool.StreamClientLimitingInterceptor(perServerGitalyRpcLimiter),
				),
				grpc.WithChainUnaryInterceptor(
					unaryClientProm,
					grpctool.UnaryClientLimitingInterceptor(globalGitalyRpcLimiter),
					grpctool.UnaryClientLimitingInterceptor(perServerGitalyRpcLimiter),
				),
			}
			opts = append(opts, dialOptions...)
			return client.DialContext(ctx, address, opts)
		}),
	), nil
}

func (a *ConfiguredApp) constructRedisClient(tp trace.TracerProvider, mp otelmetric.MeterProvider) (rueidis.Client, error) {
	cfg := a.Configuration.Redis
	dialTimeout := cfg.DialTimeout.AsDuration()
	writeTimeout := cfg.WriteTimeout.AsDuration()
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
		passwordBytes, err := os.ReadFile(cfg.PasswordFile) // nolint:govet
		if err != nil {
			return nil, err
		}
		password = string(passwordBytes)
	}
	opts := rueidis.ClientOption{
		Dialer: net.Dialer{
			Timeout: dialTimeout,
		},
		TLSConfig:        tlsConfig,
		Username:         cfg.Username,
		Password:         password,
		ClientName:       kasName,
		ConnWriteTimeout: writeTimeout,
		MaxFlushDelay:    20 * time.Microsecond,
		DisableCache:     true,
		SelectDB:         int(cfg.DatabaseIndex),
	}
	if cfg.Network == "unix" {
		opts.DialFn = redistool.UnixDialer
	}
	switch v := cfg.RedisConfig.(type) {
	case *kascfg.RedisCF_Server:
		opts.InitAddress = []string{v.Server.Address}
		if opts.TLSConfig != nil {
			opts.TLSConfig.ServerName, _, _ = strings.Cut(v.Server.Address, ":")
		}
	case *kascfg.RedisCF_Sentinel:
		opts.InitAddress = v.Sentinel.Addresses
		var sentinelPassword string
		if v.Sentinel.SentinelPasswordFile != "" {
			sentinelPasswordBytes, err := os.ReadFile(v.Sentinel.SentinelPasswordFile) // nolint:govet
			if err != nil {
				return nil, err
			}
			sentinelPassword = string(sentinelPasswordBytes)
		}
		opts.Sentinel = rueidis.SentinelOption{
			Dialer:    opts.Dialer,
			TLSConfig: opts.TLSConfig,
			MasterSet: v.Sentinel.MasterName,
			Username:  cfg.Username,
			Password:  sentinelPassword,
		}
	default:
		// This should never happen
		return nil, fmt.Errorf("unexpected Redis config type: %T", cfg.RedisConfig)
	}
	redisClient, err := rueidis.NewClient(opts)
	if err != nil {
		return nil, err
	}
	if a.isTracingEnabled() {
		// Instrument Redis client with tracing only if it's configured.
		redisClient = rueidisotel.WithClient(redisClient, rueidisotel.WithTracerProvider(tp), rueidisotel.WithMeterProvider(mp))
	}
	return redisClient, nil
}

func constructTracingExporter(ctx context.Context, tracingConfig *kascfg.TracingCF) (tracesdk.SpanExporter, error) {
	otlpEndpoint := tracingConfig.OtlpEndpoint

	u, err := url.Parse(otlpEndpoint)
	if err != nil {
		return nil, fmt.Errorf("parsing tracing url %s failed: %w", otlpEndpoint, err)
	}

	var otlpOptions []otlptracehttp.Option

	switch u.Scheme {
	case "https":
	case "http":
		otlpOptions = append(otlpOptions, otlptracehttp.WithInsecure())
	default:
		return nil, fmt.Errorf("unsupported schema of tracing url %q, only `http` and `https` are permitted", u.Scheme)
	}

	otlpOptions = append(otlpOptions, otlptracehttp.WithEndpoint(u.Host))
	otlpOptions = append(otlpOptions, otlptracehttp.WithURLPath(u.Path))

	otlpTokenSecretFile := tracingConfig.OtlpTokenSecretFile
	if otlpTokenSecretFile != nil {
		token, err := os.ReadFile(*otlpTokenSecretFile) // nolint: gosec, govet
		if err != nil {
			return nil, fmt.Errorf("unable to read OTLP token from %q: %w", *otlpTokenSecretFile, err)
		}
		token = bytes.TrimSpace(token)

		// This is just a temporary measure to allow for smooth migration from
		// Gitlab Observability UI tokens to Gitlab Access Tokens.
		// Issue: https://gitlab.com/gitlab-org/opstrace/opstrace/-/issues/2148
		//
		// The idea is simple - we try to determine the type of the token and
		// basing on it set correct HTTP headers. Gitlab
		// Observability Backend makes the decision which auth mechanism to use
		// basing on which HTTP header is present.
		headers := make(map[string]string)
		if bytes.HasPrefix(token, []byte("glpat-")) {
			headers["Private-Token"] = string(token)
		} else {
			headers[httpz.AuthorizationHeader] = fmt.Sprintf("Bearer %s", token)
		}

		otlpOptions = append(otlpOptions, otlptracehttp.WithHeaders(headers))
	}

	tlsConfig, err := tlstool.DefaultClientTLSConfigWithCACert(tracingConfig.GetOtlpCaCertificateFile())
	if err != nil {
		return nil, err
	}
	otlpOptions = append(otlpOptions, otlptracehttp.WithTLSClientConfig(tlsConfig))

	return otlptracehttp.New(ctx, otlpOptions...)
}

func (a *ConfiguredApp) constructOTELMeterProvider(r *resource.Resource, reg prometheus.Registerer) (*metricsdk.MeterProvider, func() error, error) {
	otelPromExp, err := promexp.New(promexp.WithRegisterer(reg))
	if err != nil {
		return nil, nil, err
	}
	mp := metricsdk.NewMeterProvider(metricsdk.WithReader(otelPromExp), metricsdk.WithResource(r))
	mpStop := func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return mp.Shutdown(ctx)
	}
	return mp, mpStop, nil
}

func (a *ConfiguredApp) constructOTELTracingTools(ctx context.Context, r *resource.Resource) (trace.TracerProvider, propagation.TextMapPropagator, func() error, error) {
	if !a.isTracingEnabled() {
		return trace.NewNoopTracerProvider(), propagation.NewCompositeTextMapPropagator(), func() error { return nil }, nil
	}

	// Exporter must be constructed right before TracerProvider as it's started implicitly so needs to be stopped,
	// which TracerProvider does in its Shutdown() method.
	exporter, err := constructTracingExporter(ctx, a.Configuration.Observability.Tracing)
	if err != nil {
		return nil, nil, nil, err
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithResource(r),
		tracesdk.WithBatcher(exporter),
	)
	p := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	tpStop := func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return tp.Shutdown(ctx)
	}
	return tp, p, tpStop, nil
}

func (a *ConfiguredApp) isTracingEnabled() bool {
	return a.Configuration.Observability.Tracing != nil
}

func startModules(stage stager.Stage, modules []modserver.Module) {
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
}

func constructRedisReadinessProbe(redisClient rueidis.Client) observability.Probe {
	return func(ctx context.Context) error {
		pingCmd := redisClient.B().Ping().Build()
		err := redisClient.Do(ctx, pingCmd).Error()
		if err != nil {
			return fmt.Errorf("redis: %w", err)
		}
		return nil
	}
}

func constructOTELResource() (*resource.Resource, error) {
	return resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(kasName),
			semconv.ServiceVersion(cmd.Version),
		),
	)
}

func gitlabBuildInfoGauge(m otelmetric.Meter) error {
	// Only allocate the option once
	attributes := otelmetric.WithAttributeSet(attribute.NewSet(kasVersionAttr.String(cmd.Version), kasBuiltAttr.String(cmd.BuildTime)))
	_, err := m.Int64ObservableGauge(gitlabBuildInfoGaugeMetricName,
		otelmetric.WithDescription("Current build info for this GitLab Service"),
		otelmetric.WithInt64Callback(func(ctx context.Context, observer otelmetric.Int64Observer) error {
			observer.Observe(1, attributes)
			return nil
		}),
	)
	return err
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
	_ redistool.RpcApi = (*tokenLimiterApi)(nil)
)

type tokenLimiterApi struct {
	rpcApi modserver.AgentRpcApi
}

func (a *tokenLimiterApi) Log() *zap.Logger {
	return a.rpcApi.Log()
}

func (a *tokenLimiterApi) HandleProcessingError(msg string, err error) {
	a.rpcApi.HandleProcessingError(a.rpcApi.Log(), modshared.NoAgentId, msg, err)
}

func (a *tokenLimiterApi) RequestKey() []byte {
	return api.AgentToken2key(a.rpcApi.AgentToken())
}

var (
	_ grpctool.ServerErrorReporter = (*serverErrorReporter)(nil)
)

// serverErrorReporter implements the grpctool.ServerErrorReporter interface
// in order to report unknown grpc status code errors.
// In this case the errz.ErrReporter is used as a proxy for the modserver.RpcApi
// which logs and captures errors in Sentry.
type serverErrorReporter struct {
	log         *zap.Logger
	errReporter errz.ErrReporter
}

func (r *serverErrorReporter) Report(ctx context.Context, fullMethod string, err error) {
	r.errReporter.HandleProcessingError(ctx, r.log, fmt.Sprintf("Unknown gRPC error in %q", fullMethod), err)
}
