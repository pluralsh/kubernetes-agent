package kasapp

import (
	"errors"
	"net/http"
	"strconv"
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/api"
	gapi "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/gitlab/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/cache"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/errz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/testing/mock_gitlab"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/testing/testhelpers"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	_ modserver.RpcApi = (*serverRpcApi)(nil)
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
			code:       codes.Unavailable,
			captureErr: "HTTP status code: 500",
		},
	}
	for _, tc := range tests {
		t.Run(strconv.Itoa(tc.httpStatus), func(t *testing.T) {
			log, hub, rpcApi, correlationId := setupRpcApi(t, tc.httpStatus)
			if tc.captureErr != "" {
				hub.EXPECT().
					CaptureEvent(gomock.Any()).
					Do(func(event *sentry.Event) {
						assert.Equal(t, correlationId, event.Tags[modserver.CorrelationIdSentryField])
						assert.Empty(t, event.User.ID)
						assert.Equal(t, sentry.LevelError, event.Level)
						assert.Equal(t, "gitlab.ClientError", event.Exception[0].Type)
						assert.Equal(t, "AgentInfo(): "+tc.captureErr, event.Exception[0].Value)
					})
			}
			info, err := rpcApi.AgentInfo(rpcApi.StreamCtx, log)
			assert.Equal(t, tc.code, status.Code(err))
			assert.Nil(t, info)
		})
	}
}

func TestRpcHandleProcessingError_UserError(t *testing.T) {
	log, _, rpcApi, _ := setupRpcApi(t, http.StatusInternalServerError)
	err := errz.NewUserError("boom")
	rpcApi.HandleProcessingError(log, testhelpers.AgentId, "Bla", err)
}

func TestRpcHandleProcessingError_NonUserError_AgentId(t *testing.T) {
	log, hub, rpcApi, correlationId := setupRpcApi(t, http.StatusInternalServerError)
	err := errors.New("boom")
	hub.EXPECT().
		CaptureEvent(gomock.Any()).
		Do(func(event *sentry.Event) {
			assert.Equal(t, correlationId, event.Tags[modserver.CorrelationIdSentryField])
			assert.Equal(t, strconv.FormatInt(testhelpers.AgentId, 10), event.User.ID)
			assert.Equal(t, sentry.LevelError, event.Level)
			assert.Equal(t, "*errors.errorString", event.Exception[0].Type)
			assert.Equal(t, "Bla: boom", event.Exception[0].Value)
		})
	rpcApi.HandleProcessingError(log, testhelpers.AgentId, "Bla", err)
}

func TestRpcHandleProcessingError_NonUserError_NoAgentId(t *testing.T) {
	log, hub, rpcApi, correlationId := setupRpcApi(t, http.StatusInternalServerError)
	err := errors.New("boom")
	hub.EXPECT().
		CaptureEvent(gomock.Any()).
		Do(func(event *sentry.Event) {
			assert.Equal(t, correlationId, event.Tags[modserver.CorrelationIdSentryField])
			assert.Empty(t, event.User.ID)
			assert.Equal(t, sentry.LevelError, event.Level)
			assert.Equal(t, "*errors.errorString", event.Exception[0].Type)
			assert.Equal(t, "Bla: boom", event.Exception[0].Value)
		})
	rpcApi.HandleProcessingError(log, modshared.NoAgentId, "Bla", err)
}

func setupRpcApi(t *testing.T, statusCode int) (*zap.Logger, *MockSentryHub, *serverRpcApi, string) {
	log := zaptest.NewLogger(t)
	ctrl := gomock.NewController(t)
	hub := NewMockSentryHub(ctrl)
	ctx, correlationId := testhelpers.CtxWithCorrelation(t)
	gitLabClient := mock_gitlab.SetupClient(t, gapi.AgentInfoApiPath, func(w http.ResponseWriter, r *http.Request) {
		testhelpers.AssertGetJsonRequestIsCorrect(t, r, correlationId)
		w.WriteHeader(statusCode)
	})
	ctx = api.InjectAgentMD(ctx, &api.AgentMD{
		Token: testhelpers.AgentkToken,
	})
	rpcApi := &serverRpcApi{
		RpcApiStub: modshared.RpcApiStub{
			StreamCtx: ctx,
		},
		GitLabClient:   gitLabClient,
		Hub:            hub,
		AgentInfoCache: cache.NewWithError(0, 0), // no cache!
	}
	return log, hub, rpcApi, correlationId
}
