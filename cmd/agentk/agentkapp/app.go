package agentkapp

import (
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
	gitops_agent "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/gitops/agent"
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
	"gitlab.com/gitlab-org/labkit/correlation"
	grpccorrelation "gitlab.com/gitlab-org/labkit/correlation/grpc"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zapgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/keepalive"
	"k8s.io/cli-runtime/pkg/genericclioptions"
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
	// KasAddress specifies the address of kas.
	KasAddress      string
	CACertFile      string
	TokenFile       string
	K8sClientGetter genericclioptions.RESTClientGetter
}

func (a *App) Run(ctx context.Context) (retErr error) {
	// Construct gRPC connection to gitlab-kas
	kasConn, err := a.constructKasConnection(ctx)
	if err != nil {
		return err
	}
	defer errz.SafeClose(kasConn, &retErr)

	// Internal gRPC client->listener pipe
	internalListener := grpctool.NewDialListener()

	// Construct internal gRPC server
	internalServer := a.constructInternalServer(ctx)

	// Construct connection to internal gRPC server
	internalServerConn, err := a.constructInternalServerConn(ctx, internalListener.DialContext)
	if err != nil {
		return err
	}
	defer errz.SafeClose(internalServerConn, &retErr)

	// Construct agent modules
	modules, internalModules, err := a.constructModules(internalServer, kasConn, internalServerConn)
	if err != nil {
		return err
	}
	runner := a.newModuleRunner(kasConn)
	modulesRun := runner.RegisterModules(modules)
	internalModulesRun := runner.RegisterModules(internalModules)

	// Start things up. Stages are shut down in reverse order.
	return stager.RunStages(ctx,
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
		},
	}
}

func (a *App) constructModules(internalServer *grpc.Server, kasConn, internalServerConn grpc.ClientConnInterface) ([]modagent.Module, []modagent.Module, error) {
	k8sFactory := util.NewFactory(a.K8sClientGetter)
	fTracker := newFeatureTracker(a.Log)
	accessClient := gitlab_access_rpc.NewGitlabAccessClient(kasConn)
	factories := []modagent.Factory{
		&observability_agent.Factory{
			LogLevel:            a.LogLevel,
			GrpcLogLevel:        a.GrpcLogLevel,
			DefaultGrpcLogLevel: defaultGrpcLogLevel,
		},
		&gitops_agent.Factory{},
		&starboard_vulnerability.Factory{},
		&reverse_tunnel_agent.Factory{
			InternalServerConn: internalServerConn,
		},
		&kubernetes_api_agent.Factory{},
	}
	var modules []modagent.Module
	var internalModules []modagent.Module
	serviceAccountName := os.Getenv(envVarServiceAccountName)
	for _, f := range factories {
		moduleName := f.Name()
		module, err := f.New(&modagent.Config{
			Log:       a.Log.With(logz.ModuleName(moduleName)),
			AgentMeta: a.AgentMeta,
			Api: &agentAPI{
				moduleName:     moduleName,
				client:         accessClient,
				featureTracker: fTracker,
			},
			K8sUtilFactory:     k8sFactory,
			KasConn:            kasConn,
			Server:             internalServer,
			AgentName:          agentName,
			ServiceAccountName: serviceAccountName,
		})
		if err != nil {
			return nil, nil, err
		}
		if f.UsesInternalServer() {
			internalModules = append(internalModules, module)
		} else {
			modules = append(modules, module)
		}
	}
	return modules, internalModules, nil
}

func (a *App) constructKasConnection(ctx context.Context) (*grpc.ClientConn, error) {
	tokenData, err := os.ReadFile(a.TokenFile)
	if err != nil {
		return nil, fmt.Errorf("token file: %w", err)
	}
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
			grpccorrelation.StreamClientCorrelationInterceptor(grpccorrelation.WithClientName(agentName)),
			grpctool.StreamClientValidatingInterceptor,
		),
		grpc.WithChainUnaryInterceptor(
			grpc_prometheus.UnaryClientInterceptor,
			grpccorrelation.UnaryClientCorrelationInterceptor(grpccorrelation.WithClientName(agentName)),
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
		addressToDial = u.Host
	case "grpcs":
		addressToDial = u.Host
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

func (a *App) constructInternalServer(auxCtx context.Context) *grpc.Server {
	factory := func(ctx context.Context, method string) modagent.RpcApi {
		return &agentRpcApi{
			RpcApiStub: modshared.RpcApiStub{
				Logger:    a.Log.With(logz.CorrelationId(correlation.ExtractFromContext(ctx))),
				StreamCtx: ctx,
			},
		}
	}
	return grpc.NewServer(
		grpc.StatsHandler(grpctool.NewServerMaxConnAgeStatsHandler(auxCtx, 0)),
		grpc.ChainStreamInterceptor(
			grpc_prometheus.StreamServerInterceptor,              // 1. measure all invocations
			grpccorrelation.StreamServerCorrelationInterceptor(), // 2. add correlation id
			modagent.StreamRpcApiInterceptor(factory),            // 3. inject RPC API
			grpc_validator.StreamServerInterceptor(),             // x. wrap with validator
		),
		grpc.ChainUnaryInterceptor(
			grpc_prometheus.UnaryServerInterceptor,              // 1. measure all invocations
			grpccorrelation.UnaryServerCorrelationInterceptor(), // 2. add correlation id
			modagent.UnaryRpcApiInterceptor(factory),            // 3. inject RPC API
			grpc_validator.UnaryServerInterceptor(),             // x. wrap with validator
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
			Version:      cmd.Version,
			CommitId:     cmd.Commit,
			PodNamespace: os.Getenv(envVarPodNamespace),
			PodName:      os.Getenv(envVarPodName),
		},
		K8sClientGetter: kubeConfigFlags,
	}
	c := &cobra.Command{
		Use:   "agentk",
		Short: "GitLab Agent for Kubernetes",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) (retErr error) {
			lockedSyncer := zapcore.Lock(logz.NoSync(os.Stderr))
			var err error
			a.Log, a.LogLevel, err = logger(defaultLogLevel, lockedSyncer)
			if err != nil {
				return err
			}
			defer errz.SafeCall(a.Log.Sync, &retErr)

			var grpcLog *zap.Logger
			grpcLog, a.GrpcLogLevel, err = logger(defaultGrpcLogLevel, lockedSyncer)
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

func logger(levelEnum agentcfg.LogLevelEnum, sync zapcore.WriteSyncer) (*zap.Logger, zap.AtomicLevel, error) {
	level, err := logz.LevelFromString(levelEnum.String())
	if err != nil {
		return nil, zap.NewAtomicLevel(), err
	}
	atomicLevel := zap.NewAtomicLevelAt(level)
	return logz.LoggerWithLevel(atomicLevel, sync), atomicLevel, nil
}
