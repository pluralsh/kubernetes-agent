package chartops

import (
	"context"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/gitops/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/retry"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/pkg/agentcfg"
	"go.uber.org/zap"
	"helm.sh/helm/v3/pkg/action"
	"k8s.io/apimachinery/pkg/util/wait"
)

type worker struct {
	log               *zap.Logger
	chartCfg          *agentcfg.ChartCF
	installPollConfig retry.PollConfig
	actionCfg         *action.Configuration
	objWatcher        rpc.ObjectsToSynchronizeWatcherInterface
}

func (w *worker) Run(ctx context.Context) {
	// Data flow: fetch() -> load() -> installOrUpgrade()
	desiredState := make(chan fetchedData)
	jobs := make(chan job)

	var wg wait.Group
	defer wg.Wait()           // Wait for all pipeline stages to finish
	defer close(desiredState) // Close desiredState to signal load() there is no more work to be done.
	wg.Start(func() {
		w.installOrUpgrade(jobs)
	})
	wg.Start(func() {
		defer close(jobs) // Close jobs to signal installOrUpgrade() there is no more work to be done.
		w.decode(desiredState, jobs)
	})
	w.fetch(ctx, desiredState)
}
