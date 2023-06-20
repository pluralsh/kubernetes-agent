package agent

import (
	"context"
	"fmt"

	notificationv1 "github.com/fluxcd/notification-controller/api/v1"
	sourcev1 "github.com/fluxcd/source-controller/api/v1"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/flux"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/prototool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/syncz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/agentcfg"
	"go.uber.org/zap"
	apiextensionsv1api "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

var (
	requiredFluxCrds = [...]schema.GroupResource{
		sourcev1.GroupVersion.WithResource("gitrepositories").GroupResource(),
		notificationv1.GroupVersion.WithResource("receivers").GroupResource(),
	}
)

type module struct {
	log               *zap.Logger
	k8sExtApiClient   apiextensionsv1client.ApiextensionsV1Interface
	informersFactory  func() (informers.GenericInformer, informers.GenericInformer, cache.Indexer)
	clientFactory     clientFactory
	controllerFactory controllerFactory
}

func (m *module) IsRunnableConfiguration(cfg *agentcfg.AgentConfiguration) bool {
	// NOTE: always running Flux module for now, but check in `Run()` if Flux is installed
	return true
}

func (m *module) Run(ctx context.Context, cfg <-chan *agentcfg.AgentConfiguration) error {
	if !m.isFluxInstalled(ctx) {
		m.log.Debug("Flux is not installed, skipping module. A restart is required for this to be checked again")
		<-ctx.Done()
		return nil
	}

	wh := syncz.NewProtoWorkerHolder[*agentcfg.FluxCF](
		func(config *agentcfg.FluxCF) syncz.Worker {
			return syncz.WorkerFunc(func(ctx context.Context) {
				if err := m.run(ctx, config); err != nil {
					m.log.Error("Failed to run module", logz.Error(err))
				}
			})
		},
	)
	defer wh.StopAndWait()

	for config := range cfg {
		wh.ApplyConfig(ctx, config.Flux)
	}
	return nil
}

func (m *module) run(ctx context.Context, cfg *agentcfg.FluxCF) error {
	var wg wait.Group
	defer wg.Wait()

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	gitRepositoryInformer, receiverInformer, receiverIndexer := m.informersFactory()

	cl, err := m.clientFactory(runCtx, cfg.WebhookReceiverUrl, receiverIndexer)
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

func (m *module) isFluxInstalled(ctx context.Context) bool {
	for _, crd := range requiredFluxCrds {
		ok, err := checkCRDExistsAndEstablished(ctx, m.k8sExtApiClient, crd)
		if err != nil {
			m.log.Error("Unable to check if CRD is installed", logz.K8sGroup(crd.Group), logz.Error(err))
			return false
		}
		if !ok {
			m.log.Debug("Required Flux CRD is not established", logz.K8sResource(crd.Resource))
			return false
		}
	}
	return true
}

func checkCRDExistsAndEstablished(ctx context.Context, client apiextensionsv1client.ApiextensionsV1Interface, crd schema.GroupResource) (bool, error) {
	obj, err := client.CustomResourceDefinitions().Get(ctx, crd.String(), metav1.GetOptions{})
	if err != nil {
		if kubeerrors.IsNotFound(err) {
			return false, nil
		}
		return false, fmt.Errorf("unable to get CRD %s: %w", crd.String(), err)
	}

	established := false
	for _, cond := range obj.Status.Conditions {
		switch cond.Type { // nolint:exhaustive
		case apiextensionsv1api.Established:
			if cond.Status == apiextensionsv1api.ConditionTrue {
				established = true
			}
			// we don't really care about any other conditions for now, because we don't own this CRD
			// and expect the owner to make sure it becomes established.
		}
	}
	return established, nil
}
