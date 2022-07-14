package agent

import (
	"context"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/gitops/rpc"
	"k8s.io/apimachinery/pkg/util/wait"
)

type defaultGitopsWorker struct {
	objWatcher rpc.ObjectsToSynchronizeWatcherInterface
	synchronizerConfig
}

func (w *defaultGitopsWorker) Run(ctx context.Context) {
	var wg wait.Group
	defer wg.Wait()
	desiredState := make(chan rpc.ObjectsToSynchronizeData)
	defer close(desiredState)
	wg.Start(func() {
		s := newSynchronizer(w.synchronizerConfig)
		s.run(desiredState)
	})
	req := &rpc.ObjectsToSynchronizeRequest{
		ProjectId: w.project.Id,
		Paths:     w.project.Paths,
	}
	w.objWatcher.Watch(ctx, req, func(ctx context.Context, data rpc.ObjectsToSynchronizeData) {
		select {
		case <-ctx.Done():
		case desiredState <- data:
		}
	})
}
