package chartops

import (
	"context"
	"net/http"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/gitops/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/retry"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/pkg/agentcfg"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/util/wait"
)

type worker struct {
	log               *zap.Logger
	chartCfg          *agentcfg.ChartCF
	installPollConfig retry.PollConfig
	helm              Helm
	httpClient        http.RoundTripper
	objWatcher        rpc.ObjectsToSynchronizeWatcherInterface
}

func (w *worker) Run(ctx context.Context) {
	// Data flow: fetch_and_load() -> installOrUpgrade()
	jobs := make(chan job)

	var wg wait.Group
	defer wg.Wait()   // Wait for all pipeline stages to finish
	defer close(jobs) // Close jobs to signal installOrUpgrade() there is no more work to be done.
	wg.Start(func() {
		w.installOrUpgrade(jobs)
	})
	w.fetch(ctx, jobs)
}
