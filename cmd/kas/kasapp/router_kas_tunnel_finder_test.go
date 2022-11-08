package kasapp

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/reverse_tunnel/tracker"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_reverse_tunnel_tracker"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/testhelpers"
	"go.uber.org/zap/zaptest"
)

func TestTunnelFinder_NoDialsForNonMatchingService(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tf, querier, _, _ := setupTunnelFinder(ctx, t)
	tf.fullMethod = "/service/method" // doesn't match tunnel
	gomock.InOrder(
		querier.EXPECT().
			GetTunnelsByAgentId(gomock.Any(), testhelpers.AgentId, gomock.Any()).
			Do(func(ctx context.Context, agentId int64, cb tracker.GetTunnelsByAgentIdCallback) {
				done, err := cb(tunnelInfo())
				assert.NoError(t, err)
				assert.False(t, done)
			}),
		querier.EXPECT().
			GetTunnelsByAgentId(gomock.Any(), testhelpers.AgentId, gomock.Any()).
			Do(func(ctx context.Context, agentId int64, cb tracker.GetTunnelsByAgentIdCallback) {
				cancel()
			}),
	)
	tf.Run(ctx)
	assert.Empty(t, tf.connections)
}

func TestTunnelFinder_PollStartsSingleGoroutineForUrl(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tf, querier, rpcApi, kasPool := setupTunnelFinder(ctx, t)

	cnt := 0

	querier.EXPECT().
		GetTunnelsByAgentId(gomock.Any(), testhelpers.AgentId, gomock.Any()).
		Do(func(ctx context.Context, agentId int64, cb tracker.GetTunnelsByAgentIdCallback) {
			done, err := cb(tunnelInfo())
			assert.NoError(t, err)
			assert.False(t, done)
			ti := tunnelInfo()
			ti.ConnectionId = 23
			done, err = cb(ti)
			assert.NoError(t, err)
			assert.False(t, done)
			cnt++
			if cnt == 2 {
				cancel()
			}
		}).
		Times(2)
	gomock.InOrder(
		kasPool.EXPECT().
			Dial(gomock.Any(), "grpc://pipe").
			DoAndReturn(func(ctx context.Context, targetUrl string) (grpctool.PoolConn, error) {
				<-ctx.Done() // block to simulate a long running dial
				return nil, ctx.Err()
			}),
		rpcApi.EXPECT().
			HandleProcessingError(gomock.Any(), testhelpers.AgentId, gomock.Any(), gomock.Any()),
	)

	tf.Run(ctx)
	assert.Len(t, tf.connections, 1)
}

func TestTunnelFinder_PollStartsGoroutineForEachUrl(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tf, querier, rpcApi, kasPool := setupTunnelFinder(ctx, t)

	cnt := 0

	querier.EXPECT().
		GetTunnelsByAgentId(gomock.Any(), testhelpers.AgentId, gomock.Any()).
		Do(func(ctx context.Context, agentId int64, cb tracker.GetTunnelsByAgentIdCallback) {
			done, err := cb(tunnelInfo())
			assert.NoError(t, err)
			assert.False(t, done)
			ti := tunnelInfo()
			ti.ConnectionId = 23
			ti.KasUrl = "grpc://pipe2"
			done, err = cb(ti)
			assert.NoError(t, err)
			assert.False(t, done)
			cnt++
			if cnt == 3 {
				cancel()
			}
		}).
		Times(3)
	kasPool.EXPECT().
		Dial(gomock.Any(), "grpc://pipe").
		DoAndReturn(func(ctx context.Context, targetUrl string) (grpctool.PoolConn, error) {
			<-ctx.Done() // block to simulate a long running dial
			return nil, ctx.Err()
		})
	kasPool.EXPECT().
		Dial(gomock.Any(), "grpc://pipe2").
		DoAndReturn(func(ctx context.Context, targetUrl string) (grpctool.PoolConn, error) {
			<-ctx.Done() // block to simulate a long running dial
			return nil, ctx.Err()
		})
	rpcApi.EXPECT().
		HandleProcessingError(gomock.Any(), testhelpers.AgentId, gomock.Any(), gomock.Any()).
		Times(2)
	tf.Run(ctx)
	assert.Len(t, tf.connections, 2)
}

func setupTunnelFinder(ctx context.Context, t *testing.T) (*tunnelFinder, *mock_reverse_tunnel_tracker.MockQuerier, *mock_modserver.MockRpcApi, *MockKasPool) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	querier := mock_reverse_tunnel_tracker.NewMockQuerier(ctrl)
	rpcApi := mock_modserver.NewMockRpcApi(ctrl)
	kasPool := NewMockKasPool(ctrl)

	tf := &tunnelFinder{
		log:           zaptest.NewLogger(t),
		kasPool:       kasPool,
		tunnelQuerier: querier,
		rpcApi:        rpcApi,
		fullMethod:    "/gitlab.agent.grpctool.test.Testing/RequestResponse",
		agentId:       testhelpers.AgentId,
		outgoingCtx:   ctx,
		pollConfig:    testhelpers.NewPollConfig(100 * time.Millisecond),
		foundTunnel:   make(chan readyTunnel),
		connections:   make(map[string]kasConnAttempt),
	}
	return tf, querier, rpcApi, kasPool
}
