package agent

type persistedStateTracker struct {
	// TODO: The key in the map is the workspace name for now. When we are making this production ready
	// issue: https://gitlab.com/gitlab-org/gitlab/-/issues/402758
	//       we should revisit whether we want to make the logic more resilient by including namespace
	//       in the key (because resource is in k8s uniquely identified by the combination of namespace and name)
	persistedVersion map[string]string
}

func newPersistedStateTracker() *persistedStateTracker {
	return &persistedStateTracker{persistedVersion: make(map[string]string)}
}

// returns true if the recordVersion method has been called with the same versions
func (p *persistedStateTracker) isPersisted(name, deploymentResourceVersion string) bool {
	version, ok := p.persistedVersion[name]
	if !ok {
		return false
	}
	if deploymentResourceVersion != version {
		return false
	}
	return true
}

func (p *persistedStateTracker) recordVersion(wi *WorkspaceRailsInfo) {
	p.persistedVersion[wi.Name] = wi.DeploymentResourceVersion
}

// delete removes persisted workspace from memory
func (p *persistedStateTracker) delete(name string) {
	delete(p.persistedVersion, name)
}
