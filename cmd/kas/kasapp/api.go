package kasapp

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"sync"

	"github.com/getsentry/sentry-go"
	"github.com/redis/go-redis/v9"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/errz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/logz"
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
	gitPushEventsRedisChannel = "kas_git_push_events"
)

type SentryHub interface {
	CaptureEvent(event *sentry.Event) *sentry.EventID
}

type serverApi struct {
	log                       *zap.Logger
	Hub                       SentryHub
	redisClient               redis.UniversalClient
	gitPushEventSubscriptions subscriptions
}

type subscriptions struct {
	mu  sync.Mutex
	chs []chan<- *modserver.Project
}

func (s *subscriptions) add(ch chan<- *modserver.Project) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.chs = append(s.chs, ch)
}

func (s *subscriptions) remove(ch chan<- *modserver.Project) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, c := range s.chs {
		if c == ch {
			s.chs = append(s.chs[:i], s.chs[i+1:]...)
			break
		}
	}
}

func (a *serverApi) HandleProcessingError(ctx context.Context, log *zap.Logger, agentId int64, msg string, err error) {
	handleProcessingError(ctx, a.hub, log, agentId, msg, err)
}

func (a *serverApi) hub() (SentryHub, string) {
	return a.Hub, ""
}

// OnGitPushEvent runs the given callback function for a received Git push event.
// The Git push event may come from any GitLab project and as such it's up to the
// callback to filter out the events that it's interested in.
// This particular implementation registers an unbuffered channel for the callback
// which receives the actual event from a redis subscription.
// This is mainly to unblock the redis subscription from the callback execution.
func (a *serverApi) OnGitPushEvent(ctx context.Context, callback modserver.GitPushEventCallback) {
	ch := make(chan *modserver.Project)
	defer a.gitPushEventSubscriptions.remove(ch)
	a.gitPushEventSubscriptions.add(ch)

	done := ctx.Done()
	for {
		select {
		case <-done:
			return
		case m := <-ch:
			callback(ctx, m)
		}
	}
}

func (a *serverApi) publishGitPushEvent(ctx context.Context, e *modserver.Project) error {
	payload, err := redisProtoMarshal(e)
	if err != nil {
		return fmt.Errorf("failed to marshal proto message to publish: %w", err)
	}
	return a.redisClient.Publish(ctx, gitPushEventsRedisChannel, payload).Err()
}

// subscribeGitPushEvent subscribes to the Git push event redis channel
// and will dispatch each event to the registered callbacks.
func (a *serverApi) subscribeGitPushEvent(ctx context.Context) {
	// go-redis will automatically re-connect on error
	pubsub := a.redisClient.Subscribe(ctx, gitPushEventsRedisChannel)
	defer func() {
		if err := pubsub.Close(); err != nil {
			a.HandleProcessingError(ctx, a.log, modshared.NoAgentId, fmt.Sprintf("failed to close channel %q after subscribing", gitPushEventsRedisChannel), err)
		}
	}()
	ch := pubsub.Channel()
	done := ctx.Done()
	for {
		select {
		case <-done:
			return
		case message := <-ch:
			protoMessage, err := redisProtoUnmarshal(message.Payload)
			if err != nil {
				a.HandleProcessingError(ctx, a.log, modshared.NoAgentId, fmt.Sprintf("receiver message in channel %q cannot be unmarshalled into proto message", gitPushEventsRedisChannel), err)
				continue
			}

			switch m := (protoMessage).(type) {
			case *modserver.Project:
				a.dispatchGitPushEvent(ctx, m)
			default:
				a.HandleProcessingError(
					ctx,
					a.log,
					modshared.NoAgentId,
					"GitOps: unable to handle received git push event",
					fmt.Errorf("failed to cast proto message of type %T to concrete type", m))
			}
		}
	}
}

// dispatchGitPushEvent dispatches the given `project` which is the message of the Git push event
// to all registered subscriptions registered by OnGitPushEvent.
// This particular implementation will drop events per registered callback if their
// registered channel is blocked, e.g. when the callback is too slow to handle the produced events.
// This is suboptimal, but will decouple and unblock the redis subscription from callback function's performance.
func (a *serverApi) dispatchGitPushEvent(ctx context.Context, project *modserver.Project) {
	a.gitPushEventSubscriptions.mu.Lock()
	defer a.gitPushEventSubscriptions.mu.Unlock()

	done := ctx.Done()
	for _, ch := range a.gitPushEventSubscriptions.chs {
		select {
		case <-done:
			return
		case ch <- project:
		default:
			// NOTE: if for whatever reason the subscriber isn't able to keep up with the events,
			// we just drop them for now.
			a.log.Debug("Dropping Git push event", logz.ProjectId(project.FullPath))
			continue
		}
	}
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
	tc := trace.SpanContextFromContext(ctx)
	traceId := tc.TraceID()
	if traceId.IsValid() {
		event.Tags[modserver.SentryFieldTraceId] = traceId.String()
		sampled := "false"
		if tc.IsSampled() {
			sampled = "true"
		}
		event.Tags[modserver.SentryFieldTraceSampled] = sampled
	}
	event.Transaction = transaction
	hub.CaptureEvent(event)
}

func removeRandomPort(err string) string {
	return removePortStdLibError.ReplaceAllString(err, "${1}x${2}")
}
