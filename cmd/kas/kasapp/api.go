package kasapp

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/errz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/logz"
	"gitlab.com/gitlab-org/labkit/errortracking"
	"go.uber.org/zap"
)

type serverApi struct {
	ErrorTracker errortracking.Tracker
}

func (a *serverApi) HandleProcessingError(ctx context.Context, log *zap.Logger, agentId int64, msg string, err error) {
	handleProcessingError(ctx, a.ErrorTracker, log, agentId, msg, err)
}

func handleProcessingError(ctx context.Context, errTracker errortracking.Tracker, log *zap.Logger, agentId int64, msg string, err error) {
	if grpctool.RequestCanceled(err) {
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
		logAndCapture(ctx, errTracker, log, agentId, msg, err)
	}
}

func logAndCapture(ctx context.Context, errTracker errortracking.Tracker, log *zap.Logger, agentId int64, msg string, err error) {
	// don't add logz.CorrelationIdFromContext() or logz.AgentId() here as they've been added to the logger already
	log.Error(msg, logz.Error(err))
	opts := []errortracking.CaptureOption{errortracking.WithContext(ctx)}
	if agentId != modshared.NoAgentId {
		opts = append(opts, errortracking.WithField(modshared.AgentIdErrTrackingField, strconv.FormatInt(agentId, 10)))
	}
	errTracker.Capture(fmt.Errorf("%s: %w", msg, err), opts...)
}
