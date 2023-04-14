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
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/remote_development/agent/informer"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/remote_development/agent/k8s"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/retry"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/testhelpers"
	"go.uber.org/zap/zaptest"
)

type WorkerTestSuite struct {
	suite.Suite

	runner        worker
	mockApi       *mock_modagent.MockApi
	mockInformer  *informer.MockInformer
	mockK8sClient *k8s.MockClient
}

func TestRemoteDevModuleWorker(t *testing.T) {
	suite.Run(t, new(WorkerTestSuite))
}

func (w *WorkerTestSuite) TestSuccessfulTerminationOfWorkspace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	existingWorkspaceA := w.newMockWorkspaceInAgent("workspaceA")

	w.mockInformer.Resources = map[string]*informer.ParsedWorkspace{
		existingWorkspaceA.Name: existingWorkspaceA,
	}

	// test assumes an existing running workspace that rails intends to terminate
	w.ensureWorkspaceExists(ctx, w.runner.stateTracker, w.mockK8sClient, existingWorkspaceA)

	workspaceChangesFromRails := w.generateInfoForWorkspaceChanges(existingWorkspaceA.Name, WorkspaceStateTerminated, "Running")
	w.expectRequestAndReturnReply(w.mockApi, w.generateRailsRequest(), w.generateRailsResponse(workspaceChangesFromRails))
	err := w.runner.Run(ctx)
	w.Require().NoError(err)
	w.Require().False(w.runner.terminatedTracker.isTerminated(existingWorkspaceA.Name))
	w.Require().Contains(w.runner.stateTracker.persistedVersion, existingWorkspaceA.Name)

	// simulate "Terminating" state for workspace i.e. create workspace if it doesn't already exist
	// TODO: investigate 'Terminating' state and if its possible at all after DWO removal
	// issue: https://gitlab.com/gitlab-org/gitlab/-/issues/396464
	w.ensureWorkspaceExists(ctx, w.runner.stateTracker, w.mockK8sClient, existingWorkspaceA)

	// While the workspace termination in progress, it is expected that rails will continue to request termination
	// while agent will continue to skip when generating payload for the workspace in question
	// TODO: Rails no longer will continue to request termination of a workspace which has already been requested for termination
	// issue: https://gitlab.com/gitlab-org/gitlab/-/issues/406565
	workspaceChangesFromRails = w.generateInfoForWorkspaceChanges(existingWorkspaceA.Name, WorkspaceStateTerminated, "Terminating")
	w.expectRequestAndReturnReply(w.mockApi, w.generateRailsRequest(), w.generateRailsResponse(workspaceChangesFromRails))
	err = w.runner.Run(ctx)
	w.Require().NoError(err)
	w.Require().False(w.runner.terminatedTracker.isTerminated(existingWorkspaceA.Name))
	w.Require().Contains(w.runner.stateTracker.persistedVersion, existingWorkspaceA.Name)

	// In this cycle, agent will discover that the workspace has been deleted
	// which will result in the workspace being tracked in the termination tracker
	delete(w.mockInformer.Resources, existingWorkspaceA.Name)
	workspaceChangesFromRails = w.generateInfoForWorkspaceChanges(existingWorkspaceA.Name, WorkspaceStateTerminated, "Terminating")
	w.expectRequestAndReturnReply(w.mockApi, w.generateRailsRequest(), w.generateRailsResponse(workspaceChangesFromRails))
	err = w.runner.Run(ctx)
	w.Require().NoError(err)
	w.Require().True(w.runner.terminatedTracker.isTerminated(existingWorkspaceA.Name))
	w.Require().Contains(w.runner.stateTracker.persistedVersion, existingWorkspaceA.Name)

	// In the next cycle, it is expected that successful termination will be communicated to rails
	// Rail will then acknowledge that both the desired and actual state are in sync
	expectedRailsRequest := RequestPayload{
		MessageType: MessageTypeWorkspaceUpdates,
		WorkspaceAgentInfos: []WorkspaceAgentInfo{
			{
				Name:       existingWorkspaceA.Name,
				Terminated: true,
			},
		},
	}
	workspaceChangesFromRails = w.generateInfoForWorkspaceChanges(existingWorkspaceA.Name, WorkspaceStateTerminated, WorkspaceStateTerminated)
	w.expectRequestAndReturnReply(w.mockApi, expectedRailsRequest, w.generateRailsResponse(workspaceChangesFromRails))
	err = w.runner.Run(ctx)
	w.Require().NoError(err)
	w.Require().False(w.runner.terminatedTracker.isTerminated(existingWorkspaceA.Name))
	w.Require().NotContains(w.runner.stateTracker.persistedVersion, existingWorkspaceA.Name)

	// Eventually, no more messages for the terminated workspace should be exchanged between agent & rails
	w.expectRequestAndReturnReply(w.mockApi, w.generateRailsRequest(), w.generateRailsResponse())
	err = w.runner.Run(ctx)
	w.Require().NoError(err)
	w.Require().False(w.runner.terminatedTracker.isTerminated(existingWorkspaceA.Name))
	w.Require().NotContains(w.runner.stateTracker.persistedVersion, existingWorkspaceA.Name)
	w.Require().False(w.mockK8sClient.NamespaceExists(ctx, existingWorkspaceA.Namespace))
}

