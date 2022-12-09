package errz

import (
	"context"
	"errors"
	"io"

	"go.uber.org/zap"
)

// ErrReporter provides a way to report errors.
type ErrReporter interface {
	// HandleProcessingError can be used to handle errors occurring while processing a request.
	// If err is a (or wraps a) errz.UserError, it might be handled specially.
	HandleProcessingError(ctx context.Context, log *zap.Logger, msg string, err error)
}

func ContextDone(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}

func SafeClose(toClose io.Closer, err *error) {
	if toClose == nil {
		return
	}
	SafeCall(toClose.Close, err)
}

func SafeCall(toCall func() error, err *error) {
	e := toCall()
	if *err == nil {
		*err = e
	}
}
