package kasapp

import (
	"errors"
	"net/http"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/api"
	gapi "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/gitlab/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/cache"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/errz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/testing/matcher"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/testing/mock_errtracker"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/testing/mock_gitlab"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/testing/testhelpers"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	_ modserver.Api = (*serverApi)(nil)
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
			log, errTracker, rpcApi := setupRpcApi(t, tc.httpStatus)
			if tc.captureErr != "" {
				errTracker.EXPECT().
					Capture(matcher.ErrorEq("AgentInfo(): "+tc.captureErr), gomock.Any())
			}
			info, err := rpcApi.AgentInfo(rpcApi.StreamCtx, log)
			assert.Equal(t, tc.code, status.Code(err))
			assert.Nil(t, info)
		})
	}
}

func TestRpcHandleProcessingError_UserError(t *testing.T) {
	log, _, rpcApi := setupRpcApi(t, http.StatusInternalServerError)
	err := errz.NewUserError("boom")
	rpcApi.HandleProcessingError(log, 123, "Bla", err)
}

func TestRpcHandleProcessingError_NonUserError_AgentId(t *testing.T) {
	log, errTracker, rpcApi := setupRpcApi(t, http.StatusInternalServerError)
	err := errors.New("boom")
	errTracker.EXPECT().
		Capture(matcher.ErrorEq("Bla: boom"), gomock.Len(2))
	rpcApi.HandleProcessingError(log, 123, "Bla", err)
}

func TestRpcHandleProcessingError_NonUserError_NoAgentId(t *testing.T) {
	log, errTracker, rpcApi := setupRpcApi(t, http.StatusInternalServerError)
	err := errors.New("boom")
	errTracker.EXPECT().
		Capture(matcher.ErrorEq("Bla: boom"), gomock.Len(1))
	rpcApi.HandleProcessingError(log, modshared.NoAgentId, "Bla", err)
}

func setupRpcApi(t *testing.T, statusCode int) (*zap.Logger, *mock_errtracker.MockTracker, *serverRpcApi) {
	log := zaptest.NewLogger(t)
	ctrl := gomock.NewController(t)
	errTracker := mock_errtracker.NewMockTracker(ctrl)
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
		ErrorTracker:   errTracker,
		AgentInfoCache: cache.NewWithError(0, 0), // no cache!
	}
	return log, errTracker, rpcApi
}
