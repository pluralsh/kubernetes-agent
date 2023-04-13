package informer

import "context"

type Informer interface {
	Start(ctx context.Context) error
	List() []*ParsedWorkspace
}

/*
ParsedWorkspace is our internal view on the full unstructured.Unstructured k8s resource.
We use ParsedWorkspace for logic like checking if the latest change has been persisted
or whether we terminated the workspace
*/
type ParsedWorkspace struct {
	Name              string
	Namespace         string
	ResourceVersion   string
	K8sDeploymentInfo map[string]interface{}
}
