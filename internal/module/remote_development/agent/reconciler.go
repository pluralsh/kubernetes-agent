package agent

//goland:noinspection GoSnakeCaseUsage
import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"time"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/remote_development/agent/k8s"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/remote_development/agent/util"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/errz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/httpz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/ioz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/retry"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/errors"
)

type MessageType string
type WorkspaceUpdateType string
type TerminationProgress string
type ErrorType string

const (
	WorkspaceUpdateTypePartial WorkspaceUpdateType = "partial"
	WorkspaceUpdateTypeFull    WorkspaceUpdateType = "full"

	WorkspaceStateTerminated string = "Terminated"
	WorkspaceStateError      string = "Error"

	TerminationProgressTerminating TerminationProgress = "Terminating"
	TerminationProgressTerminated  TerminationProgress = "Terminated"

	ErrorTypeApplier ErrorType = "applier"
)

// reconciler is equipped to and responsible for carrying
// one cycle of reconciliation when reconciler.Run() is invoked
type reconciler struct {
	log          *zap.Logger
	agentId      int64
	api          modagent.Api
	pollConfig   retry.PollConfigFactory
	stateTracker *persistedStateTracker
	informer     informer
	k8sClient    k8s.Client

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

	// applierErrorTracker is used when applying k8s configs for a workspace to the cluster.
	// It tracks the asynchronous errors received and serves as the single source of truth when
	// reporting these errors to Rails
	applierErrorTracker *errorTracker

	// versionCounter is an atomic counter that creates monotonically increasing values for use
	// as versions by the applierErrorTracker. It is monotonically increasing in order to determine
	// if one version is older than another and thereby only watch for errors corresponding to the
	// latest operation
	versionCounter atomic.Uint64
}

// TODO: revisit all request and response types, and make more strongly typed if possible
// issue: https://gitlab.com/gitlab-org/gitlab/-/issues/396882

type WorkspaceAgentInfo struct {
	Name                    string                 `json:"name"`
	Namespace               string                 `json:"namespace,omitempty"`
	LatestK8sDeploymentInfo map[string]interface{} `json:"latest_k8s_deployment_info,omitempty"`
	TerminationProgress     TerminationProgress    `json:"termination_progress,omitempty"`
	ErrorDetails            *ErrorDetails          `json:"error_details,omitempty"`
}

type ErrorDetails struct {
	ErrorType    ErrorType `json:"error_type,omitempty"`
	ErrorMessage string    `json:"error_message,omitempty"`
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

	applierErrorsSnapshot := r.applierErrorTracker.createSnapshot()

	// Load and the info on the workspaces from the informer. Skip it if the persisted state in
	// rails is already up-to-date for the workspace
	workspaceAgentInfos := r.generateWorkspaceAgentInfos(applierErrorsSnapshot)

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
		// TODO: discuss whether/how to deal with errors returned by rails when processing individual workspaces
		//  for ex. there may be a case where no workspace is found for a given combination of workspace name &
		//  workspace
		// issue: https://gitlab.com/gitlab-org/gitlab/-/issues/409807

		// increment version to guarantee ordering
		version := r.versionCounter.Add(1)

		errCh := r.applyWorkspaceChanges(ctx, workspaceRailsInfo, applierErrorsSnapshot)
		r.applierErrorTracker.watchForLatestErrors(ctx, workspaceRailsInfo.Name, workspaceRailsInfo.Namespace, version, errCh)
	}

	r.hasFullSyncRunBefore = true
	return nil
}

