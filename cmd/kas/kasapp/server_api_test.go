package kasapp

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/errz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/testhelpers"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/event"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"k8s.io/apimachinery/pkg/util/wait"
)

var (
	_ modserver.Api = (*serverApi)(nil)
)

func TestHandleProcessingError_UserError(t *testing.T) {
	ctx, log, _, apiObj, _ := setupApi(t)
	err := errz.NewUserError("boom")
	apiObj.HandleProcessingError(ctx, log, testhelpers.AgentId, "Bla", err)
}

func TestHandleProcessingError_NonUserError_AgentId(t *testing.T) {
	ctx, log, hub, apiObj, traceId := setupApi(t)
	err := errors.New("boom")
	hub.EXPECT().
		CaptureEvent(gomock.Any()).
		Do(func(event *sentry.Event) {
			assert.Equal(t, traceId.String(), event.Tags[modserver.SentryFieldTraceId])
			assert.Equal(t, strconv.FormatInt(testhelpers.AgentId, 10), event.User.ID)
			assert.Equal(t, sentry.LevelError, event.Level)
			assert.Equal(t, "*errors.errorString", event.Exception[0].Type)
			assert.Equal(t, "Bla: boom", event.Exception[0].Value)
		})
	apiObj.HandleProcessingError(ctx, log, testhelpers.AgentId, "Bla", err)
}

func TestHandleProcessingError_NonUserError_NoAgentId_NoTraceId(t *testing.T) {
	_, log, hub, apiObj, _ := setupApi(t)
	err := errors.New("boom")
	hub.EXPECT().
		CaptureEvent(gomock.Any()).
		Do(func(event *sentry.Event) {
			assert.NotContains(t, event.Tags, modserver.SentryFieldTraceId)
			assert.Empty(t, event.User.ID)
			assert.Equal(t, sentry.LevelError, event.Level)
			assert.Equal(t, "*errors.errorString", event.Exception[0].Type)
			assert.Equal(t, "Bla: boom", event.Exception[0].Value)
		})
	apiObj.HandleProcessingError(context.Background(), log, modshared.NoAgentId, "Bla", err)
}

func setupApi(t *testing.T) (context.Context, *zap.Logger, *MockSentryHub, *serverApi, trace.TraceID) {
	log := zaptest.NewLogger(t)
	ctrl := gomock.NewController(t)
	hub := NewMockSentryHub(ctrl)
	ctx, traceId := testhelpers.CtxWithSpanContext(t)
	apiObj := newServerApi(log, hub, nil)
	return ctx, log, hub, apiObj, traceId
}

func TestRemoveRandomPort(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "",
			expected: "",
		},
		{
			input:    "bla",
			expected: "bla",
		},
		{
			input:    "read tcp 10.222.67.20:40272->10.216.1.45:11443: read: connection reset by peer",
			expected: "read tcp 10.222.67.20:x->10.216.1.45:11443: read: connection reset by peer",
		},
		{
			input:    "some error with ip and port 10.222.67.20:40272: bla",
			expected: "some error with ip and port 10.222.67.20:40272: bla",
		},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			actual := removeRandomPort(tc.input)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestServerApi_GitPushEvent_DispatchingMultiple(t *testing.T) {
	// GIVEN
	var wg wait.Group
	defer wg.Wait()

	a := newServerApi(zaptest.NewLogger(t), nil, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ev := &event.GitPushEvent{}
	// recorder for callback hits
	rec1 := make(chan struct{})
	rec2 := make(chan struct{})
	subscriber1 := func(_ context.Context, e *event.GitPushEvent) {
		assert.Same(t, ev, e)
		close(rec1)
	}
	subscriber2 := func(_ context.Context, e *event.GitPushEvent) {
		assert.Same(t, ev, e)
		close(rec2)
	}

	// WHEN
	// starting multiple subscribers
	wg.Start(func() {
		a.OnGitPushEvent(ctx, subscriber1)
	})
	wg.Start(func() {
		a.OnGitPushEvent(ctx, subscriber2)
	})

	// give the OnGitPushEvent goroutines time to be scheduled and registered
	assert.Eventually(t, func() bool {
		a.gitPushEventSubscriptions.mu.Lock()
		defer a.gitPushEventSubscriptions.mu.Unlock()
		return len(a.gitPushEventSubscriptions.chs) == 2
	}, time.Minute, time.Millisecond)

	// dispatch a single git push event
	a.dispatchGitPushEvent(ctx, ev)

	// THEN
	<-rec1
	<-rec2
}

func TestServerApi_GitPushEventSubscriptions(t *testing.T) {
	var s subscriptions

	ch1 := make(chan<- *event.GitPushEvent)
	ch2 := make(chan<- *event.GitPushEvent)
	ch3 := make(chan<- *event.GitPushEvent)

	s.add(ch1)
	s.add(ch2)
	s.add(ch3)

	assert.Equal(t, ch1, s.chs[0])
	assert.Equal(t, ch2, s.chs[1])
	assert.Equal(t, ch3, s.chs[2])

	s.remove(ch2)

	assert.Equal(t, ch1, s.chs[0])
	assert.Equal(t, ch3, s.chs[1])
	assert.Nil(t, s.chs[:3:3][2])

	s.remove(ch1)
	s.remove(ch3)
	assert.Nil(t, s.chs[:3:3][0])
	assert.Nil(t, s.chs[:3:3][1])
	assert.Empty(t, s.chs)
}
