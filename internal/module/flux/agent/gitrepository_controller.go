package agent

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	notificationv1 "github.com/fluxcd/notification-controller/api/v1"
	"github.com/fluxcd/pkg/apis/meta"
	sourcev1 "github.com/fluxcd/source-controller/api/v1"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/logz"
	"go.uber.org/zap"
	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	applycorev1 "k8s.io/client-go/applyconfigurations/core/v1"
	applymetav1 "k8s.io/client-go/applyconfigurations/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/informers"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const (
	objectNamePrefix = "gitlab-"
	// See https://kubernetes.io/docs/reference/labels-annotations-taints/#app-kubernetes-io-managed-by
	managedByAnnotationKey   = "app.kubernetes.io/managed-by"
	managedByAnnotationValue = "gitlab"
	// Annotation key that has a GitLab project full path as its value, e.g. `gitlab-org/gitlab`.
	projectAnnotationKey   = "agent.gitlab.com/project"
	agentIdAnnotationKey   = "agent.gitlab.com/id"
	receiverSecretInterval = 5 * time.Minute
)

type gitRepositoryController struct {
	log                                 *zap.Logger
	api                                 modagent.Api
	agentId                             int64
	gitLabExternalUrl                   url.URL
	gitRepositoryInformerCacheHasSynced cache.InformerSynced
	gitRepositoryLister                 cache.GenericLister
	receiverInformerCacheHasSynced      cache.InformerSynced
	projectReconciler                   projectReconciler
	receiverApiClient                   dynamic.NamespaceableResourceInterface
	corev1ApiClient                     v1.CoreV1Interface
	workqueue                           workqueue.RateLimitingInterface
}

func newGitRepositoryController(
	ctx context.Context,
	log *zap.Logger,
	api modagent.Api,
	agentId int64,
	gitLabExternalUrl url.URL,
	gitRepositoryInformer informers.GenericInformer,
	receiverInformer informers.GenericInformer,
	projectReconciler projectReconciler,
	receiverApiClient dynamic.NamespaceableResourceInterface,
	corev1ApiClient v1.CoreV1Interface) (*gitRepositoryController, error) {
	queue := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "GitRepositories")

	gitRepositorySharedInformer := gitRepositoryInformer.Informer()
	receiverSharedInformer := receiverInformer.Informer()

	c := &gitRepositoryController{
		log:                                 log,
		api:                                 api,
		agentId:                             agentId,
		gitLabExternalUrl:                   gitLabExternalUrl,
		gitRepositoryInformerCacheHasSynced: gitRepositorySharedInformer.HasSynced,
		gitRepositoryLister:                 gitRepositoryInformer.Lister(),
		receiverInformerCacheHasSynced:      receiverSharedInformer.HasSynced,
		projectReconciler:                   projectReconciler,
		receiverApiClient:                   receiverApiClient,
		corev1ApiClient:                     corev1ApiClient,
		workqueue:                           queue,
	}

	// register for GitRepository informer events
	_, err := gitRepositorySharedInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			c.log.Debug("Handling add of GitRepository")
			c.enqueue(obj)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			newU := newObj.(*unstructured.Unstructured)
			oldU := oldObj.(*unstructured.Unstructured)
			if oldU.GetResourceVersion() == newU.GetResourceVersion() {
				return
			}
			c.log.Debug("Handling update of GitRepository")
			c.enqueue(newObj)
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add event handlers for GitRepository resources: %w", err)
	}

	// register for Receiver informer events
	_, err = receiverSharedInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			c.log.Debug("Handling add of Receiver")
			c.handleReceiverObj(ctx, obj)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			newU := newObj.(*unstructured.Unstructured)
			oldU := oldObj.(*unstructured.Unstructured)
			if oldU.GetResourceVersion() == newU.GetResourceVersion() {
				return
			}
			c.log.Debug("Handling update of Receiver")
			c.handleReceiverObj(ctx, newObj)
		},
		DeleteFunc: func(obj interface{}) {
			c.log.Debug("Handling delete of Receiver")
			c.handleReceiverObj(ctx, obj)
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add event handlers for Receiver resources: %w", err)
	}
	return c, nil
}

