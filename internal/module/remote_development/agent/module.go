package agent

import (
	"context"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/remote_development"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/errz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/pkg/agentcfg"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type remoteDevWorker interface {
	Run(context.Context) error
}

type module struct {
	log           *zap.Logger
	api           modagent.Api
	workerFactory func(ctx context.Context, cfg *agentcfg.RemoteCF) (remoteDevWorker, error)
}

func (m *module) Run(ctx context.Context, cfg <-chan *agentcfg.AgentConfiguration) error {
	var activeTask stoppableTask
	defer func() {
		if activeTask != nil {
			activeTask.StopAndWait()
		}
	}()
	var existingConfig *agentcfg.RemoteCF

	for config := range cfg {
		// This loop reacts to configuration changes stopping and starting workers.

		// If the config has not changed don't do anything
		if proto.Equal(config.RemoteDevelopment, existingConfig) {
			continue
		}

		if activeTask != nil {
			activeTask.StopAndWait()
			activeTask = nil
			existingConfig = nil
		}

		if config.RemoteDevelopment != nil && config.RemoteDevelopment.Enabled {
			m.log.Debug("Remote Development is enabled")

			// Set new configuration
			w, err := m.workerFactory(ctx, config.RemoteDevelopment)
			if err != nil {
				m.api.HandleProcessingError(ctx, m.log, modshared.NoAgentId, "Error starting worker", err)
				continue
			}
			existingConfig = config.RemoteDevelopment

			activeTask = newStoppableTask(ctx, func(workerCtx context.Context) {
				err := w.Run(workerCtx)
				if err != nil && !errz.ContextDone(err) {
					m.log.Error("Error running worker", logz.Error(err))
				}
			})
		} else {
			m.log.Debug("Remote Development is disabled")
		}
	}
	return nil
}

//goland:noinspection GoUnusedParameter
func (m *module) DefaultAndValidateConfiguration(config *agentcfg.AgentConfiguration) error {
	return nil
}

func (m *module) Name() string {
	return remote_development.ModuleName
}
