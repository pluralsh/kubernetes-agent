package agentkapp

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"time"

	"github.com/go-logr/zapr"
	"github.com/spf13/pflag"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/cmd"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/agentk"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/agentrpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/api/apiutil"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/tools/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/tools/wstunnel"
	grpccorrelation "gitlab.com/gitlab-org/labkit/correlation/grpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/klog/v2"
	"nhooyr.io/websocket"
)

const (
	defaultRefreshConfigurationRetryPeriod    = 10 * time.Second
	defaultGetObjectsToSynchronizeRetryPeriod = 10 * time.Second

	defaultMaxMessageSize = 10 * 1024 * 1024
	correlationClientName = "gitlab-agent"
)

type App struct {
	Log      *zap.Logger
	LogLevel zap.AtomicLevel
	// KasAddress specifies the address of kas.
	KasAddress      string
	TokenFile       string
	K8sClientGetter resource.RESTClientGetter
}

func (a *App) Run(ctx context.Context) error {
	defer a.Log.Sync() // nolint: errcheck
	// Kubernetes uses klog so here we pipe all logs from it to our logger via an adapter.
	klog.SetLogger(zapr.NewLogger(a.Log))
	restConfig, err := a.K8sClientGetter.ToRESTConfig()
	if err != nil {
		return fmt.Errorf("ToRESTConfig: %v", err)
	}
	tokenData, err := ioutil.ReadFile(a.TokenFile)
	if err != nil {
		return fmt.Errorf("token file: %v", err)
	}
	conn, err := a.kasConnection(ctx, string(tokenData))
	if err != nil {
		return err
	}
	defer conn.Close() // nolint: errcheck
	agent := agentk.New(agentk.Config{
		Log:       a.Log,
		LogLevel:  a.LogLevel,
		KasClient: agentrpc.NewKasClient(conn),
		EngineFactory: &agentk.DefaultGitOpsEngineFactory{
			KubeClientConfig: restConfig,
		},
		K8sClientGetter:                    a.K8sClientGetter,
		RefreshConfigurationRetryPeriod:    defaultRefreshConfigurationRetryPeriod,
		GetObjectsToSynchronizeRetryPeriod: defaultGetObjectsToSynchronizeRetryPeriod,
	})
	return agent.Run(ctx)
}

func (a *App) kasConnection(ctx context.Context, token string) (*grpc.ClientConn, error) {
	u, err := url.Parse(a.KasAddress)
	if err != nil {
		return nil, fmt.Errorf("invalid kas address: %v", err)
	}
	opts := []grpc.DialOption{
		grpc.WithUserAgent(fmt.Sprintf("agentk/%s/%s", cmd.Version, cmd.Commit)),
		// keepalive.ClientParameters must be specified at least as large as what is allowed by the
		// server-side grpc.KeepaliveEnforcementPolicy
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                20 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.WithChainStreamInterceptor(
			grpccorrelation.StreamClientCorrelationInterceptor(grpccorrelation.WithClientName(correlationClientName)),
		),
		grpc.WithChainUnaryInterceptor(
			grpccorrelation.UnaryClientCorrelationInterceptor(grpccorrelation.WithClientName(correlationClientName)),
		),
	}
	var addressToDial string
	// "grpcs" is the only scheme where encryption is done by gRPC.
	// "wss" is secure too but gRPC cannot know that, so we tell it it's not.
	secure := u.Scheme == "grpcs"
	switch u.Scheme {
	case "ws", "wss":
		addressToDial = a.KasAddress
		opts = append(opts, grpc.WithContextDialer(wstunnel.DialerForGRPC(defaultMaxMessageSize, &websocket.DialOptions{
			// TODO
		})))
	case "grpc":
		addressToDial = u.Host
	//case "grpcs":
	// TODO https://gitlab.com/gitlab-org/cluster-integration/gitlab-agent/-/issues/7
	default:
		return nil, fmt.Errorf("unsupported scheme in GitLab Kubernetes Agent Server address: %q", u.Scheme)
	}
	if !secure {
		opts = append(opts, grpc.WithInsecure())
	}
	opts = append(opts, grpc.WithPerRPCCredentials(apiutil.NewTokenCredentials(token, !secure)))
	conn, err := grpc.DialContext(ctx, addressToDial, opts...)
	if err != nil {
		return nil, fmt.Errorf("gRPC.dial: %v", err)
	}
	return conn, nil
}

func NewFromFlags(flagset *pflag.FlagSet, arguments []string) (cmd.Runnable, error) {
	log, level := logger()
	app := &App{
		Log:      log,
		LogLevel: level,
	}
	flagset.StringVar(&app.KasAddress, "kas-address", "", "GitLab Kubernetes Agent Server address")
	flagset.StringVar(&app.TokenFile, "token-file", "", "File with access token")
	kubeConfigFlags := genericclioptions.NewConfigFlags(true)
	kubeConfigFlags.AddFlags(flagset)
	if err := flagset.Parse(arguments); err != nil {
		return nil, err
	}
	app.K8sClientGetter = kubeConfigFlags
	return app, nil
}

func logger() (*zap.Logger, zap.AtomicLevel) {
	level := zap.NewAtomicLevelAt(agentk.DefaultLogLevel)
	return logz.LoggerWithLevel(level), level
}
