package agent

import (
	"context"
	"fmt"
	"time"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/remote_development"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/remote_development/agent/informer"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/remote_development/agent/k8s"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/retry"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/pkg/agentcfg"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
)

const (
	interval      = 10 * time.Second
	initBackoff   = 10 * time.Second
	maxBackoff    = time.Minute
	resetDuration = 2 * time.Minute
	backoffFactor = 2.0
	jitter        = 1.0

	agentIdLabelSelector = "remotedevelopment.gitlab/agent-id"
	resyncDuration       = 5 * time.Minute
)

var (
	deploymentGVR = schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "deployments",
	}
)

type Factory struct {
}

func (f *Factory) New(config *modagent.Config) (modagent.Module, error) {
	restConfig, err := config.K8sUtilFactory.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	client, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	k8sClient, err := k8s.New(config.Log, config.K8sUtilFactory)
	if err != nil {
		return nil, err
	}

	pollFactory := retry.NewPollConfigFactory(interval, retry.NewExponentialBackoffFactory(
		initBackoff,
		maxBackoff,
		resetDuration,
		backoffFactor,
		jitter,
	))

	return &module{
		log: config.Log,
		api: config.Api,
		workerFactory: func(ctx context.Context, cfg *agentcfg.RemoteCF) (remoteDevWorker, error) {
			agentId, err := config.Api.GetAgentId(ctx)
			if err != nil {
				return nil, err
			}

			factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(client, resyncDuration, corev1.NamespaceAll, func(opts *metav1.ListOptions) {
				opts.LabelSelector = fmt.Sprintf("%s=%d", agentIdLabelSelector, agentId)
			})
			inf, err := informer.NewK8sInformer(config.Log, factory.ForResource(deploymentGVR).Informer())
			if err != nil {
				return nil, err
			}
			return &worker{
				log:               config.Log,
				agentId:           agentId,
				api:               config.Api,
				pollConfig:        pollFactory,
				pollFunction:      retry.PollWithBackoff,
				stateTracker:      newPersistedStateTracker(),
				terminatedTracker: newPersistedTerminatedWorkspacesTracker(),
				informer:          inf,
				k8sClient:         k8sClient,
			}, nil
		},
	}, nil
}

func (f *Factory) Name() string {
	return remote_development.ModuleName
}

func (f *Factory) StartStopPhase() modshared.ModuleStartStopPhase {
	return modshared.ModuleStartBeforeServers
}