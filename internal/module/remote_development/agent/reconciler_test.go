package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/remote_development/agent/k8s"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/testhelpers"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"
)

type ReconcilerTestSuite struct {
	suite.Suite

	runner        reconciler
	mockApi       *mock_modagent.MockApi
	mockInformer  *mockInformer
	mockK8sClient *k8s.MockClient
}

func TestRemoteDevModuleReconciler(t *testing.T) {
	suite.Run(t, new(ReconcilerTestSuite))
}

func (r *ReconcilerTestSuite) TestSuccessfulTerminationOfWorkspace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	existingWorkspaceA := r.newMockWorkspaceInAgent("workspaceA")
	trackerKeyForWorkspace := terminationTrackerKey{
		name:      existingWorkspaceA.Name,
		namespace: existingWorkspaceA.Namespace,
	}

	r.mockInformer.Resources = map[string]*parsedWorkspace{
		existingWorkspaceA.Name: existingWorkspaceA,
	}

	// test assumes an existing running workspace that rails intends to terminate
	r.ensureWorkspaceExists(ctx, r.runner.stateTracker, r.mockK8sClient, existingWorkspaceA)

	workspaceChangesFromRails := r.generateInfoForWorkspaceChanges(existingWorkspaceA.Name, WorkspaceStateTerminated, "Running")
	r.expectRequestAndReturnReply(r.mockApi, r.generateRailsRequest(WorkspaceUpdateTypeFull), r.generateRailsResponse(workspaceChangesFromRails))
	err := r.runner.Run(ctx)
	r.NoError(err)
	r.ensureWorkspaceHasTerminationProgress(existingWorkspaceA.Name, existingWorkspaceA.Namespace, TerminationProgressTerminating)
	r.Contains(r.runner.stateTracker.persistedVersion, existingWorkspaceA.Name)

	// simulate "Terminating" state for workspace i.e. create workspace if it doesn't already exist
	r.ensureWorkspaceExists(ctx, r.runner.stateTracker, r.mockK8sClient, existingWorkspaceA)

	// In the next partial sync, and until the workspace is removed from the cluster, agentk will keep sending workspace
	// info with termination progress as Terminating
	workspaceChangesFromRails = r.generateInfoForWorkspaceChanges(existingWorkspaceA.Name, WorkspaceStateTerminated, "Terminating")
	r.expectRequestAndReturnReply(
		r.mockApi,
		r.generateRailsRequest(WorkspaceUpdateTypePartial, r.agentInfoWithTerminationProgress(existingWorkspaceA, TerminationProgressTerminating)),
		r.generateRailsResponse(workspaceChangesFromRails),
	)
	err = r.runner.Run(ctx)
	r.NoError(err)
	r.ensureWorkspaceHasTerminationProgress(existingWorkspaceA.Name, existingWorkspaceA.Namespace, TerminationProgressTerminating)
	r.Contains(r.runner.stateTracker.persistedVersion, existingWorkspaceA.Name)

	// In this cycle, agent will discover that the workspace has been deleted which will result in the workspace being
	// removed from all the trackers after a successful exchange with rails
	delete(r.mockInformer.Resources, existingWorkspaceA.Name)
	workspaceChangesFromRails = r.generateInfoForWorkspaceChanges(existingWorkspaceA.Name, WorkspaceStateTerminated, WorkspaceStateTerminated)
	r.expectRequestAndReturnReply(
		r.mockApi,
		r.generateRailsRequest(WorkspaceUpdateTypePartial, r.agentInfoWithTerminationProgress(existingWorkspaceA, TerminationProgressTerminated)),
		r.generateRailsResponse(workspaceChangesFromRails),
	)
	err = r.runner.Run(ctx)
	r.NoError(err)
	r.NotContains(r.runner.stateTracker.persistedVersion, existingWorkspaceA.Name)
	r.NotContains(r.runner.terminationTracker, trackerKeyForWorkspace)

	// In the next cycle, no more information for the terminated workspace will be shared with rails as there
	// no information about this workspace either in agent's internal state nor in the cluster
	r.expectRequestAndReturnReply(r.mockApi, r.generateRailsRequest(WorkspaceUpdateTypePartial), r.generateRailsResponse())
	err = r.runner.Run(ctx)
	r.NoError(err)
	r.NotContains(r.runner.terminationTracker, trackerKeyForWorkspace)
	r.NotContains(r.runner.stateTracker.persistedVersion, existingWorkspaceA.Name)
	r.False(r.mockK8sClient.NamespaceExists(ctx, existingWorkspaceA.Namespace))
}

