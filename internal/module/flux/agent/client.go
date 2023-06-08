package agent

import (
	"context"
	"fmt"
	"io"
	"reflect"
	"time"

	notificationv1 "github.com/fluxcd/notification-controller/api/v1"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/flux/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/retry"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/syncz"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
)

const (
	// projectReceiverIndex the name for the index on the FluxCD notification/Receiver objects
	// that maps them to GitLab project paths.
	projectReceiverIndex                  = "project"
	defaultReconciliationDebounceDuration = 10 * time.Second
)

// client represents the client part of the Flux agent module.
type client struct {
	log                            *zap.Logger
	agentApi                       modagent.Api
	agentId                        int64
	fluxGitLabClient               rpc.GitLabFluxClient
	pollCfgFactory                 retry.PollConfigFactory
	receiverIndexer                cache.Indexer
	reconcileTrigger               reconcileTrigger
	updateProjectsToReconcileC     chan []string
	reconciliationDebounceDuration time.Duration
}

type clientFactory func(ctx context.Context, url string, receiverIndexer cache.Indexer) (*client, error)

// projectReconciler is an interface to start reconciliation
// of projects available to an underlying indexer
type projectReconciler interface {
	// ReconcileIndexedProjects starts the reconciliation of whatever
	// projects are in the underlying index.
	// This method only starts the reconciliation - which is an asynchronous process
	// and effectively is run in another goroutine.
	ReconcileIndexedProjects(ctx context.Context)
}

// newClient adds an index on the given indexer and returns a new client.
// The added index maps GitLab project paths to receiver objects, see createProjectIndex
func newClient(log *zap.Logger, agentApi modagent.Api, agentId int64, fluxGitLabClient rpc.GitLabFluxClient, pollCfgFactory retry.PollConfigFactory, receiverIndexer cache.Indexer, reconcileTrigger reconcileTrigger) (*client, error) {
	err := addProjectIndex(receiverIndexer)
	if err != nil {
		return nil, err
	}
	updateProjectsToReconcileC := make(chan []string)
	return &client{
		log:                            log,
		agentApi:                       agentApi,
		agentId:                        agentId,
		fluxGitLabClient:               fluxGitLabClient,
		pollCfgFactory:                 pollCfgFactory,
		receiverIndexer:                receiverIndexer,
		reconcileTrigger:               reconcileTrigger,
		updateProjectsToReconcileC:     updateProjectsToReconcileC,
		reconciliationDebounceDuration: defaultReconciliationDebounceDuration,
	}, nil
}

// addProjectIndex adds a new index to the given indexer.
// The added index maps GitLab project paths (like `gitlab-org/gitlab`) to FluxCD notification/Receiver objects,
// which in turn trace back to FluxCD source/GitRepository objects.
// The index is created based on the projectAnnotationKey annotation on the notification/Receiver object
// which are being indexed.
func addProjectIndex(receiverIndexer cache.Indexer) error {
	err := receiverIndexer.AddIndexers(map[string]cache.IndexFunc{
		projectReceiverIndex: func(obj interface{}) ([]string, error) {
			u, ok := obj.(*unstructured.Unstructured)
			if !ok {
				return nil, fmt.Errorf("failed to cast object of type %T into *unstructured.Unstructured", obj)
			}

			project, ok := u.GetAnnotations()[projectAnnotationKey]
			if !ok {
				// NOTE: this is not an issue at this point, because it may very well be that this
				// receiver doesn't (yet) have the annotation. This function will eventually
				// be called again once it has it.
				return nil, nil
			}

			return []string{project}, nil
		},
	})
	if err != nil {
		return fmt.Errorf("unable to add %s indexer: %w", projectReceiverIndex, err)
	}
	return nil
}

// RunProjectReconciliation runs a new reconciliation process
// for the latest projects.
// A reconciliation for a new set of projects can be started using
// the ReconcileIndexedProjects method.
// There is only ever one reconciliation process running at the same
// time and a call to ReconcileIndexedProjects will terminate
// a potentially previously started reconciliation process.
// The c.reconciliationDebounceDuration controls how long to debounce
// before starting a new reconciliation for received projects.
func (c *client) RunProjectReconciliation(ctx context.Context) {
	done := ctx.Done()

	wh := syncz.NewWorkerHolder[[]string](
		func(projectsToReconcile []string) syncz.Worker {
			return syncz.WorkerFunc(func(ctx context.Context) {
				c.reconcileProjects(ctx, projectsToReconcile)
			})
		},
		isEqualProjectSets,
	)
	defer wh.StopAndWait()

	debounceTimer := time.NewTimer(time.Hour)
	debounceTimer.Stop()
	defer debounceTimer.Stop()

	var lastProjects []string

	for {
		select {
		case <-done:
			// Shutdown
			return // nolint:govet
		case projects := <-c.updateProjectsToReconcileC:
			lastProjects = projects
			if !debounceTimer.Stop() {
				select {
				case <-debounceTimer.C:
				default:
				}
			}
			debounceTimer.Reset(c.reconciliationDebounceDuration)
		case <-debounceTimer.C:
			wh.ApplyConfig(ctx, lastProjects)
		}
	}
}

