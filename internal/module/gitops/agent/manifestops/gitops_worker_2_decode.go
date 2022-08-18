package manifestops

import (
	"context"
	"fmt"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/gitops/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/logz"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/cli-utils/pkg/common"
	"sigs.k8s.io/cli-utils/pkg/inventory"
)

func (w *defaultGitopsWorker) decode(desiredState <-chan rpc.ObjectsToSynchronizeData, jobs chan<- applyJob) {
	var jobCancel context.CancelFunc
	defer func() {
		if jobCancel != nil {
			jobCancel()
		}
	}()

	d := syncDecoder{
		restClientGetter: w.restClientGetter,
		defaultNamespace: w.project.DefaultNamespace,
	}

	p := retryPipeline{
		inputCh:      desiredState,
		outputCh:     jobs,
		retryBackoff: w.decodeRetryPolicy,
		process: func(input inputT) (outputT, processResult) {
			objs, err := d.Decode(input.Sources)
			if err != nil {
				w.log.Error("Failed to decode GitOps objects", logz.Error(err), logz.CommitId(input.CommitId))
				return outputT{}, backoff
			}
			invObj, objs, err := w.splitObjects(input.ProjectId, objs)
			if err != nil {
				w.log.Error("Failed to locate inventory object in GitOps objects", logz.Error(err), logz.CommitId(input.CommitId))
				return outputT{}, done
			}
			if jobCancel != nil {
				jobCancel() // Cancel running/pending job ASAP
			}
			newJob := applyJob{
				commitId: input.CommitId,
				invInfo:  inventory.WrapInventoryInfoObj(invObj),
				objects:  objs,
			}
			newJob.ctx, jobCancel = context.WithCancel(context.Background()) // nolint: govet
			return newJob, success
		},
	}
	p.run()
}

func (w *defaultGitopsWorker) splitObjects(projectId int64, objs []*unstructured.Unstructured) (*unstructured.Unstructured, []*unstructured.Unstructured, error) {
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
		return w.defaultInventoryObjTemplate(projectId), resources, nil
	case 1:
		return invs[0], resources, nil
	default:
		return nil, nil, fmt.Errorf("expecting zero or one inventory object, found %d", len(invs))
	}
}

func (w *defaultGitopsWorker) defaultInventoryObjTemplate(projectId int64) *unstructured.Unstructured {
	id := inventoryId(w.agentId, projectId)
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"name":      "inventory-" + id,
				"namespace": w.project.DefaultNamespace,
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
