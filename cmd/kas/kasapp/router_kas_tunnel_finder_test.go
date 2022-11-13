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
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/testhelpers"
	"go.uber.org/zap/zaptest"
)

const (
	selfAddr = "grpc://self"
)

func TestTunnelFinder_PollStartsSingleGoroutineForUrl(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tf, querier, rpcApi, kasPool := setupTunnelFinder(ctx, t)

	gomock.InOrder(
		querier.EXPECT().
			KasUrlsByAgentId(gomock.Any(), testhelpers.AgentId, gomock.Any()).
			Do(func(ctx context.Context, agentId int64, cb tracker.KasUrlsByAgentIdCallback) {
				done, err := cb(kasUrlPipe)
				assert.NoError(t, err)
				assert.False(t, done)
			}),
		querier.EXPECT().
			KasUrlsByAgentId(gomock.Any(), testhelpers.AgentId, gomock.Any()).
			Do(func(ctx context.Context, agentId int64, cb tracker.KasUrlsByAgentIdCallback) {
				done, err := cb(kasUrlPipe)
				assert.NoError(t, err)
				assert.False(t, done)
				cancel()
			}),
	)
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
	assert.Len(t, tf.connections, 2)
	assert.Contains(t, tf.connections, selfAddr)
	assert.Contains(t, tf.connections, "grpc://pipe")
}

func TestTunnelFinder_PollStartsGoroutineForEachUrl(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tf, querier, rpcApi, kasPool := setupTunnelFinder(ctx, t)

	gomock.InOrder(
		querier.EXPECT().
			KasUrlsByAgentId(gomock.Any(), testhelpers.AgentId, gomock.Any()).
			Do(func(ctx context.Context, agentId int64, cb tracker.KasUrlsByAgentIdCallback) {
				done, err := cb(kasUrlPipe)
				assert.NoError(t, err)
				assert.False(t, done)
				done, err = cb("grpc://pipe2")
				assert.NoError(t, err)
				assert.True(t, done)
			}),
		querier.EXPECT().
			KasUrlsByAgentId(gomock.Any(), testhelpers.AgentId, gomock.Any()).
			Do(func(ctx context.Context, agentId int64, cb tracker.KasUrlsByAgentIdCallback) {
				done, err := cb(kasUrlPipe)
				assert.NoError(t, err)
				assert.False(t, done)
				done, err = cb("grpc://pipe2")
				assert.NoError(t, err)
				assert.False(t, done)
				cancel()
			}),
	)
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
	assert.Len(t, tf.connections, 3)
	assert.Contains(t, tf.connections, selfAddr)
	assert.Contains(t, tf.connections, "grpc://pipe")
	assert.Contains(t, tf.connections, "grpc://pipe2")
}

func setupTunnelFinder(ctx context.Context, t *testing.T) (*tunnelFinder, *mock_reverse_tunnel_tracker.MockQuerier, *mock_modserver.MockRpcApi, *mock_rpc.MockPoolInterface) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	querier := mock_reverse_tunnel_tracker.NewMockQuerier(ctrl)
	rpcApi := mock_modserver.NewMockRpcApi(ctrl)
	kasPool := mock_rpc.NewMockPoolInterface(ctrl)

	tf := &tunnelFinder{
		log:              zaptest.NewLogger(t),
		kasPool:          kasPool,
		tunnelQuerier:    querier,
		rpcApi:           rpcApi,
		fullMethod:       "/gitlab.agent.grpctool.test.Testing/RequestResponse",
		ownPrivateApiUrl: selfAddr,
		agentId:          testhelpers.AgentId,
		outgoingCtx:      ctx,
		pollConfig:       testhelpers.NewPollConfig(100 * time.Millisecond),
		foundTunnel:      make(chan readyTunnel),
		connections:      make(map[string]kasConnAttempt),
	}
	gomock.InOrder(
		kasPool.EXPECT().
			Dial(gomock.Any(), selfAddr).
			DoAndReturn(func(ctx context.Context, targetUrl string) (grpctool.PoolConn, error) {
				<-ctx.Done() // block to simulate a long running dial
				return nil, ctx.Err()
			}),
		rpcApi.EXPECT().
			HandleProcessingError(gomock.Any(), testhelpers.AgentId, gomock.Any(), gomock.Any()),
	)
	return tf, querier, rpcApi, kasPool
}
