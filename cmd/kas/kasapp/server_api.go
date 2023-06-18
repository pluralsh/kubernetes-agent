package kasapp

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/redis/rueidis"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/errz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/retry"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/syncz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/event"
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
)

const (
	redisAttemptInterval = 50 * time.Millisecond
	redisInitBackoff     = 100 * time.Millisecond
	redisMaxBackoff      = 10 * time.Second
	redisResetDuration   = 20 * time.Second
	redisBackoffFactor   = 2.0
	redisJitter          = 1.0

	gitPushEventsRedisChannel = "kas_git_push_events2"
)

type SentryHub interface {
	CaptureEvent(event *sentry.Event) *sentry.EventID
}

type serverApi struct {
	log             *zap.Logger
	Hub             SentryHub
	redisClient     rueidis.Client
	gitPushEvent    syncz.Subscriptions[*event.GitPushEvent]
	redisPollConfig retry.PollConfigFactory
}

func newServerApi(log *zap.Logger, hub SentryHub, redisClient rueidis.Client) *serverApi {
	return &serverApi{
		log:         log,
		Hub:         hub,
		redisClient: redisClient,
		redisPollConfig: retry.NewPollConfigFactory(redisAttemptInterval, retry.NewExponentialBackoffFactory(
			redisInitBackoff,
			redisMaxBackoff,
			redisResetDuration,
			redisBackoffFactor,
			redisJitter,
		)),
	}
}

func (a *serverApi) HandleProcessingError(ctx context.Context, log *zap.Logger, agentId int64, msg string, err error) {
	handleProcessingError(ctx, a.hub, log, agentId, msg, err)
}

func (a *serverApi) hub() (SentryHub, string) {
	return a.Hub, ""
}

func (a *serverApi) OnGitPushEvent(ctx context.Context, cb syncz.EventCallback[*event.GitPushEvent]) {
	a.gitPushEvent.On(ctx, cb)
}

func (a *serverApi) publishGitPushEvent(ctx context.Context, e *event.GitPushEvent) error {
	payload, err := redisProtoMarshal(e)
	if err != nil {
		return fmt.Errorf("failed to marshal proto message to publish: %w", err)
	}
	publishCmd := a.redisClient.B().Publish().Channel(gitPushEventsRedisChannel).Message(string(payload)).Build()
	return a.redisClient.Do(ctx, publishCmd).Error()
}

// subscribeGitPushEvent subscribes to the Git push event redis channel
// and will dispatch each event to the registered callbacks.
func (a *serverApi) subscribeGitPushEvent(ctx context.Context) {
	_ = retry.PollWithBackoff(ctx, a.redisPollConfig(), func(ctx context.Context) (error, retry.AttemptResult) {
		subCmd := a.redisClient.B().Subscribe().Channel(gitPushEventsRedisChannel).Build()
		err := a.redisClient.Receive(ctx, subCmd, func(msg rueidis.PubSubMessage) {
			protoMessage, err := redisProtoUnmarshal(msg.Message)
			if err != nil {
				a.HandleProcessingError(ctx, a.log, modshared.NoAgentId, fmt.Sprintf("receiver message in channel %q cannot be unmarshalled into proto message", gitPushEventsRedisChannel), err)
				return
			}

			switch e := (protoMessage).(type) {
			case *event.GitPushEvent:
				a.gitPushEvent.Dispatch(ctx, e)
			default:
				a.HandleProcessingError(
					ctx,
					a.log,
					modshared.NoAgentId,
					"GitOps: unable to handle received git push event",
					fmt.Errorf("failed to cast proto message of type %T to concrete type", e))
			}
		})
		switch err { // nolint:errorlint
		case nil, context.Canceled, context.DeadlineExceeded:
			return nil, retry.ContinueImmediately
		default:
			a.HandleProcessingError(ctx, a.log, modshared.NoAgentId, "Error handling Redis SUBSCRIBE", err)
			return nil, retry.Backoff
		}
	})
}

func redisProtoMarshal(m proto.Message) ([]byte, error) {
	a, err := anypb.New(m) // use Any to capture type information so that a value can be instantiated in redisProtoUnmarshal()
	if err != nil {
		return nil, err
	}
	return proto.Marshal(a)
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

	e := sentry.NewEvent()
	if agentId != modshared.NoAgentId {
		e.User.ID = strconv.FormatInt(agentId, 10)
	}
	e.Level = sentry.LevelError
	e.Exception = []sentry.Exception{
		{
			Type:       reflect.TypeOf(err).String(),
			Value:      fmt.Sprintf("%s: %s", msg, errStr),
			Stacktrace: sentry.ExtractStacktrace(err),
		},
	}
	tc := trace.SpanContextFromContext(ctx)
	traceId := tc.TraceID()
	if traceId.IsValid() {
		e.Tags[modserver.SentryFieldTraceId] = traceId.String()
		sampled := "false"
		if tc.IsSampled() {
			sampled = "true"
		}
		e.Tags[modserver.SentryFieldTraceSampled] = sampled
	}
	e.Transaction = transaction
	hub.CaptureEvent(e)
}

func removeRandomPort(err string) string {
	return removePortStdLibError.ReplaceAllString(err, "${1}x${2}")
}
