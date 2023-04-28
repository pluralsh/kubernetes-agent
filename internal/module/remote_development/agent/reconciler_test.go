package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/remote_development/agent/k8s"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/retry"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/testhelpers"
	"go.uber.org/zap/zaptest"
)

type ReconcilerTestSuite struct {
	suite.Suite

	runner        reconciler
	mockApi       *mock_modagent.MockApi
	mockInformer  *mockInformer
	mockK8sClient *k8s.MockClient
}

func TestRemoteDevModuleWorker(t *testing.T) {
	suite.Run(t, new(ReconcilerTestSuite))
}

func (r *ReconcilerTestSuite) TestSuccessfulTerminationOfWorkspace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	existingWorkspaceA := r.newMockWorkspaceInAgent("workspaceA")

	r.mockInformer.Resources = map[string]*parsedWorkspace{
		existingWorkspaceA.Name: existingWorkspaceA,
	}

	// test assumes an existing running workspace that rails intends to terminate
	r.ensureWorkspaceExists(ctx, r.runner.stateTracker, r.mockK8sClient, existingWorkspaceA)

	workspaceChangesFromRails := r.generateInfoForWorkspaceChanges(existingWorkspaceA.Name, WorkspaceStateTerminated, "Running")
	r.expectRequestAndReturnReply(r.mockApi, r.generateRailsRequest(WorkspaceUpdateTypeFull), r.generateRailsResponse(workspaceChangesFromRails))
	err := r.runner.Run(ctx)
	r.Require().NoError(err)
	r.Require().False(r.runner.terminatedTracker.isTerminated(existingWorkspaceA.Name))
	r.Require().Contains(r.runner.stateTracker.persistedVersion, existingWorkspaceA.Name)

	// simulate "Terminating" state for workspace i.e. create workspace if it doesn't already exist
	// TODO: investigate 'Terminating' state and if its possible at all after DWO removal
	// issue: https://gitlab.com/gitlab-org/gitlab/-/issues/396464
	r.ensureWorkspaceExists(ctx, r.runner.stateTracker, r.mockK8sClient, existingWorkspaceA)

	// While the workspace termination in progress, it is expected that rails will continue to request termination
	// while agent will continue to skip when generating payload for the workspace in question
	// TODO: Rails no longer will continue to request termination of a workspace which has already been requested for termination
	// issue: https://gitlab.com/gitlab-org/gitlab/-/issues/406565
	workspaceChangesFromRails = r.generateInfoForWorkspaceChanges(existingWorkspaceA.Name, WorkspaceStateTerminated, "Terminating")
	r.expectRequestAndReturnReply(r.mockApi, r.generateRailsRequest(WorkspaceUpdateTypePartial), r.generateRailsResponse(workspaceChangesFromRails))
	err = r.runner.Run(ctx)
	r.Require().NoError(err)
	r.Require().False(r.runner.terminatedTracker.isTerminated(existingWorkspaceA.Name))
	r.Require().Contains(r.runner.stateTracker.persistedVersion, existingWorkspaceA.Name)

	// In this cycle, agent will discover that the workspace has been deleted
	// which will result in the workspace being tracked in the termination tracker
	delete(r.mockInformer.Resources, existingWorkspaceA.Name)
	workspaceChangesFromRails = r.generateInfoForWorkspaceChanges(existingWorkspaceA.Name, WorkspaceStateTerminated, "Terminating")
	r.expectRequestAndReturnReply(r.mockApi, r.generateRailsRequest(WorkspaceUpdateTypePartial), r.generateRailsResponse(workspaceChangesFromRails))
	err = r.runner.Run(ctx)
	r.Require().NoError(err)
	r.Require().True(r.runner.terminatedTracker.isTerminated(existingWorkspaceA.Name))
	r.Require().Contains(r.runner.stateTracker.persistedVersion, existingWorkspaceA.Name)

	// In the next cycle, it is expected that successful termination will be communicated to rails
	// Rail will then acknowledge that both the desired and actual state are in sync
	expectedRailsRequest := RequestPayload{
		MessageType: MessageTypeWorkspaceUpdates,
		UpdateType:  WorkspaceUpdateTypePartial,
		WorkspaceAgentInfos: []WorkspaceAgentInfo{
			{
				Name:       existingWorkspaceA.Name,
				Terminated: true,
			},
		},
	}
	workspaceChangesFromRails = r.generateInfoForWorkspaceChanges(existingWorkspaceA.Name, WorkspaceStateTerminated, WorkspaceStateTerminated)
	r.expectRequestAndReturnReply(r.mockApi, expectedRailsRequest, r.generateRailsResponse(workspaceChangesFromRails))
	err = r.runner.Run(ctx)
	r.Require().NoError(err)
	r.Require().False(r.runner.terminatedTracker.isTerminated(existingWorkspaceA.Name))
	r.Require().NotContains(r.runner.stateTracker.persistedVersion, existingWorkspaceA.Name)

	// Eventually, no more messages for the terminated workspace should be exchanged between agent & rails
	r.expectRequestAndReturnReply(r.mockApi, r.generateRailsRequest(WorkspaceUpdateTypePartial), r.generateRailsResponse())
	err = r.runner.Run(ctx)
	r.Require().NoError(err)
	r.Require().False(r.runner.terminatedTracker.isTerminated(existingWorkspaceA.Name))
	r.Require().NotContains(r.runner.stateTracker.persistedVersion, existingWorkspaceA.Name)
	r.Require().False(r.mockK8sClient.NamespaceExists(ctx, existingWorkspaceA.Namespace))
}

