package agent

//goland:noinspection GoSnakeCaseUsage
import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/remote_development/agent/k8s"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/errz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/retry"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/agentcfg"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/errors"
)

type MessageType string
type WorkspaceUpdateType string
type TerminationProgress string

const (
	WorkspaceUpdateTypePartial WorkspaceUpdateType = "partial"
	WorkspaceUpdateTypeFull    WorkspaceUpdateType = "full"

	WorkspaceStateTerminated       string              = "Terminated"
	TerminationProgressTerminating TerminationProgress = "Terminating"
	TerminationProgressTerminated  TerminationProgress = "Terminated"
)

// reconciler is equipped to and responsible for carrying
// one cycle of reconciliation when reconciler.Run() is invoked
type reconciler struct {
	log          *zap.Logger
	agentId      int64
	api          modagent.Api
	pollConfig   retry.PollConfigFactory
	pollFunction func(ctx context.Context, cfg retry.PollConfig, f retry.PollWithBackoffCtxFunc) error
	stateTracker *persistedStateTracker
	informer     informer
	k8sClient    k8s.Client
	config       *agentcfg.RemoteCF

	// This is used to determine whether the reconciliation cycle corresponds to
	// a full or partial sync. When a reconciler runs for the first time, this will be false
	// indicating a full sync. After the full sync successfully completes, this will be
	// indicating partial sync for subsequent cycles for the same reconciler
	hasFullSyncRunBefore bool

	// terminatingTracker tracks all workspaces for whom termination has been initiated but
	// still exist in the cluster and therefore considered "Terminating". These workspaces are
	// removed from the tracker once they are removed from the cluster and their Terminated
	// status is persisted in Rails
	terminatingTracker persistedTerminatingWorkspacesTracker
}

// TODO: revisit all request and response types, and make more strongly typed if possible
// issue: https://gitlab.com/gitlab-org/gitlab/-/issues/396882

type WorkspaceAgentInfo struct {
	Name                    string                 `json:"name"`
	Namespace               string                 `json:"namespace,omitempty"`
	LatestK8sDeploymentInfo map[string]interface{} `json:"latest_k8s_deployment_info,omitempty"`
	TerminationProgress     TerminationProgress    `json:"termination_progress,omitempty"`
}

type RequestPayload struct {
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
	workspaceAgentInfos := r.generateWorkspaceAgentInfos()

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

func (r *reconciler) generateWorkspaceAgentInfos() []WorkspaceAgentInfo {
	parsedWorkspaces := r.informer.List()
	// "workspaceAgentInfos" is constructed by looping over "parsedWorkspaces" and "terminatingTracker".
	// It can remain a nil slice because GitLab already has the latest version of all workspaces,
	// However, we want it to be an empty(0-length) slice. Hence, initializing it.
	// TODO: add a test case - https://gitlab.com/gitlab-org/gitlab/-/issues/407554
	workspaceAgentInfos := make([]WorkspaceAgentInfo, 0)

	// nonTerminatedWorkspaces is a set of all workspaces in the cluster returned by the informer
	// it is compared with the workspaces in the terminatingTracker (considered "Terminating") to determine
	// which works have been removed from the cluster and can be deemed "Terminated"
	nonTerminatedWorkspaces := make(map[string]struct{})

	for _, workspace := range parsedWorkspaces {
		// any workspace returned by the informer is deemed as having not been terminated irrespective of its status
		// workspaces that have been terminated completely will be absent from this set
		nonTerminatedWorkspaces[workspace.Name] = struct{}{}

		// if Rails already knows about the latest version of the resource, don't send the info again
		if r.stateTracker.isPersisted(workspace.Name, workspace.ResourceVersion) {
			r.log.Debug("Skipping sending workspace info. GitLab already has the latest version", logz.WorkspaceName(workspace.Name))
			continue
		}

		workspaceAgentInfos = append(workspaceAgentInfos, WorkspaceAgentInfo{
			Name:                    workspace.Name,
			Namespace:               workspace.Namespace,
			LatestK8sDeploymentInfo: workspace.K8sDeploymentInfo,
		})
	}

	// For each workspace that has been scheduled for termination, check if it exists in the cluster
	// If it is missing from the cluster, it can be considered as Terminated otherwise Terminating
	//
	// TODO this implementation will repeatedly send workspaces that have to be completely removed from the
	//	cluster and are still considered Terminating. This may need to be optimized in case it becomes an issue
	// issue: https://gitlab.com/gitlab-org/gitlab/-/issues/408844
	for workspaceName := range r.terminatingTracker {
		_, existsInCluster := nonTerminatedWorkspaces[workspaceName]

		terminationProgress := TerminationProgressTerminating
		if !existsInCluster {
			terminationProgress = TerminationProgressTerminated
		}

		r.log.Debug("Sending termination progress workspace info workspace",
			logz.WorkspaceName(workspaceName),
			logz.WorkspaceTerminationProgress(string(terminationProgress)),
		)
		workspaceAgentInfos = append(workspaceAgentInfos, WorkspaceAgentInfo{
			Name:                workspaceName,
			TerminationProgress: terminationProgress,
		})
	}

	return workspaceAgentInfos
}

func (r *reconciler) performWorkspaceUpdateRequestToRailsApi(
	ctx context.Context,
	updateType WorkspaceUpdateType,
	workspaceAgentInfos []WorkspaceAgentInfo,
) (workspaceRailsInfos []*WorkspaceRailsInfo, retError error) {
	// Do the POST request to the Rails API
	r.log.Debug("Making GitLab request")
	var requestPayload = RequestPayload{
		UpdateType:          updateType,
		WorkspaceAgentInfos: workspaceAgentInfos,
	}
	// below code is from internal/module/starboard_vulnerability/agent/reporter.go
	resp, err := r.api.MakeGitLabRequest(ctx, "/reconcile",
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
	if r.terminatingTracker.isTerminating(name) && actualState == WorkspaceStateTerminated {
		r.log.Debug("ActualState=Terminated, so deleting it from persistedTerminatingWorkspacesTracker", logz.WorkspaceNamespace(namespace))
		r.terminatingTracker.delete(name)
		r.stateTracker.delete(name)
		return nil
	}

	if !r.k8sClient.NamespaceExists(ctx, namespace) {
		// nothing to as the workspace has already absent from the cluster
		return nil
	}

	r.log.Debug("Namespace for terminated workspace still exists, so deleting the namespace", logz.WorkspaceNamespace(namespace))
	err := r.k8sClient.DeleteNamespace(ctx, namespace)
	if err != nil {
		return fmt.Errorf("failed to terminate workspace by deleting namespace: %w", err)
	}

	r.terminatingTracker.add(name)

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
