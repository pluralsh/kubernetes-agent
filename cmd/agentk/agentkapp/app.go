package agentkapp

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/ash2k/stager"
	"github.com/go-logr/zapr"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/spf13/cobra"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/cmd"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/agent_configuration/rpc"
	gitlab_access_rpc "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/gitlab_access/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/gitops/agent/chartops"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/gitops/agent/manifestops"
	kubernetes_api_agent "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/kubernetes_api/agent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modshared"
	observability_agent "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/observability/agent"
	reverse_tunnel_agent "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/reverse_tunnel/agent"
	starboard_vulnerability "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/starboard_vulnerability/agent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/errz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/httpz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/retry"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/tlstool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/wstunnel"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/pkg/agentcfg"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zapgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/keepalive"
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes/scheme"
	client_core_v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	"k8s.io/kubectl/pkg/cmd/util"
	"nhooyr.io/websocket"
)

const (
	defaultLogLevel     agentcfg.LogLevelEnum = 0 // whatever is 0 is the default value
	defaultGrpcLogLevel                       = agentcfg.LogLevelEnum_error

	defaultMaxMessageSize = 10 * 1024 * 1024
	agentName             = "gitlab-agent"

	envVarPodNamespace       = "POD_NAMESPACE"
	envVarPodName            = "POD_NAME"
	envVarServiceAccountName = "SERVICE_ACCOUNT_NAME"

	getConfigurationInitBackoff   = 10 * time.Second
	getConfigurationMaxBackoff    = 5 * time.Minute
	getConfigurationResetDuration = 10 * time.Minute
	getConfigurationBackoffFactor = 2.0
	getConfigurationJitter        = 1.0
)

type App struct {
	Log          *zap.Logger
	LogLevel     zap.AtomicLevel
	GrpcLogLevel zap.AtomicLevel
	AgentMeta    *modshared.AgentMeta
	AgentId      *AgentIdHolder
	// KasAddress specifies the address of kas.
	KasAddress         string
	ServiceAccountName string
	CACertFile         string
	TokenFile          string
	K8sClientGetter    genericclioptions.RESTClientGetter
}

func (a *App) Run(ctx context.Context) (retErr error) {
	// TODO Tracing
	tp := trace.NewNoopTracerProvider()
	p := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})

	// Construct gRPC connection to gitlab-kas
	kasConn, err := a.constructKasConnection(ctx, tp, p)
	if err != nil {
		return err
	}
	defer errz.SafeClose(kasConn, &retErr)

	// Internal gRPC client->listener pipe
	internalListener := grpctool.NewDialListener()

	// Construct internal gRPC server
	internalServer := a.constructInternalServer(ctx, tp, p)

	// Construct connection to internal gRPC server
	internalServerConn, err := a.constructInternalServerConn(ctx, internalListener.DialContext)
	if err != nil {
		return err
	}
	defer errz.SafeClose(internalServerConn, &retErr)

	// Construct Kubernetes tools.
	k8sFactory := util.NewFactory(a.K8sClientGetter)
	kubeClient, err := k8sFactory.KubernetesClientSet()
	if err != nil {
		return err
	}

	// Construct event recorder
	eventBroadcaster := record.NewBroadcaster()
	eventRecorder := eventBroadcaster.NewRecorder(scheme.Scheme, core_v1.EventSource{Component: agentName})

	// Construct leader runner
	lr := newLeaderRunner(&leaseLeaderElector{
		namespace: a.AgentMeta.PodNamespace,
		name: func(ctx context.Context) (string, error) {
			id, err := a.AgentId.get(ctx) // nolint: govet
			if err != nil {
				return "", err
			}
			// We use agent id as part of lock name so that agentk Pods of different id don't compete with
			// each other. Only Pods with same agent id should compete for a lock. Put differently, agentk Pods
			// with same agent id have the same lock name but with different id have different lock name.
			return fmt.Sprintf("agent-%d-lock", id), nil
		},
		identity:           a.AgentMeta.PodName,
		coordinationClient: kubeClient.CoordinationV1(),
		eventRecorder:      eventRecorder,
	})

	// Construct agent modules
	modules, internalModules, err := a.constructModules(internalServer, kasConn, internalServerConn, k8sFactory, lr)
	if err != nil {
		return err
	}
	runner := a.newModuleRunner(kasConn)
	modulesRun := runner.RegisterModules(modules)
	internalModulesRun := runner.RegisterModules(internalModules)

	// Start events processing pipeline.
	loggingWatch := eventBroadcaster.StartStructuredLogging(0)
	defer loggingWatch.Stop()
	eventBroadcaster.StartRecordingToSink(&client_core_v1.EventSinkImpl{Interface: kubeClient.CoreV1().Events("")})
	defer eventBroadcaster.Shutdown()

	// Start things up. Stages are shut down in reverse order.
	return stager.RunStages(ctx,
		func(stage stager.Stage) {
			stage.Go(func(ctx context.Context) error {
				// Start leader runner.
				lr.Run(ctx)
				return nil
			})
		},
		func(stage stager.Stage) {
			// Start modules.
			stage.Go(modulesRun)
		},
		func(stage stager.Stage) {
			// Start internal gRPC server. It is used by internal modules, so it is shut down after them.
			a.startInternalServer(stage, internalServer, internalListener)
		},
		func(stage stager.Stage) {
			// Start modules that use internal server.
			stage.Go(internalModulesRun)
		},
		func(stage stager.Stage) {
			// Start configuration refresh.
			stage.Go(runner.RunConfigurationRefresh)
		},
	)
}

