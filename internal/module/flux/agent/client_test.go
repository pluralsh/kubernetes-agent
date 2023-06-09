package agent

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	notificationv1 "github.com/fluxcd/notification-controller/api/v1"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/flux/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/retry"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/matcher"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_k8s"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_modagent"
	"go.uber.org/zap/zaptest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/utils/clock"
)

func TestClient_isEqualProjectSets(t *testing.T) {
	testcases := []struct {
		name                     string
		projects                 []string
		cachedProjects           []string
		expectedIsUpdateRequired bool
	}{
		{
			name:                     "no projects",
			projects:                 []string{},
			cachedProjects:           []string{},
			expectedIsUpdateRequired: true,
		},
		{
			name:                     "new projects on empty cache",
			projects:                 []string{"foo"},
			cachedProjects:           []string{},
			expectedIsUpdateRequired: false,
		},
		{
			name:                     "removed project on existing cache",
			projects:                 []string{},
			cachedProjects:           []string{"foo"},
			expectedIsUpdateRequired: false,
		},
		{
			name:                     "same projects in same order",
			projects:                 []string{"foo", "bar"},
			cachedProjects:           []string{"foo", "bar"},
			expectedIsUpdateRequired: true,
		},
		{
			name:                     "same projects in different order",
			projects:                 []string{"foo", "bar"},
			cachedProjects:           []string{"bar", "foo"},
			expectedIsUpdateRequired: true,
		},
		{
			name:                     "with duplicates",
			projects:                 []string{"foo", "foo", "bar"},
			cachedProjects:           []string{"foo", "bar"},
			expectedIsUpdateRequired: true,
		},
		{
			name:                     "with duplicates in cache",
			projects:                 []string{"foo", "bar"},
			cachedProjects:           []string{"foo", "foo", "bar"},
			expectedIsUpdateRequired: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// WHEN
			actualIsEqual := isEqualProjectSets(tc.cachedProjects, tc.projects)

			// THEN
			assert.Equal(t, tc.expectedIsUpdateRequired, actualIsEqual)
		})
	}
}

