package chartops

import (
	"net"
	"net/http"
	"time"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/gitops"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/gitops/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/retry"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/tlstool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/agentcfg"
	"go.uber.org/zap"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/kube"
	"helm.sh/helm/v3/pkg/registry"
	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"
)

const (
	getObjectsToSynchronizeInitBackoff   = 10 * time.Second
	getObjectsToSynchronizeMaxBackoff    = 5 * time.Minute
	getObjectsToSynchronizeResetDuration = 10 * time.Minute
	getObjectsToSynchronizeBackoffFactor = 2.0
	getObjectsToSynchronizeJitter        = 1.0

	defaultReinstallInterval = 5 * time.Minute
	installInitBackoff       = 10 * time.Second
	installMaxBackoff        = time.Minute
	installResetDuration     = 2 * time.Minute
	installBackoffFactor     = 2.0
	installJitter            = 1.0
)

type Factory struct {
}

func (f *Factory) IsProducingLeaderModules() bool {
	return true
}

func (f *Factory) New(config *modagent.Config) (modagent.Module, error) {
	clientset, err := config.K8sUtilFactory.KubernetesClientSet()
	if err != nil {
		return nil, err
	}
	// TODO support debug, credentials, output writer
	registryClient, err := registry.NewClient(
		registry.ClientOptEnableCache(true),
	)
	if err != nil {
		return nil, err
	}
	coreV1client := clientset.CoreV1()
	return &module{
		log: config.Log,
		workerFactory: &workerFactory{
			log: config.Log,
			helm: func(log *zap.Logger, chartCfg *agentcfg.ChartCF) Helm {
				infof := log.Sugar().Infof
				d := driver.NewSecrets(coreV1client.Secrets(*chartCfg.Namespace))
				d.Log = infof
				return &HelmActions{
					ActionCfg: &action.Configuration{
						RESTClientGetter: config.K8sUtilFactory,
						Releases: &storage.Storage{
							Driver:     d,
							MaxHistory: int(*chartCfg.MaxHistory),
							Log:        infof,
						},
						KubeClient: &kube.Client{
							Factory:   config.K8sUtilFactory,
							Log:       infof,
							Namespace: *chartCfg.Namespace,
						},
						RegistryClient: registryClient,
						Capabilities:   nil, // Empty to re-discover supported APIs.
						Log:            infof,
					},
				}
			},
			httpClient: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				TLSClientConfig:       tlstool.DefaultClientTLSConfig(),
				MaxIdleConns:          10,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ResponseHeaderTimeout: 20 * time.Second,
				ForceAttemptHTTP2:     true,
			},
			gitopsClient: rpc.NewGitopsClient(config.KasConn),
			installPollConfig: retry.NewPollConfigFactory(defaultReinstallInterval, retry.NewExponentialBackoffFactory(
				installInitBackoff,
				installMaxBackoff,
				installResetDuration,
				installBackoffFactor,
				installJitter,
			)),
			watchPollConfig: retry.NewPollConfigFactory(0, retry.NewExponentialBackoffFactory(
				getObjectsToSynchronizeInitBackoff,
				getObjectsToSynchronizeMaxBackoff,
				getObjectsToSynchronizeResetDuration,
				getObjectsToSynchronizeBackoffFactor,
				getObjectsToSynchronizeJitter,
			)),
		},
	}, nil
}

func (f *Factory) Name() string {
	return gitops.AgentChartModuleName
}

func (f *Factory) StartStopPhase() modshared.ModuleStartStopPhase {
	return modshared.ModuleStartBeforeServers
}
