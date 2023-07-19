package agent

import (
	"context"
	"time"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/remote_development"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/errz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/prototool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/syncz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/agentcfg"
	"go.uber.org/zap"
)

const (
	defaultFullSyncInterval    = 1 * time.Hour
	defaultPartialSyncInterval = 10 * time.Second
)

type remoteDevReconciler interface {
	Run(context.Context) error
	Stop()
}

type module struct {
	log               *zap.Logger
	api               modagent.Api
	reconcilerFactory func(ctx context.Context) (remoteDevReconciler, error)
}

func (m *module) Run(ctx context.Context, cfg <-chan *agentcfg.AgentConfiguration) error {
	wh := syncz.NewProtoWorkerHolder[*agentcfg.RemoteDevelopmentCF](
		func(config *agentcfg.RemoteDevelopmentCF) syncz.Worker {
			if !config.Enabled {
				return syncz.WorkerFunc(func(ctx context.Context) {
					// nop worker
				})
			}

			return syncz.WorkerFunc(func(ctx context.Context) {
				m.log.Debug("Remote Development - starting reconciler run")
				defer m.log.Debug("Remote Development - reconciler run ended")

				w := &worker{
					log:                 m.log,
					api:                 m.api,
					fullSyncInterval:    config.GetFullSyncInterval().AsDuration(),
					partialSyncInterval: config.GetPartialSyncInterval().AsDuration(),
					reconcilerFactory:   m.reconcilerFactory,
				}

				err := w.Run(ctx)
				if err != nil && !errz.ContextDone(err) {
					m.log.Error("Error running reconciler", logz.Error(err))
				}
			})
		},
	)
	defer wh.StopAndWait()

	// This loop reacts to configuration changes stopping and starting workers.
	for config := range cfg {
		wh.ApplyConfig(ctx, config.RemoteDevelopment)
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
