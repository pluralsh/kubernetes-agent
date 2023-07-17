package agent

import (
	"context"
	"errors"
	"net/url"
	"testing"

	notificationv1 "github.com/fluxcd/notification-controller/api/v1"
	sourcev1 "github.com/fluxcd/source-controller/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_k8s"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_modagent"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"
	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	_ controller = &gitRepositoryController{}
)

func TestGitRepositoryController_getProjectPathFromRepositoryUrl(t *testing.T) {
	testcases := []struct {
		name             string
		fullUrl          string
		expectedFullPath string
	}{
		{
			name:             "HTTPS Url with .git extension",
			fullUrl:          "https://gitlab.com/gitlab-org/cluster-integration/gitlab-agent.git",
			expectedFullPath: "gitlab-org/cluster-integration/gitlab-agent",
		},
		{
			name:             "HTTPS Url without .git extension",
			fullUrl:          "https://gitlab.com/gitlab-org/cluster-integration/gitlab-agent",
			expectedFullPath: "gitlab-org/cluster-integration/gitlab-agent",
		},
		{
			name:             "SSH Url with .git extension",
			fullUrl:          "ssh://git@gitlab.com/gitlab-org/cluster-integration/gitlab-agent.git",
			expectedFullPath: "gitlab-org/cluster-integration/gitlab-agent",
		},
		{
			name:             "SSH Url without .git extension",
			fullUrl:          "ssh://git@gitlab.com/gitlab-org/cluster-integration/gitlab-agent",
			expectedFullPath: "gitlab-org/cluster-integration/gitlab-agent",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// WHEN
			actualFullPath, err := getProjectPathFromRepositoryUrl(tc.fullUrl)

			// THEN
			require.NoError(t, err)
			assert.Equal(t, tc.expectedFullPath, actualFullPath)
		})
	}
}

func TestGitRepositoryController_ProcessNextItemWithInvalidObject(t *testing.T) {
	// GIVEN
	ctrl := gomock.NewController(t)
	mockAgentApi := mock_modagent.NewMockApi(ctrl)
	mockWorkqueue := mock_k8s.NewMockRateLimitingWorkqueue(ctrl)
	c := &gitRepositoryController{
		log:       zaptest.NewLogger(t),
		api:       mockAgentApi,
		workqueue: mockWorkqueue,
	}

	notAString := 42

	// setup mock expectations
	mockWorkqueue.EXPECT().Get().Return(notAString, false)
	mockWorkqueue.EXPECT().Forget(notAString)
	mockWorkqueue.EXPECT().Done(notAString)
	mockAgentApi.EXPECT().HandleProcessingError(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())

	// WHEN
	_ = c.processNextItem(context.Background())
}

func TestGitRepositoryController_ReconcileWithInvalidKeyError(t *testing.T) {
	// GIVEN
	c := &gitRepositoryController{}

	// WHEN
	res := c.reconcile(context.Background(), "foo/bar/too-much")

	// THEN
	assert.Equal(t, Error, res.status)
	assert.ErrorContains(t, res.error, "invalid key format")
}

func TestGitRepositoryController_RetryIfUnableToGetObjectToReconcile(t *testing.T) {
	// GIVEN
	ctrl := gomock.NewController(t)
	mockGitRepositoryLister := mock_k8s.NewMockGenericLister(ctrl)
	mockNamespaceLister := mock_k8s.NewMockGenericNamespaceLister(ctrl)
	c := &gitRepositoryController{
		gitRepositoryLister: mockGitRepositoryLister,
	}

	// setup mock expectations
	mockGitRepositoryLister.EXPECT().ByNamespace("namespace").Return(mockNamespaceLister)
	mockNamespaceLister.EXPECT().Get("name").Return(nil, errors.New("test"))

	// WHEN
	res := c.reconcile(context.Background(), "namespace/name")

	// THEN
	assert.Equal(t, RetryRateLimited, res.status)
	assert.ErrorContains(t, res.error, "unable to list GitRepository object namespace/name")
}

