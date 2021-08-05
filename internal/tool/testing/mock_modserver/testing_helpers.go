package mock_modserver

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/testing/testhelpers"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc/metadata"
)

func IncomingCtx(t *testing.T, rpcApi modserver.RpcApi) context.Context {
	creds := grpctool.NewTokenCredentials(testhelpers.AgentkToken, false)
	meta, err := creds.GetRequestMetadata(context.Background())
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	ctx = metadata.NewIncomingContext(ctx, metadata.New(meta))
	agentMD, err := grpctool.AgentMDFromRawContext(ctx)
	require.NoError(t, err)
	ctx = api.InjectAgentMD(ctx, agentMD)
	ctx = grpctool.InjectLogger(ctx, zaptest.NewLogger(t))
	ctx = grpctool.InjectRpcApi(ctx, rpcApi)

	return ctx
}
