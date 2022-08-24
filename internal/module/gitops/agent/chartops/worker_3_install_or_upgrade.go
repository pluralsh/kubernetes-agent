package chartops

import (
	"context"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/errz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/retry"
	"go.uber.org/zap"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage/driver"
)

type job struct {
	ctx   context.Context
	log   *zap.Logger
	chart *chart.Chart
	vals  map[string]interface{}
}

func (w *worker) installOrUpgrade(jobs <-chan job) {
	for job := range jobs {
		_ = retry.PollWithBackoff(job.ctx, w.installPollConfig, func(ctx context.Context) (error, retry.AttemptResult) {
			job.log.Info("Synchronizing chart")
			err := w.doJob(ctx, job)
			if err != nil {
				if errz.ContextDone(err) {
					job.log.Info("Synchronization was canceled", logz.Error(err))
				} else {
					job.log.Warn("Synchronization failed", logz.Error(err))
				}
				return nil, retry.Backoff
			}
			job.log.Info("Chart synchronized")
			return nil, retry.Continue
		})
	}
}

func (w *worker) doJob(ctx context.Context, job job) error {
	isInstalled, err := w.isInstalled()
	if err != nil {
		return err
	}
	var rel *release.Release
	if isInstalled {
		upgrade := action.NewUpgrade(w.actionCfg)
		upgrade.Namespace = *w.chartCfg.Namespace
		upgrade.MaxHistory = int(*w.chartCfg.MaxHistory)
		rel, err = upgrade.RunWithContext(ctx, w.chartCfg.ReleaseName, job.chart, job.vals)
	} else {
		install := action.NewInstall(w.actionCfg)
		install.Namespace = *w.chartCfg.Namespace
		install.ReleaseName = w.chartCfg.ReleaseName
		rel, err = install.RunWithContext(ctx, job.chart, job.vals)

	}
	_ = rel // TODO anything we want to do with the release?
	return err
}

func (w *worker) isInstalled() (bool, error) {
	histClient := action.NewHistory(w.actionCfg)
	_, err := histClient.Run(w.chartCfg.ReleaseName)
	switch err { // nolint: errorlint
	case driver.ErrReleaseNotFound:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, err
	}
}