func (a *App) newModuleRunner(kasConn *grpc.ClientConn) *moduleRunner {
	return &moduleRunner{
		log: a.Log,
		configurationWatcher: &rpc.ConfigurationWatcher{
			Log:       a.Log,
			AgentMeta: a.AgentMeta,
			Client:    rpc.NewAgentConfigurationClient(kasConn),
			PollConfig: retry.NewPollConfigFactory(0, retry.NewExponentialBackoffFactory(
				getConfigurationInitBackoff,
				getConfigurationMaxBackoff,
				getConfigurationResetDuration,
				getConfigurationBackoffFactor,
				getConfigurationJitter,
			)),
			ConfigPreProcessor: func(data rpc.ConfigurationData) error {
				return a.AgentId.set(data.Config.AgentId)
			},
		},
	}
}

func (a *App) constructModules(internalServer *grpc.Server, kasConn, internalServerConn grpc.ClientConnInterface,
	k8sFactory util.Factory, lr *leaderRunner) ([]modagent.Module, []modagent.Module, error) {
	accessClient := gitlab_access_rpc.NewGitlabAccessClient(kasConn)
	factories := []modagent.Factory{
		&observability_agent.Factory{
			LogLevel:            a.LogLevel,
			GrpcLogLevel:        a.GrpcLogLevel,
			DefaultGrpcLogLevel: defaultGrpcLogLevel,
		},
		&manifestops.Factory{},
		&chartops.Factory{},
		&starboard_vulnerability.Factory{},
		&reverse_tunnel_agent.Factory{
			InternalServerConn: internalServerConn,
		},
		&kubernetes_api_agent.Factory{},
	}
	var modules []modagent.Module
	var internalModules []modagent.Module
	for _, f := range factories {
		moduleName := f.Name()
		module, err := f.New(&modagent.Config{
			Log:       a.Log.With(logz.ModuleName(moduleName)),
			AgentMeta: a.AgentMeta,
			Api: &agentAPI{
				moduleName: moduleName,
				agentId:    a.AgentId,
				client:     accessClient,
			},
			K8sUtilFactory:     k8sFactory,
			KasConn:            kasConn,
			Server:             internalServer,
			AgentName:          agentName,
			ServiceAccountName: a.ServiceAccountName,
		})
		if err != nil {
			return nil, nil, err
		}
		module = lr.MaybeWrapModule(module)
		if f.UsesInternalServer() {
			internalModules = append(internalModules, module)
		} else {
			modules = append(modules, module)
		}
	}
	return modules, internalModules, nil
}