// Run runs the reconciliation loop of this controller.
// New reconciliation requests can be enqueued using the enqueue
// method. The reconciliation loop runs a single worker,
// but this may be changed easily in the future.
func (c *gitRepositoryController) Run(ctx context.Context) {
	var wg wait.Group
	// this wait group has strictly to be the last thing to wait for,
	// the queue must be shutdown before.
	defer wg.Wait()
	// making the sure the work queue is being stopped when shutting down
	defer c.workqueue.ShutDown()

	c.log.Info("Starting GitRepository controller")
	defer c.log.Info("Stopped GitRepository controller")

	c.log.Debug("Waiting for GitRepository informer caches to sync")

	if ok := cache.WaitForCacheSync(ctx.Done(), c.gitRepositoryInformerCacheHasSynced, c.receiverInformerCacheHasSynced); !ok {
		// NOTE: context was canceled and we can just return
		return
	}

	c.log.Debug("Starting GitRepository worker")
	wg.Start(func() {
		// this is a long-running function that continuously
		// processes the items from the work queue.
		for c.processNextItem(ctx) {
		}
	})

	c.log.Debug("Started GitRepository worker")
	<-ctx.Done()
	c.log.Debug("Shutting down GitRepository worker")
}

// processNextItem processes the next item in the work queue.
// It returns false when the work queue was shutdown otherwise true (even for errors so that the loop continues)
func (c *gitRepositoryController) processNextItem(ctx context.Context) bool {
	// get next item to process
	obj, shutdown := c.workqueue.Get()
	if shutdown {
		return false
	}

	err := c.processItem(ctx, obj)
	if err != nil {
		c.api.HandleProcessingError(ctx, c.log.With(logz.ObjectKey(obj)), c.agentId, "Failed to reconcile GitRepository", err)
	}
	return true
}

// processItem processes a single given item.
func (c *gitRepositoryController) processItem(ctx context.Context, obj interface{}) error {
	defer c.workqueue.Done(obj)
	var key string
	var ok bool
	// The keys should be strings in the form of `namespace/name`
	if key, ok = obj.(string); !ok {
		// We have to forget the item, because it is actually invalid
		// and there is no chance for it to become valid.
		c.workqueue.Forget(obj)
		return fmt.Errorf("expected string in workqueue but got %#v", obj)
	}

	result := c.reconcile(ctx, key)
	switch result.status {
	case RetryRateLimited:
		// put back to work queue to handle transient errors
		c.workqueue.AddRateLimited(obj)
		return result.error
	case Error:
		c.workqueue.Forget(obj)
		return result.error
	case Success:
		c.log.Debug("Successfully reconciled GitRepository", logz.NamespacedName(key))
		c.workqueue.Forget(obj)
	}
	return nil
}

// reconcile reconciles a single GitRepository object specified with the key argument.
// The key must be in the format of `namespace/name` (GitRepositories are namespaced)
// and references an object.
// This reconcile may be call on events for any kinds of objects, but the key
// argument must always reference a GitRepository. This is common in controllers
// that manage more than one resource type - mostly those are created for the
// resource given by the key.
// In this GitRepository controller case these additional objects are
// Receiver and Secret.
// A reconcile will create or update the Receiver and Secret resource
// required for the GitRepository at hand.
// If the given GitRepository object does no longer exists the reconciliation is stopped,
// but this is normal behavior because reconcile may have been called because of Receiver
// events (as explained above).
// The hostname of the GitRepository.Spec.URL must match the configured GitLab External URL
// in order for this controller to take any action on the GitRepository.
// If the hostname doesn't match the object is left untouched and the reconciliation request
// is dropped.
func (c *gitRepositoryController) reconcile(ctx context.Context, key string) reconciliationResult {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return reconciliationResult{status: Error, error: fmt.Errorf("invalid key format: %q, should be in `namespace/name` format", key)}
	}

	obj, err := c.gitRepositoryLister.ByNamespace(namespace).Get(name)
	if err != nil {
		if kubeerrors.IsNotFound(err) {
			c.log.Debug("Queued GitRepository no longer exists, dropping it", logz.NamespacedName(key))
			return reconciliationResult{status: Success}
		}
		return reconciliationResult{status: RetryRateLimited, error: fmt.Errorf("unable to list GitRepository object %s", key)}
	}

	u, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return reconciliationResult{status: Error, error: fmt.Errorf("received GitRepository object %s cannot be parsed to unstructured data", key)}
	}

	var gitRepository sourcev1.GitRepository
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), &gitRepository)
	if err != nil {
		return reconciliationResult{status: Error, error: fmt.Errorf("unable to convert unstructured object to GitRepository: %w", err)}
	}

	// check if the hosts of the GitRepository URL matches the GitLab external URL that this agent is connected to.
	// If not, then we don't reconcile this GitRepository and leave it untouched
	grUrl, err := url.Parse(gitRepository.Spec.URL)
	if err != nil {
		return reconciliationResult{status: Error, error: fmt.Errorf("unable to parse GitRepository URL %q: %w", gitRepository.Spec.URL, err)}
	}

	if c.gitLabExternalUrl.Hostname() != grUrl.Hostname() {
		c.log.Debug("Dropping reconciliation for GitRepository that is not on configured GitLab host", logz.NamespacedName(key), logz.Url(c.gitLabExternalUrl.Hostname()), logz.Url(grUrl.Hostname()))
		return reconciliationResult{status: Success}
	}

	gitRepositoryGitLabPath, err := getProjectPathFromRepositoryUrl(gitRepository.Spec.URL)
	if err != nil {
		return reconciliationResult{status: Error, error: fmt.Errorf("unable to extract GitLab project path from URL: %w", err)}
	}

	c.log.Debug("Reconciling GitRepository", logz.NamespacedName(key), logz.GitRepositoryUrl(gitRepository.Spec.URL), logz.ProjectId(gitRepositoryGitLabPath))

	// reconcile the Secret required for the Receiver
	secret := newWebhookReceiverSecret(c.agentId, &gitRepository)
	if err = c.reconcileWebhookReceiverSecret(ctx, secret); err != nil {
		return reconciliationResult{status: RetryRateLimited, error: err}
	}

	// reconcile the actual Receiver
	receiver := newWebhookReceiver(c.agentId, &gitRepository, gitRepositoryGitLabPath, *secret.Name)
	if err = c.reconcileWebhookReceiver(ctx, receiver); err != nil {
		return reconciliationResult{status: RetryRateLimited, error: err}
	}

	return reconciliationResult{status: Success}
}

