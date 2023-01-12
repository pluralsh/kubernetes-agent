package kasapp

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
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
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/cmd"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/api"
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
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/ioz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/metric"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/prototool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/redistool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/retry"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/tlstool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/pkg/kascfg"
	"gitlab.com/gitlab-org/gitaly/v15/client"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	_ "google.golang.org/grpc/encoding/gzip" // Install the gzip compressor
)

const (
	routingAttemptInterval   = 50 * time.Millisecond
	routingInitBackoff       = 100 * time.Millisecond
	routingMaxBackoff        = 1 * time.Second
	routingResetDuration     = 10 * time.Second
	routingBackoffFactor     = 2.0
	routingJitter            = 1.0
	routingTunnelFindTimeout = 20 * time.Second

	authSecretLength = 32

	kasName = "gitlab-kas"

	kasRoutingMetricName              = "k8s_api_proxy_routing_duration_seconds"
	kasRoutingStatusLabelName         = "status"
	kasRoutingStatusSuccessLabelValue = "success"
	kasRoutingStatusErrorLabelValue   = "error"
)

type ConfiguredApp struct {
	Log               *zap.Logger
	Configuration     *kascfg.ConfigurationFile
	OwnPrivateApiUrl  string
	OwnPrivateApiHost string
}

func (a *ConfiguredApp) Run(ctx context.Context) (retErr error) {
	if a.OwnPrivateApiUrl == "" {
		return fmt.Errorf("%s environment variable is required so that kas instance is accessible to other kas instances", envVarOwnPrivateApiUrl)
	}
	// Metrics
	// TODO use an independent registry
	// reg := prometheus.NewPedanticRegistry()
	registerer := prometheus.DefaultRegisterer
	gatherer := prometheus.DefaultGatherer
	//goCollector := prometheus.NewGoCollector()
	err := metric.Register(registerer, gitlabBuildInfoGauge())
	if err != nil {
		return err
	}

	// Probe Registry
	probeRegistry := observability.NewProbeRegistry()

	// Tracing
	tp, p, tpStop, err := a.constructTracingTools(ctx)
	if err != nil {
		return err
	}
	defer func() {
		tpErr := tpStop()
		if retErr == nil && tpErr != nil {
			retErr = tpErr
		}
	}()

	// GitLab REST client
	gitLabClient, err := a.constructGitLabClient()
	if err != nil {
		return err
	}

	// Sentry
	sentryHub, err := a.constructSentryHub(p)
	if err != nil {
		return fmt.Errorf("error tracker: %w", err)
	}

	srvApi := &serverApi{Hub: sentryHub}
	errRep := modshared.ApiToErrReporter(srvApi)

	// Redis
	redisClient, err := a.constructRedisClient()
	if err != nil {
		return err
	}
	probeRegistry.RegisterReadinessProbe("redis", constructRedisReadinessProbe(redisClient))

	// RPC API factory
	rpcApiFactory, agentRpcApiFactory := a.constructRpcApiFactory(errRep, sentryHub, gitLabClient, redisClient)

	// Server for handling agentk requests
	agentSrv, err := newAgentServer(a.Log, a.Configuration, tp, redisClient, agentRpcApiFactory, probeRegistry) // nolint: contextcheck
	if err != nil {
		return fmt.Errorf("agent server: %w", err)
	}

	// Server for handling external requests e.g. from GitLab
	apiSrv, err := newApiServer(a.Log, a.Configuration, tp, p, rpcApiFactory, probeRegistry) // nolint: contextcheck
	if err != nil {
		return fmt.Errorf("API server: %w", err)
	}

	// Server for handling API requests from other kas instances
	privateApiSrv, err := newPrivateApiServer(a.Log, errRep, a.Configuration, tp, p, rpcApiFactory, a.OwnPrivateApiUrl, a.OwnPrivateApiHost, probeRegistry) // nolint: contextcheck
	if err != nil {
		return fmt.Errorf("private API server: %w", err)
	}

	// Construct internal gRPC server
	internalSrv, err := newInternalServer(tp, p, rpcApiFactory, probeRegistry) // nolint: contextcheck
	if err != nil {
		return err
	}
	defer errz.SafeClose(internalSrv.inMemConn, &retErr)

	// Reverse gRPC tunnel tracker
	tunnelTracker := a.constructTunnelTracker(srvApi, redisClient)

	// Tunnel registry
	tunnelRegistry, err := reverse_tunnel.NewTunnelRegistry(a.Log, errRep, tunnelTracker, a.OwnPrivateApiUrl)
	if err != nil {
		return err
	}
	defer tunnelRegistry.Stop() // nolint: contextcheck

	// Kas to agentk router
	kasToAgentRouter, err := a.constructKasToAgentRouter(
		tunnelTracker,
		tunnelRegistry,
		a.OwnPrivateApiUrl,
		internalSrv.server,
		privateApiSrv,
		privateApiSrv.kasPool,
		srvApi,
		registerer)
	if err != nil {
		return err
	}

	// Agent tracker
	agentTracker := a.constructAgentTracker(errRep, redisClient)

	// Usage tracker
	usageTracker := usage_metrics.NewUsageTracker()

	// Gitaly client
	gitalyClientPool := a.constructGitalyPool(tp, p)
	defer errz.SafeClose(gitalyClientPool, &retErr)

	// Module factories
	factories := []modserver.Factory{
		&observability_server.Factory{
			Gatherer: gatherer,
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
	poolWrapper := &gitaly.Pool{
		ClientPool: gitalyClientPool,
	}
	meterProvider := global.MeterProvider() // TODO
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
			Registerer:       registerer,
			UsageTracker:     usageTracker,
			AgentServer:      agentSrv.server,
			ApiServer:        apiSrv.server,
			RegisterAgentApi: kasToAgentRouter.RegisterAgentApi,
			AgentConn:        internalSrv.inMemConn,
			Gitaly:           poolWrapper,
			TraceProvider:    tp,
			TracePropagator:  p,
			MeterProvider:    meterProvider, // TODO https://gitlab.com/gitlab-org/cluster-integration/gitlab-agent/-/issues/359
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
			// Tunnel tracker should be stopped ASAP to unregister this kas replica from Redis
			// to avoid getting more requests on the private API server.
			stage.Go(tunnelTracker.Run)
			startModules(stage, afterServersModules)
		},
	)
}