func (r *ReconcilerTestSuite) TestSuccessfulCreationOfWorkspace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	workspace := "workspaceA"
	currentWorkspaceState := r.newMockWorkspaceInAgent(workspace)
	trackerKeyForWorkspace := terminationTrackerKey{
		name:      currentWorkspaceState.Name,
		namespace: currentWorkspaceState.Namespace,
	}

	// step1: expect nothing in rails req, get creation req in rails resp => expect changes to be applied
	r.expectRequestAndReturnReply(r.mockApi, r.generateRailsRequest(WorkspaceUpdateTypeFull), r.generateRailsResponse(r.generateInfoForNewWorkspace(workspace)))
	err := r.runner.Run(ctx)
	r.NoError(err)
	r.NotContains(r.runner.terminationTracker, trackerKeyForWorkspace)
	r.Contains(r.runner.stateTracker.persistedVersion, workspace)

	// step2: simulate transition to "Starting" step (modify resource version in informer), expect rails req to contain update
	r.updateMockWorkspaceStateInInformer(r.mockInformer, currentWorkspaceState)

	//ensure rails acks the latest persisted version
	workspaceChangesFromRails := r.generateInfoForWorkspaceChanges(workspace, "Running", "Starting")
	workspaceChangesFromRails.DeploymentResourceVersion = currentWorkspaceState.ResourceVersion

	r.expectRequestAndReturnReply(
		r.mockApi,
		r.generateRailsRequest(WorkspaceUpdateTypePartial, r.agentInfoForNonTerminatedWorkspace(currentWorkspaceState)),
		r.generateRailsResponse(workspaceChangesFromRails),
	)
	err = r.runner.Run(ctx)
	r.NoError(err)
	r.NotContains(r.runner.terminationTracker, trackerKeyForWorkspace)
	r.Contains(r.runner.stateTracker.persistedVersion, workspace)

	// step3: simulate transition to "Running" step(modify resource version in informer), expect rails req to contain update
	r.updateMockWorkspaceStateInInformer(r.mockInformer, currentWorkspaceState)
	workspaceChangesFromRails = r.generateInfoForWorkspaceChanges(workspace, "Running", "Running")
	workspaceChangesFromRails.DeploymentResourceVersion = currentWorkspaceState.ResourceVersion

	r.expectRequestAndReturnReply(
		r.mockApi,
		r.generateRailsRequest(WorkspaceUpdateTypePartial, r.agentInfoForNonTerminatedWorkspace(currentWorkspaceState)),
		r.generateRailsResponse(workspaceChangesFromRails),
	)
	err = r.runner.Run(ctx)
	r.NoError(err)
	r.NotContains(r.runner.terminationTracker, trackerKeyForWorkspace)
	r.Contains(r.runner.stateTracker.persistedVersion, workspace)

	// step4: nothing changes in resource, expect rails req to not contain workspace metadata, expect metadata to be present in tracker (but not in terminated tracker)
	r.expectRequestAndReturnReply(r.mockApi, r.generateRailsRequest(WorkspaceUpdateTypePartial), r.generateRailsResponse())
	err = r.runner.Run(ctx)
	r.NoError(err)
	r.NotContains(r.runner.terminationTracker, trackerKeyForWorkspace)
	r.Contains(r.runner.stateTracker.persistedVersion, workspace)
}