func (r *ReconcilerTestSuite) TestSuccessfulCreationOfWorkspace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	workspace := "workspaceA"

	// step1: expect nothing in rails req, get creation req in rails resp => expect changes to be applied
	r.expectRequestAndReturnReply(r.mockApi, r.generateRailsRequest(WorkspaceUpdateTypeFull), r.generateRailsResponse(r.generateInfoForNewWorkspace(workspace)))
	err := r.runner.Run(ctx)
	r.Require().NoError(err)
	r.Require().False(r.runner.terminatedTracker.isTerminated(workspace))
	r.Require().Contains(r.runner.stateTracker.persistedVersion, workspace)

	// step2: simulate transition to "Starting" step (modify resource version in informer), expect rails req to contain update
	currentWorkspaceState := r.newMockWorkspaceInAgent(workspace)
	r.updateMockWorkspaceStateInInformer(r.mockInformer, currentWorkspaceState)

	//ensure rails acks the latest persisted version
	workspaceChangesFromRails := r.generateInfoForWorkspaceChanges(workspace, "Running", "Starting")
	workspaceChangesFromRails.DeploymentResourceVersion = currentWorkspaceState.ResourceVersion

	r.expectRequestAndReturnReply(r.mockApi, r.generateRailsRequest(WorkspaceUpdateTypePartial, currentWorkspaceState), r.generateRailsResponse(workspaceChangesFromRails))
	err = r.runner.Run(ctx)
	r.Require().NoError(err)
	r.Require().False(r.runner.terminatedTracker.isTerminated(workspace))
	r.Require().Contains(r.runner.stateTracker.persistedVersion, workspace)

	// step3: simulate transition to "Running" step(modify resource version in informer), expect rails req to contain update
	r.updateMockWorkspaceStateInInformer(r.mockInformer, currentWorkspaceState)
	workspaceChangesFromRails = r.generateInfoForWorkspaceChanges(workspace, "Running", "Running")
	workspaceChangesFromRails.DeploymentResourceVersion = currentWorkspaceState.ResourceVersion

	r.expectRequestAndReturnReply(r.mockApi, r.generateRailsRequest(WorkspaceUpdateTypePartial, currentWorkspaceState), r.generateRailsResponse(workspaceChangesFromRails))
	err = r.runner.Run(ctx)
	r.Require().NoError(err)
	r.Require().False(r.runner.terminatedTracker.isTerminated(workspace))
	r.Require().Contains(r.runner.stateTracker.persistedVersion, workspace)

	// step4: nothing changes in resource, expect rails req to not contain workspace metadata, expect metadata to be present in tracker (but not in terminated tracker)
	r.expectRequestAndReturnReply(r.mockApi, r.generateRailsRequest(WorkspaceUpdateTypePartial), r.generateRailsResponse())
	err = r.runner.Run(ctx)
	r.Require().NoError(err)
	r.Require().False(r.runner.terminatedTracker.isTerminated(workspace))
	r.Require().Contains(r.runner.stateTracker.persistedVersion, workspace)
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

	if !mockK8sClient.NamespaceExists(ctx, existingWorkspaceA.Namespace) {
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
		r.Require().NoError(err)

		err = json.Unmarshal(raw, &request)
		r.Require().NoError(err)

		return request
	}

	mockApi.EXPECT().
		MakeGitLabRequest(gomock.Any(), "/", gomock.Any()).Times(1).
		DoAndReturn(func(ctx context.Context, path string, opts ...modagent.GitLabRequestOption) (*modagent.GitLabResponse, error) {
			requestConfig, err := modagent.ApplyRequestOptions(opts)
			r.Require().NoError(err)

			requestBody := extractRequestPayload(requestConfig.Body)
			r.Require().Equal(expectedRequest, requestBody)

			var body []byte

			body, err = json.Marshal(response)
			r.Require().NoError(err)

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

func (r *ReconcilerTestSuite) generateRailsRequest(updateType WorkspaceUpdateType, workspaces ...*parsedWorkspace) RequestPayload {
	// "infos" is constructed by looping over "workspaces".
	// It can remain a nil slice if "workspaces" is a 0-length slice.
	// However, we want it to be an empty(0-length) slice. Hence, initializing it.
	infos := make([]WorkspaceAgentInfo, 0)

	for _, workspace := range workspaces {
		info := WorkspaceAgentInfo{
			Name:      workspace.Name,
			Namespace: workspace.Namespace,
		}

		infos = append(infos, info)
	}
	return RequestPayload{
		MessageType:         MessageTypeWorkspaceUpdates,
		UpdateType:          updateType,
		WorkspaceAgentInfos: infos,
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

	r.runner = reconciler{
		log:        zaptest.NewLogger(r.T()),
		agentId:    testhelpers.AgentId,
		api:        r.mockApi,
		pollConfig: testhelpers.NewPollConfig(time.Second),
		pollFunction: func(ctx context.Context, cfg retry.PollConfig, f retry.PollWithBackoffCtxFunc) error {
			err, _ := f(ctx)
			return err
		},
		stateTracker:      newPersistedStateTracker(),
		terminatedTracker: newPersistedTerminatedWorkspacesTracker(),
		informer:          r.mockInformer,
		k8sClient:         r.mockK8sClient,
	}
}
