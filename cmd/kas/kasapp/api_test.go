package kasapp

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/errz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/testing/matcher"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/testing/mock_errtracker"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/testing/testhelpers"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

var (
	_ modserver.Api = (*serverApi)(nil)
)

func TestHandleProcessingError_UserError(t *testing.T) {
	ctx, log, _, apiObj := setupApi(t)
	err := errz.NewUserError("boom")
	apiObj.HandleProcessingError(ctx, log, 123, "Bla", err)
}

func TestHandleProcessingError_NonUserError_AgentId(t *testing.T) {
	ctx, log, errTracker, apiObj := setupApi(t)
	err := errors.New("boom")
	errTracker.EXPECT().
		Capture(matcher.ErrorEq("Bla: boom"), gomock.Len(2))
	apiObj.HandleProcessingError(ctx, log, 123, "Bla", err)
}

func TestHandleProcessingError_NonUserError_NoAgentId(t *testing.T) {
	ctx, log, errTracker, apiObj := setupApi(t)
	err := errors.New("boom")
	errTracker.EXPECT().
		Capture(matcher.ErrorEq("Bla: boom"), gomock.Len(1))
	apiObj.HandleProcessingError(ctx, log, modshared.NoAgentId, "Bla", err)
}

func setupApi(t *testing.T) (context.Context, *zap.Logger, *mock_errtracker.MockTracker, *serverApi) {
	log := zaptest.NewLogger(t)
	ctrl := gomock.NewController(t)
	errTracker := mock_errtracker.NewMockTracker(ctrl)
	ctx, _ := testhelpers.CtxWithCorrelation(t)
	apiObj := &serverApi{
		ErrorTracker: errTracker,
	}
	return ctx, log, errTracker, apiObj
}