func TestGitRepositoryController_DropNotExistingObjectToReconcileWithSuccess(t *testing.T) {
	// GIVEN
	ctrl := gomock.NewController(t)
	mockGitRepositoryLister := mock_k8s.NewMockGenericLister(ctrl)
	mockNamespaceLister := mock_k8s.NewMockGenericNamespaceLister(ctrl)
	c := &gitRepositoryController{
		log:                 zaptest.NewLogger(t),
		gitRepositoryLister: mockGitRepositoryLister,
	}

	// setup mock expectations
	mockGitRepositoryLister.EXPECT().ByNamespace("namespace").Return(mockNamespaceLister)
	mockNamespaceLister.EXPECT().Get("name").Return(nil, kubeerrors.NewNotFound(schema.GroupResource{}, "test"))

	// WHEN
	res := c.reconcile(context.Background(), "namespace/name")

	// THEN
	assert.Equal(t, Success, res.status)
}

func TestGitRepositoryController_DropGitRepositoryNotInConfiguredGitLabWithSuccess(t *testing.T) {
	// GIVEN
	ctrl := gomock.NewController(t)
	mockGitRepositoryLister := mock_k8s.NewMockGenericLister(ctrl)
	mockNamespaceLister := mock_k8s.NewMockGenericNamespaceLister(ctrl)
	c := &gitRepositoryController{
		log:                 zaptest.NewLogger(t),
		gitRepositoryLister: mockGitRepositoryLister,
		gitLabExternalUrl:   url.URL{Scheme: "https", Host: "another-host.example.com"},
	}

	// setup mock expectations
	mockGitRepositoryLister.EXPECT().ByNamespace("namespace").Return(mockNamespaceLister)
	mockNamespaceLister.EXPECT().Get("name").Return(getTestGitRepositoryAsRuntimeObject(t), nil)

	// WHEN
	res := c.reconcile(context.Background(), "namespace/name")

	// THEN
	assert.Equal(t, Success, res.status)
}

func TestGitRepositoryController_SuccessfullyCreateReceiverAndSecretForGitRepository(t *testing.T) {
	// GIVEN
	ctrl := gomock.NewController(t)
	mockGitRepositoryLister := mock_k8s.NewMockGenericLister(ctrl)
	mockNamespaceLister := mock_k8s.NewMockGenericNamespaceLister(ctrl)
	mockReceiverApiClient := mock_k8s.NewMockNamespaceableResourceInterface(ctrl)
	mockNamespacedReceiverApiClient := mock_k8s.NewMockResourceInterface(ctrl)
	mockCoreV1ApiClient := mock_k8s.NewMockCoreV1Interface(ctrl)
	mockSecretsApiClient := mock_k8s.NewMockSecretInterface(ctrl)
	c := &gitRepositoryController{
		log:                 zaptest.NewLogger(t),
		gitRepositoryLister: mockGitRepositoryLister,
		receiverApiClient:   mockReceiverApiClient,
		corev1ApiClient:     mockCoreV1ApiClient,
		gitLabExternalUrl:   url.URL{Scheme: "https", Host: "gitlab.example.com:8080"},
	}

	// setup mock expectations
	mockGitRepositoryLister.EXPECT().ByNamespace("namespace").Return(mockNamespaceLister)
	mockNamespaceLister.EXPECT().Get("name").Return(getTestGitRepositoryAsRuntimeObject(t), nil)

	// Secret
	mockCoreV1ApiClient.EXPECT().Secrets("namespace").Return(mockSecretsApiClient)
	mockSecretsApiClient.EXPECT().Apply(gomock.Any(), gomock.Any(), gomock.Any())

	// Receiver
	mockReceiverApiClient.EXPECT().Namespace("namespace").Return(mockNamespacedReceiverApiClient)
	mockNamespacedReceiverApiClient.EXPECT().Apply(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())

	// WHEN
	res := c.reconcile(context.Background(), "namespace/name")

	// THEN
	assert.Equal(t, Success, res.status)
}

func TestGitRepositoryController_RetryOnSecretReconciliationFailureForGitRepository(t *testing.T) {
	// GIVEN
	ctrl := gomock.NewController(t)
	mockGitRepositoryLister := mock_k8s.NewMockGenericLister(ctrl)
	mockNamespaceLister := mock_k8s.NewMockGenericNamespaceLister(ctrl)
	mockReceiverApiClient := mock_k8s.NewMockNamespaceableResourceInterface(ctrl)
	mockCoreV1ApiClient := mock_k8s.NewMockCoreV1Interface(ctrl)
	mockSecretsApiClient := mock_k8s.NewMockSecretInterface(ctrl)
	c := &gitRepositoryController{
		log:                 zaptest.NewLogger(t),
		gitRepositoryLister: mockGitRepositoryLister,
		receiverApiClient:   mockReceiverApiClient,
		corev1ApiClient:     mockCoreV1ApiClient,
		gitLabExternalUrl:   url.URL{Scheme: "https", Host: "gitlab.example.com:8080"},
	}

	// setup mock expectations
	mockGitRepositoryLister.EXPECT().ByNamespace("namespace").Return(mockNamespaceLister)
	mockNamespaceLister.EXPECT().Get("name").Return(getTestGitRepositoryAsRuntimeObject(t), nil)

	// Secret
	mockCoreV1ApiClient.EXPECT().Secrets("namespace").Return(mockSecretsApiClient)
	mockSecretsApiClient.EXPECT().Apply(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("expected apply failure for secret")).Times(1)

	// WHEN
	res := c.reconcile(context.Background(), "namespace/name")

	// THEN
	assert.Equal(t, RetryRateLimited, res.status)
}