func (a *ConfiguredApp) constructRpcApiFactory(errRep errz.ErrReporter, sentryHub *sentry.Hub, gitLabClient gitlab.ClientInterface, redisClient redis.UniversalClient) (modserver.RpcApiFactory, modserver.AgentRpcApiFactory) {
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
			gapi.IsCacheableError,
		),
	}
	return f.New, fAgent.New
}

func (a *ConfiguredApp) constructKasToAgentRouter(tunnelQuerier tracker.Querier, tunnelFinder reverse_tunnel.TunnelFinder,
	ownPrivateApiUrl string, internalServer, privateApiServer grpc.ServiceRegistrar, kasPool grpctool.PoolInterface,
	api modshared.Api, registerer prometheus.Registerer) (kasRouter, error) {
	gatewayKasVisitor, err := grpctool.NewStreamVisitor(&GatewayKasResponse{})
	if err != nil {
		return nil, err
	}
	kasRoutingDuration := constructKasRoutingDurationHistogram()
	err = metric.Register(registerer, kasRoutingDuration)
	if err != nil {
		return nil, err
	}
	pollConfig := retry.NewPollConfigFactory(routingAttemptInterval, retry.NewExponentialBackoffFactory(
		routingInitBackoff,
		routingMaxBackoff,
		routingResetDuration,
		routingBackoffFactor,
		routingJitter,
	))
	return &router{
		kasPool:                   kasPool,
		tunnelQuerier:             tracker.NewAggregatingQuerier(a.Log, tunnelQuerier, api, pollConfig),
		tunnelFinder:              tunnelFinder,
		ownPrivateApiUrl:          ownPrivateApiUrl,
		pollConfig:                pollConfig,
		internalServer:            internalServer,
		privateApiServer:          privateApiServer,
		gatewayKasVisitor:         gatewayKasVisitor,
		kasRoutingDurationSuccess: kasRoutingDuration.WithLabelValues(kasRoutingStatusSuccessLabelValue),
		kasRoutingDurationError:   kasRoutingDuration.WithLabelValues(kasRoutingStatusErrorLabelValue),
		tunnelFindTimeout:         routingTunnelFindTimeout,
	}, nil
}

func (a *ConfiguredApp) constructAgentTracker(errRep errz.ErrReporter, redisClient redis.UniversalClient) agent_tracker.Tracker {
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

func (a *ConfiguredApp) constructTunnelTracker(api modshared.Api, redisClient redis.UniversalClient) tracker.Tracker {
	cfg := a.Configuration
	return tracker.NewRedisTracker(
		a.Log,
		api,
		redisClient,
		cfg.Redis.KeyPrefix+":tunnel_tracker2",
		cfg.Agent.RedisConnInfoTtl.AsDuration(),
		cfg.Agent.RedisConnInfoRefresh.AsDuration(),
		cfg.Agent.RedisConnInfoGc.AsDuration(),
		a.OwnPrivateApiUrl,
	)
}

func (a *ConfiguredApp) constructSentryHub(p propagation.TextMapPropagator) (*sentry.Hub, error) {
	s := a.Configuration.Observability.Sentry
	dialer := net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	sentryClient, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:         s.Dsn, // empty DSN disables Sentry transport
		SampleRate:  1,     // no sampling
		Release:     fmt.Sprintf("%s/%s", cmd.Version, cmd.Commit),
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
			}, otelhttp.WithPropagators(p),
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

func (a *ConfiguredApp) constructGitLabClient() (*gitlab.Client, error) {
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
		gitlab.WithTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})),
		gitlab.WithUserAgent(kasServerName()),
		gitlab.WithTLSConfig(clientTLSConfig),
		gitlab.WithRateLimiter(rate.NewLimiter(
			rate.Limit(cfg.Gitlab.ApiRateLimit.RefillRatePerSecond),
			int(cfg.Gitlab.ApiRateLimit.BucketSize),
		)),
	), nil
}

