package agent

// persistedTerminatingWorkspacesTracker is a set to track workspaces that exist
// in Terminating state
type persistedTerminatingWorkspacesTracker map[string]struct{}

func newPersistedTerminatingWorkspacesTracker() persistedTerminatingWorkspacesTracker {
	return make(map[string]struct{})
}

func (p persistedTerminatingWorkspacesTracker) isTerminating(workspaceName string) bool {
	_, ok := p[workspaceName]
	return ok
}

func (p persistedTerminatingWorkspacesTracker) add(workspaceName string) {
	p[workspaceName] = struct{}{}
}

func (p persistedTerminatingWorkspacesTracker) delete(workspaceName string) {
	delete(p, workspaceName)
}