func TestGitRepositoryController_RetryOnReceiverReconciliationFailureForGitRepository(t *testing.T) {
	// GIVEN
	ctrl := gomock.NewController(t)
	mockGitRepositoryLister := mock_k8s.NewMockGenericLister(ctrl)
	mockNamespaceLister := mock_k8s.NewMockGenericNamespaceLister(ctrl)
	mockReceiverApiClient := mock_k8s.NewMockNamespaceableResourceInterface(ctrl)
	mockNamespacedReceiverApiClient := mock_k8s.NewMockResourceInterface(ctrl)
	mockCoreV1ApiClient := mock_k8s.NewMockCoreV1Interface(ctrl)
	mockSecretsApiClient := mock_k8s.NewMockSecretInterface(ctrl)
	c := &gitRepositoryController{
		log:                 zaptest.NewLogger(t),
		gitRepositoryLister: mockGitRepositoryLister,
		receiverApiClient:   mockReceiverApiClient,
		corev1ApiClient:     mockCoreV1ApiClient,
		gitLabExternalUrl:   url.URL{Scheme: "https", Host: "gitlab.example.com:8080"},
	}

	// setup mock expectations
	mockGitRepositoryLister.EXPECT().ByNamespace("namespace").Return(mockNamespaceLister)
	mockNamespaceLister.EXPECT().Get("name").Return(getTestGitRepositoryAsRuntimeObject(t), nil)

	// Secret
	mockCoreV1ApiClient.EXPECT().Secrets("namespace").Return(mockSecretsApiClient)
	mockSecretsApiClient.EXPECT().Apply(gomock.Any(), gomock.Any(), gomock.Any())

	// Receiver
	mockReceiverApiClient.EXPECT().Namespace("namespace").Return(mockNamespacedReceiverApiClient)
	mockNamespacedReceiverApiClient.EXPECT().Apply(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("expected apply failure for receiver")).Times(1)

	// WHEN
	res := c.reconcile(context.Background(), "namespace/name")

	// THEN
	assert.Equal(t, RetryRateLimited, res.status)
}

func TestGitRepositoryController_IgnoreConflictOnReceiverReconciliationFailureForGitRepository(t *testing.T) {
	// GIVEN
	ctrl := gomock.NewController(t)
	mockGitRepositoryLister := mock_k8s.NewMockGenericLister(ctrl)
	mockNamespaceLister := mock_k8s.NewMockGenericNamespaceLister(ctrl)
	mockReceiverApiClient := mock_k8s.NewMockNamespaceableResourceInterface(ctrl)
	mockNamespacedReceiverApiClient := mock_k8s.NewMockResourceInterface(ctrl)
	mockCoreV1ApiClient := mock_k8s.NewMockCoreV1Interface(ctrl)
	mockSecretsApiClient := mock_k8s.NewMockSecretInterface(ctrl)
	c := &gitRepositoryController{
		log:                 zaptest.NewLogger(t),
		gitRepositoryLister: mockGitRepositoryLister,
		receiverApiClient:   mockReceiverApiClient,
		corev1ApiClient:     mockCoreV1ApiClient,
		gitLabExternalUrl:   url.URL{Scheme: "https", Host: "gitlab.example.com:8080"},
	}

	// setup mock expectations
	mockGitRepositoryLister.EXPECT().ByNamespace("namespace").Return(mockNamespaceLister)
	mockNamespaceLister.EXPECT().Get("name").Return(getTestGitRepositoryAsRuntimeObject(t), nil)

	// Secret
	mockCoreV1ApiClient.EXPECT().Secrets("namespace").Return(mockSecretsApiClient)
	mockSecretsApiClient.EXPECT().Apply(gomock.Any(), gomock.Any(), gomock.Any())

	// Receiver
	mockReceiverApiClient.EXPECT().Namespace("namespace").Return(mockNamespacedReceiverApiClient)
	mockNamespacedReceiverApiClient.EXPECT().Apply(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, kubeerrors.NewConflict(schema.GroupResource{}, "test", errors.New("conflict"))).Times(1)

	// WHEN
	res := c.reconcile(context.Background(), "namespace/name")

	// THEN
	assert.Equal(t, Success, res.status)
}

