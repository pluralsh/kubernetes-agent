package agent

import (
	"context"
	"fmt"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/flux"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/prototool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/agentcfg"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

type module struct {
	log               *zap.Logger
	informersFactory  func() (informers.GenericInformer, informers.GenericInformer, cache.Indexer)
	clientFactory     clientFactory
	controllerFactory controllerFactory
}

func (m *module) IsRunnableConfiguration(cfg *agentcfg.AgentConfiguration) bool {
	m.log.Debug("Not running module for now ...")
	return false
}

func (m *module) Run(ctx context.Context, cfg <-chan *agentcfg.AgentConfiguration) error {
	var wg wait.Group
	defer wg.Wait()

	var configCtx context.Context
	var cancel func()
	maybeCancel := func() {
		if cancel != nil {
			cancel()
		}
	}
	defer maybeCancel()

	for config := range cfg {
		// stop previous runs if any
		maybeCancel()
		wg.Wait()

		// create new context for module run with a specific configuration
		configCtx, cancel = context.WithCancel(ctx) // nolint:govet

		wg.Start(func() {
			if err := m.run(configCtx, config); err != nil {
				m.log.Error("failed to run module", logz.Error(err))
			}
		})
	}
	return nil // nolint:govet
}

func (m *module) run(ctx context.Context, cfg *agentcfg.AgentConfiguration) error {
	var wg wait.Group
	defer wg.Wait()

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	gitRepositoryInformer, receiverInformer, receiverIndexer := m.informersFactory()

	cl, err := m.clientFactory(runCtx, cfg.Flux.WebhookReceiverUrl, receiverIndexer)
	if err != nil {
		return fmt.Errorf("unable to create receiver: %w", err)
	}
	c, err := m.controllerFactory(runCtx, gitRepositoryInformer, receiverInformer, cl)
	if err != nil {
		return fmt.Errorf("unable to start controller: %w", err)
	}

	wg.StartWithChannel(runCtx.Done(), gitRepositoryInformer.Informer().Run)
	wg.StartWithChannel(runCtx.Done(), receiverInformer.Informer().Run)
	wg.StartWithContext(runCtx, cl.RunProjectReconciliation)

	c.Run(runCtx)
	return nil
}

func (m *module) DefaultAndValidateConfiguration(cfg *agentcfg.AgentConfiguration) error {
	prototool.NotNil(&cfg.Flux)
	prototool.String(&cfg.Flux.WebhookReceiverUrl, defaultServiceApiBaseUrl)
	return nil
}

func (m *module) Name() string {
	return flux.ModuleName
}
