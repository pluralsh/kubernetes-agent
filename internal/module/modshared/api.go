package modshared

import (
	"context"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/retry"
	"go.uber.org/zap"
)

const (
	NoAgentId int64 = 0
)

// Api provides the API for the module to use.
type Api interface {
	// HandleProcessingError can be used to handle errors occurring while processing a request.
	// If err is a (or wraps a) errz.UserError, it might be handled specially.
	HandleProcessingError(ctx context.Context, log *zap.Logger, agentId int64, msg string, err error)
}

// RpcApi provides the API for the module's gRPC handlers to use.
type RpcApi interface {
	// Log returns a logger to use in the context of the request being processed.
	Log() *zap.Logger
	// HandleProcessingError can be used to handle errors occurring while processing a request.
	// If err is a (or wraps a) errz.UserError, it might be handled specially.
	HandleProcessingError(log *zap.Logger, agentId int64, msg string, err error)
	// HandleIoError can be used to handle I/O error produced by gRPC Send(), Recv() methods or any other I/O error.
	// It returns an error, compatible with gRPC status package.
	HandleIoError(log *zap.Logger, msg string, err error) error
	// PollWithBackoff runs f every duration given by BackoffManager.
	//
	// PollWithBackoff should be used by the top-level polling, so that it can be gracefully interrupted
	// by the server when necessary. E.g. when stream is nearing it's max connection age or program needs to
	// be shut down.
	// If sliding is true, the period is computed after f runs. If it is false then
	// period includes the runtime for f.
	// It returns when:
	// - stream's context is cancelled or max connection age has been reached. nil is returned in this case.
	// - f returns Done. error from f is returned in this case.
	PollWithBackoff(cfg retry.PollConfig, f retry.PollWithBackoffFunc) error
}