func TestGitRepositoryController_ReceiverObjUpdateChangeTriggersProjectReconciliation(t *testing.T) {
	// GIVEN
	ctrl := gomock.NewController(t)
	mockGitRepositoryLister := mock_k8s.NewMockGenericLister(ctrl)
	mockNamespaceLister := mock_k8s.NewMockGenericNamespaceLister(ctrl)
	mockProjectReconciler := NewMockprojectReconciler(ctrl)
	mockWorkqueue := mock_k8s.NewMockRateLimitingWorkqueue(ctrl)
	c := &gitRepositoryController{
		log:                 zaptest.NewLogger(t),
		gitRepositoryLister: mockGitRepositoryLister,
		projectReconciler:   mockProjectReconciler,
		workqueue:           mockWorkqueue,
		gitLabExternalUrl:   url.URL{Scheme: "https", Host: "gitlab.example.com:8080"},
	}

	// setup mock expectations
	mockGitRepositoryLister.EXPECT().ByNamespace("namespace").Return(mockNamespaceLister)
	mockNamespaceLister.EXPECT().Get("name").Return(getTestGitRepositoryAsRuntimeObject(t), nil)
	mockWorkqueue.EXPECT().Add(gomock.Any())

	mockProjectReconciler.EXPECT().ReconcileIndexedProjects(gomock.Any())

	// WHEN
	c.handleReceiverObj(context.Background(), getTestReceiverAsInterface())
}

func TestGitRepositoryController_ReceiverObjUpdateChangeTriggersProjectReconciliationForDelete(t *testing.T) {
	// GIVEN
	ctrl := gomock.NewController(t)
	mockGitRepositoryLister := mock_k8s.NewMockGenericLister(ctrl)
	mockNamespaceLister := mock_k8s.NewMockGenericNamespaceLister(ctrl)
	mockProjectReconciler := NewMockprojectReconciler(ctrl)
	mockWorkqueue := mock_k8s.NewMockRateLimitingWorkqueue(ctrl)
	c := &gitRepositoryController{
		log:                 zaptest.NewLogger(t),
		gitRepositoryLister: mockGitRepositoryLister,
		projectReconciler:   mockProjectReconciler,
		workqueue:           mockWorkqueue,
		gitLabExternalUrl:   url.URL{Scheme: "https", Host: "gitlab.example.com:8080"},
	}

	// setup mock expectations
	mockGitRepositoryLister.EXPECT().ByNamespace("namespace").Return(mockNamespaceLister)
	mockNamespaceLister.EXPECT().Get("name").Return(nil, kubeerrors.NewNotFound(schema.GroupResource{}, "test"))

	mockProjectReconciler.EXPECT().ReconcileIndexedProjects(gomock.Any())

	// WHEN
	c.handleReceiverObj(context.Background(), getTestReceiverAsInterface())
}

func getTestGitRepositoryAsRuntimeObject(t *testing.T) runtime.Object {
	gitRepository := &sourcev1.GitRepository{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "namespace",
			Name:      "name",
		},
		Spec: sourcev1.GitRepositorySpec{
			URL: "https://gitlab.example.com/some-org/some-repo.git",
		},
	}
	o, err := runtime.DefaultUnstructuredConverter.ToUnstructured(gitRepository)
	assert.NoError(t, err)
	u := &unstructured.Unstructured{Object: o}
	return u
}

func getTestReceiver() *notificationv1.Receiver {
	isController := true
	return &notificationv1.Receiver{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "namespace",
			Name:      "gitlab-name",
			OwnerReferences: []metav1.OwnerReference{{
				APIVersion: sourcev1.GroupVersion.String(),
				Kind:       sourcev1.GitRepositoryKind,
				Name:       "name",
				Controller: &isController,
			}},
		},
	}
}

func getTestReceiverAsInterface() interface{} {
	var o metav1.Object // nolint:gosimple
	o = getTestReceiver()
	return o
}
