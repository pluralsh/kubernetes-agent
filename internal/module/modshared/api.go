package modshared

import (
	"context"

	"gitlab.com/gitlab-org/labkit/errortracking"
	"go.uber.org/zap"
)

const (
	NoAgentId int64 = 0

	AgentIdErrTrackingField = "gitlab.AgentId"
)

// Api provides the API for the module to use.
type Api interface {
	errortracking.Tracker
	// HandleProcessingError can be used to handle errors occurring while processing a request.
	// If err is a (or wraps a) errz.UserError, it might be handled specially.
	HandleProcessingError(ctx context.Context, log *zap.Logger, agentId int64, msg string, err error)
}
