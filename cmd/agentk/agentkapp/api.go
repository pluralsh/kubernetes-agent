package agentkapp

import (
	"context"
	"errors"
	"io"
	"net/url"
	"sync"

	"github.com/pluralsh/kuberentes-agent/internal/tool/errz"
	"github.com/pluralsh/kuberentes-agent/internal/tool/grpctool"
	"github.com/pluralsh/kuberentes-agent/internal/tool/logz"
	"go.uber.org/zap"
)

// agentAPI is an implementation of modagent.API.
type agentAPI struct {
	moduleName        string
	agentId           *ValueHolder[int64]
	gitLabExternalUrl *ValueHolder[url.URL]
}

func (a *agentAPI) GetAgentId(ctx context.Context) (int64, error) {
	return a.agentId.get(ctx)
}

func (a *agentAPI) GetGitLabExternalUrl(ctx context.Context) (url.URL, error) {
	return a.gitLabExternalUrl.get(ctx)
}

func (a *agentAPI) TryGetAgentId() (int64, bool) {
	return a.agentId.tryGet()
}

func (a *agentAPI) HandleProcessingError(ctx context.Context, log *zap.Logger, agentId int64, msg string, err error) {
	handleProcessingError(ctx, log, agentId, msg, err)
}

func handleProcessingError(ctx context.Context, log *zap.Logger, agentId int64, msg string, err error) { // nolint:unparam
	if grpctool.RequestCanceledOrTimedOut(err) {
		// An error caused by context signalling done
		return
	}
	var ue errz.UserError
	isUserError := errors.As(err, &ue)
	if isUserError {
		// TODO Don't log it, send it somewhere the user can see it https://gitlab.com/gitlab-org/gitlab/-/issues/277323
		// Log at Info for now.
		log.Info(msg, logz.Error(err))
	} else {
		// don't add logz.TraceIdFromContext(ctx) here as it's been added to the logger already
		log.Error(msg, logz.Error(err))
	}
}

type cancelingReadCloser struct {
	io.ReadCloser
	cancel context.CancelFunc
}

func (c cancelingReadCloser) Close() error {
	// First close the pipe, then cancel the context.
	// Doing it the other way around can result in another goroutine calling CloseWithError() with the
	// cancellation error.
	err := c.ReadCloser.Close()
	c.cancel()
	return err
}

// valueOrError holds a value or an error. The first value/error is persisted, rest are discarded.
// onError is a callback that is called on each SetError() invocation, regardless of whether a value or an error
// has already been set.
// Thread safe.
type valueOrError[T any] struct {
	mu      sync.Mutex
	locker  *sync.Cond
	onError func(error) error
	value   T
	err     error
	isSet   bool
}

func newValueOrError[T any](onError func(error) error) *valueOrError[T] {
	v := &valueOrError[T]{
		onError: onError,
	}
	v.locker = sync.NewCond(&v.mu)
	return v
}

func (v *valueOrError[T]) SetValue(value T) {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.isSet {
		return
	}
	v.isSet = true
	v.value = value
	v.locker.Broadcast()
}

func (v *valueOrError[T]) SetError(err error) {
	v.mu.Lock()
	defer v.mu.Unlock()
	err = v.onError(err)
	if v.isSet {
		return
	}
	v.isSet = true
	v.err = err
	v.locker.Broadcast()
}

// Wait returns a value or an error, blocking the caller until one of them is set.
func (v *valueOrError[T]) Wait() (T, error) {
	v.mu.Lock()
	defer v.mu.Unlock()
	if !v.isSet {
		v.locker.Wait()
	}
	return v.value, v.err
}

type onceReadCloser struct {
	io.ReadCloser
	once     sync.Once
	closeErr error
}

func (oc *onceReadCloser) Close() error {
	oc.once.Do(oc.close)
	return oc.closeErr
}

func (oc *onceReadCloser) close() {
	oc.closeErr = oc.ReadCloser.Close()
}
