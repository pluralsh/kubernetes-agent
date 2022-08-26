package agent

import (
	"context"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/gitops/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/retry"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/pkg/agentcfg"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/resource"
	"sigs.k8s.io/cli-utils/pkg/apply"
	"sigs.k8s.io/cli-utils/pkg/apply/event"
	"sigs.k8s.io/cli-utils/pkg/common"
	"sigs.k8s.io/cli-utils/pkg/inventory"
	"sigs.k8s.io/cli-utils/pkg/object"
	"sigs.k8s.io/cli-utils/pkg/object/validation"
)

const (
	dryRunStrategyNone   = "none"
	dryRunStrategyClient = "client"
	dryRunStrategyServer = "server"

	prunePropagationPolicyOrphan     = "orphan"
	prunePropagationPolicyBackground = "background"
	prunePropagationPolicyForeground = "foreground"

	inventoryPolicyMustMatch          = "must_match"
	inventoryPolicyAdoptIfNoInventory = "adopt_if_no_inventory"
	inventoryPolicyAdoptAll           = "adopt_all"
)

var (
	dryRunStrategyMapping = map[string]common.DryRunStrategy{
		dryRunStrategyNone:   common.DryRunNone,
		dryRunStrategyClient: common.DryRunClient,
		dryRunStrategyServer: common.DryRunServer,
	}
	prunePropagationPolicyMapping = map[string]metav1.DeletionPropagation{
		prunePropagationPolicyOrphan:     metav1.DeletePropagationOrphan,
		prunePropagationPolicyBackground: metav1.DeletePropagationBackground,
		prunePropagationPolicyForeground: metav1.DeletePropagationForeground,
	}
	inventoryPolicyMapping = map[string]inventory.Policy{
		inventoryPolicyMustMatch:          inventory.PolicyMustMatch,
		inventoryPolicyAdoptIfNoInventory: inventory.PolicyAdoptIfNoInventory,
		inventoryPolicyAdoptAll:           inventory.PolicyAdoptAll,
	}
)

type Applier interface {
	Run(ctx context.Context, invInfo inventory.Info, objects object.UnstructuredSet, options apply.ApplierOptions) <-chan event.Event
}

type GitopsWorkerFactory interface {
	New(int64, *agentcfg.ManifestProjectCF) GitopsWorker
}

type GitopsWorker interface {
	Run(context.Context)
}

type defaultGitopsWorkerFactory struct {
	log               *zap.Logger
	applier           Applier
	restClientGetter  resource.RESTClientGetter
	gitopsClient      rpc.GitopsClient
	watchPollConfig   retry.PollConfigFactory
	applierPollConfig retry.PollConfigFactory
	decodeRetryPolicy retry.BackoffManagerFactory
}

func (f *defaultGitopsWorkerFactory) New(agentId int64, project *agentcfg.ManifestProjectCF) GitopsWorker {
	l := f.log.With(logz.ProjectId(project.Id))
	return &defaultGitopsWorker{
		log:               l,
		agentId:           agentId,
		project:           project,
		applier:           f.applier,
		restClientGetter:  f.restClientGetter,
		applierPollConfig: f.applierPollConfig(),
		applyOptions: apply.ApplierOptions{
			ServerSideOptions: common.ServerSideOptions{
				// It's supported since Kubernetes 1.16, so there should be no reason not to use it.
				// https://kubernetes.io/docs/reference/using-api/server-side-apply/
				ServerSideApply: true,
				// GitOps repository is the source of truth and that's what we are applying, so overwrite any conflicts.
				// https://kubernetes.io/docs/reference/using-api/server-side-apply/#conflicts
				ForceConflicts: true,
				// https://kubernetes.io/docs/reference/using-api/server-side-apply/#field-management
				FieldManager: "agentk",
			},
			ReconcileTimeout:       project.ReconcileTimeout.AsDuration(),
			EmitStatusEvents:       true,
			NoPrune:                !project.GetPrune(),
			DryRunStrategy:         f.mapDryRunStrategy(project.DryRunStrategy),
			PrunePropagationPolicy: f.mapPrunePropagationPolicy(project.PrunePropagationPolicy),
			PruneTimeout:           project.PruneTimeout.AsDuration(),
			InventoryPolicy:        f.mapInventoryPolicy(project.InventoryPolicy),
			ValidationPolicy:       validation.ExitEarly,
		},
		decodeRetryPolicy: f.decodeRetryPolicy(),
		objWatcher: &rpc.ObjectsToSynchronizeWatcher{
			Log:          l,
			GitopsClient: f.gitopsClient,
			PollConfig:   f.watchPollConfig,
		},
	}
}

func (f *defaultGitopsWorkerFactory) mapDryRunStrategy(strategy string) common.DryRunStrategy {
	ret, ok := dryRunStrategyMapping[strategy]
	if !ok {
		// This shouldn't happen because we've checked the value in DefaultAndValidateConfiguration().
		// Just being extra cautious.
		f.log.Sugar().Errorf("Invalid dry-run strategy: %q, using client dry-run for safety - NO CHANGES WILL BE APPLIED!", strategy)
		ret = common.DryRunClient
	}
	return ret
}

func (f *defaultGitopsWorkerFactory) mapPrunePropagationPolicy(policy string) metav1.DeletionPropagation {
	ret, ok := prunePropagationPolicyMapping[policy]
	if !ok {
		// This shouldn't happen because we've checked the value in DefaultAndValidateConfiguration().
		// Just being extra cautious.
		f.log.Sugar().Errorf("Invalid prune propagation policy: %q, defaulting to %s", policy, metav1.DeletePropagationForeground)
		ret = metav1.DeletePropagationForeground
	}
	return ret
}

func (f *defaultGitopsWorkerFactory) mapInventoryPolicy(policy string) inventory.Policy {
	ret, ok := inventoryPolicyMapping[policy]
	if !ok {
		// This shouldn't happen because we've checked the value in DefaultAndValidateConfiguration().
		// Just being extra cautious.
		f.log.Sugar().Errorf("Invalid inventory policy: %q, defaulting to 'must match'", policy)
		ret = inventory.PolicyMustMatch
	}
	return ret
}
