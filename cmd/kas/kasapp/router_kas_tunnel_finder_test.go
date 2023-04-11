package kasapp

import (
	"context"
	"sync"
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

	var wg sync.WaitGroup
	wg.Add(2)

	gomock.InOrder(
		kasPool.EXPECT().
			Dial(gomock.Any(), selfAddr).
			DoAndReturn(func(ctx context.Context, targetUrl string) (grpctool.PoolConn, error) {
				wg.Done()
				<-ctx.Done() // block to simulate a long running dial
				return nil, ctx.Err()
			}),
		rpcApi.EXPECT().
			HandleProcessingError(gomock.Any(), testhelpers.AgentId, gomock.Any(), gomock.Any()),
	)
	querier.EXPECT().
		PollKasUrlsByAgentId(gomock.Any(), testhelpers.AgentId, gomock.Any()).
		Do(func(ctx context.Context, agentId int64, cb tracker.PollKasUrlsByAgentIdCallback) {
			done := cb(true, kasUrlPipe)
			assert.False(t, done)
			done = cb(true, kasUrlPipe)
			assert.False(t, done)
			wg.Wait()
			cancel()
			<-ctx.Done()
		})
	gomock.InOrder(
		kasPool.EXPECT().
			Dial(gomock.Any(), "grpc://pipe").
			DoAndReturn(func(ctx context.Context, targetUrl string) (grpctool.PoolConn, error) {
				wg.Done()
				<-ctx.Done() // block to simulate a long running dial
				return nil, ctx.Err()
			}),
		rpcApi.EXPECT().
			HandleProcessingError(gomock.Any(), testhelpers.AgentId, gomock.Any(), gomock.Any()),
	)

	_, err := tf.Find(ctx)
	assert.Same(t, context.Canceled, err)
	assert.Len(t, tf.connections, 2)
	assert.Contains(t, tf.connections, selfAddr)
	assert.Contains(t, tf.connections, "grpc://pipe")
}

func TestTunnelFinder_PollStartsGoroutineForEachUrl(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tf, querier, rpcApi, kasPool := setupTunnelFinder(ctx, t)

	var wg sync.WaitGroup
	wg.Add(3)

	gomock.InOrder(
		kasPool.EXPECT().
			Dial(gomock.Any(), selfAddr).
			DoAndReturn(func(ctx context.Context, targetUrl string) (grpctool.PoolConn, error) {
				wg.Done()
				<-ctx.Done() // block to simulate a long running dial
				return nil, ctx.Err()
			}),
		rpcApi.EXPECT().
			HandleProcessingError(gomock.Any(), testhelpers.AgentId, gomock.Any(), gomock.Any()),
	)
	querier.EXPECT().
		PollKasUrlsByAgentId(gomock.Any(), testhelpers.AgentId, gomock.Any()).
		Do(func(ctx context.Context, agentId int64, cb tracker.PollKasUrlsByAgentIdCallback) {
			done := cb(true, kasUrlPipe)
			assert.False(t, done)
			done = cb(false, "grpc://pipe2")
			assert.False(t, done)
			done = cb(true, kasUrlPipe)
			assert.False(t, done)
			done = cb(false, "grpc://pipe2")
			assert.False(t, done)
			wg.Wait()
			cancel()
			<-ctx.Done()
		})
	kasPool.EXPECT().
		Dial(gomock.Any(), "grpc://pipe").
		DoAndReturn(func(ctx context.Context, targetUrl string) (grpctool.PoolConn, error) {
			wg.Done()
			<-ctx.Done() // block to simulate a long running dial
			return nil, ctx.Err()
		})
	kasPool.EXPECT().
		Dial(gomock.Any(), "grpc://pipe2").
		DoAndReturn(func(ctx context.Context, targetUrl string) (grpctool.PoolConn, error) {
			wg.Done()
			<-ctx.Done() // block to simulate a long running dial
			return nil, ctx.Err()
		})
	rpcApi.EXPECT().
		HandleProcessingError(gomock.Any(), testhelpers.AgentId, gomock.Any(), gomock.Any()).
		Times(2)
	_, err := tf.Find(ctx)
	assert.Same(t, context.Canceled, err)
	assert.Len(t, tf.connections, 3)
	assert.Contains(t, tf.connections, selfAddr)
	assert.Contains(t, tf.connections, "grpc://pipe")
	assert.Contains(t, tf.connections, "grpc://pipe2")
}

func setupTunnelFinder(ctx context.Context, t *testing.T) (*tunnelFinder, *mock_reverse_tunnel_tracker.MockPollingQuerier, *mock_modserver.MockRpcApi, *mock_rpc.MockPoolInterface) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	querier := mock_reverse_tunnel_tracker.NewMockPollingQuerier(ctrl)
	rpcApi := mock_modserver.NewMockRpcApi(ctrl)
	kasPool := mock_rpc.NewMockPoolInterface(ctrl)

	tf := newTunnelFinder(
		zaptest.NewLogger(t),
		kasPool,
		querier,
		rpcApi,
		"/gitlab.agent.grpctool.test.Testing/RequestResponse",
		selfAddr,
		testhelpers.AgentId,
		ctx,
		testhelpers.NewPollConfig(100*time.Millisecond),
	)
	return tf, querier, rpcApi, kasPool
}
