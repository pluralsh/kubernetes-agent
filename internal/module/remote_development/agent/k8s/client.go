package k8s

import (
	"context"
	"errors"
	"io"
	"strings"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/logz"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubectl/pkg/cmd/util"
	"sigs.k8s.io/cli-utils/pkg/apply"
	"sigs.k8s.io/cli-utils/pkg/apply/event"
	"sigs.k8s.io/cli-utils/pkg/common"
	"sigs.k8s.io/cli-utils/pkg/inventory"
	"sigs.k8s.io/cli-utils/pkg/kstatus/watcher"
	"sigs.k8s.io/cli-utils/pkg/object/validation"
)

var (
	noInventoryFoundErr       = errors.New("no inventory found")
	noOwningInventoryFoundErr = errors.New("no owning inventory found")
)

type K8sClient struct {
	log       *zap.Logger
	clientset kubernetes.Interface
	applier   *apply.Applier
	factory   util.Factory
}

// applierInfo contains the information that is needed to run an applier command to Kubernetes.
// It contains the inventory object and the objects tracked by that inventory.
type applierInfo struct {
	invInfo *unstructured.Unstructured
	objects []*unstructured.Unstructured
}

func New(log *zap.Logger, factory util.Factory) (*K8sClient, error) {
	clientset, err := factory.KubernetesClientSet()
	if err != nil {
		return nil, err
	}

	invClient, err := inventory.ClusterClientFactory{
		StatusPolicy: inventory.StatusPolicyNone,
	}.NewClient(factory)
	if err != nil {
		return nil, err
	}

	applier, err := apply.NewApplierBuilder().
		WithFactory(factory).
		WithInventoryClient(invClient).
		Build()
	if err != nil {
		return nil, err
	}

	return &K8sClient{
		log:       log,
		clientset: clientset,
		applier:   applier,
		factory:   factory,
	}, nil
}

func (k *K8sClient) NamespaceExists(ctx context.Context, name string) bool {
	client := k.clientset.CoreV1().Namespaces()
	_, err := client.Get(ctx, name, metav1.GetOptions{})
	return err == nil
}

func (k *K8sClient) CreateNamespace(ctx context.Context, name string) error {
	client := k.clientset.CoreV1().Namespaces()
	nsSpec := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	_, err := client.Create(ctx, nsSpec, metav1.CreateOptions{})
	return err
}

func (k *K8sClient) DeleteNamespace(ctx context.Context, name string) error {
	client := k.clientset.CoreV1().Namespaces()
	return client.Delete(ctx, name, metav1.DeleteOptions{})
}

func (k *K8sClient) Apply(ctx context.Context, config string) error {
	objs, err := k.decode(strings.NewReader(config))
	if err != nil {
		return err
	}

	parsedApplierInfo, err := k.groupObjectsByInventory(objs)
	if err != nil {
		return err
	}

	for _, applierInfo := range parsedApplierInfo {
		// process work - apply to cluster
		k.log.Debug("Applying work to cluster", logz.InventoryName(applierInfo.invInfo.GetName()))
		applierOptions := apply.ApplierOptions{
			ServerSideOptions: common.ServerSideOptions{
				ServerSideApply: true,
				ForceConflicts:  true,
				FieldManager:    modagent.FieldManager,
			},
			ReconcileTimeout:         0,
			EmitStatusEvents:         true,
			PruneTimeout:             0,
			ValidationPolicy:         validation.ExitEarly,
			WatcherRESTScopeStrategy: watcher.RESTScopeNamespace,
		}
		events := k.applier.Run(ctx, inventory.WrapInventoryInfoObj(applierInfo.invInfo), applierInfo.objects, applierOptions)

		k.log.Debug("Applied work to cluster", logz.InventoryName(applierInfo.invInfo.GetName()))

		// We listen on the events channel in goroutine
		// not to block this function's execution
		go func() {
			// TODO: Add logging as a part of processing the events above
			// issue: https://gitlab.com/gitlab-org/gitlab/-/issues/397001
			for e := range events {
				k.log.Debug("Applied event", applyEvent(e))
			}
		}()
	}

	return nil
}

func (k *K8sClient) decode(data io.Reader) ([]*unstructured.Unstructured, error) {
	// 1. parse in local mode to retrieve objects.
	builder := resource.NewBuilder(k.factory).
		ContinueOnError().
		Flatten().
		Unstructured().
		Local()

	builder.Stream(data, "main")

	result := builder.Do()
	var objs []*unstructured.Unstructured
	err := result.Visit(func(info *resource.Info, err error) error {
		if err != nil {
			return err
		}
		objs = append(objs, info.Object.(*unstructured.Unstructured))
		return nil
	})
	if err != nil {
		return nil, err
	}

	return objs, nil
}

func (k *K8sClient) groupObjectsByInventory(objs []*unstructured.Unstructured) (map[string]*applierInfo, error) {
	info := map[string]*applierInfo{}
	for _, obj := range objs {
		if inventory.IsInventoryObject(obj) {
			info[obj.GetName()] = &applierInfo{
				invInfo: obj,
				objects: make([]*unstructured.Unstructured, 0),
			}
		}
	}
	if len(info) == 0 {
		return nil, noInventoryFoundErr
	}

	for _, obj := range objs {
		if inventory.IsInventoryObject(obj) {
			continue
		}
		annotations := obj.GetAnnotations()
		key, ok := annotations[inventory.OwningInventoryKey]
		if !ok {
			return nil, noOwningInventoryFoundErr
		}
		info[key].objects = append(info[key].objects, obj)
	}

	return info, nil
}

func applyEvent(event event.Event) zap.Field {
	return zap.Inline(zapcore.ObjectMarshalerFunc(func(encoder zapcore.ObjectEncoder) error {
		encoder.AddString(logz.ApplyEvent, event.String())
		return nil
	}))
}
