package kasapp

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/getsentry/sentry-go"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/errz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/logz"
	"gitlab.com/gitlab-org/labkit/correlation"
	"go.uber.org/zap"
)

type SentryHub interface {
	CaptureEvent(event *sentry.Event) *sentry.EventID
}

type serverApi struct {
	Hub SentryHub
}

func (a *serverApi) HandleProcessingError(ctx context.Context, log *zap.Logger, agentId int64, msg string, err error) {
	handleProcessingError(ctx, a.Hub, log, agentId, msg, err)
}

func handleProcessingError(ctx context.Context, hub SentryHub, log *zap.Logger, agentId int64, msg string, err error) {
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
		logAndCapture(ctx, hub, log, agentId, msg, err)
	}
}

func logAndCapture(ctx context.Context, hub SentryHub, log *zap.Logger, agentId int64, msg string, err error) {
	log.Error(msg, logz.Error(err))
	event := sentry.NewEvent()
	if agentId != modshared.NoAgentId {
		event.User.ID = strconv.FormatInt(agentId, 10)
	}
	event.Level = sentry.LevelError
	event.Exception = []sentry.Exception{
		{
			Type:       reflect.TypeOf(err).String(),
			Value:      fmt.Sprintf("%s: %v", msg, err),
			Stacktrace: sentry.ExtractStacktrace(err),
		},
	}
	correlationID := correlation.ExtractFromContext(ctx)
	if correlationID != "" {
		event.Tags[modserver.CorrelationIdSentryField] = correlationID
	}
	hub.CaptureEvent(event)
}
