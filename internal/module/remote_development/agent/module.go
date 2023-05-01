package agent

import (
	"context"
	"time"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/remote_development"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/errz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/prototool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/agentcfg"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const (
	defaultFullSyncInterval    = 1 * time.Hour
	defaultPartialSyncInterval = 10 * time.Second
)

// remote dev module is expected to only run on the leader agentk replica
// as such module is expected to implement modagent.LeaderModule and
// the following has been added to ensure this compliance
var _ modagent.LeaderModule = (*module)(nil)

type remoteDevReconciler interface {
	Run(context.Context) error
	Stop()
}

type module struct {
	log               *zap.Logger
	api               modagent.Api
	reconcilerFactory func(ctx context.Context, cfg *agentcfg.RemoteCF) (remoteDevReconciler, error)
}

func (m *module) IsRunnableConfiguration(cfg *agentcfg.AgentConfiguration) bool {
	return cfg.RemoteDevelopment.Enabled
}

func (m *module) Run(ctx context.Context, cfg <-chan *agentcfg.AgentConfiguration) error {
	var activeTask stoppableTask
	defer func() {
		if activeTask != nil {
			activeTask.StopAndWait()
		}
	}()
	var latestConfig *agentcfg.RemoteCF

	for config := range cfg {
		// This loop reacts to configuration changes stopping and starting workers.

		// If the config has not changed don't do anything
		if proto.Equal(config.RemoteDevelopment, latestConfig) {
			continue
		}

		if activeTask != nil {
			activeTask.StopAndWait()
			activeTask = nil
			latestConfig = nil
		}

		latestConfig = config.RemoteDevelopment

		activeTask = newStoppableTask(ctx, func(moduleCtx context.Context) {
			m.log.Debug("Remote Development - starting reconciler run")

			w := &worker{
				log:               m.log,
				api:               m.api,
				reconcilerFactory: m.reconcilerFactory,
			}

			err := w.StartReconciliation(moduleCtx, latestConfig)
			if err != nil && !errz.ContextDone(err) {
				m.log.Error("Error running reconciler", logz.Error(err))
			}
			m.log.Debug("Remote Development - reconciler run ended")
		})
	}
	return nil
}

//goland:noinspection GoUnusedParameter
func (m *module) DefaultAndValidateConfiguration(config *agentcfg.AgentConfiguration) error {
	prototool.NotNil(&config.RemoteDevelopment)

	// config.RemoteDevelopment.Enabled will default to false if not provided which is expected
	prototool.Duration(&config.RemoteDevelopment.PartialSyncInterval, defaultPartialSyncInterval)
	prototool.Duration(&config.RemoteDevelopment.FullSyncInterval, defaultFullSyncInterval)

	return nil
}

func (m *module) Name() string {
	return remote_development.ModuleName
}