func (r *reconciler) applyWorkspaceChanges(ctx context.Context, workspaceRailsInfo *WorkspaceRailsInfo, applierErrorsSnapshot map[errorTrackerKey]operationState) <-chan error {
	r.stateTracker.recordVersion(workspaceRailsInfo)

	name := workspaceRailsInfo.Name
	namespace := workspaceRailsInfo.Namespace
	key := errorTrackerKey{
		name:      name,
		namespace: namespace,
	}

	// If an error was reported to rails successfully, it can now be safely cleaned up from the tracker.
	// we can do a cleanup if the version in the tracker matches the version of
	// error sent to Rails
	if state, hasEntry := applierErrorsSnapshot[key]; hasEntry && state.err != nil {
		r.applierErrorTracker.deleteErrorIfVersion(name, namespace, state.version)
	}

	// When desired state is Terminated, trigger workspace deletion and exit early
	// to avoid processing/applying the workspace config
	if workspaceRailsInfo.DesiredState == WorkspaceStateTerminated {
		// Handle Terminated state (delete the namespace and workspace) and return
		err := r.handleDesiredStateIsTerminated(ctx, name, namespace, workspaceRailsInfo.ActualState)
		if err != nil {
			// TODO: this is technically not an applier error. A separate error type should be identified for it
			// and reported accordingly
			// issue: https://gitlab.com/gitlab-org/gitlab/-/issues/420966
			applierErr := fmt.Errorf("error when handling terminated state for workspace %s: %w", name, err)
			return util.ToAsync(applierErr)
		}
		// we don't want to continue by creating namespace if we just terminated the workspace
		return util.ToAsync[error](nil)
	}

	// Desired state is not Terminated, so continue to handle workspace creation and config apply if needed

	// Create namespace if needed
	namespaceExists := r.k8sClient.NamespaceExists(ctx, namespace)
	if !namespaceExists {
		err := r.k8sClient.CreateNamespace(ctx, namespace)
		if err != nil && !errors.IsAlreadyExists(err) {
			// TODO: this is technically not an applier error. A separate error type should be identified for it
			// and reported accordingly
			// issue: https://gitlab.com/gitlab-org/gitlab/-/issues/420966
			applierErr := fmt.Errorf("error creating namespace %s: %w", namespace, err)
			return util.ToAsync(applierErr)
		}
	}

	// Apply workspace config if one was provided in the workspaceRailsInfo
	if workspaceRailsInfo.ConfigToApply != "" {
		return r.k8sClient.Apply(ctx, workspaceRailsInfo.ConfigToApply)
	}
	return util.ToAsync[error](nil)
}

func (r *reconciler) generateWorkspaceAgentInfos(applierErrorsSnapshot map[errorTrackerKey]operationState) []WorkspaceAgentInfo {
	parsedWorkspaces := r.informer.List()

	// unterminatedWorkspaces is a set of all workspaces in the cluster returned by the informer
	// it is compared with the workspaces in the terminatingTracker (considered "Terminating") to determine
	// which ones have been removed from the cluster and can be deemed "Terminated"
	unterminatedWorkspaces := make(map[string]struct{})
	for _, workspace := range parsedWorkspaces {
		// any workspace returned by the informer is deemed as having not been terminated irrespective of its status
		// workspaces that have been terminated completely will be absent from this set
		unterminatedWorkspaces[workspace.Name] = struct{}{}
	}

	workspaceInfos := r.collectInfoForUnpersistedWorkspaces(parsedWorkspaces)
	r.enrichRailsPayloadWithWorkspaceTerminationProgress(workspaceInfos, unterminatedWorkspaces)
	r.enrichRailsPayloadWithApplierErrorDetails(workspaceInfos, applierErrorsSnapshot)

	// "result" is constructed by taking the values collected in workspaceInfos which in turn
	// is prepared by starting off with unpersisted workspace data and enrich it with termination progress
	// and error details
	result := make([]WorkspaceAgentInfo, 0, len(workspaceInfos))
	for _, agentInfo := range workspaceInfos {
		result = append(result, agentInfo)
	}

	return result
}

func (r *reconciler) collectInfoForUnpersistedWorkspaces(existingWorkspaces []*parsedWorkspace) map[workspaceInfoKey]WorkspaceAgentInfo {
	workspaceAgentInfos := make(map[workspaceInfoKey]WorkspaceAgentInfo)

	for _, workspace := range existingWorkspaces {
		// if Rails already knows about the latest version of the resource, don't send the info again
		if r.stateTracker.isPersisted(workspace.Name, workspace.ResourceVersion) {
			r.log.Debug("Skipping sending workspace info. GitLab already has the latest version", logz.WorkspaceName(workspace.Name))
			continue
		}

		key := workspaceInfoKey{
			Name:      workspace.Name,
			Namespace: workspace.Namespace,
		}
		workspaceAgentInfos[key] = WorkspaceAgentInfo{
			Name:                    workspace.Name,
			Namespace:               workspace.Namespace,
			LatestK8sDeploymentInfo: workspace.K8sDeploymentInfo,
		}
	}

	return workspaceAgentInfos
}