func (w *WorkerTestSuite) TestSuccessfulCreationOfWorkspace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	workspace := "workspaceA"

	// step1: expect nothing in rails req, get creation req in rails resp => expect changes to be applied
	w.expectRequestAndReturnReply(w.mockApi, w.generateRailsRequest(), w.generateRailsResponse(w.generateInfoForNewWorkspace(workspace)))
	err := w.runner.Run(ctx)
	w.Require().NoError(err)
	w.Require().False(w.runner.terminatedTracker.isTerminated(workspace))
	w.Require().Contains(w.runner.stateTracker.persistedVersion, workspace)

	// step2: simulate transition to "Starting" step (modify resource version in informer), expect rails req to contain update
	currentWorkspaceState := w.newMockWorkspaceInAgent(workspace)
	w.updateMockWorkspaceStateInInformer(w.mockInformer, currentWorkspaceState)

	//ensure rails acks the latest persisted version
	workspaceChangesFromRails := w.generateInfoForWorkspaceChanges(workspace, "Running", "Starting")
	workspaceChangesFromRails.DeploymentResourceVersion = currentWorkspaceState.ResourceVersion

	w.expectRequestAndReturnReply(w.mockApi, w.generateRailsRequest(currentWorkspaceState), w.generateRailsResponse(workspaceChangesFromRails))
	err = w.runner.Run(ctx)
	w.Require().NoError(err)
	w.Require().False(w.runner.terminatedTracker.isTerminated(workspace))
	w.Require().Contains(w.runner.stateTracker.persistedVersion, workspace)

	// step3: simulate transition to "Running" step(modify resource version in informer), expect rails req to contain update
	w.updateMockWorkspaceStateInInformer(w.mockInformer, currentWorkspaceState)
	workspaceChangesFromRails = w.generateInfoForWorkspaceChanges(workspace, "Running", "Running")
	workspaceChangesFromRails.DeploymentResourceVersion = currentWorkspaceState.ResourceVersion

	w.expectRequestAndReturnReply(w.mockApi, w.generateRailsRequest(currentWorkspaceState), w.generateRailsResponse(workspaceChangesFromRails))
	err = w.runner.Run(ctx)
	w.Require().NoError(err)
	w.Require().False(w.runner.terminatedTracker.isTerminated(workspace))
	w.Require().Contains(w.runner.stateTracker.persistedVersion, workspace)

	// step4: nothing changes in resource, expect rails req to not contain workspace metadata, expect metadata to be present in tracker (but not in terminated tracker)
	w.expectRequestAndReturnReply(w.mockApi, w.generateRailsRequest(), w.generateRailsResponse())
	err = w.runner.Run(ctx)
	w.Require().NoError(err)
	w.Require().False(w.runner.terminatedTracker.isTerminated(workspace))
	w.Require().Contains(w.runner.stateTracker.persistedVersion, workspace)
}

func (w *WorkerTestSuite) updateMockWorkspaceStateInInformer(mockInformer *informer.MockInformer, workspace *informer.ParsedWorkspace) {
	workspace.ResourceVersion = workspace.ResourceVersion + "1"

	mockInformer.Resources[workspace.Name] = workspace
}

