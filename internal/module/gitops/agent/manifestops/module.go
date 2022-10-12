package manifestops

import (
	"context"
	"fmt"
	"time"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/gitops"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/gitops/agent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/prototool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/pkg/agentcfg"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	defaultGitOpsManifestNamespace = metav1.NamespaceDefault
	defaultGitOpsManifestPathGlob  = "**/*.{yaml,yml,json}"
	defaultDryRunStrategy          = dryRunStrategyNone
	defaultPrune                   = true
	defaultPruneTimeout            = time.Hour
	defaultReconcileTimeout        = time.Hour
	defaultPrunePropagationPolicy  = prunePropagationPolicyForeground
	defaultInventoryPolicy         = inventoryPolicyMustMatch
)

type module struct {
	log           *zap.Logger
	workerFactory agent.WorkerFactory
}

func (m *module) IsRunnableConfiguration(cfg *agentcfg.AgentConfiguration) bool {
	return cfg.Gitops != nil && len(cfg.Gitops.ManifestProjects) > 0
}

func (m *module) Run(ctx context.Context, cfg <-chan *agentcfg.AgentConfiguration) error {
	wm := agent.NewWorkerManager(m.log, m.workerFactory)
	defer wm.StopAllWorkers()
	for config := range cfg {
		err := wm.ApplyConfiguration(config.AgentId, config) // nolint: contextcheck
		if err != nil {
			m.log.Error("Failed to apply manifest projects configuration", logz.Error(err))
			continue
		}
	}
	return nil
}

func (m *module) DefaultAndValidateConfiguration(config *agentcfg.AgentConfiguration) error {
	prototool.NotNil(&config.Gitops)
	for _, project := range config.Gitops.ManifestProjects {
		err := applyDefaultsToManifestProject(project)
		if err != nil {
			return fmt.Errorf("project %s: %w", project.Id, err)
		}
	}
	return nil
}

func applyDefaultsToManifestProject(project *agentcfg.ManifestProjectCF) error {
	prototool.String(&project.DefaultNamespace, defaultGitOpsManifestNamespace)
	if len(project.Paths) == 0 {
		project.Paths = []*agentcfg.PathCF{
			{
				Glob: defaultGitOpsManifestPathGlob,
			},
		}
	}
	prototool.Duration(&project.ReconcileTimeout, defaultReconcileTimeout)
	prototool.String(&project.DryRunStrategy, defaultDryRunStrategy)
	if _, ok := dryRunStrategyMapping[project.DryRunStrategy]; !ok {
		return fmt.Errorf("invalid dry-run strategy: %q", project.DryRunStrategy)
	}
	if project.PruneOneof == nil {
		project.PruneOneof = &agentcfg.ManifestProjectCF_Prune{Prune: defaultPrune}
	}
	prototool.Duration(&project.PruneTimeout, defaultPruneTimeout)
	prototool.String(&project.PrunePropagationPolicy, defaultPrunePropagationPolicy)
	if _, ok := prunePropagationPolicyMapping[project.PrunePropagationPolicy]; !ok {
		return fmt.Errorf("invalid prune propagation policy: %q", project.PrunePropagationPolicy)
	}
	prototool.String(&project.InventoryPolicy, defaultInventoryPolicy)
	if _, ok := inventoryPolicyMapping[project.InventoryPolicy]; !ok {
		return fmt.Errorf("invalid inventory policy: %q", project.InventoryPolicy)
	}
	return nil
}

func (m *module) Name() string {
	return gitops.AgentManifestModuleName
}