func (r *ReconcilerTestSuite) TestSuccessfulReportingOfApplierErrors() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	workspace := "workspaceA"
	currentWorkspaceState := r.newMockWorkspaceInAgent(workspace)
	applierError := errors.New("applier error")

	// step1: expect nothing in rails req, get creation req in rails resp => changes should FAIL to apply
	r.expectRequestAndReturnReply(
		r.mockApi,
		r.generateRailsRequest(WorkspaceUpdateTypeFull),
		r.generateRailsResponse(
			r.generateInfoForWorkspaceChanges(
				workspace,
				"Running",
				"Creating",
			),
		),
	)
	r.ensureK8sClientReturnsApplierError(applierError)

	err := r.runner.Run(ctx)
	r.runner.applierErrorTracker.waitForErrors()

	r.NoError(err)
	r.Contains(r.runner.stateTracker.persistedVersion, workspace)
	r.Contains(r.runner.applierErrorTracker.store, errorTrackerKey{
		name:      currentWorkspaceState.Name,
		namespace: currentWorkspaceState.Namespace,
	})

	// step 2: expect applier error details in Rails request and ensure tracked error is purged
	// from the tracker
	r.expectRequestAndReturnReply(
		r.mockApi,
		r.generateRailsRequest(
			WorkspaceUpdateTypePartial,
			r.agentInfoWithApplierErrors(currentWorkspaceState, applierError),
		),
		r.generateRailsResponse(
			r.generateInfoForWorkspaceInErrorState(workspace),
		),
	)

	err = r.runner.Run(ctx)
	r.runner.applierErrorTracker.waitForErrors()

	r.NoError(err)
	r.Contains(r.runner.stateTracker.persistedVersion, workspace)
	r.NotContains(r.runner.applierErrorTracker.store, errorTrackerKey{
		name:      currentWorkspaceState.Name,
		namespace: currentWorkspaceState.Namespace,
	})
}

func (r *ReconcilerTestSuite) TestTerminationOfATerminatedWorkspace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	terminatedWorkspace := &parsedWorkspace{
		Name:      "terminated-workspace",
		Namespace: "terminated-workspace-namespace",
	}
	trackerKey := terminationTrackerKey{
		name:      terminatedWorkspace.Name,
		namespace: terminatedWorkspace.Namespace,
	}

	// In the first cycle, the terminated workspace will be tracked by terminationTracker
	workspaceChangesFromRails := r.generateInfoForWorkspaceChanges(terminatedWorkspace.Name, WorkspaceStateTerminated, "Running")
	r.expectRequestAndReturnReply(r.mockApi, r.generateRailsRequest(WorkspaceUpdateTypeFull), r.generateRailsResponse(workspaceChangesFromRails))

	err := r.runner.Run(ctx)
	r.NoError(err)
	r.ensureWorkspaceHasTerminationProgress(terminatedWorkspace.Name, terminatedWorkspace.Namespace, TerminationProgressTerminated)

	// In the next partial sync, agent is expected to not find the workspace in the cluster and consequently
	// inform Rails that the workspace has already been deleted
	r.expectRequestAndReturnReply(
		r.mockApi,
		r.generateRailsRequest(WorkspaceUpdateTypePartial, r.agentInfoWithTerminationProgress(terminatedWorkspace, TerminationProgressTerminated)),
		r.generateRailsResponse(r.generateInfoForWorkspaceChanges(terminatedWorkspace.Name, WorkspaceStateTerminated, WorkspaceStateTerminated)),
	)
	err = r.runner.Run(ctx)
	r.NoError(err)
	r.NotContains(r.runner.terminationTracker, trackerKey)
}

