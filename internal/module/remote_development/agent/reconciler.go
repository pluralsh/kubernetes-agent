package agent

//goland:noinspection GoSnakeCaseUsage
import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/errors"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/remote_development/agent/k8s"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/errz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/retry"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/pkg/agentcfg"
)

type MessageType string
type WorkspaceUpdateType string

const (
	MessageTypeWorkspaceUpdates MessageType = "workspace_updates"

	WorkspaceUpdateTypePartial WorkspaceUpdateType = "partial"
	WorkspaceUpdateTypeFull    WorkspaceUpdateType = "full"

	WorkspaceStateTerminated string = "Terminated"
)

// reconciler is equipped to and responsible for carrying
// one cycle of reconciliation when reconciler.Run() is invoked
type reconciler struct {
	log               *zap.Logger
	agentId           int64
	api               modagent.Api
	pollConfig        retry.PollConfigFactory
	pollFunction      func(ctx context.Context, cfg retry.PollConfig, f retry.PollWithBackoffCtxFunc) error
	stateTracker      *persistedStateTracker
	terminatedTracker persistedTerminatedWorkspacesTracker
	informer          informer
	k8sClient         k8s.Client
	config            *agentcfg.RemoteCF

	// This is used to determine whether the reconciliation cycle corresponds to
	// a full or partial sync. When a reconciler runs for the first time, this will be false
	// indicating a full sync. After the full sync successfully completes, this will be
	// indicating partial sync for subsequent cycles for the same reconciler
	hasFullSyncRunBefore bool
}

// TODO: revisit all request and response types, and make more strongly typed if possible
// issue: https://gitlab.com/gitlab-org/gitlab/-/issues/396882

type WorkspaceAgentInfo struct {
	Name                    string                 `json:"name"`
	Namespace               string                 `json:"namespace,omitempty"`
	LatestK8sDeploymentInfo map[string]interface{} `json:"latest_k8s_deployment_info,omitempty"`
	Terminated              bool                   `json:"terminated,omitempty"`
}

type RequestPayload struct {
	MessageType         MessageType          `json:"message_type"`
	UpdateType          WorkspaceUpdateType  `json:"update_type"`
	WorkspaceAgentInfos []WorkspaceAgentInfo `json:"workspace_agent_infos"`
}

type WorkspaceRailsInfo struct {
	Name                      string `json:"name"`
	Namespace                 string `json:"namespace"`
	DeploymentResourceVersion string `json:"deployment_resource_version,omitempty"`
	ActualState               string `json:"actual_state,omitempty"`
	DesiredState              string `json:"desired_state,omitempty"`
	ConfigToApply             string `json:"config_to_apply,omitempty"`
}

type ResponsePayload struct {
	WorkspaceRailsInfos []*WorkspaceRailsInfo `json:"workspace_rails_infos"`
}

func (r *reconciler) Run(ctx context.Context) error {
	r.log.Debug("Running reconciliation loop")
	defer r.log.Debug("Reconciliation loop ended")

	// Load and the info on the workspaces from the informer. Skip it if the persisted state in
	// rails is already up-to-date for the workspace
	workspaceAgentInfos, err := r.generateWorkspaceAgentInfos()
	if err != nil {
		return err
	}

	updateType := WorkspaceUpdateTypePartial
	if !r.hasFullSyncRunBefore {
		updateType = WorkspaceUpdateTypeFull
	}

	// Submit the workspace update request to the Rails API. Sends the latest AgentInfos in the request,
	// and receives the latest RailsInfos in the response.
	workspaceRailsInfos, err := r.performWorkspaceUpdateRequestToRailsApi(ctx, updateType, workspaceAgentInfos)
	if err != nil {
		return err
	}

	// Workspace update request was completed successfully, now process any RailsInfos received in the response
	for _, workspaceRailsInfo := range workspaceRailsInfos {
		err = r.applyWorkspaceChanges(ctx, workspaceRailsInfo)
		if err != nil {
			// TODO: how to report this back to rails?
			// issue: https://gitlab.com/gitlab-org/gitlab/-/issues/397001
			r.api.HandleProcessingError(ctx, r.log, r.agentId, "Error when applying workspace info", err)
		}
	}

	r.hasFullSyncRunBefore = true
	return nil
}

