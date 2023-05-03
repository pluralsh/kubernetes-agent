package agent

import (
	"context"
	"errors"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/logz"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/tools/cache"
)

// https://caiorcferreira.github.io/post/the-kubernetes-dynamic-client/

type k8sInformer struct {
	informer       cache.SharedIndexInformer
	log            *zap.Logger
	backgroundTask stoppableTask
}

func newK8sInformer(log *zap.Logger, informer cache.SharedIndexInformer) (*k8sInformer, error) {
	_, err := informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			// Handler logic
			u := obj.(*unstructured.Unstructured)
			log.Debug("Received add event", extractEventFields(u)...)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			// Handler logic
			newU := newObj.(*unstructured.Unstructured)
			oldU := oldObj.(*unstructured.Unstructured)
			if equality.Semantic.DeepEqual(oldU, newU) {
				return
			}
			log.Debug("Received update event", extractEventFields(newU)...)
		},
		DeleteFunc: func(obj interface{}) {
			// Handler logic
			switch u := obj.(type) {
			case *unstructured.Unstructured:
				log.Debug("Received delete event", extractEventFields(u)...)
			default:
				log.Debug("Received unknown delete event")
			}
		},
	})
	if err != nil {
		return nil, err
	}
	return &k8sInformer{
		informer: informer,
		log:      log,
	}, nil
}

func extractEventFields(event *unstructured.Unstructured) []zap.Field {
	return []zap.Field{
		logz.WorkspaceNamespace(event.GetNamespace()),
		logz.WorkspaceName(event.GetName()),
	}
}

func (i *k8sInformer) Start(ctx context.Context) error {
	i.backgroundTask = newStoppableTask(ctx, func(ctx context.Context) {
		i.informer.Run(ctx.Done())
	})

	isSynced := cache.WaitForCacheSync(ctx.Done(), i.informer.HasSynced)

	if !isSynced {
		return errors.New("failed to sync informer during init")
	}

	return nil
}

func (i *k8sInformer) List() []*parsedWorkspace {
	list := i.informer.GetIndexer().List()
	result := make([]*parsedWorkspace, 0, len(list))

	for _, raw := range list {
		result = append(result, i.parseUnstructuredToWorkspace(raw.(*unstructured.Unstructured)))
	}

	return result
}

func (i *k8sInformer) parseUnstructuredToWorkspace(rawWorkspace *unstructured.Unstructured) *parsedWorkspace {
	return &parsedWorkspace{
		Name:              rawWorkspace.GetName(),
		Namespace:         rawWorkspace.GetNamespace(),
		ResourceVersion:   rawWorkspace.GetResourceVersion(),
		K8sDeploymentInfo: rawWorkspace.Object,
	}
}

func (i *k8sInformer) Stop() {
	if i.backgroundTask != nil {
		i.backgroundTask.StopAndWait()
		i.log.Info("informer stopped")
	}
}
