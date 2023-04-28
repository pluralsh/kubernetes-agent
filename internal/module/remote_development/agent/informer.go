package agent

import "context"

type informer interface {
	Start(ctx context.Context) error
	List() []*parsedWorkspace
	Stop()
}

/*
parsedWorkspace is our internal view on the full unstructured.Unstructured k8s resource.
We use parsedWorkspace for logic like checking if the latest change has been persisted
or whether we terminated the workspace
*/
type parsedWorkspace struct {
	Name              string
	Namespace         string
	ResourceVersion   string
	K8sDeploymentInfo map[string]interface{}
}