func (r *reconciler) applyWorkspaceChanges(ctx context.Context, workspaceRailsInfo *WorkspaceRailsInfo) error {
	r.stateTracker.recordVersion(workspaceRailsInfo)

	name := workspaceRailsInfo.Name
	namespace := workspaceRailsInfo.Namespace

	// When desired state is Terminated, trigger workspace deletion and exit early
	// to avoid processing/applying the workspace config
	if workspaceRailsInfo.DesiredState == WorkspaceStateTerminated {
		// Handle Terminated state (delete the namespace and workspace) and return
		err := r.handleDesiredStateIsTerminated(ctx, name, namespace, workspaceRailsInfo.ActualState)
		if err != nil {
			return fmt.Errorf("error when handling terminated state for workspace %s: %w", name, err)
		}
		// we don't want to continue by creating namespace if we just terminated the workspace
		return nil
	}

	// Desired state is not Terminated, so continue to handle workspace creation and config apply if needed

	// Create namespace if needed
	namespaceExists := r.k8sClient.NamespaceExists(ctx, namespace)
	if !namespaceExists {
		err := r.k8sClient.CreateNamespace(ctx, namespace)
		if err != nil && !errors.IsAlreadyExists(err) {
			return fmt.Errorf("error creating namespace %s: %w", namespace, err)
		}
	}

	// Apply workspace config if one was provided in the workspaceRailsInfo
	if workspaceRailsInfo.ConfigToApply != "" {
		err := r.k8sClient.Apply(ctx, workspaceRailsInfo.Namespace, workspaceRailsInfo.ConfigToApply)
		if err != nil {
			return fmt.Errorf("error applying workspace config (namespace %s, workspace name %s): %w", namespace, name, err)
		}
	}
	return nil
}

func (r *reconciler) generateWorkspaceAgentInfos() ([]WorkspaceAgentInfo, error) {
	parsedWorkspaces := r.informer.List()
	// "workspaceAgentInfos" is constructed by looping over "parsedWorkspaces" and "terminatedTracker".
	// It can remain a nil slice because GitLab already has the latest version of all workspaces,
	// However, we want it to be an empty(0-length) slice. Hence, initializing it.
	// TODO: add a test case - https://gitlab.com/gitlab-org/gitlab/-/issues/407554
	workspaceAgentInfos := make([]WorkspaceAgentInfo, 0)

	for _, workspace := range parsedWorkspaces {
		// if Rails already knows about the latest version of the resource, don't send the info again
		if r.stateTracker.isPersisted(workspace.Name, workspace.ResourceVersion) {
			r.log.Debug("Skipping sending workspace info. GitLab already has the latest version", logz.WorkspaceName(workspace.Name))
			continue
		}

		if r.terminatedTracker.isTerminated(workspace.Name) {
			// TODO: Instead of returning an error and skipping the rest of the informer items, this should instead be returned in the `Error` field of the WorkspaceAgentInfo
			// issue: https://gitlab.com/gitlab-org/gitlab/-/issues/397001
			return nil, fmt.Errorf("invalid state for workspace, workspace exists in terminatedTracker but still exists in informer: %s", workspace.Name)
		}

		workspaceAgentInfos = append(workspaceAgentInfos, WorkspaceAgentInfo{
			Name:                    workspace.Name,
			Namespace:               workspace.Namespace,
			LatestK8sDeploymentInfo: workspace.K8sDeploymentInfo,
		})
	}

	// Add any already-deleted workspaces that are in the persistedTerminatedWorkspacesTracker
	// In this case we send a minimal WorkspaceAgentInfo with only the name and Terminated = true
	// TODO encapsulate this functionality so we don't have to expose the map which is a detail of terminated tracker implementation
	// issue: https://gitlab.com/gitlab-org/gitlab/-/issues/402758
	for terminatedWorkspaceName := range r.terminatedTracker {
		r.log.Debug("Sending workspace info for already-terminated workspace", logz.WorkspaceName(terminatedWorkspaceName))
		workspaceAgentInfos = append(workspaceAgentInfos, WorkspaceAgentInfo{
			Name:       terminatedWorkspaceName,
			Terminated: true,
		})
	}

	return workspaceAgentInfos, nil
}

