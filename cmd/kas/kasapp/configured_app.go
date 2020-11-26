package kasapp

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"time"

	"cloud.google.com/go/profiler"
	"github.com/ash2k/stager"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/cmd"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/agentrpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/api/apiutil"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/gitaly"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/gitlab"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/kas"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/redis"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/tool/metric"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/tool/tlstool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/tool/tracing"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/tool/wstunnel"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/pkg/kascfg"
	"gitlab.com/gitlab-org/gitaly/client"
	grpccorrelation "gitlab.com/gitlab-org/labkit/correlation/grpc"
	"gitlab.com/gitlab-org/labkit/errortracking"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	_ "google.golang.org/grpc/encoding/gzip" // Install the gzip compressor
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/stats"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	authSecretLength      = 32
	defaultMaxMessageSize = 10 * 1024 * 1024

	correlationClientName     = "gitlab-kas"
	tracingServiceName        = "gitlab-kas"
	googleProfilerServiceName = "gitlab-kas"
)

type ConfiguredApp struct {
	Configuration *kascfg.ConfigurationFile
	Log           *zap.Logger
}

func (a *ConfiguredApp) Run(ctx context.Context) error {
	// Metrics
	// TODO use an independent registry with https://gitlab.com/gitlab-org/cluster-integration/gitlab-agent/-/issues/32
	// reg := prometheus.NewPedanticRegistry()
	reg := prometheus.DefaultRegisterer
	gatherer := prometheus.DefaultGatherer
	ssh := metric.ServerStatsHandler()
	csh := metric.ClientStatsHandler()
	//goCollector := prometheus.NewGoCollector()
	cleanup, err := metric.Register(reg, ssh, csh)
	if err != nil {
		return err
	}
	defer cleanup()

	// Start things up
	st := stager.New()
	a.startGoogleProfiler(st)
	a.startObservabilityServer(st, gatherer, reg)
	a.startGrpcServer(st, reg, ssh, csh)
	return st.Run(ctx)
}

func kasUserAgent() string {
	return fmt.Sprintf("gitlab-kas/%s/%s", cmd.Version, cmd.Commit)
}

func (a *ConfiguredApp) constructErrorTracker() (errortracking.Tracker, error) {
	s := a.Configuration.Observability.Sentry
	if s.Dsn == "" {
		return nopTracker{}, nil
	}
	a.Log.Debug("Initializing Sentry error tracking", logz.SentryDSN(s.Dsn), logz.SentryEnv(s.Environment))
	tracker, err := errortracking.NewTracker(
		errortracking.WithSentryDSN(s.Dsn),
		errortracking.WithVersion(kasUserAgent()),
		errortracking.WithSentryEnvironment(s.Environment),
	)
	if err != nil {
		return nil, err
	}
	return tracker, nil
}

func (a *ConfiguredApp) startGoogleProfiler(st stager.Stager) {
	cfg := a.Configuration.Observability.GoogleProfiler
	if !cfg.Enabled {
		return
	}
	stage := st.NextStage()
	stage.Go(func(ctx context.Context) error {
		config := profiler.Config{
			Service:        googleProfilerServiceName,
			ServiceVersion: cmd.Version,
			MutexProfiling: true, // like in LabKit
			ProjectID:      cfg.ProjectId,
		}
		var opts []option.ClientOption
		if cfg.CredentialsFile != "" {
			opts = append(opts, option.WithCredentialsFile(cfg.CredentialsFile))
		}
		err := profiler.Start(config, opts...)
		if err != nil {
			return fmt.Errorf("google profiler: %v", err)
		}
		return nil
	})
}

func (a *ConfiguredApp) startObservabilityServer(st stager.Stager, gatherer prometheus.Gatherer, registerer prometheus.Registerer) {
	cfg := a.Configuration.Observability
	stage := st.NextStage()
	stage.Go(func(ctx context.Context) error {
		lis, err := net.Listen(cfg.Listen.Network.String(), cfg.Listen.Address)
		if err != nil {
			return err
		}
		defer lis.Close() // nolint: errcheck

		a.Log.Info("Observability endpoint is up",
			logz.NetNetworkFromAddr(lis.Addr()),
			logz.NetAddressFromAddr(lis.Addr()),
		)

		metricSrv := &metric.Server{
			Name:                  kasUserAgent(),
			Listener:              lis,
			PrometheusUrlPath:     cfg.Prometheus.UrlPath,
			LivenessProbeUrlPath:  cfg.LivenessProbe.UrlPath,
			ReadinessProbeUrlPath: cfg.ReadinessProbe.UrlPath,
			Gatherer:              gatherer,
			Registerer:            registerer,
		}
		return metricSrv.Run(ctx)
	})
}