func (a *App) constructKasConnection(ctx context.Context, tp trace.TracerProvider, p propagation.TextMapPropagator) (*grpc.ClientConn, error) {
	tokenData, err := os.ReadFile(a.TokenFile)
	if err != nil {
		return nil, fmt.Errorf("token file: %w", err)
	}
	tokenData = bytes.TrimSuffix(tokenData, []byte{'\n'})
	tlsConfig, err := tlstool.DefaultClientTLSConfigWithCACert(a.CACertFile)
	if err != nil {
		return nil, err
	}
	u, err := url.Parse(a.KasAddress)
	if err != nil {
		return nil, fmt.Errorf("invalid gitlab-kas address: %w", err)
	}
	userAgent := fmt.Sprintf("%s/%s/%s", agentName, a.AgentMeta.Version, a.AgentMeta.CommitId)
	opts := []grpc.DialOption{
		grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
		grpc.WithUserAgent(userAgent),
		// keepalive.ClientParameters must be specified at least as large as what is allowed by the
		// server-side grpc.KeepaliveEnforcementPolicy
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			// kas allows min 20 seconds, trying to stay below 60 seconds (typical load-balancer timeout) and
			// above kas' server keepalive Time so that kas pings the client sometimes. This helps mitigate
			// reverse-proxies' enforced server response timeout.
			Time:                55 * time.Second,
			PermitWithoutStream: true,
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
	var addressToDial string
	// "grpcs" is the only scheme where encryption is done by gRPC.
	// "wss" is secure too but gRPC cannot know that, so we tell it it's not.
	secure := u.Scheme == "grpcs"
	switch u.Scheme {
	case "ws", "wss":
		addressToDial = a.KasAddress
		dialer := net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}
		opts = append(opts, grpc.WithContextDialer(wstunnel.DialerForGRPC(defaultMaxMessageSize, &websocket.DialOptions{
			HTTPClient: &http.Client{
				Transport: &http.Transport{
					Proxy:                 http.ProxyFromEnvironment,
					DialContext:           dialer.DialContext,
					TLSClientConfig:       tlsConfig,
					MaxIdleConns:          10,
					IdleConnTimeout:       90 * time.Second,
					TLSHandshakeTimeout:   10 * time.Second,
					ResponseHeaderTimeout: 20 * time.Second,
					ExpectContinueTimeout: 20 * time.Second,
				},
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			},
			HTTPHeader: http.Header{
				httpz.UserAgentHeader: []string{userAgent},
			},
		})))
	case "grpc":
		addressToDial = grpcHostWithPort(u)
	case "grpcs":
		addressToDial = grpcHostWithPort(u)
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	default:
		return nil, fmt.Errorf("unsupported scheme in GitLab Kubernetes Agent Server address: %q", u.Scheme)
	}
	if !secure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	opts = append(opts, grpc.WithPerRPCCredentials(grpctool.NewTokenCredentials(api.AgentToken(tokenData), !secure)))
	conn, err := grpc.DialContext(ctx, addressToDial, opts...)
	if err != nil {
		return nil, fmt.Errorf("gRPC.dial: %w", err)
	}
	return conn, nil
}

func (a *App) constructInternalServer(auxCtx context.Context, tp trace.TracerProvider, p propagation.TextMapPropagator) *grpc.Server {
	factory := func(ctx context.Context, method string) modagent.RpcApi {
		return &agentRpcApi{
			RpcApiStub: modshared.RpcApiStub{
				Logger:    a.Log.With(logz.TraceIdFromContext(ctx)),
				StreamCtx: ctx,
			},
		}
	}
	return grpc.NewServer(
		grpc.StatsHandler(grpctool.NewServerMaxConnAgeStatsHandler(auxCtx, 0)),
		grpc.ChainStreamInterceptor(
			grpc_prometheus.StreamServerInterceptor,                                                        // 1. measure all invocations
			otelgrpc.StreamServerInterceptor(otelgrpc.WithTracerProvider(tp), otelgrpc.WithPropagators(p)), // 2. trace
			modagent.StreamRpcApiInterceptor(factory),                                                      // 3. inject RPC API
			grpc_validator.StreamServerInterceptor(),                                                       // x. wrap with validator
		),
		grpc.ChainUnaryInterceptor(
			grpc_prometheus.UnaryServerInterceptor,                                                        // 1. measure all invocations
			otelgrpc.UnaryServerInterceptor(otelgrpc.WithTracerProvider(tp), otelgrpc.WithPropagators(p)), // 2. trace
			modagent.UnaryRpcApiInterceptor(factory),                                                      // 3. inject RPC API
			grpc_validator.UnaryServerInterceptor(),                                                       // x. wrap with validator
		),
	)
}

