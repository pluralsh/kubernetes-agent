package manifestops

import (
	"context"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/gitops/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/agentcfg"
)

func (w *worker) watch(ctx context.Context, desiredState chan<- rpc.ObjectsToSynchronizeData) {
	req := &rpc.ObjectsToSynchronizeRequest{
		ProjectId: *w.project.Id,
		Ref:       rpc.NewRpcRef(w.project.Ref),
		Paths:     configPathsToRpcPaths(w.project.Paths),
	}
	w.objWatcher.Watch(ctx, req, func(ctx context.Context, data rpc.ObjectsToSynchronizeData) {
		select {
		case <-ctx.Done():
		case desiredState <- data:
		}
	})
}

func configPathsToRpcPaths(paths []*agentcfg.PathCF) []*rpc.PathCF {
	p := make([]*rpc.PathCF, 0, len(paths))
	for _, path := range paths {
		p = append(p, &rpc.PathCF{Path: &rpc.PathCF_Glob{Glob: path.Glob}})
	}
	return p
}
