package kasapp

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"

	"github.com/getsentry/sentry-go"
	"github.com/redis/go-redis/v9"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modserver/notifications"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/errz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/logz"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

var (
	// Help deduplicate errors like:
	//    read tcp 10.222.67.20:40272->10.216.1.45:11443: read: connection reset by peer
	// by removing the random egress port.

	ipv4 = `(?:\d+\.){3}\d+`

	removePortStdLibError = regexp.MustCompile(`(` + ipv4 + `:)\d+(->` + ipv4 + `:\d+)`)

	_ notifications.Subscriber = &serverApi{}
)

type SentryHub interface {
	CaptureEvent(event *sentry.Event) *sentry.EventID
}

type serverApi struct {
	log         *zap.Logger
	Hub         SentryHub
	redisClient redis.UniversalClient
}

func (a *serverApi) HandleProcessingError(ctx context.Context, log *zap.Logger, agentId int64, msg string, err error) {
	handleProcessingError(ctx, a.hub, log, agentId, msg, err)
}

func (a *serverApi) hub() (SentryHub, string) {
	return a.Hub, ""
}

// TODO: optimize open subscriptions by multiplexing callbacks and only subscribe to redis once.
func (a *serverApi) Subscribe(ctx context.Context, channel string, callback notifications.Callback) {
	// go-redis will automatically re-connect on error
	pubsub := a.redisClient.Subscribe(ctx, channel)
	defer func() {
		if err := pubsub.Close(); err != nil {
			a.HandleProcessingError(ctx, a.log, modshared.NoAgentId, fmt.Sprintf("failed to close channel %q after subscribing", channel), err)
		}
	}()
	ch := pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			return
		case message := <-ch:
			protoMessage, err := redisProtoUnmarshal(message.Payload)
			if err != nil {
				a.HandleProcessingError(ctx, a.log, modshared.NoAgentId, fmt.Sprintf("receiver message in channel %q cannot be unmarshalled into proto message", channel), err)
				continue
			}
			callback(ctx, protoMessage)
		}
	}
}

func redisProtoUnmarshal(payload string) (proto.Message, error) {
	var a anypb.Any
	err := proto.Unmarshal([]byte(payload), &a)
	if err != nil {
		return nil, err
	}
	return a.UnmarshalNew()
}

func handleProcessingError(ctx context.Context, hub func() (SentryHub, string), log *zap.Logger, agentId int64, msg string, err error) {
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
		h, transaction := hub()
		logAndCapture(ctx, h, transaction, log, agentId, msg, err)
	}
}

func logAndCapture(ctx context.Context, hub SentryHub, transaction string, log *zap.Logger, agentId int64, msg string, err error) {
	log.Error(msg, logz.Error(err))

	errStr := removeRandomPort(err.Error())

	event := sentry.NewEvent()
	if agentId != modshared.NoAgentId {
		event.User.ID = strconv.FormatInt(agentId, 10)
	}
	event.Level = sentry.LevelError
	event.Exception = []sentry.Exception{
		{
			Type:       reflect.TypeOf(err).String(),
			Value:      fmt.Sprintf("%s: %s", msg, errStr),
			Stacktrace: sentry.ExtractStacktrace(err),
		},
	}
	traceId := trace.SpanContextFromContext(ctx).TraceID()
	if traceId.IsValid() {
		event.Tags[modserver.TraceIdSentryField] = traceId.String()
	}
	event.Transaction = transaction
	hub.CaptureEvent(event)
}

func removeRandomPort(err string) string {
	return removePortStdLibError.ReplaceAllString(err, "${1}x${2}")
}
