package agent

// terminatingWorkspacesTrackerKey is used as a key within persistedTerminatingWorkspacesTracker
// to uniquely identify a combination of workspace name and namespace that must be tracked
type terminatingWorkspacesTrackerKey struct {
	name      string
	namespace string
}

// persistedTerminatingWorkspacesTracker is a set to track workspaces that exist
// in Terminating state
type persistedTerminatingWorkspacesTracker map[terminatingWorkspacesTrackerKey]struct{}

func newPersistedTerminatingWorkspacesTracker() persistedTerminatingWorkspacesTracker {
	return make(map[terminatingWorkspacesTrackerKey]struct{})
}

func (p persistedTerminatingWorkspacesTracker) isTerminating(name string, namespace string) bool {
	key := terminatingWorkspacesTrackerKey{
		name:      name,
		namespace: namespace,
	}
	_, ok := p[key]
	return ok
}

func (p persistedTerminatingWorkspacesTracker) add(name string, namespace string) {
	key := terminatingWorkspacesTrackerKey{
		name:      name,
		namespace: namespace,
	}
	p[key] = struct{}{}
}

func (p persistedTerminatingWorkspacesTracker) delete(name string, namespace string) {
	key := terminatingWorkspacesTrackerKey{
		name:      name,
		namespace: namespace,
	}
	delete(p, key)
}