// enqueue adds the given object to the controller work queue for processing
// The given object in the obj argument must be a GitRepository resource
// even if enqueue was called because of an event in another resource.
// Most likely that other resource is owned by said GitRepository.
func (c *gitRepositoryController) enqueue(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		c.log.Error("Unable to enqueue object, because key cannot be retrieved from object", logz.Error(err))
		return
	}
	c.workqueue.Add(key)
}

// handleReceiverObj handles informer events for the Receiver object given in the obj argument.
// If that Receiver object is owned by a GitRepository that GitRepository is enqueued
// for reconciliation by this controller.
// No matter the outcome of handleReceiverObj the projects to reconcile will always
// be updated to the current indexed Receiver objects.
func (c *gitRepositoryController) handleReceiverObj(ctx context.Context, obj interface{}) {
	var object metav1.Object
	var ok bool

	defer c.projectReconciler.ReconcileIndexedProjects(ctx)

	if object, ok = obj.(metav1.Object); !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			c.api.HandleProcessingError(ctx, c.log, c.agentId, "Failed to handle Receiver object", errors.New("unable to decode object, invalid type"))
			return
		}
		object, ok = tombstone.Obj.(metav1.Object)
		if !ok {
			c.api.HandleProcessingError(ctx, c.log, c.agentId, "Failed to handle Receiver object", errors.New("unable to decode tombstone object, invalid type"))
			return
		}
		c.log.Debug("Recovered deleted object", logz.NamespacedName(types.NamespacedName{Name: object.GetName(), Namespace: object.GetNamespace()}.String()))
	}

	ownerRef := metav1.GetControllerOf(object)
	if ownerRef == nil {
		return
	}

	// If this object is not owned by a GitRepository, we should not do anything more with it.
	gv, err := schema.ParseGroupVersion(ownerRef.APIVersion)
	if err != nil {
		c.api.HandleProcessingError(ctx, c.log, c.agentId, fmt.Sprintf("Failed to parse Receiver owner group version %q", ownerRef.APIVersion), err)
		return
	}

	if gv.Group != sourcev1.GroupVersion.Group || ownerRef.Kind != sourcev1.GitRepositoryKind {
		return
	}

	gitRepository, err := c.gitRepositoryLister.ByNamespace(object.GetNamespace()).Get(ownerRef.Name)
	if err != nil {
		if kubeerrors.IsNotFound(err) {
			c.log.Debug("Ignoring orphaned Receiver object", logz.NamespacedName(types.NamespacedName{Name: object.GetName(), Namespace: object.GetNamespace()}.String()))
		} else {
			c.api.HandleProcessingError(ctx, c.log, c.agentId, "Failed to handle Receiver object", errors.New("unable to get owner reference of Receiver"))
		}
		return
	}

	c.enqueue(gitRepository)
}

