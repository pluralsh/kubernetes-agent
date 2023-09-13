package mock_modserver

import (
	"context"
	"testing"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/testhelpers"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/entity"
	"go.uber.org/zap/zaptest"
)

func IncomingAgentCtx(t *testing.T, rpcApi *MockAgentRpcApi) context.Context {
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
	ctx = modserver.InjectAgentRpcApi(ctx, rpcApi)

	return ctx
}

func AgentMeta() *entity.AgentMeta {
	return &entity.AgentMeta{
		Version:      "v1.2.3",
		CommitId:     "32452345",
		PodNamespace: "ns1",
		PodName:      "n1",
	}
}
