package agent

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"
	rdutil "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/remote_development/agent/util"
)

type ErrorTrackerTestSuite struct {
	suite.Suite

	tracker *errorTracker

	testError      error
	testWorkspace  string
	testNamespace  string
	testTrackerKey errorTrackerKey
}

func TestRemoteDevModuleErrorTracker(t *testing.T) {
	suite.Run(t, new(ErrorTrackerTestSuite))
}

func (e *ErrorTrackerTestSuite) TestErrorTracking_ErrorIsPublished() {
	ctx := context.Background()

	asyncErr, unblockPublish := e.newBlockedAsyncError(e.testError)
	e.NotContains(e.tracker.createSnapshot(), e.testTrackerKey)

	v1 := uint64(1)
	e.tracker.watchForLatestErrors(ctx, e.testWorkspace, e.testNamespace, v1, asyncErr)

	// since publishing of error is blocked, verify state of tracker
	snapshot := e.tracker.createSnapshot()
	e.Contains(snapshot, e.testTrackerKey)

	actualState := snapshot[e.testTrackerKey]
	e.Equal(v1, actualState.version)
	e.NoError(actualState.err)

	// resume error publish, wait for it to be received and verify state of tracker again
	unblockPublish()
	e.tracker.waitForErrors()

	snapshot = e.tracker.createSnapshot()
	e.Contains(snapshot, e.testTrackerKey)

	actualState = snapshot[e.testTrackerKey]
	e.Equal(v1, actualState.version)
	e.ErrorIs(actualState.err, e.testError)
}

func (e *ErrorTrackerTestSuite) TestErrorTracking_NilError() {
	ctx := context.Background()

	// blocked value is created so that tracker state can be verified before publish
	asyncNilError, unblockPublish := e.newBlockedAsyncError(nil)

	e.NotContains(e.tracker.createSnapshot(), e.testTrackerKey)

	v1 := uint64(1)
	e.tracker.watchForLatestErrors(ctx, e.testWorkspace, e.testNamespace, v1, asyncNilError)

	// since publishing of nil is blocked, verify state of tracker
	snapshot := e.tracker.createSnapshot()
	e.Contains(snapshot, e.testTrackerKey)
	e.Equal(operationState{
		version: v1,
		err:     nil,
	}, snapshot[e.testTrackerKey])

	// unblock the async nil and verify final state of tracker
	unblockPublish()
	e.tracker.waitForErrors()

	snapshot = e.tracker.createSnapshot()
	e.NotContains(snapshot, e.testTrackerKey)
}

func (e *ErrorTrackerTestSuite) TestErrorTracking_NoErrors() {
	ctx := context.Background()

	asyncOp1 := rdutil.RunWithAsyncResult(func(_ chan<- error) {
		// no errors are published to the write-only channel
	})

	e.NotContains(e.tracker.createSnapshot(), e.testTrackerKey)

	v1 := uint64(1)
	e.tracker.watchForLatestErrors(ctx, e.testWorkspace, e.testNamespace, v1, asyncOp1)

	// wait for errors(if any) and verify final state of tracker
	e.tracker.waitForErrors()

	snapshot := e.tracker.createSnapshot()
	e.NotContains(snapshot, e.testTrackerKey)
}

func (e *ErrorTrackerTestSuite) TestErrorTracking_MultipleErrors() {
	ctx := context.Background()

	err1 := errors.New("applier error 1")
	err2 := errors.New("applier error 2")

	asyncOp1 := rdutil.RunWithAsyncResult(func(ch chan<- error) {
		ch <- err1
		ch <- err2
	})

	v1 := uint64(1)
	e.tracker.watchForLatestErrors(ctx, e.testWorkspace, e.testNamespace, v1, asyncOp1)
	e.tracker.waitForErrors()

	snapshot := e.tracker.createSnapshot()
	e.Contains(snapshot, e.testTrackerKey)

	storedState := snapshot[e.testTrackerKey]
	e.NotNil(storedState)
	e.Error(storedState.err)

	// verify that the received error is a combined error
	e.ErrorIs(storedState.err, err1)
	e.ErrorIs(storedState.err, err2)
}

func (e *ErrorTrackerTestSuite) TestErrorTracking_MultipleNilErrors() {
	ctx := context.Background()

	asyncOp1 := rdutil.RunWithAsyncResult(func(ch chan<- error) {
		ch <- nil
		ch <- nil
	})

	v1 := uint64(1)
	e.tracker.watchForLatestErrors(ctx, e.testWorkspace, e.testNamespace, v1, asyncOp1)
	e.tracker.waitForErrors()

	snapshot := e.tracker.createSnapshot()
	e.NotContains(snapshot, e.testTrackerKey)
}

