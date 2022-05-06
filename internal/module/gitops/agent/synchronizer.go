package agent

import (
	"context"
	"fmt"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/gitops/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/retry"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/pkg/agentcfg"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/cli-runtime/pkg/resource"
	"sigs.k8s.io/cli-utils/pkg/apply"
	"sigs.k8s.io/cli-utils/pkg/common"
	"sigs.k8s.io/cli-utils/pkg/inventory"
)

// synchronizerConfig holds configuration for a synchronizer.
type synchronizerConfig struct {
	log               *zap.Logger
	agentId           int64
	project           *agentcfg.ManifestProjectCF
	applier           Applier
	restMapper        meta.RESTMapper
	restClientGetter  resource.RESTClientGetter
	applierPollConfig retry.PollConfig
	applyOptions      apply.ApplierOptions
}

type synchronizer struct {
	synchronizerConfig
}

func newSynchronizer(config synchronizerConfig) *synchronizer {
	return &synchronizer{
		synchronizerConfig: config,
	}
}

func (s *synchronizer) run(desiredState <-chan rpc.ObjectsToSynchronizeData) {
	jobs := make(chan syncJob)
	sw := syncWorker{
		log:               s.log,
		applier:           s.applier,
		applierPollConfig: s.applierPollConfig,
		applyOptions:      s.applyOptions,
	}
	var wg wait.Group
	defer wg.Wait()   // Wait for sw to exit
	defer close(jobs) // Close jobs to signal sw there is no more work to be done
	wg.Start(func() {
		sw.Run(jobs) // Start sw
	})

	var (
		jobsCh    chan syncJob
		newJob    syncJob
		jobCancel context.CancelFunc
	)
	defer func() {
		if jobCancel != nil {
			jobCancel()
		}
	}()

	d := syncDecoder{
		restMapper:       s.restMapper,
		restClientGetter: s.restClientGetter,
		defaultNamespace: s.project.DefaultNamespace,
	}

	for {
		select {
		case state, ok := <-desiredState:
			if !ok {
				return // nolint: govet
			}
			objs, err := d.Decode(state.Sources)
			if err != nil {
				s.log.Error("Failed to decode GitOps objects", logz.Error(err), logz.CommitId(state.CommitId))
				continue
			}
			invObj, objs, err := s.splitObjects(state.ProjectId, objs)
			if err != nil {
				s.log.Error("Failed to locate inventory object in GitOps objects", logz.Error(err), logz.CommitId(state.CommitId))
				continue
			}
			if jobCancel != nil {
				jobCancel() // Cancel running/pending job ASAP
			}
			newJob = syncJob{
				commitId: state.CommitId,
				invInfo:  inventory.WrapInventoryInfoObj(invObj),
				objects:  objs,
			}
			newJob.ctx, jobCancel = context.WithCancel(context.Background()) // nolint: govet
			jobsCh = jobs                                                    // Enable select case
		case jobsCh <- newJob: // Try to send new job to syncWorker. This case is active only when jobsCh is not nil
			// Success!
			newJob = syncJob{} // Erase contents to help GC
			jobsCh = nil       // Disable this select case (send to nil channel blocks forever)
		}
	}
}

func (s *synchronizer) splitObjects(projectId int64, objs []*unstructured.Unstructured) (*unstructured.Unstructured, []*unstructured.Unstructured, error) {
	invs := make([]*unstructured.Unstructured, 0, 1)
	resources := make([]*unstructured.Unstructured, 0, len(objs))
	for _, obj := range objs {
		if inventory.IsInventoryObject(obj) {
			invs = append(invs, obj)
		} else {
			resources = append(resources, obj)
		}
	}
	switch len(invs) {
	case 0:
		return s.defaultInventoryObjTemplate(projectId), resources, nil
	case 1:
		return invs[0], resources, nil
	default:
		return nil, nil, fmt.Errorf("expecting zero or one inventory object, found %d", len(invs))
	}
}

func (s *synchronizer) defaultInventoryObjTemplate(projectId int64) *unstructured.Unstructured {
	id := inventoryId(s.agentId, projectId)
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"name":      "inventory-" + id,
				"namespace": s.project.DefaultNamespace,
				"labels": map[string]interface{}{
					common.InventoryLabel: id,
				},
			},
		},
	}
}

func inventoryId(agentId, projectId int64) string {
	return fmt.Sprintf("%d-%d", agentId, projectId)
}
