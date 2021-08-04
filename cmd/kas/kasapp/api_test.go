package kasapp

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/gitlab/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modshared"
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
	_ modserver.API = (*serverAPI)(nil)
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
			ctx, log, errTracker, apiObj := setupApi(t, tc.httpStatus)
			if tc.captureErr != "" {
				errTracker.EXPECT().
					Capture(matcher.ErrorEq("GetAgentInfo(): "+tc.captureErr), gomock.Any())
			}
			info, err := apiObj.GetAgentInfo(ctx, log, testhelpers.AgentkToken)
			assert.Equal(t, tc.code, status.Code(err))
			assert.Nil(t, info)
		})
	}
}

func TestHandleProcessingError_UserError(t *testing.T) {
	ctx, log, _, apiObj := setupApi(t, http.StatusInternalServerError)
	err := errz.NewUserError("boom")
	apiObj.HandleProcessingError(ctx, log, 123, "Bla", err)
}

func TestHandleProcessingError_NonUserError_AgentId(t *testing.T) {
	ctx, log, errTracker, apiObj := setupApi(t, http.StatusInternalServerError)
	err := errors.New("boom")
	errTracker.EXPECT().
		Capture(matcher.ErrorEq("Bla: boom"), gomock.Len(2))
	apiObj.HandleProcessingError(ctx, log, 123, "Bla", err)
}

func TestHandleProcessingError_NonUserError_NoAgentId(t *testing.T) {
	ctx, log, errTracker, apiObj := setupApi(t, http.StatusInternalServerError)
	err := errors.New("boom")
	errTracker.EXPECT().
		Capture(matcher.ErrorEq("Bla: boom"), gomock.Len(1))
	apiObj.HandleProcessingError(ctx, log, modshared.NoAgentId, "Bla", err)
}

func setupApi(t *testing.T, statusCode int) (context.Context, *zap.Logger, *mock_errtracker.MockTracker, *serverAPI) {
	log := zaptest.NewLogger(t)
	ctrl := gomock.NewController(t)
	errTracker := mock_errtracker.NewMockTracker(ctrl)
	ctx, correlationId := testhelpers.CtxWithCorrelation(t)
	gitLabClient := mock_gitlab.SetupClient(t, api.AgentInfoApiPath, func(w http.ResponseWriter, r *http.Request) {
		testhelpers.AssertGetJsonRequestIsCorrect(t, r, correlationId)
		w.WriteHeader(statusCode)
	})
	apiObj := newAPI(apiConfig{
		GitLabClient:           gitLabClient,
		ErrorTracker:           errTracker,
		AgentInfoCacheTtl:      0, // no cache!
		AgentInfoCacheErrorTtl: 0,
	})
	return ctx, log, errTracker, apiObj
}
