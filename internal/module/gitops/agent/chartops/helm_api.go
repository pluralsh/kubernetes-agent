package chartops

import (
	"context"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
)

type ChartValues map[string]interface{}

type Helm interface {
	History(name string) ([]*release.Release, error)
	Upgrade(ctx context.Context, name string, chart *chart.Chart, vals ChartValues, cfg UpgradeConfig) (*release.Release, error)
	Install(ctx context.Context, chart *chart.Chart, vals ChartValues, cfg InstallConfig) (*release.Release, error)
}

type UpgradeConfig struct {
	// Namespace is the namespace in which this operation should be performed.
	Namespace string
	// MaxHistory limits the maximum number of revisions saved per release
	MaxHistory int
}

type InstallConfig struct {
	// Namespace is the namespace in which this operation should be performed.
	Namespace   string
	ReleaseName string
}

type HelmActions struct {
	ActionCfg *action.Configuration
}

func (h *HelmActions) History(name string) ([]*release.Release, error) {
	return action.NewHistory(h.ActionCfg).Run(name)
}

func (h *HelmActions) Upgrade(ctx context.Context, name string, chart *chart.Chart, vals ChartValues, cfg UpgradeConfig) (*release.Release, error) {
	upgrade := action.NewUpgrade(h.ActionCfg)
	upgrade.Namespace = cfg.Namespace
	upgrade.MaxHistory = cfg.MaxHistory
	return upgrade.RunWithContext(ctx, name, chart, vals)
}

func (h *HelmActions) Install(ctx context.Context, chart *chart.Chart, vals ChartValues, cfg InstallConfig) (*release.Release, error) {
	install := action.NewInstall(h.ActionCfg)
	install.Namespace = cfg.Namespace
	install.ReleaseName = cfg.ReleaseName
	return install.RunWithContext(ctx, chart, vals)
}