func (r *ReconcilerTestSuite) TestSerializationOfRailsRequests() {
	workspace := &parsedWorkspace{
		Name:            "workspace1",
		Namespace:       "namespace1",
		ResourceVersion: "123",
		K8sDeploymentInfo: map[string]interface{}{
			"a": 1,
		},
	}

	tests := []struct {
		testCase              string
		request               RequestPayload
		expectedSerialization string
	}{
		{
			testCase: "partial sync payload with latest k8s info",
			request: r.generateRailsRequest(
				WorkspaceUpdateTypePartial,
				WorkspaceAgentInfo{
					Name:                    workspace.Name,
					Namespace:               workspace.Namespace,
					LatestK8sDeploymentInfo: workspace.K8sDeploymentInfo,
				},
			),
			expectedSerialization: "{\"update_type\":\"partial\",\"workspace_agent_infos\":[{\"name\":\"workspace1\",\"namespace\":\"namespace1\",\"latest_k8s_deployment_info\":{\"a\":1}}]}",
		},
		{
			testCase: "partial sync payload with only terminating progress",
			request: r.generateRailsRequest(
				WorkspaceUpdateTypePartial,
				r.agentInfoWithTerminationProgress(workspace, TerminationProgressTerminating),
			),
			expectedSerialization: "{\"update_type\":\"partial\",\"workspace_agent_infos\":[{\"name\":\"workspace1\",\"namespace\":\"namespace1\",\"termination_progress\":\"Terminating\"}]}",
		},
		{
			testCase: "partial sync payload with applier errors",
			request: r.generateRailsRequest(
				WorkspaceUpdateTypePartial,
				r.agentInfoWithApplierErrors(workspace, errors.New("applierError")),
			),
			expectedSerialization: "{\"update_type\":\"partial\",\"workspace_agent_infos\":[{\"name\":\"workspace1\",\"namespace\":\"namespace1\",\"error_details\":{\"error_type\":\"applier\",\"error_message\":\"applierError\"}}]}",
		},
		{
			testCase:              "full sync payload",
			request:               r.generateRailsRequest(WorkspaceUpdateTypeFull),
			expectedSerialization: "{\"update_type\":\"full\",\"workspace_agent_infos\":[]}",
		},
	}

	for _, tc := range tests {
		raw, err := json.Marshal(tc.request)

		r.NoError(err)
		r.Equal(tc.expectedSerialization, string(raw))
	}
}

func (r *ReconcilerTestSuite) ensureWorkspaceHasTerminationProgress(name, namespace string, expectedProgress TerminationProgress) {
	trackerKey := terminationTrackerKey{
		name:      name,
		namespace: namespace,
	}

	r.Contains(r.runner.terminationTracker, trackerKey)
	r.Equal(expectedProgress, r.runner.terminationTracker[trackerKey])
}
func (r *ReconcilerTestSuite) ensureK8sClientReturnsApplierError(err error) {
	r.mockK8sClient.MockError = err
}

func (r *ReconcilerTestSuite) updateMockWorkspaceStateInInformer(mockInformer *mockInformer, workspace *parsedWorkspace) {
	workspace.ResourceVersion = workspace.ResourceVersion + "1"

	mockInformer.Resources[workspace.Name] = workspace
}

func (r *ReconcilerTestSuite) ensureWorkspaceExists(ctx context.Context, stateTracker *persistedStateTracker, mockK8sClient *k8s.MockClient, existingWorkspaceA *parsedWorkspace) {
	if _, ok := stateTracker.persistedVersion[existingWorkspaceA.Name]; !ok {
		stateTracker.recordVersion(&WorkspaceRailsInfo{
			Name:                      existingWorkspaceA.Name,
			Namespace:                 existingWorkspaceA.Namespace,
			DeploymentResourceVersion: existingWorkspaceA.ResourceVersion,
		})
	}

	namespaceExists, err := mockK8sClient.NamespaceExists(ctx, existingWorkspaceA.Namespace)
	r.NoError(err)
	if !namespaceExists {
		_ = mockK8sClient.CreateNamespace(ctx, existingWorkspaceA.Namespace)
	}
}

func (r *ReconcilerTestSuite) generateRailsResponse(infos ...*WorkspaceRailsInfo) ResponsePayload {
	return ResponsePayload{WorkspaceRailsInfos: infos}
}

func (r *ReconcilerTestSuite) expectRequestAndReturnReply(mockApi *mock_modagent.MockApi, expectedRequest RequestPayload, response ResponsePayload) {
	extractRequestPayload := func(dataReader io.ReadCloser) RequestPayload {
		var request RequestPayload

		raw, err := io.ReadAll(dataReader)
		r.NoError(err)

		err = json.Unmarshal(raw, &request)
		r.NoError(err)

		return request
	}

	mockApi.EXPECT().
		MakeGitLabRequest(gomock.Any(), "/reconcile", gomock.Any()).Times(1).
		DoAndReturn(func(ctx context.Context, path string, opts ...modagent.GitLabRequestOption) (*modagent.GitLabResponse, error) {
			requestConfig, err := modagent.ApplyRequestOptions(opts)
			r.NoError(err)

			requestBody := extractRequestPayload(requestConfig.Body)
			r.Equal(expectedRequest, requestBody)

			var body []byte

			body, err = json.Marshal(response)
			r.NoError(err)

			return &modagent.GitLabResponse{
				StatusCode: http.StatusCreated,
				Body:       io.NopCloser(bytes.NewReader(body)),
			}, nil
		})
}

