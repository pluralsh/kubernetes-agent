package manifestops

import (
	"context"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/gitops/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/retry"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/pkg/agentcfg"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/cli-runtime/pkg/resource"
	"sigs.k8s.io/cli-utils/pkg/apply"
)

type defaultGitopsWorker struct {
	log               *zap.Logger
	agentId           int64
	project           *agentcfg.ManifestProjectCF
	applier           Applier
	restClientGetter  resource.RESTClientGetter
	applierPollConfig retry.PollConfig
	applyOptions      apply.ApplierOptions
	decodeRetryPolicy retry.BackoffManager
	objWatcher        rpc.ObjectsToSynchronizeWatcherInterface
}

func (w *defaultGitopsWorker) Run(ctx context.Context) {
	// Data flow: watch() -> decode() -> apply()
	desiredState := make(chan rpc.ObjectsToSynchronizeData)
	jobs := make(chan applyJob)

	var wg wait.Group
	defer wg.Wait()           // Wait for all pipeline stages to finish
	defer close(desiredState) // Close desiredState to signal decode() there is no more work to be done.
	wg.Start(func() {
		w.apply(jobs)
	})
	wg.Start(func() {
		defer close(jobs) // Close jobs to signal apply() there is no more work to be done.
		w.decode(desiredState, jobs)
	})
	w.watch(ctx, desiredState)
}