func (a *ConfiguredApp) startGrpcServer(st stager.Stager, registerer prometheus.Registerer, ssh, csh stats.Handler) {
	stage := st.NextStage()
	stage.Go(func(ctx context.Context) error {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		cfg := a.Configuration

		gitLabUrl, err := url.Parse(cfg.Gitlab.Address)
		if err != nil {
			return err
		}
		// TLS cert for talking to GitLab/Workhorse.
		clientTLSConfig, err := tlstool.DefaultClientTLSConfigWithCACert(cfg.Gitlab.CaCertificateFile)
		if err != nil {
			return err
		}
		// Secret for JWT signing
		decodedAuthSecret, err := a.loadAuthSecret()
		if err != nil {
			return fmt.Errorf("authentication secret: %v", err)
		}
		// Tracing
		tracer, closer, err := tracing.ConstructTracer(tracingServiceName, cfg.Observability.Tracing.ConnectionString)
		if err != nil {
			return fmt.Errorf("tracing: %v", err)
		}
		defer closer.Close() // nolint: errcheck

		// Sentry
		tracker, err := a.constructErrorTracker()
		if err != nil {
			return fmt.Errorf("error tracker: %v", err)
		}

		// gRPC listener
		lis, err := net.Listen(cfg.Agent.Listen.Network.String(), cfg.Agent.Listen.Address)
		if err != nil {
			return err
		}
		defer lis.Close() // nolint: errcheck

		a.Log.Info("Listening for agentk connections",
			logz.NetNetworkFromAddr(lis.Addr()),
			logz.NetAddressFromAddr(lis.Addr()),
			logz.IsWebSocket(cfg.Agent.Listen.Websocket),
		)

		if cfg.Agent.Listen.Websocket {
			wsWrapper := wstunnel.ListenerWrapper{
				// TODO set timeouts
				ReadLimit: defaultMaxMessageSize,
			}
			lis = wsWrapper.Wrap(lis)
		}

		userAgent := kasUserAgent()
		gitalyClientPool := constructGitalyPool(cfg.Gitaly, csh, tracer, userAgent)
		defer gitalyClientPool.Close() // nolint: errcheck
		gitLabClient := &gitlab.RateLimitingClient{
			Delegate: gitlab.NewClient(
				gitLabUrl,
				decodedAuthSecret,
				gitlab.WithCorrelationClientName(correlationClientName),
				gitlab.WithUserAgent(userAgent),
				gitlab.WithTracer(tracer),
				gitlab.WithLogger(a.Log),
				gitlab.WithTLSConfig(clientTLSConfig),
			),
			Limiter: rate.NewLimiter(
				rate.Limit(cfg.Gitlab.ApiRateLimit.RefillRatePerSecond),
				int(cfg.Gitlab.ApiRateLimit.BucketSize),
			),
		}
		gitLabCachingClient := gitlab.NewCachingClient(gitLabClient, gitlab.CacheOptions{
			CacheTTL:      cfg.Agent.InfoCacheTtl.AsDuration(),
			CacheErrorTTL: cfg.Agent.InfoCacheErrorTtl.AsDuration(),
		}, gitlab.CacheOptions{
			CacheTTL:      cfg.Agent.Gitops.ProjectInfoCacheTtl.AsDuration(),
			CacheErrorTTL: cfg.Agent.Gitops.ProjectInfoCacheErrorTtl.AsDuration(),
		})

		connectionMaxAge := cfg.Agent.Limits.ConnectionMaxAge.AsDuration()
		srv, cleanup, err := kas.NewServer(kas.Config{
			Log: a.Log,
			GitalyPool: &gitaly.Pool{
				ClientPool: gitalyClientPool,
			},
			GitLabClient:                   gitLabCachingClient,
			Registerer:                     registerer,
			ErrorTracker:                   tracker,
			AgentConfigurationPollPeriod:   cfg.Agent.Configuration.PollPeriod.AsDuration(),
			GitopsPollPeriod:               cfg.Agent.Gitops.PollPeriod.AsDuration(),
			UsageReportingPeriod:           cfg.Observability.UsageReportingPeriod.AsDuration(),
			MaxConfigurationFileSize:       cfg.Agent.Limits.MaxConfigurationFileSize,
			MaxGitopsManifestFileSize:      cfg.Agent.Limits.MaxGitopsManifestFileSize,
			MaxGitopsTotalManifestFileSize: cfg.Agent.Limits.MaxGitopsTotalManifestFileSize,
			MaxGitopsNumberOfPaths:         cfg.Agent.Limits.MaxGitopsNumberOfPaths,
			MaxGitopsNumberOfFiles:         cfg.Agent.Limits.MaxGitopsNumberOfFiles,
			ConnectionMaxAge:               connectionMaxAge,
		})
		if err != nil {
			return fmt.Errorf("kas.NewServer: %v", err)
		}
		defer cleanup()

		// TODO construct independent metrics interceptors with https://gitlab.com/gitlab-org/cluster-integration/gitlab-agent/-/issues/32
		grpcStreamServerInterceptors := []grpc.StreamServerInterceptor{
			grpc_prometheus.StreamServerInterceptor, // This one should be the first one to measure all invocations
			apiutil.StreamAgentMetaInterceptor(),    // This one should be the second one to ensure agent presents a token
			grpccorrelation.StreamServerCorrelationInterceptor(grpccorrelation.WithoutPropagation()),
			grpc_opentracing.StreamServerInterceptor(grpc_opentracing.WithTracer(tracer)),
			grpctool.StreamServerCtxAugmentingInterceptor(grpctool.JoinContexts(ctx)),
		}
		grpcUnaryServerInterceptors := []grpc.UnaryServerInterceptor{
			grpc_prometheus.UnaryServerInterceptor, // This one should be the first one to measure all invocations
			apiutil.UnaryAgentMetaInterceptor(),    // This one should be the second one to ensure agent presents a token
			grpccorrelation.UnaryServerCorrelationInterceptor(grpccorrelation.WithoutPropagation()),
			grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(tracer)),
			grpctool.UnaryServerCtxAugmentingInterceptor(grpctool.JoinContexts(ctx)),
		}
		if cfg.Redis != nil {
			redisCfg := &redis.Config{
				// Url is parsed below
				Password:       cfg.Redis.Password,
				MaxIdle:        int32(cfg.Redis.MaxIdle),
				MaxActive:      int32(cfg.Redis.MaxActive),
				ReadTimeout:    cfg.Redis.ReadTimeout.AsDuration(),
				WriteTimeout:   cfg.Redis.WriteTimeout.AsDuration(),
				KeepAlive:      cfg.Redis.Keepalive.AsDuration(),
				SentinelMaster: cfg.Redis.SentinelMaster,
				// Sentinels is parsed below
			}
			if cfg.Redis.Url != "" {
				redisCfg.URL, err = url.Parse(cfg.Redis.Url)
				if err != nil {
					return fmt.Errorf("kas.redis.NewPool: redis.url is not a valid URL: %v", err)
				}
			}
			for i, addr := range cfg.Redis.Sentinels {
				u, err := url.Parse(addr)
				if err != nil {
					return fmt.Errorf("kas.redis.NewPool: redis.sentinels[%d] is not a valid URL: %v", i, err)
				}
				redisCfg.Sentinels = append(redisCfg.Sentinels, u)
			}
			redisPool := redis.NewPool(redisCfg)
			agentConnectionLimiter := redis.NewTokenLimiter(
				a.Log,
				redisPool,
				cfg.Agent.Limits.RedisKeyPrefix,
				uint64(cfg.Agent.Limits.ConnectionsPerTokenPerMinute),
				func(ctx context.Context) string { return string(apiutil.AgentTokenFromContext(ctx)) },
			)
			grpcStreamServerInterceptors = append(grpcStreamServerInterceptors, grpctool.StreamServerLimitingInterceptor(agentConnectionLimiter))
			grpcUnaryServerInterceptors = append(grpcUnaryServerInterceptors, grpctool.UnaryServerLimitingInterceptor(agentConnectionLimiter))
		}

		serverOpts := []grpc.ServerOption{
			grpc.StatsHandler(ssh),
			grpc.ChainStreamInterceptor(grpcStreamServerInterceptors...),
			grpc.ChainUnaryInterceptor(grpcUnaryServerInterceptors...),
			grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
				MinTime:             20 * time.Second,
				PermitWithoutStream: true,
			}),
			grpc.KeepaliveParams(keepalive.ServerParameters{
				// MaxConnectionAge should be below connectionMaxAge so that when kas closes a long running response
				// stream, gRPC will close the underlying connection. -20% to account for jitter (see doc for the field)
				// and ensure it's somewhat below connectionMaxAge.
				// See https://github.com/grpc/grpc-go/blob/v1.33.1/internal/transport/http2_server.go#L949-L1047 to better understand how this all works.
				MaxConnectionAge: time.Duration(0.8 * float64(connectionMaxAge)),
				// Give pending RPCs plenty of time to complete.
				// In practice it will happen in 10-30% of connectionMaxAge time (see above).
				MaxConnectionAgeGrace: connectionMaxAge,
				// trying to stay below 60 seconds (typical load-balancer timeout)
				Time: 50 * time.Second,
			}),
		}

		certFile := cfg.Agent.Listen.CertificateFile
		keyFile := cfg.Agent.Listen.KeyFile
		switch {
		case certFile != "" && keyFile != "":
			config, err := tlstool.DefaultServerTLSConfig(certFile, keyFile)
			if err != nil {
				return err
			}
			serverOpts = append(serverOpts, grpc.Creds(credentials.NewTLS(config)))
		case certFile == "" && keyFile == "":
		default:
			return fmt.Errorf("both certificate_file (%s) and key_file (%s) must be either set or not set", certFile, keyFile)
		}

		grpcServer := grpc.NewServer(serverOpts...)
		agentrpc.RegisterKasServer(grpcServer, srv)

		var wg wait.Group
		defer wg.Wait() // wait for grpcServer to shutdown
		defer cancel()  // cancel ctx
		wg.StartWithContext(ctx, srv.Run)
		wg.Start(func() {
			<-ctx.Done() // can be cancelled because Serve() failed or because main ctx was cancelled
			grpcServer.GracefulStop()
		})
		return grpcServer.Serve(lis)
	})
}