func (r *ReconcilerTestSuite) generateInfoForNewWorkspace(name string) *WorkspaceRailsInfo {
	return &WorkspaceRailsInfo{
		Name:                      name,
		Namespace:                 name + "-namespace",
		DeploymentResourceVersion: "",
		ActualState:               "Creating",
		DesiredState:              "Running",
		ConfigToApply:             "",
	}
}

func (r *ReconcilerTestSuite) generateInfoForWorkspaceChanges(name, desiredState, actualState string) *WorkspaceRailsInfo {
	return &WorkspaceRailsInfo{
		Name:                      name,
		Namespace:                 name + "-namespace",
		DeploymentResourceVersion: "1",
		ActualState:               actualState,
		DesiredState:              desiredState,
		ConfigToApply:             "test",
	}
}

func (r *ReconcilerTestSuite) generateInfoForWorkspaceInErrorState(name string) *WorkspaceRailsInfo {
	return &WorkspaceRailsInfo{
		Name:                      name,
		Namespace:                 name + "-namespace",
		DeploymentResourceVersion: "1",
		ActualState:               WorkspaceStateError,
		DesiredState:              "Running",
		ConfigToApply:             "",
	}
}

func (r *ReconcilerTestSuite) generateRailsRequest(updateType WorkspaceUpdateType, agentInfos ...WorkspaceAgentInfo) RequestPayload {
	// agentInfos may be a nil slice. However, we want it to be an empty(0-length) slice. Hence, the explicit initialization.
	if len(agentInfos) == 0 {
		agentInfos = make([]WorkspaceAgentInfo, 0)
	}

	return RequestPayload{
		UpdateType:          updateType,
		WorkspaceAgentInfos: agentInfos,
	}
}

func (r *ReconcilerTestSuite) agentInfoForNonTerminatedWorkspace(workspace *parsedWorkspace) WorkspaceAgentInfo {
	return WorkspaceAgentInfo{
		Name:      workspace.Name,
		Namespace: workspace.Namespace,
	}
}

func (r *ReconcilerTestSuite) agentInfoWithApplierErrors(workspace *parsedWorkspace, err error) WorkspaceAgentInfo {
	return WorkspaceAgentInfo{
		Name:      workspace.Name,
		Namespace: workspace.Namespace,
		ErrorDetails: &ErrorDetails{
			ErrorType:    ErrorTypeApplier,
			ErrorMessage: err.Error(),
		},
	}
}

func (r *ReconcilerTestSuite) agentInfoWithTerminationProgress(workspace *parsedWorkspace, progress TerminationProgress) WorkspaceAgentInfo {
	return WorkspaceAgentInfo{
		Name:                workspace.Name,
		Namespace:           workspace.Namespace,
		TerminationProgress: progress,
	}
}

func (r *ReconcilerTestSuite) newMockWorkspaceInAgent(name string) *parsedWorkspace {
	return &parsedWorkspace{
		Name:              name,
		Namespace:         name + "-namespace",
		ResourceVersion:   "1",
		K8sDeploymentInfo: nil,
	}
}

func (r *ReconcilerTestSuite) SetupTest() {
	ctrl := gomock.NewController(r.T())
	r.mockApi = mock_modagent.NewMockApi(ctrl)

	r.mockK8sClient = k8s.NewMockClient()
	r.mockInformer = newMockInformer()

	// this should ideally be called once per run
	//  however, since each test may have multiple runs, this is just put here for simplicity
	r.mockApi.EXPECT().GetAgentId(gomock.Any()).AnyTimes()

	// this may be called in cases where the reconciliation loop results in applier errors
	r.mockApi.EXPECT().HandleProcessingError(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	r.runner = reconciler{
		log:                 zaptest.NewLogger(r.T()),
		agentId:             testhelpers.AgentId,
		api:                 r.mockApi,
		pollConfig:          testhelpers.NewPollConfig(time.Second),
		stateTracker:        newPersistedStateTracker(),
		informer:            r.mockInformer,
		k8sClient:           r.mockK8sClient,
		terminationTracker:  newTerminationTracker(),
		applierErrorTracker: newErrorTracker(),
	}
}
