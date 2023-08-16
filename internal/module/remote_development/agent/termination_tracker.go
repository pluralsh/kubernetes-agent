package agent

// terminationTrackerKey is used as a key within terminationTracker
// to uniquely identify a combination of workspace name and namespace that must be tracked
type terminationTrackerKey struct {
	name      string
	namespace string
}

// terminationTracker is a set to track workspaces that exist in Terminating/Terminated state
type terminationTracker map[terminationTrackerKey]TerminationProgress

func newTerminationTracker() terminationTracker {
	return make(map[terminationTrackerKey]TerminationProgress)
}

func (t terminationTracker) add(name string, namespace string, progress TerminationProgress) {
	key := terminationTrackerKey{
		name:      name,
		namespace: namespace,
	}
	t[key] = progress
}

func (t terminationTracker) delete(name string, namespace string) {
	key := terminationTrackerKey{
		name:      name,
		namespace: namespace,
	}
	delete(t, key)
}
