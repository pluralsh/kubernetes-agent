package agent

// TODO: This could be improved to be a map of "Terminating|Terminated" values instead of just a
//
//	 issue: https://gitlab.com/gitlab-org/gitlab/-/issues/402758
//		boolean. This would allow us to send an ActualState of Terminating and have the Rails
//		UI be able to display Terminating instead of the last persisted actual_state.
type persistedTerminatedWorkspacesTracker map[string]struct{}

func newPersistedTerminatedWorkspacesTracker() persistedTerminatedWorkspacesTracker {
	return make(map[string]struct{})
}

// returns true if the workspaceName is in the map
func (p persistedTerminatedWorkspacesTracker) isTerminated(workspaceName string) bool {
	_, ok := p[workspaceName]
	return ok
}

func (p persistedTerminatedWorkspacesTracker) add(workspaceName string) {
	p[workspaceName] = struct{}{}
}

func (p persistedTerminatedWorkspacesTracker) delete(workspaceName string) {
	delete(p, workspaceName)
}
