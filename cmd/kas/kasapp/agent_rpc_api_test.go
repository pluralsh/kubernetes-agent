package kasapp

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/api"
	gapi "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/gitlab/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/cache"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/errz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_cache"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_gitlab"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/testhelpers"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	_ modserver.RpcApi             = (*serverRpcApi)(nil)
	_ modserver.RpcApiFactory      = (*serverRpcApiFactory)(nil).New
	_ modserver.AgentRpcApi        = (*serverAgentRpcApi)(nil)
	_ modserver.AgentRpcApiFactory = (*serverAgentRpcApiFactory)(nil).New
)

func TestGetAgentInfo_Errors(t *testing.T) {
	tests := []struct {
		httpStatus int
		code       codes.Code
		captureErr string
	}{
		{
			httpStatus: http.StatusForbidden,
			code:       codes.PermissionDenied,
		},
		{
			httpStatus: http.StatusUnauthorized,
			code:       codes.Unauthenticated,
		},
		{
			httpStatus: http.StatusNotFound,
			code:       codes.NotFound,
		},
		{
			httpStatus: http.StatusInternalServerError,
			captureErr: "HTTP status code: 500 for path /api/v4/internal/kubernetes/agent_info",
			code:       codes.Unavailable,
		},
		{
			httpStatus: http.StatusBadGateway,
			captureErr: "HTTP status code: 502 for path /api/v4/internal/kubernetes/agent_info",
			code:       codes.Unavailable,
		},
		{
			httpStatus: http.StatusServiceUnavailable,
			captureErr: "HTTP status code: 503 for path /api/v4/internal/kubernetes/agent_info",
			code:       codes.Unavailable,
		},
	}
	for _, tc := range tests {
		t.Run(strconv.Itoa(tc.httpStatus), func(t *testing.T) {
			ctx, log, hub, rpcApi, traceId := setupAgentRpcApi(t, tc.httpStatus)
			if tc.captureErr != "" {
				hub.EXPECT().
					CaptureEvent(gomock.Any()).
					Do(func(event *sentry.Event) {
						assert.Equal(t, traceId.String(), event.Tags[modserver.TraceIdSentryField])
						assert.Empty(t, event.User.ID)
						assert.Equal(t, sentry.LevelError, event.Level)
						assert.Equal(t, "*gitlab.ClientError", event.Exception[0].Type)
						assert.Equal(t, "AgentInfo(): "+tc.captureErr, event.Exception[0].Value)
					})
			}
			info, err := rpcApi.AgentInfo(ctx, log)
			assert.Equal(t, tc.code, status.Code(err))
			assert.Nil(t, info)
		})
	}
}

func TestRpcHandleProcessingError_UserError(t *testing.T) {
	_, log, _, rpcApi, _ := setupAgentRpcApi(t, http.StatusInternalServerError)
	err := errz.NewUserError("boom")
	rpcApi.HandleProcessingError(log, testhelpers.AgentId, "Bla", err)
}

func TestRpcHandleProcessingError_NonUserError_AgentId(t *testing.T) {
	_, log, hub, rpcApi, traceId := setupAgentRpcApi(t, http.StatusInternalServerError)
	err := errors.New("boom")
	hub.EXPECT().
		CaptureEvent(gomock.Any()).
		Do(func(event *sentry.Event) {
			assert.Equal(t, traceId.String(), event.Tags[modserver.TraceIdSentryField])
			assert.Equal(t, strconv.FormatInt(testhelpers.AgentId, 10), event.User.ID)
			assert.Equal(t, sentry.LevelError, event.Level)
			assert.Equal(t, "*errors.errorString", event.Exception[0].Type)
			assert.Equal(t, "Bla: boom", event.Exception[0].Value)
		})
	rpcApi.HandleProcessingError(log, testhelpers.AgentId, "Bla", err)
}

func TestRpcHandleProcessingError_NonUserError_NoAgentId(t *testing.T) {
	_, log, hub, rpcApi, traceId := setupAgentRpcApi(t, http.StatusInternalServerError)
	err := errors.New("boom")
	hub.EXPECT().
		CaptureEvent(gomock.Any()).
		Do(func(event *sentry.Event) {
			assert.Equal(t, traceId.String(), event.Tags[modserver.TraceIdSentryField])
			assert.Empty(t, event.User.ID)
			assert.Equal(t, sentry.LevelError, event.Level)
			assert.Equal(t, "*errors.errorString", event.Exception[0].Type)
			assert.Equal(t, "Bla: boom", event.Exception[0].Value)
		})
	rpcApi.HandleProcessingError(log, modshared.NoAgentId, "Bla", err)
}

func setupAgentRpcApi(t *testing.T, statusCode int) (context.Context, *zap.Logger, *MockSentryHub, *serverAgentRpcApi, trace.TraceID) {
	log := zaptest.NewLogger(t)
	ctrl := gomock.NewController(t)
	hub := NewMockSentryHub(ctrl)
	errCacher := mock_cache.NewMockErrCacher[api.AgentToken](ctrl)
	ctx, traceId := testhelpers.CtxWithSpanContext(t)
	gitLabClient := mock_gitlab.SetupClient(t, gapi.AgentInfoApiPath, func(w http.ResponseWriter, r *http.Request) {
		testhelpers.AssertGetJsonRequestIsCorrect(t, r, traceId)
		w.WriteHeader(statusCode)
	})
	sra := &serverRpcApi{
		RpcApiStub: modshared.RpcApiStub{
			Logger:    log,
			StreamCtx: ctx,
		},
		sentryHubRoot: sentry.NewHub(nil, sentry.NewScope()),
		service:       "svc",
		method:        "method",
	}
	sra.hub() // so that Once fires here and doesn't overwrite our mock later
	sra.sentryHub = hub

	rpcApi := &serverAgentRpcApi{
		RpcApi:         sra,
		Token:          testhelpers.AgentkToken,
		GitLabClient:   gitLabClient,
		AgentInfoCache: cache.NewWithError[api.AgentToken, *api.AgentInfo](0, 0, errCacher, func(err error) bool { return false }), // no cache!
	}
	return ctx, log, hub, rpcApi, traceId
}
