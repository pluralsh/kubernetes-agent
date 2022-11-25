package chartops

import (
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/gitops/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/retry"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/pkg/agentcfg"
	"go.uber.org/zap"
	"helm.sh/helm/v3/pkg/action"
)

type workerFactory struct {
	log               *zap.Logger
	actionCfg         func(log *zap.Logger, chartCfg *agentcfg.ChartCF) *action.Configuration
	gitopsClient      rpc.GitopsClient
	installPollConfig retry.PollConfigFactory
	watchPollConfig   retry.PollConfigFactory
}

func (f *workerFactory) New(agentId int64, source modagent.WorkSource[*agentcfg.ChartCF]) modagent.Worker {
	chartCfg := source.Configuration()
	l := f.log.With(logz.WorkerId(source.ID()), logz.AgentId(agentId))
	return &worker{
		log:               l,
		chartCfg:          chartCfg,
		installPollConfig: f.installPollConfig(),
		actionCfg:         f.actionCfg(l, chartCfg),
		objWatcher: &rpc.ObjectsToSynchronizeWatcher{
			Log:          l,
			GitopsClient: f.gitopsClient,
			PollConfig:   f.watchPollConfig,
		},
	}
}

func (f *workerFactory) SourcesFromConfiguration(cfg *agentcfg.AgentConfiguration) []modagent.WorkSource[*agentcfg.ChartCF] {
	res := make([]modagent.WorkSource[*agentcfg.ChartCF], 0, len(cfg.Gitops.Charts))
	for _, chart := range cfg.Gitops.Charts {
		res = append(res, (*manifestSource)(chart))
	}
	return res
}

type manifestSource agentcfg.ChartCF

func (s *manifestSource) ID() string {
	return *s.Namespace + "/" + s.ReleaseName
}

func (s *manifestSource) Configuration() *agentcfg.ChartCF {
	return (*agentcfg.ChartCF)(s)
}
