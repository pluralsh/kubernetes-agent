package informer

import (
	"context"
	"errors"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/logz"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/tools/cache"
)

// https://caiorcferreira.github.io/post/the-kubernetes-dynamic-client/

type K8sInformer struct {
	informer cache.SharedIndexInformer
}

func NewK8sInformer(log *zap.Logger, informer cache.SharedIndexInformer) (*K8sInformer, error) {
	_, err := informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			// Handler logic
			u := obj.(*unstructured.Unstructured)
			log.Debug("Received add event", logz.WorkspaceNamespace(u.GetNamespace()), logz.WorkspaceName(u.GetName()))
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			// Handler logic
			newU := newObj.(*unstructured.Unstructured)
			oldU := oldObj.(*unstructured.Unstructured)
			if equality.Semantic.DeepEqual(oldU, newU) {
				return
			}
			log.Debug("Received update event", logz.WorkspaceNamespace(newU.GetNamespace()), logz.WorkspaceName(newU.GetName()))
		},
		DeleteFunc: func(obj interface{}) {
			// Handler logic
			switch u := obj.(type) {
			case *unstructured.Unstructured:
				log.Debug("Received delete event", logz.WorkspaceNamespace(u.GetNamespace()), logz.WorkspaceName(u.GetName()))
			default:
				log.Debug("Received unknown delete event")
			}
		},
	})
	if err != nil {
		return nil, err
	}
	return &K8sInformer{
		informer: informer,
	}, nil
}

func (i *K8sInformer) Start(ctx context.Context) error {

	go i.informer.Run(ctx.Done())

	isSynced := cache.WaitForCacheSync(ctx.Done(), i.informer.HasSynced)

	if !isSynced {
		return errors.New("failed to sync informer during init")
	}

	return nil
}

func (i *K8sInformer) List() []*ParsedWorkspace {
	list := i.informer.GetIndexer().List()
	result := make([]*ParsedWorkspace, 0, len(list))

	for _, raw := range list {
		result = append(result, i.parseUnstructuredToWorkspace(raw.(*unstructured.Unstructured)))
	}

	return result
}

func (i *K8sInformer) parseUnstructuredToWorkspace(rawWorkspace *unstructured.Unstructured) *ParsedWorkspace {
	return &ParsedWorkspace{
		Name:              rawWorkspace.GetName(),
		Namespace:         rawWorkspace.GetNamespace(),
		ResourceVersion:   rawWorkspace.GetResourceVersion(),
		K8sDeploymentInfo: rawWorkspace.Object,
	}
}