func (c *gitRepositoryController) reconcileWebhookReceiver(ctx context.Context, receiver *notificationv1.Receiver) error {
	namespacedReceiverApiClient := c.receiverApiClient.Namespace(receiver.Namespace)

	o, err := runtime.DefaultUnstructuredConverter.ToUnstructured(receiver)
	if err != nil {
		return fmt.Errorf("failed to convert Receiver %s/%s to unstructured object: %w", receiver.Namespace, receiver.Name, err)
	}
	u := &unstructured.Unstructured{Object: o}

	if _, err = namespacedReceiverApiClient.Apply(ctx, receiver.Name, u, metav1.ApplyOptions{FieldManager: modagent.FieldManager, Force: true}); err != nil {
		if kubeerrors.IsConflict(err) {
			c.log.Debug("Unable to apply Receiver, because there is a newer version of it available", logz.NamespacedName(u.GetNamespace()+u.GetName()))
			return nil
		}
		return fmt.Errorf("failed to apply Receiver: %w", err)
	}
	return nil
}

func (c *gitRepositoryController) reconcileWebhookReceiverSecret(ctx context.Context, secret *applycorev1.SecretApplyConfiguration) error {
	secrets := c.corev1ApiClient.Secrets(*secret.Namespace)
	if _, err := secrets.Apply(ctx, secret, metav1.ApplyOptions{FieldManager: modagent.FieldManager, Force: true}); err != nil {
		if kubeerrors.IsConflict(err) {
			c.log.Debug("Unable to apply Secret, because there is a newer version of it available", logz.NamespacedName(*secret.Namespace+*secret.Name))
			return nil
		}
		return fmt.Errorf("failed to apply Secret for Receiver: %w", err)
	}
	return nil
}

func newWebhookReceiver(agentId int64, repository *sourcev1.GitRepository, project string, secretName string) *notificationv1.Receiver {
	return &notificationv1.Receiver{
		TypeMeta: metav1.TypeMeta{
			Kind:       notificationv1.ReceiverKind,
			APIVersion: notificationv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      objectWithPrefix(repository.Name),
			Namespace: repository.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(repository, repository.GroupVersionKind()),
			},
			Annotations: map[string]string{
				managedByAnnotationKey: managedByAnnotationValue,
				agentIdAnnotationKey:   strconv.FormatInt(agentId, 10),
				projectAnnotationKey:   project,
			},
		},
		Spec: notificationv1.ReceiverSpec{
			// FIXME: we should use `generic-hmac` so that there is proper authentication to the hook.
			Type:     notificationv1.GenericReceiver,
			Interval: &metav1.Duration{Duration: receiverSecretInterval},
			Resources: []notificationv1.CrossNamespaceObjectReference{
				{
					Kind:      repository.Kind,
					Name:      repository.Name,
					Namespace: repository.Namespace,
				},
			},
			SecretRef: meta.LocalObjectReference{
				Name: secretName,
			},
		},
	}
}

func newWebhookReceiverSecret(agentId int64, repository *sourcev1.GitRepository) *applycorev1.SecretApplyConfiguration {
	return applycorev1.Secret(objectWithPrefix(repository.Name), repository.Namespace).
		WithOwnerReferences(
			applymetav1.OwnerReference().
				WithAPIVersion(repository.GroupVersionKind().GroupVersion().String()).
				WithKind(repository.Kind).
				WithName(repository.Name).
				WithUID(repository.GetUID()).
				WithBlockOwnerDeletion(true).
				WithController(true)).
		WithAnnotations(map[string]string{
			managedByAnnotationKey: managedByAnnotationValue,
			agentIdAnnotationKey:   strconv.FormatInt(agentId, 10),
		}).
		WithData(map[string][]byte{"token": {}})
}

func objectWithPrefix(name string) string {
	return objectNamePrefix + name
}

// getProjectPathFromRepositoryUrl converts a full HTTP(S) or SSH Url into a GitLab full project path
// Flux does not support the SCP-like syntax for SSH and requires a correct SSH Url.
// See https://fluxcd.io/flux/components/source/gitrepositories/#url
func getProjectPathFromRepositoryUrl(fullUrl string) (string, error) {
	u, err := url.Parse(fullUrl)
	if err != nil {
		return "", err
	}

	fullPath := strings.TrimLeft(u.Path, "/")
	path := strings.TrimSuffix(fullPath, ".git")
	return path, nil
}
