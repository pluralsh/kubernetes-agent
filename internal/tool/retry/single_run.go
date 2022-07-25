package retry

import "sync/atomic"

const (
	notRunning int32 = 0
	running    int32 = 1
)

type SingleRun struct {
	isRunning int32
}

// Run runs f in a goroutine if it's not running already i.e. was never ran or has returned.
func (s *SingleRun) Run(f func()) {
	if !atomic.CompareAndSwapInt32(&s.isRunning, notRunning, running) {
		// Already running.
		return
	}
	// Not running, run it.
	go func() {
		defer atomic.StoreInt32(&s.isRunning, notRunning)
		f()
	}()
}

// IsRunning returns whether a goroutine is running at the moment.
// Normally this method should not be used as it's inherently racy, but can be useful in tests.
func (s *SingleRun) IsRunning() bool {
	return atomic.LoadInt32(&s.isRunning) == running
}