// ReconcileIndexedProjects starts the reconciliation of the latest indexed projects.
// The index must be created using createProjectIndex using the newClient factory function.
// This method only *starts* the reconciliation, but actual process runs within
// RunProjectReconciliation.
func (c *client) ReconcileIndexedProjects(ctx context.Context) {
	projects := c.receiverIndexer.ListIndexFuncValues(projectReceiverIndex)
	c.log.Debug("Reconcile project update", logz.ProjectsToReconcile(projects))

	select {
	case <-ctx.Done():
	case c.updateProjectsToReconcileC <- projects:
	}
}

// isEqualProjectSets returns true if the given project sets are equal.
// The order and possible duplicates don't matter.
func isEqualProjectSets(projects1, projects2 []string) bool {
	uniqueProjects := func(ps []string) map[string]struct{} {
		us := make(map[string]struct{}, len(ps))
		for _, p := range ps {
			us[p] = struct{}{}
		}
		return us
	}
	ux := uniqueProjects(projects1)
	uy := uniqueProjects(projects2)
	return reflect.DeepEqual(ux, uy)
}

// reconcileProjects makes an API call to the server to wait for reconciliation updates of a set of projects.
// Once one of these projects is updated it triggers the associated FluxCD notification/Receiver webhook.
func (c *client) reconcileProjects(ctx context.Context, projects []string) {
	c.log.Debug("Started watching projects for reconciliation", logz.ProjectsToReconcile(projects))
	defer c.log.Debug("Stopped watching projects for reconciliation", logz.ProjectsToReconcile(projects))

	_ = retry.PollWithBackoff(ctx, c.pollCfgFactory(), func(ctx context.Context) (error, retry.AttemptResult) {
		rpcClient, err := c.fluxGitLabClient.ReconcileProjects(ctx, &rpc.ReconcileProjectsRequest{Project: rpc.ReconcileProjectsFromSlice(projects)})
		if err != nil {
			c.agentApi.HandleProcessingError(ctx, c.log, c.agentId, "Failed to reconcile projects", err)
			return nil, retry.Backoff
		}

		for {
			resp, err := rpcClient.Recv()
			if err != nil {
				if err == io.EOF { // nolint:errorlint
					// server closed connection, retrying
					return nil, retry.ContinueImmediately
				}
				if grpctool.RequestCanceled(err) {
					// request was canceled
					c.log.Debug("ReconcileProjects request has been canceled, backing off and awaiting cancellation")
					return nil, retry.Backoff
				}
				c.agentApi.HandleProcessingError(ctx, c.log, c.agentId, "Failed to receive project to reconcile", err)
				return nil, retry.Backoff
			}

			c.reconcileProject(ctx, resp.Project.Id)
		}
	})
}

// reconcileProject reconciles a single project by triggering the FluxCD notification/Receiver webhooks
// indexed with this project.
// The Receiver object must have:
// - a projectAnnotationKey that matches the given project
// - a webhook path
// ... in order to be triggered.
func (c *client) reconcileProject(ctx context.Context, project string) {
	log := c.log.With(logz.ProjectId(project))
	objs, err := c.receiverIndexer.ByIndex(projectReceiverIndex, project)
	if err != nil {
		log.Error("Unable to get Receivers for project", logz.Error(err))
		return
	}

	for _, obj := range objs {
		u := obj.(*unstructured.Unstructured)
		var nr notificationv1.Receiver
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), &nr)
		if err != nil {
			log.Error("Unable to convert unstructured object to Receiver", logz.Error(err))
			continue
		}

		if p := nr.Annotations[projectAnnotationKey]; project != p || nr.Status.WebhookPath == "" {
			continue
		}

		err = c.reconcileTrigger.reconcile(ctx, nr.Status.WebhookPath)
		if err != nil {
			log.Error("Unable to trigger Receiver", logz.Error(err))
			continue
		}
	}
}