// This verifies cases where the tracker watches for errors using a lower version
// after initiating a watch with a higher version for the same workspace / namespace combination
func (e *ErrorTrackerTestSuite) TestErrorTracking_OutOfOrderVersions() {
	ctx := context.Background()

	testError2 := errors.New("some other error")
	asyncOp2, unblockOp2 := e.newBlockedAsyncError(testError2)
	v2 := uint64(2)
	e.tracker.watchForLatestErrors(ctx, e.testWorkspace, e.testNamespace, v2, asyncOp2)

	asyncOp1, unblockOp1 := e.newBlockedAsyncError(e.testError)
	v1 := uint64(1)
	e.tracker.watchForLatestErrors(ctx, e.testWorkspace, e.testNamespace, v1, asyncOp1)

	// at this point, verify that only entry for version=2 is being tracked
	snapshot := e.tracker.createSnapshot()
	e.Len(snapshot, 1)
	e.Contains(snapshot, e.testTrackerKey)

	actualState := snapshot[e.testTrackerKey]
	e.Equal(v2, actualState.version)
	e.NoError(actualState.err)

	// unblock the publishing of async errors
	unblockOp1()
	unblockOp2()
	e.tracker.waitForErrors()

	snapshot = e.tracker.createSnapshot()
	e.Contains(snapshot, e.testTrackerKey)
	actualState = snapshot[e.testTrackerKey]
	e.Equal(v2, actualState.version)
	e.ErrorIs(actualState.err, testError2)
}

func (e *ErrorTrackerTestSuite) TestErrorTracking_ConflictingErrors() {
	ctx := context.Background()

	asyncOp1, unblockOp1 := e.newBlockedAsyncError(e.testError)

	v1 := uint64(1)
	e.tracker.watchForLatestErrors(ctx, e.testWorkspace, e.testNamespace, v1, asyncOp1)

	testError2 := errors.New("some other error")
	asyncOp2, unblockOp2 := e.newBlockedAsyncError(testError2)

	v2 := uint64(2)
	e.tracker.watchForLatestErrors(ctx, e.testWorkspace, e.testNamespace, v2, asyncOp2)

	// at this point, the track should stop tracking version=1
	snapshot := e.tracker.createSnapshot()
	e.Len(snapshot, 1)
	e.Contains(snapshot, e.testTrackerKey)
	actualState := snapshot[e.testTrackerKey]
	e.Equal(v2, actualState.version)
	e.NoError(actualState.err)

	// unblock the publishing of async errors
	unblockOp1()
	unblockOp2()
	e.tracker.waitForErrors()

	snapshot = e.tracker.createSnapshot()
	e.Contains(snapshot, e.testTrackerKey)

	actualState = snapshot[e.testTrackerKey]
	e.Equal(v2, actualState.version)
	e.ErrorIs(actualState.err, testError2)
}

func (e *ErrorTrackerTestSuite) TestErrorTracking_ConflictsWithEventualSuccess() {
	ctx := context.Background()

	asyncOp1, unblockOp1 := e.newBlockedAsyncError(e.testError)

	v1 := uint64(1)
	e.tracker.watchForLatestErrors(ctx, e.testWorkspace, e.testNamespace, v1, asyncOp1)

	asyncOp2, unblockOp2 := e.newBlockedAsyncError(nil)

	v2 := uint64(2)
	e.tracker.watchForLatestErrors(ctx, e.testWorkspace, e.testNamespace, v2, asyncOp2)

	// at this point, the track should stop tracking version=1
	snapshot := e.tracker.createSnapshot()
	e.Len(snapshot, 1)
	e.Contains(snapshot, e.testTrackerKey)
	actualState := snapshot[e.testTrackerKey]
	e.Equal(v2, actualState.version)
	e.NoError(actualState.err)

	// unblock the publishing of async errors
	unblockOp1()
	unblockOp2()
	e.tracker.waitForErrors()

	snapshot = e.tracker.createSnapshot()
	e.NotContains(snapshot, e.testTrackerKey)
	e.Len(snapshot, 0)
}

// newBlockedAsyncError returns a blocked channel that returns the passed error once the
// returned unblock function is invoked. The channel is closed automatically after.
func (e *ErrorTrackerTestSuite) newBlockedAsyncError(err error) (result <-chan error, unblock func()) {
	var wg sync.WaitGroup

	wg.Add(1)

	return rdutil.RunWithAsyncResult(func(ch chan<- error) {
		wg.Wait()

		ch <- err
	}), wg.Done
}

func (e *ErrorTrackerTestSuite) SetupTest() {
	e.tracker = newErrorTracker()

	e.testError = errors.New("applier error")
	e.testWorkspace = "test-workspace"
	e.testNamespace = "test-namespace"
	e.testTrackerKey = errorTrackerKey{
		name:      e.testWorkspace,
		namespace: e.testNamespace,
	}
}
