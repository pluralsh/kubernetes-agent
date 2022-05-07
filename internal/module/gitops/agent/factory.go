package agent

import (
	"time"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/gitops"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/gitops/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/retry"
	"sigs.k8s.io/cli-utils/pkg/apply"
	"sigs.k8s.io/cli-utils/pkg/inventory"
)

const (
	getObjectsToSynchronizeInitBackoff   = 10 * time.Second
	getObjectsToSynchronizeMaxBackoff    = 5 * time.Minute
	getObjectsToSynchronizeResetDuration = 10 * time.Minute
	getObjectsToSynchronizeBackoffFactor = 2.0
	getObjectsToSynchronizeJitter        = 1.0

	defaultReapplyInterval = 5 * time.Minute
	applierInitBackoff     = 10 * time.Second
	applierMaxBackoff      = time.Minute
	applierResetDuration   = time.Minute
	applierBackoffFactor   = 2.0
	applierJitter          = 1.0
)

type Factory struct {
}

func (f *Factory) New(config *modagent.Config) (modagent.Module, error) {
	invClient, err := inventory.ClusterClientFactory{
		StatusPolicy: inventory.StatusPolicyNone,
	}.NewClient(config.K8sUtilFactory)
	if err != nil {
		return nil, err
	}
	applier, err := apply.NewApplierBuilder().
		WithFactory(config.K8sUtilFactory).
		WithInventoryClient(invClient).
		Build()
	if err != nil {
		return nil, err
	}
	restMapper, err := config.K8sUtilFactory.ToRESTMapper()
	if err != nil {
		return nil, err
	}
	return &module{
		log: config.Log,
		workerFactory: &defaultGitopsWorkerFactory{
			log:              config.Log,
			applier:          applier,
			restMapper:       restMapper,
			restClientGetter: config.K8sUtilFactory,
			gitopsClient:     rpc.NewGitopsClient(config.KasConn),
			watchPollConfig: retry.NewPollConfigFactory(0, retry.NewExponentialBackoffFactory(
				getObjectsToSynchronizeInitBackoff,
				getObjectsToSynchronizeMaxBackoff,
				getObjectsToSynchronizeResetDuration,
				getObjectsToSynchronizeBackoffFactor,
				getObjectsToSynchronizeJitter,
			)),
			applierPollConfig: retry.NewPollConfigFactory(defaultReapplyInterval, retry.NewExponentialBackoffFactory(
				applierInitBackoff,
				applierMaxBackoff,
				applierResetDuration,
				applierBackoffFactor,
				applierJitter,
			)),
		},
	}, nil
}

func (f *Factory) Name() string {
	return gitops.ModuleName
}