func (r *reconciler) enrichRailsPayloadWithApplierErrorDetails(payload map[workspaceInfoKey]WorkspaceAgentInfo, errorsSnapshot map[errorTrackerKey]operationState) {
	// for each entry in the error snapshot, either a new payload has to be
	// created or the error details must be merged into the existing payload
	// being sent over to rails
	for trackerKey, state := range errorsSnapshot {
		if state.err == nil {
			// this case will occur if the config is in the process of being
			// applied and the error status is not yet known
			continue
		}

		key := workspaceInfoKey{
			Name:      trackerKey.name,
			Namespace: trackerKey.namespace,
		}

		var agentInfo WorkspaceAgentInfo
		var exists bool

		agentInfo, exists = payload[key]
		if !exists {
			agentInfo = WorkspaceAgentInfo{
				Name:      trackerKey.name,
				Namespace: trackerKey.namespace,
			}
		}

		agentInfo.ErrorDetails = &ErrorDetails{
			ErrorType:    ErrorTypeApplier,
			ErrorMessage: state.err.Error(),
		}
		payload[key] = agentInfo
	}
}

func (r *reconciler) enrichRailsPayloadWithWorkspaceTerminationProgress(payload map[workspaceInfoKey]WorkspaceAgentInfo, unterminatedWorkspaces map[string]struct{}) {
	// For each workspace that has been scheduled for termination, check if it exists in the cluster
	// If it is missing from the cluster, it can be considered as Terminated otherwise Terminating
	//
	// TODO this implementation will repeatedly send workspaces that have to be completely removed from the
	//	cluster and are still considered Terminating. This may need to be optimized in case it becomes an issue
	// issue: https://gitlab.com/gitlab-org/gitlab/-/issues/408844
	for entry := range r.terminatingTracker {
		_, existsInCluster := unterminatedWorkspaces[entry.name]

		terminationProgress := TerminationProgressTerminating
		if !existsInCluster {
			terminationProgress = TerminationProgressTerminated
		}

		r.log.Debug("Sending termination progress workspace info workspace",
			logz.WorkspaceName(entry.name),
			logz.WorkspaceTerminationProgress(string(terminationProgress)),
		)

		key := workspaceInfoKey{
			Name:      entry.name,
			Namespace: entry.namespace,
		}

		var agentInfo WorkspaceAgentInfo
		var exists bool

		agentInfo, exists = payload[key]
		if !exists {
			agentInfo = WorkspaceAgentInfo{
				Name:      entry.name,
				Namespace: entry.namespace,
			}
		}

		agentInfo.TerminationProgress = terminationProgress
		payload[key] = agentInfo
	}
}

func (r *reconciler) performWorkspaceUpdateRequestToRailsApi(
	ctx context.Context,
	updateType WorkspaceUpdateType,
	workspaceAgentInfos []WorkspaceAgentInfo,
) (workspaceRailsInfos []*WorkspaceRailsInfo, retError error) {
	// Do the POST request to the Rails API
	r.log.Debug("Making GitLab request")

	if workspaceAgentInfos == nil {
		// In case there is nothing to report to rails, populate the payload with an empty slice
		// to be explicit about the intent
		// TODO: add a test case - https://gitlab.com/gitlab-org/gitlab/-/issues/407554
		workspaceAgentInfos = []WorkspaceAgentInfo{}
	}

	startTime := time.Now()
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
	r.log.Debug("Made request to the Rails API",
		logz.StatusCode(resp.StatusCode),
		logz.RequestId(resp.Header.Get(httpz.RequestIdHeader)),
		logz.DurationInMilliseconds(time.Since(startTime)),
	)

	defer errz.SafeClose(resp.Body, &retError)
	if resp.StatusCode != http.StatusCreated {
		_ = ioz.DiscardData(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var responsePayload ResponsePayload
	err = json.Unmarshal(body, &responsePayload)
	if err != nil {
		return nil, fmt.Errorf("error parsing response body: %w", err)
	}

	r.log.Debug(
		"Read body from the Rails API",
		logz.PayloadSizeInBytes(len(body)),
		logz.WorkspaceDataCount(len(responsePayload.WorkspaceRailsInfos)),
	)

	return responsePayload.WorkspaceRailsInfos, nil
}

func (r *reconciler) handleDesiredStateIsTerminated(ctx context.Context, name string, namespace string, actualState string) error {
	if r.terminatingTracker.isTerminating(name, namespace) && actualState == WorkspaceStateTerminated {
		r.log.Debug("ActualState=Terminated, so deleting it from persistedTerminatingWorkspacesTracker", logz.WorkspaceNamespace(namespace))
		r.terminatingTracker.delete(name, namespace)
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

	r.terminatingTracker.add(name, namespace)

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