func (r *reconciler) performWorkspaceUpdateRequestToRailsApi(
	ctx context.Context,
	updateType WorkspaceUpdateType,
	workspaceAgentInfos []WorkspaceAgentInfo,
) (workspaceRailsInfos []*WorkspaceRailsInfo, retError error) {
	// Do the POST request to the Rails API
	r.log.Debug("Making GitLab request")
	var requestPayload = RequestPayload{
		MessageType:         MessageTypeWorkspaceUpdates,
		UpdateType:          updateType,
		WorkspaceAgentInfos: workspaceAgentInfos,
	}
	// below code is from internal/module/starboard_vulnerability/agent/reporter.go
	resp, err := r.api.MakeGitLabRequest(ctx, "/",
		modagent.WithRequestMethod(http.MethodPost),
		modagent.WithJsonRequestBody(requestPayload),
	) // nolint: govet
	if err != nil {
		return nil, fmt.Errorf("error making api request: %w", err)
	}
	r.log.Debug("Made request to the Rails API", logz.StatusCode(resp.StatusCode))

	defer errz.SafeClose(resp.Body, &retError)
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	r.log.Debug("Read body from the Rails API", zap.String("body", string(body)))

	var responsePayload ResponsePayload
	err = json.Unmarshal(body, &responsePayload)
	if err != nil {
		return nil, fmt.Errorf("error parsing response body: %w", err)
	}
	return responsePayload.WorkspaceRailsInfos, nil
}

func (r *reconciler) handleDesiredStateIsTerminated(ctx context.Context, name string, namespace string, actualState string) error {
	if r.terminatedTracker.isTerminated(name) && actualState == WorkspaceStateTerminated {
		r.log.Debug("ActualState=Terminated, so deleting it from persistedTerminatedWorkspacesTracker", logz.WorkspaceNamespace(namespace))
		r.terminatedTracker.delete(name)
		r.stateTracker.delete(name)
		return nil
	}

	if r.k8sClient.NamespaceExists(ctx, namespace) {
		r.log.Debug("Namespace for terminated workspace still exists, so deleting the namespace", logz.WorkspaceNamespace(namespace))
		err := r.k8sClient.DeleteNamespace(ctx, namespace)
		if err != nil {
			return fmt.Errorf("failed to terminate workspace by deleting namespace: %w", err)
		}
		// TODO: Rails no longer will continue to request termination of a workspace which has already been requested for termination
		// issue: https://gitlab.com/gitlab-org/gitlab/-/issues/406565
		return nil
	}

	r.log.Debug("Namespace no longer exists, sending Actual State as terminated", logz.WorkspaceName(name))
	r.terminatedTracker.add(name)

	return nil
}

func (r *reconciler) Stop() {
	// this will be invoked at the end of each full cycle with the outgoing reconciler and its informer being stopped
	// and new reconcilers (with new informers) created from scratch. The underlying principle is to prevent a re-use of
	// reconciler state when a full sync occurs to prevent issues due to corruption of internal state
	// However, this decision can be revisited in the future in of any issues
	// issue: https://gitlab.com/gitlab-org/gitlab/-/issues/404748
	r.informer.Stop()
}
