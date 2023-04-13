package agent

import (
	"context"

	"k8s.io/apimachinery/pkg/util/wait"
)

type stoppableTask interface {
	StopAndWait()
}

// TODO: enhance this to be able to capture and return error from the task being run
//
//	issue: https://gitlab.com/gitlab-org/gitlab/-/issues/404773
type simpleStoppableTask struct {
	cancel context.CancelFunc
	wg     wait.Group
}

func (t *simpleStoppableTask) StopAndWait() {
	t.cancel()
	t.wg.Wait()
}

func newStoppableTask(ctx context.Context, fn func(ctx context.Context)) stoppableTask {
	taskCtx, cancel := context.WithCancel(ctx)

	task := &simpleStoppableTask{
		cancel: cancel,
	}

	task.wg.StartWithContext(taskCtx, fn)

	return task
}