func (a *ConfiguredApp) loadAuthSecret() ([]byte, error) {
	encodedAuthSecret, err := ioutil.ReadFile(a.Configuration.Gitlab.AuthenticationSecretFile) // nolint: gosec
	if err != nil {
		return nil, fmt.Errorf("read file: %v", err)
	}
	decodedAuthSecret := make([]byte, authSecretLength)

	n, err := base64.StdEncoding.Decode(decodedAuthSecret, encodedAuthSecret)
	if err != nil {
		return nil, fmt.Errorf("decoding: %v", err)
	}
	if n != authSecretLength {
		return nil, fmt.Errorf("decoding: expecting %d bytes, was %d", authSecretLength, n)
	}
	return decodedAuthSecret, nil
}

func constructGitalyPool(g *kascfg.GitalyCF, csh stats.Handler, tracer opentracing.Tracer, userAgent string) *client.Pool {
	globalGitalyRpcLimiter := rate.NewLimiter(
		rate.Limit(g.GlobalApiRateLimit.RefillRatePerSecond),
		int(g.GlobalApiRateLimit.BucketSize),
	)
	return client.NewPoolWithOptions(
		client.WithDialOptions(
			grpc.WithUserAgent(userAgent),
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
					grpccorrelation.StreamClientCorrelationInterceptor(grpccorrelation.WithClientName(correlationClientName)),
					grpc_opentracing.StreamClientInterceptor(grpc_opentracing.WithTracer(tracer)),
					grpctool.StreamClientLimitingInterceptor(globalGitalyRpcLimiter),
					grpctool.StreamClientLimitingInterceptor(perServerGitalyRpcLimiter),
				),
				grpc.WithChainUnaryInterceptor(
					grpc_prometheus.UnaryClientInterceptor,
					grpccorrelation.UnaryClientCorrelationInterceptor(grpccorrelation.WithClientName(correlationClientName)),
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

// nopTracker is the state of the art error tracking facility.
type nopTracker struct {
}

func (n nopTracker) Capture(err error, opts ...errortracking.CaptureOption) {
}
