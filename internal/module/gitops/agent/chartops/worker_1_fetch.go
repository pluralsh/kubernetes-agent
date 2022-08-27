package chartops

import (
	"context"
	"fmt"
	"strings"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/gitops/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/pkg/agentcfg"
	"go.uber.org/zap"
	"helm.sh/helm/v3/pkg/chart/loader"
)

type fetchedData struct {
	log   *zap.Logger
	files []*loader.BufferedFile
}

func (w *worker) fetch(ctx context.Context, desiredState chan<- fetchedData) {
	switch src := w.chartCfg.Source.Source.(type) {
	case *agentcfg.ChartSourceCF_Project:
		w.fetchFromGitLabRepo(ctx, src.Project, desiredState)
	default:
		// Should never happen.
		panic(fmt.Errorf("unknown source type: %T", w.chartCfg.Source.Source))
	}
}

func (w *worker) fetchFromGitLabRepo(ctx context.Context, src *agentcfg.ChartProjectSourceCF, desiredState chan<- fetchedData) {
	pathWithSlash := strings.TrimSuffix(src.Path, "/") + "/"
	req := &rpc.ObjectsToSynchronizeRequest{
		ProjectId: src.Id,
		Paths:     []*agentcfg.PathCF{{Glob: pathWithSlash + "**"}}, // TODO validate it's a path, not a glob?
	}
	w.objWatcher.Watch(ctx, req, func(ctx context.Context, data rpc.ObjectsToSynchronizeData) {
		select {
		case <-ctx.Done():
		case desiredState <- w.data2fetchedData(data, pathWithSlash):
		}
	})
}

func (w *worker) data2fetchedData(data rpc.ObjectsToSynchronizeData, pathWithSlash string) fetchedData {
	files := make([]*loader.BufferedFile, 0, len(data.Sources))
	for _, source := range data.Sources {
		files = append(files, &loader.BufferedFile{
			Name: strings.TrimPrefix(source.Name, pathWithSlash),
			Data: source.Data,
		})
	}
	return fetchedData{
		log:   w.log.With(logz.CommitId(data.CommitId)),
		files: files,
	}
}
