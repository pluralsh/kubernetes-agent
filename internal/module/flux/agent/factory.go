package agent

import (
	"context"
	"net/http"
	"time"

	notificationv1 "github.com/fluxcd/notification-controller/api/v1"
	sourcev1 "github.com/fluxcd/source-controller/api/v1"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/flux"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/flux/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/retry"
	apiextensionsv1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/transport"
)

const (
	// resyncDuration defines the duration for the shared informer cache resync interval.
	resyncDuration = 10 * time.Minute

	reconcileProjectsInitBackoff   = 10 * time.Second
	reconcileProjectsMaxBackoff    = 5 * time.Minute
	reconcileProjectsResetDuration = 10 * time.Minute
	reconcileProjectsBackoffFactor = 2.0
	reconcileProjectsJitter        = 1.0
)

type Factory struct {
}

func (f *Factory) IsProducingLeaderModules() bool {
	return true
}

func (f *Factory) New(config *modagent.Config) (modagent.Module, error) {
	restConfig, err := config.K8sUtilFactory.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	dynamicClient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	extApiClient, err := apiextensionsv1client.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	receiverClient := dynamicClient.Resource(notificationv1.GroupVersion.WithResource("receivers"))

	kubeApiUrl, _, err := defaultServerUrlFor(restConfig)
	if err != nil {
		return nil, err
	}
	transportCfg, err := restConfig.TransportConfig()
	if err != nil {
		return nil, err
	}
	kubeApiRoundTripper, err := transport.New(transportCfg)
	if err != nil {
		return nil, err
	}

	return &module{
		log:             config.Log,
		k8sExtApiClient: extApiClient,
		informersFactory: func() (informers.GenericInformer, informers.GenericInformer, cache.Indexer) {
			informerFactory := dynamicinformer.NewDynamicSharedInformerFactory(dynamicClient, resyncDuration)
			gitRepositoryInformer := informerFactory.ForResource(sourcev1.GroupVersion.WithResource("gitrepositories"))
			receiverInformer := informerFactory.ForResource(notificationv1.GroupVersion.WithResource("receivers"))
			receiverIndexer := receiverInformer.Informer().GetIndexer()
			return gitRepositoryInformer, receiverInformer, receiverIndexer
		},
		clientFactory: func(ctx context.Context, cfgUrl string, receiverIndexer cache.Indexer) (*client, error) {
			agentId, err := config.Api.GetAgentId(ctx)
			if err != nil {
				return nil, err
			}

			rt, err := newGitRepositoryReconcileTrigger(cfgUrl, kubeApiUrl, kubeApiRoundTripper, http.DefaultTransport)
			if err != nil {
				return nil, err
			}

			return newClient(
				config.Log,
				config.Api,
				agentId,
				rpc.NewGitLabFluxClient(config.KasConn),
				retry.NewPollConfigFactory(0, retry.NewExponentialBackoffFactory(
					reconcileProjectsInitBackoff, reconcileProjectsMaxBackoff, reconcileProjectsResetDuration, reconcileProjectsBackoffFactor, reconcileProjectsJitter),
				),
				receiverIndexer,
				rt,
			)
		},
		controllerFactory: func(ctx context.Context, gitRepositoryInformer informers.GenericInformer, receiverInformer informers.GenericInformer, projectReconciler projectReconciler) (controller, error) {
			agentId, err := config.Api.GetAgentId(ctx)
			if err != nil {
				return nil, err
			}
			gitLabExternalUrl, err := config.Api.GetGitLabExternalUrl(ctx)
			if err != nil {
				return nil, err
			}

			return newGitRepositoryController(ctx, config.Log, config.Api, agentId, gitLabExternalUrl, gitRepositoryInformer, receiverInformer, projectReconciler, receiverClient, clientset.CoreV1())
		},
	}, nil
}

func (f *Factory) Name() string {
	return flux.ModuleName
}

func (f *Factory) StartStopPhase() modshared.ModuleStartStopPhase {
	return modshared.ModuleStartBeforeServers
}
