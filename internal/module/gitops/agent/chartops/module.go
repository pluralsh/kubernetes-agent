package chartops

import (
	"context"
	"fmt"
	"time"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/gitops"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/prototool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/agentcfg"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	defaultChartNamespace      = metav1.NamespaceDefault
	defaultChartMaxHistory     = 1 // no history for now as it's not very useful.
	defaultUrlValueMaxFileSize = 1024 * 1024
	defaultUrlValuePollPeriod  = time.Minute
)

type module struct {
	log           *zap.Logger
	workerFactory *workerFactory
}

func (m *module) IsRunnableConfiguration(cfg *agentcfg.AgentConfiguration) bool {
	return cfg.Gitops != nil && len(cfg.Gitops.Charts) > 0
}

func (m *module) Run(ctx context.Context, cfg <-chan *agentcfg.AgentConfiguration) error {
	wm := modagent.NewWorkerManager[*agentcfg.ChartCF](m.log, m.workerFactory)
	defer wm.StopAllWorkers()
	for config := range cfg {
		err := wm.ApplyConfiguration(config.AgentId, config) // nolint: contextcheck
		if err != nil {
			m.log.Error("Failed to apply chart synchronization configuration", logz.Error(err))
			continue
		}
	}
	return nil
}

func (m *module) DefaultAndValidateConfiguration(config *agentcfg.AgentConfiguration) error {
	prototool.NotNil(&config.Gitops)
	for _, chart := range config.Gitops.Charts {
		prototool.StringPtr(&chart.Namespace, defaultChartNamespace)
		prototool.Int32Ptr(&chart.MaxHistory, defaultChartMaxHistory)
		proj := chart.Source.GetProject() // may be nil
		for _, val := range chart.Values {
			switch src := val.From.(type) {
			case *agentcfg.ChartValuesCF_File:
				file := src.File
				if file.ProjectId == nil { // values from file without project ID
					if proj == nil { // we are not fetching chart from a project
						return fmt.Errorf("from_file %q values must have project id specified when not fetching chart from a project", file.File)
					}
					file.ProjectId = &proj.Id
				}
				if file.Ref == nil && proj != nil {
					// ref was not specified, but we are fetching from a project, so use its ref. It may be nil.
					file.Ref = proj.Ref
				}
			case *agentcfg.ChartValuesCF_Url:
				prototool.Duration(&src.Url.PollPeriod, defaultUrlValuePollPeriod)
				prototool.Uint32Ptr(&src.Url.MaxFileSize, defaultUrlValueMaxFileSize)
			}
		}
	}
	return nil
}

func (m *module) Name() string {
	return gitops.AgentChartModuleName
}