func (w *WorkerTestSuite) ensureWorkspaceExists(ctx context.Context, stateTracker *persistedStateTracker, mockK8sClient *k8s.MockClient, existingWorkspaceA *informer.ParsedWorkspace) {
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

func (w *WorkerTestSuite) generateRailsResponse(infos ...*WorkspaceRailsInfo) ResponsePayload {
	return ResponsePayload{WorkspaceRailsInfos: infos}
}

func (w *WorkerTestSuite) expectRequestAndReturnReply(mockApi *mock_modagent.MockApi, expectedRequest RequestPayload, response ResponsePayload) {
	extractRequestPayload := func(dataReader io.ReadCloser) RequestPayload {
		var request RequestPayload

		raw, err := io.ReadAll(dataReader)
		w.Require().NoError(err)

		err = json.Unmarshal(raw, &request)
		w.Require().NoError(err)

		return request
	}

	mockApi.EXPECT().
		MakeGitLabRequest(gomock.Any(), "/", gomock.Any()).Times(1).
		DoAndReturn(func(ctx context.Context, path string, opts ...modagent.GitLabRequestOption) (*modagent.GitLabResponse, error) {
			requestConfig, err := modagent.ApplyRequestOptions(opts)
			w.Require().NoError(err)

			requestBody := extractRequestPayload(requestConfig.Body)
			w.Require().Equal(expectedRequest, requestBody)

			var body []byte

			body, err = json.Marshal(response)
			w.Require().NoError(err)

			return &modagent.GitLabResponse{
				StatusCode: http.StatusCreated,
				Body:       io.NopCloser(bytes.NewReader(body)),
			}, nil
		})
}

func (w *WorkerTestSuite) generateInfoForNewWorkspace(name string) *WorkspaceRailsInfo {
	return &WorkspaceRailsInfo{
		Name:                      name,
		Namespace:                 name + "-namespace",
		DeploymentResourceVersion: "",
		ActualState:               "Creating",
		DesiredState:              "Running",
		ConfigToApply:             "",
	}
}

func (w *WorkerTestSuite) generateInfoForWorkspaceChanges(name, desiredState, actualState string) *WorkspaceRailsInfo {
	return &WorkspaceRailsInfo{
		Name:                      name,
		Namespace:                 name + "-namespace",
		DeploymentResourceVersion: "1",
		ActualState:               actualState,
		DesiredState:              desiredState,
		ConfigToApply:             "test",
	}
}

func (w *WorkerTestSuite) generateRailsRequest(workspaces ...*informer.ParsedWorkspace) RequestPayload {
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
		WorkspaceAgentInfos: infos,
	}
}

func (w *WorkerTestSuite) newMockWorkspaceInAgent(name string) *informer.ParsedWorkspace {
	return &informer.ParsedWorkspace{
		Name:              name,
		Namespace:         name + "-namespace",
		ResourceVersion:   "1",
		K8sDeploymentInfo: nil,
	}
}

func (w *WorkerTestSuite) SetupTest() {
	ctrl := gomock.NewController(w.T())
	w.mockApi = mock_modagent.NewMockApi(ctrl)

	w.mockK8sClient = k8s.NewMockClient()
	w.mockInformer = informer.NewMockInformer()

	// this should ideally be called once per run
	//  however, since each test may have multiple runs, this is just put here for simplicity
	w.mockApi.EXPECT().GetAgentId(gomock.Any()).AnyTimes()

	w.runner = worker{
		log:        zaptest.NewLogger(w.T()),
		agentId:    testhelpers.AgentId,
		api:        w.mockApi,
		pollConfig: testhelpers.NewPollConfig(time.Second),
		pollFunction: func(ctx context.Context, cfg retry.PollConfig, f retry.PollWithBackoffCtxFunc) error {
			err, _ := f(ctx)
			return err
		},
		stateTracker:      newPersistedStateTracker(),
		terminatedTracker: newPersistedTerminatedWorkspacesTracker(),
		informer:          w.mockInformer,
		k8sClient:         w.mockK8sClient,
	}
}
