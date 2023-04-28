package chartops

import (
	"context"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/errz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/retry"
	"go.uber.org/zap"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage/driver"
)

type job struct {
	ctx    context.Context
	log    *zap.Logger
	chart  *chart.Chart
	values ChartValues
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
		rel, err = w.helm.Upgrade(ctx, w.chartCfg.ReleaseName, job.chart, job.values, UpgradeConfig{
			Namespace:  *w.chartCfg.Namespace,
			MaxHistory: int(*w.chartCfg.MaxHistory),
		})
	} else {
		rel, err = w.helm.Install(ctx, job.chart, job.values, InstallConfig{
			Namespace:   *w.chartCfg.Namespace,
			ReleaseName: w.chartCfg.ReleaseName,
		})
	}
	_ = rel // TODO anything we want to do with the release?
	return err
}

func (w *worker) isInstalled() (bool, error) {
	_, err := w.helm.History(w.chartCfg.ReleaseName)
	switch err { // nolint: errorlint
	case driver.ErrReleaseNotFound:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, err
	}
}