func (a *ConfiguredApp) constructGitalyPool(tp trace.TracerProvider, p propagation.TextMapPropagator) *client.Pool {
	g := a.Configuration.Gitaly
	globalGitalyRpcLimiter := rate.NewLimiter(
		rate.Limit(g.GlobalApiRateLimit.RefillRatePerSecond),
		int(g.GlobalApiRateLimit.BucketSize),
	)
	return client.NewPoolWithOptions(
		client.WithDialOptions(
			grpc.WithUserAgent(kasServerName()),
			// Don't put interceptors here as order is important. Put them below.
		),
		client.WithDialer(func(ctx context.Context, address string, dialOptions []grpc.DialOption) (*grpc.ClientConn, error) {
			perServerGitalyRpcLimiter := rate.NewLimiter(
				rate.Limit(g.PerServerApiRateLimit.RefillRatePerSecond),
				int(g.PerServerApiRateLimit.BucketSize))
			opts := []grpc.DialOption{
				grpc.WithChainStreamInterceptor(
					grpc_prometheus.StreamClientInterceptor,
					otelgrpc.StreamClientInterceptor(otelgrpc.WithTracerProvider(tp), otelgrpc.WithPropagators(p)),
					grpctool.StreamClientLimitingInterceptor(globalGitalyRpcLimiter),
					grpctool.StreamClientLimitingInterceptor(perServerGitalyRpcLimiter),
				),
				grpc.WithChainUnaryInterceptor(
					grpc_prometheus.UnaryClientInterceptor,
					otelgrpc.UnaryClientInterceptor(otelgrpc.WithTracerProvider(tp), otelgrpc.WithPropagators(p)),
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

func (a *ConfiguredApp) fetchOtlpOptions() ([]otlptracehttp.Option, error) {
	tracingConfig := a.Configuration.GetObservability().GetTracing()
	otlpEndpoint := tracingConfig.GetOtlpEndpoint()
	otlpTokenPath := tracingConfig.GetOtlpTokenSecretFile()
	otlpCAPath := tracingConfig.GetOtlpCaCertificateFile()

	if otlpEndpoint == "" {
		return nil, nil
	}

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

	if otlpTokenPath == "" {
		return nil, errors.New("OTLP Token path must be specified")
	}
	token, err := os.ReadFile(otlpTokenPath) // nolint: gosec
	if err != nil {
		return nil, fmt.Errorf("unable to read OTLP token from %q: %w", otlpTokenPath, err)
	}

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", bytes.TrimSpace(token)),
	}
	otlpOptions = append(otlpOptions, otlptracehttp.WithHeaders(headers))

	tlsConfig, err := tlstool.DefaultClientTLSConfigWithCACert(otlpCAPath)
	if err != nil {
		return nil, err
	}
	otlpOptions = append(otlpOptions, otlptracehttp.WithTLSClientConfig(tlsConfig))

	return otlpOptions, nil
}

func (a *ConfiguredApp) constructTracingTools(ctx context.Context) (trace.TracerProvider, propagation.TextMapPropagator, func() error /* stop */, error) {
	otlpOptions, err := a.fetchOtlpOptions()
	if err != nil {
		return nil, nil, nil, err
	}

	if len(otlpOptions) == 0 {
		return trace.NewNoopTracerProvider(), propagation.NewCompositeTextMapPropagator(), func() error { return nil }, nil
	}

	exporter, err := otlptracehttp.New(ctx, otlpOptions...)
	if err != nil {
		return nil, nil, nil, err
	}

	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(kasName),
			semconv.ServiceVersionKey.String(cmd.Version),
		),
	)
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

func constructKasRoutingDurationHistogram() *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    kasRoutingMetricName,
		Help:    "The time it takes the routing kas to find a suitable tunnel in seconds",
		Buckets: prometheus.ExponentialBuckets(time.Millisecond.Seconds(), 4, 9), // 9 buckets of milliseconds as seconds (1,4,16,64,256,1k,4k,16k,64k)
	}, []string{kasRoutingStatusLabelName})
}

func constructRedisReadinessProbe(redisClient redis.UniversalClient) observability.Probe {
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
	const GitlabBuildInfoGaugeMetricName = "gitlab_build_info"
	buildInfoGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: GitlabBuildInfoGaugeMetricName,
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
