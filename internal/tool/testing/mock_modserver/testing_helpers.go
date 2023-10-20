package mock_modserver

import (
	"context"
	"testing"

	"github.com/pluralsh/kuberentes-agent/internal/module/modserver"
	"github.com/pluralsh/kuberentes-agent/internal/tool/testing/testhelpers"
	"github.com/pluralsh/kuberentes-agent/pkg/entity"
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
