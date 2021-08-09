package mock_modserver

import (
	"context"
	"testing"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/testing/testhelpers"
	"go.uber.org/zap/zaptest"
)

func IncomingCtx(t *testing.T, rpcApi *MockRpcApi) context.Context {
	rpcApi.EXPECT().
		AgentToken().
		Return(testhelpers.AgentkToken).
		AnyTimes()
	rpcApi.EXPECT().
		Log().
		Return(zaptest.NewLogger(t)).
		AnyTimes()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	ctx = modserver.InjectRpcApi(ctx, rpcApi)

	return ctx
}