func TestClient_OnlyRestartReconcilingIndexedProjectsWhenNecessary(t *testing.T) {
	testcases := []struct {
		name                     string
		projects                 []string
		cachedProjects           []string
		expectedIsUpdateRequired bool
	}{
		{
			name:                     "no projects",
			projects:                 []string{},
			cachedProjects:           []string{},
			expectedIsUpdateRequired: false,
		},
		{
			name:                     "new projects on empty cache",
			projects:                 []string{"foo"},
			cachedProjects:           []string{},
			expectedIsUpdateRequired: true,
		},
		{
			name:                     "removed project on existing cache",
			projects:                 []string{},
			cachedProjects:           []string{"foo"},
			expectedIsUpdateRequired: true,
		},
		{
			name:                     "same projects in same order",
			projects:                 []string{"foo", "bar"},
			cachedProjects:           []string{"foo", "bar"},
			expectedIsUpdateRequired: false,
		},
		{
			name:                     "same projects in different order",
			projects:                 []string{"foo", "bar"},
			cachedProjects:           []string{"bar", "foo"},
			expectedIsUpdateRequired: false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// GIVEN
			var wg wait.Group
			defer wg.Wait()

			ctrl := gomock.NewController(t)
			mockReceiverIndexer := mock_k8s.NewMockIndexer(ctrl)
			mockGitLabFluxClient := NewMockGitLabFluxClient(ctrl)
			mockAgentApi := mock_modagent.NewMockApi(ctrl)
			ch := make(chan []string, 1)
			c := client{
				log:                            zaptest.NewLogger(t),
				agentApi:                       mockAgentApi,
				receiverIndexer:                mockReceiverIndexer,
				fluxGitLabClient:               mockGitLabFluxClient,
				updateProjectsToReconcileC:     ch,
				pollCfgFactory:                 retry.NewPollConfigFactory(1*time.Hour, func() retry.BackoffManager { return &fakeBackoff{1 * time.Hour} }),
				reconciliationDebounceDuration: 0,
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// setup mock expectations
			gomock.InOrder(
				mockReceiverIndexer.EXPECT().
					ListIndexFuncValues(projectReceiverIndex).
					Return(tc.cachedProjects),
				mockReceiverIndexer.EXPECT().
					ListIndexFuncValues(projectReceiverIndex).
					Return(tc.projects),
			)

			next := make(chan struct{})
			done := make(chan struct{})

			c1 := mockGitLabFluxClient.EXPECT().
				ReconcileProjects(gomock.Any(), matcher.ProtoEq(nil, &rpc.ReconcileProjectsRequest{Project: rpc.ReconcileProjectsFromSlice(tc.cachedProjects)})).
				DoAndReturn(func(_, _ interface{}, _ ...interface{}) (rpc.GitLabFlux_ReconcileProjectsClient, error) {
					close(next)
					return nil, errors.New("just for testing, it's okay")
				})
			c2 := mockAgentApi.EXPECT().HandleProcessingError(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
			calls := []*gomock.Call{c1, c2}

			// we need this to abort the PollWithBackoff in reconcileProjects eventually
			if tc.expectedIsUpdateRequired {
				c3 := mockGitLabFluxClient.EXPECT().
					ReconcileProjects(gomock.Any(), matcher.ProtoEq(nil, &rpc.ReconcileProjectsRequest{Project: rpc.ReconcileProjectsFromSlice(tc.projects)})).
					DoAndReturn(func(_, _ interface{}, _ ...interface{}) (rpc.GitLabFlux_ReconcileProjectsClient, error) {
						close(done)
						return nil, errors.New("just for testing, it's okay")
					})
				c4 := mockAgentApi.EXPECT().HandleProcessingError(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
				calls = append(calls, c3, c4)
			}
			gomock.InOrder(calls...)

			// start reconciliation ...
			wg.StartWithContext(ctx, c.RunProjectReconciliation)

			// WHEN
			c.ReconcileIndexedProjects(ctx)
			<-next
			c.ReconcileIndexedProjects(ctx)
			// THEN
			if tc.expectedIsUpdateRequired {
				<-done
			} else {
				time.Sleep(100 * time.Millisecond) // wait a bit to ensure there are no unexpected calls.
			}
		})
	}
}

type fakeBackoff struct {
	t time.Duration
}

func (b *fakeBackoff) Backoff() clock.Timer {
	return clock.RealClock{}.NewTimer(b.t)
}

func TestClient_RestartsProjectReconciliationOnProjectsUpdate(t *testing.T) {
	// GIVEN
	var wg wait.Group
	defer wg.Wait()

	ctrl := gomock.NewController(t)
	mockGitLabFluxClient := NewMockGitLabFluxClient(ctrl)
	mockReceiverIndexer := mock_k8s.NewMockIndexer(ctrl)
	mockAgentApi := mock_modagent.NewMockApi(ctrl)
	ch := make(chan []string)
	c := client{
		log:                            zaptest.NewLogger(t),
		agentApi:                       mockAgentApi,
		fluxGitLabClient:               mockGitLabFluxClient,
		receiverIndexer:                mockReceiverIndexer,
		updateProjectsToReconcileC:     ch,
		pollCfgFactory:                 retry.NewPollConfigFactory(1*time.Hour, func() retry.BackoffManager { return &fakeBackoff{1 * time.Hour} }),
		reconciliationDebounceDuration: 0,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	firstProjects := []string{"foo", "bar"}
	secondProjects := []string{"tst", "baz"}

	// setup mock expectations
	// we need this to abort the PollWithBackoff in reconcileProjects eventually
	mockAgentApi.EXPECT().HandleProcessingError(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(2)
	mockGitLabFluxClient.EXPECT().
		ReconcileProjects(gomock.Any(), &rpc.ReconcileProjectsRequest{Project: rpc.ReconcileProjectsFromSlice(firstProjects)}).
		Return(nil, errors.New("just for testing, it's okay"))
	mockGitLabFluxClient.EXPECT().
		ReconcileProjects(gomock.Any(), &rpc.ReconcileProjectsRequest{Project: rpc.ReconcileProjectsFromSlice(secondProjects)}).
		DoAndReturn(func(_, _ interface{}, _ ...interface{}) (rpc.GitLabFlux_ReconcileProjectsClient, error) {
			cancel()
			return nil, errors.New("just for testing, it's okay")
		})

	// WHEN
	// start reconciliation ...
	wg.StartWithContext(ctx, c.RunProjectReconciliation)

	// start with first set of projects
	ch <- firstProjects

	// give some time to start reconciliation after debounce
	time.Sleep(1 * time.Second)

	// update to the second set of projects
	ch <- secondProjects

	// give some time to start reconciliation after debounce
	time.Sleep(1 * time.Second)

	// THEN
	// we cancel the context in the mock function - that's all we need to know regarding execution.
	<-ctx.Done()
}

func TestClient_SuccessfullyReconcileProjects(t *testing.T) {
	// GIVEN
	var wg wait.Group
	defer wg.Wait()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// projects to reconcile
	projects := []string{"foo", "bar"}

	c, _, _, mockReceiverIndexer, mockReconcileTrigger := setupClientForProjectReconciliation(t, projects, "foo")

	receiverObjs := []*notificationv1.Receiver{{
		ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{projectAnnotationKey: "foo"}},
		Status:     notificationv1.ReceiverStatus{WebhookPath: "/some/webhook/path"},
	}}
	mockReceiverIndexer.EXPECT().
		ByIndex(projectReceiverIndex, "foo").
		Return(receiversToUnstructuredInterfaceSlice(t, receiverObjs), nil)
	mockReconcileTrigger.EXPECT().
		reconcile(gomock.Any(), "/some/webhook/path").
		DoAndReturn(func(_, _ interface{}) error {
			cancel()
			return nil
		})

	// WHEN
	wg.Start(func() {
		c.reconcileProjects(ctx, projects)
	})

	// THEN
	// we cancel the ctx in reconcile mock function to signal that the last call we expected has been executed.
	<-ctx.Done()
}

func TestClient_ReconcileProjectsWithoutAnyInIndex(t *testing.T) {
	// GIVEN
	var wg wait.Group
	defer wg.Wait()

	// projects to reconcile
	projects := []string{"foo", "bar"}

	c, _, _, mockReceiverIndexer, mockReconcileTrigger := setupClientForProjectReconciliation(t, projects, "foo")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mockReceiverIndexer.EXPECT().ByIndex(projectReceiverIndex, "foo").DoAndReturn(func(string, string) ([]interface{}, error) {
		cancel()
		return []interface{}{}, nil
	})
	mockReconcileTrigger.EXPECT().reconcile(gomock.Any(), gomock.Any()).Times(0)

	// WHEN
	wg.Start(func() {
		c.reconcileProjects(ctx, projects)
	})

	// THEN
	// we cancel the context in the ByIndex mock function
	<-ctx.Done()
}

func TestClient_ReconcileProjectsReceiverWithoutProjectAnnotationIsIgnored(t *testing.T) {
	// GIVEN
	var wg wait.Group
	defer wg.Wait()

	// projects to reconcile
	projects := []string{"foo", "bar"}

	c, _, _, mockReceiverIndexer, mockReconcileTrigger := setupClientForProjectReconciliation(t, projects, "foo")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	receiverObjs := []*notificationv1.Receiver{{
		ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{}},
		Status:     notificationv1.ReceiverStatus{WebhookPath: "/some/webhook/path"},
	}}
	mockReceiverIndexer.EXPECT().
		ByIndex(projectReceiverIndex, "foo").
		DoAndReturn(func(string, string) ([]interface{}, error) {
			cancel()
			return receiversToUnstructuredInterfaceSlice(t, receiverObjs), nil
		})
	mockReconcileTrigger.EXPECT().reconcile(gomock.Any(), gomock.Any()).Times(0)

	// WHEN
	wg.Start(func() {
		c.reconcileProjects(ctx, projects)
	})

	// THEN
	// we cancel the context in the ByIndex mock function
	<-ctx.Done()
}

func TestClient_ReconcileProjectsReceiverWithoutWebhookPathIsIgnored(t *testing.T) {
	// GIVEN
	var wg wait.Group
	defer wg.Wait()

	// projects to reconcile
	projects := []string{"foo", "bar"}

	c, _, _, mockReceiverIndexer, mockReconcileTrigger := setupClientForProjectReconciliation(t, projects, "foo")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	receiverObjs := []*notificationv1.Receiver{{
		ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{projectAnnotationKey: "foo"}},
		Status:     notificationv1.ReceiverStatus{},
	}}
	mockReceiverIndexer.EXPECT().
		ByIndex(projectReceiverIndex, "foo").
		DoAndReturn(func(string, string) ([]interface{}, error) {
			cancel()
			return receiversToUnstructuredInterfaceSlice(t, receiverObjs), nil
		})
	mockReconcileTrigger.EXPECT().reconcile(gomock.Any(), gomock.Any()).Times(0)

	// WHEN
	wg.Start(func() {
		c.reconcileProjects(ctx, projects)
	})

	// THEN
	// we cancel the context in the ByIndex mock function
	<-ctx.Done()
}

func setupClientForProjectReconciliation(t *testing.T, projects []string, projectToReconcile string) (*client, *MockGitLabFluxClient, *MockGitLabFlux_ReconcileProjectsClient, *mock_k8s.MockIndexer, *MockreconcileTrigger) { // nolint:unparam
	ctrl := gomock.NewController(t)
	mockGitLabFluxClient := NewMockGitLabFluxClient(ctrl)
	mockRpcClient := NewMockGitLabFlux_ReconcileProjectsClient(ctrl)
	mockReceiverIndexer := mock_k8s.NewMockIndexer(ctrl)
	mockReconcileTrigger := NewMockreconcileTrigger(ctrl)
	ch := make(chan []string)
	c := &client{
		log:                        zaptest.NewLogger(t),
		fluxGitLabClient:           mockGitLabFluxClient,
		receiverIndexer:            mockReceiverIndexer,
		reconcileTrigger:           mockReconcileTrigger,
		updateProjectsToReconcileC: ch,
		pollCfgFactory:             retry.NewPollConfigFactory(1*time.Hour, func() retry.BackoffManager { return &fakeBackoff{1 * time.Hour} }),
	}

	// setup mock expectations
	mockGitLabFluxClient.EXPECT().
		ReconcileProjects(gomock.Any(), &rpc.ReconcileProjectsRequest{Project: rpc.ReconcileProjectsFromSlice(projects)}).
		Return(mockRpcClient, nil)

	mockRpcClient.EXPECT().Recv().Return(&rpc.ReconcileProjectsResponse{Project: &rpc.Project{Id: projectToReconcile}}, nil)
	mockRpcClient.EXPECT().Recv().Return(nil, io.EOF)

	return c, mockGitLabFluxClient, mockRpcClient, mockReceiverIndexer, mockReconcileTrigger
}

func receiversToUnstructuredInterfaceSlice(t *testing.T, receivers []*notificationv1.Receiver) []interface{} {
	var objs = make([]interface{}, len(receivers))
	for i, o := range receivers {
		u, err := runtime.DefaultUnstructuredConverter.ToUnstructured(o)
		assert.NoError(t, err)
		objs[i] = &unstructured.Unstructured{Object: u}
	}
	return objs
}
