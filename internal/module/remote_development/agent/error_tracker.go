package agent

import (
	"context"
	"errors"
	"sync"

	"k8s.io/apimachinery/pkg/util/wait"
)

type errorTracker struct {
	mx    sync.Mutex
	store map[errorTrackerKey]operationState
	wg    wait.Group
}

type errorTrackerKey struct {
	name      string
	namespace string
}

// operationState indicates the state of an async operation as it is being watched by the errorTracker.
type operationState struct {
	// version indicates the version of the operation for which the state is being tracked
	version uint64

	// err contains the error available at the end of an async operation
	err error
}

func newErrorTracker() *errorTracker {
	return &errorTracker{store: make(map[errorTrackerKey]operationState)}
}

func (t *errorTracker) deleteErrorIfVersion(name string, namespace string, version uint64) {
	key := errorTrackerKey{
		name:      name,
		namespace: namespace,
	}

	t.mx.Lock()
	defer t.mx.Unlock()

	existingState, exists := t.store[key]

	if exists && existingState.version == version {
		// only delete entries if they correspond to the version passed in the
		// function call. This check prevents cases where entries for a more
		// recent operation(with a higher version) may be overwritten with
		// error details of an older operation(with a lower version)
		delete(t.store, key)
	}
}

// saveErrorIfVersion will record error state only if the provided version
// matches the version in an existing record. In every other case, nothing
// will be written
func (t *errorTracker) saveErrorIfVersion(name string, namespace string, err error, version uint64) {
	key := errorTrackerKey{
		name:      name,
		namespace: namespace,
	}

	t.mx.Lock()
	defer t.mx.Unlock()

	existingState, exists := t.store[key]
	if !exists {
		// Do NOT write anything if nothing exists
		return
	}

	if existingState.version != version {
		// this check is added so to prevent overwrite of errors
		// for entries with mismatched versions
		return
	}
	t.store[key] = operationState{
		version: existingState.version,
		err:     err,
	}
}

// watchForLatestErrors will watch the provided non-nil channel for errors and record them asynchronously. If multiple watches are created for
// different versions for the same workspace & namespace, only the errors corresponding to the latest version are tracked while entries for earlier
// versions are discarded. If multiple errors are published in the channel, they are all collected first and then stored in the tracker
func (t *errorTracker) watchForLatestErrors(ctx context.Context, name string, namespace string, version uint64, errorCh <-chan error) {
	t.markEntryWithVersion(name, namespace, version)

	t.wg.StartWithContext(ctx, func(ctx context.Context) {
		var allErrors []error
		for err := range errorCh {
			allErrors = append(allErrors, err)
		}

		err := errors.Join(allErrors...)
		if err != nil {
			// at least one error was received on the channel, so it must be saved
			t.saveErrorIfVersion(name, namespace, err, version)
		} else {
			// no error was received on the channel and so the
			// entry for provided version can be safely evicted
			t.deleteErrorIfVersion(name, namespace, version)
		}
	})
}

// waitForErrors waits on existing error channels being watched to finish publishing errors(if any)
func (t *errorTracker) waitForErrors() {
	t.wg.Wait()
}

// markEntryWithVersion will create an entry in the store if and only if
// no entry exists for the particular key OR an entry exists with an older version.
// If the version of an existing entry is higher (not equal) than the version passed in the function
// call, then the writes are skipped and nothing is updated
func (t *errorTracker) markEntryWithVersion(workspace string, namespace string, version uint64) {
	key := errorTrackerKey{
		name:      workspace,
		namespace: namespace,
	}

	t.mx.Lock()
	defer t.mx.Unlock()

	existingState, exists := t.store[key]

	if exists && existingState.version > version {
		// writes should be skipped here as an older version MUST not
		// overwrite the entries corresponding to a newer version
		return
	}

	t.store[key] = operationState{
		version: version,
		err:     nil,
	}
}

func (t *errorTracker) createSnapshot() map[errorTrackerKey]operationState {
	snapshot := make(map[errorTrackerKey]operationState)

	t.mx.Lock()
	defer t.mx.Unlock()

	for key, state := range t.store {
		snapshot[key] = state
	}

	return snapshot
}