func (a *App) startInternalServer(stage stager.Stage, internalServer *grpc.Server, internalListener net.Listener) {
	grpctool.StartServer(stage, internalServer, func() (net.Listener, error) {
		return internalListener, nil
	})
}

func (a *App) constructInternalServerConn(ctx context.Context, dialContext func(ctx context.Context, addr string) (net.Conn, error)) (*grpc.ClientConn, error) {
	return grpc.DialContext(ctx, "pipe",
		grpc.WithContextDialer(dialContext),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.ForceCodec(grpctool.RawCodec{})),
	)
}

func NewCommand() *cobra.Command {
	kubeConfigFlags := genericclioptions.NewConfigFlags(true)
	a := App{
		AgentMeta: &modshared.AgentMeta{
			Version:  cmd.Version,
			CommitId: cmd.Commit,
		},
		AgentId:            NewAgentIdHolder(),
		ServiceAccountName: os.Getenv(envVarServiceAccountName),
		K8sClientGetter:    kubeConfigFlags,
	}
	c := &cobra.Command{
		Use:   "agentk",
		Short: "GitLab Agent for Kubernetes",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) (retErr error) {
			podNs := os.Getenv(envVarPodNamespace)
			if podNs == "" {
				return fmt.Errorf("%s environment variable is required but is empty", envVarPodNamespace)
			}
			podName := os.Getenv(envVarPodName)
			if podName == "" {
				return fmt.Errorf("%s environment variable is required but is empty", envVarPodName)
			}
			a.AgentMeta.PodNamespace = podNs
			a.AgentMeta.PodName = podName
			lockedSyncer := zapcore.Lock(logz.NoSync(os.Stderr))
			var err error
			a.Log, a.LogLevel, err = a.logger(defaultLogLevel, lockedSyncer)
			if err != nil {
				return err
			}
			defer errz.SafeCall(a.Log.Sync, &retErr)

			var grpcLog *zap.Logger
			grpcLog, a.GrpcLogLevel, err = a.logger(defaultGrpcLogLevel, lockedSyncer)
			if err != nil {
				return err
			}
			defer errz.SafeCall(grpcLog.Sync, &retErr)

			grpclog.SetLoggerV2(zapgrpc.NewLogger(grpcLog)) // pipe gRPC logs into zap
			// Kubernetes uses klog so here we pipe all logs from it to our logger via an adapter.
			klog.SetLogger(zapr.NewLogger(a.Log))

			return a.Run(cmd.Context())
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	f := c.Flags()
	f.StringVar(&a.KasAddress, "kas-address", "", "GitLab Kubernetes Agent Server address")
	f.StringVar(&a.CACertFile, "ca-cert-file", "", "Optional file with X.509 certificate authority certificate in PEM format")
	f.StringVar(&a.TokenFile, "token-file", "", "File with access token")
	kubeConfigFlags.AddFlags(f)
	cobra.CheckErr(c.MarkFlagRequired("kas-address"))
	cobra.CheckErr(c.MarkFlagRequired("token-file"))
	return c
}

func grpcHostWithPort(u *url.URL) string {
	port := u.Port()
	if port != "" {
		return u.Host
	}
	switch u.Scheme {
	case "grpc":
		return net.JoinHostPort(u.Host, "80")
	case "grpcs":
		return net.JoinHostPort(u.Host, "443")
	default:
		// Function called with unknown scheme, just return the original host.
		return u.Host
	}
}
