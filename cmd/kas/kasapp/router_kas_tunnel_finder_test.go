package kasapp

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/reverse_tunnel/tracker"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_reverse_tunnel_tracker"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/testhelpers"
	"go.uber.org/mock/gomock"
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
	gomock.InOrder(
		querier.EXPECT().
			CachedKasUrlsByAgentId(testhelpers.AgentId),
		querier.EXPECT().
			PollKasUrlsByAgentId(gomock.Any(), testhelpers.AgentId, gomock.Any()).
			Do(func(ctx context.Context, agentId int64, cb tracker.PollKasUrlsByAgentIdCallback) {
				cb([]string{kasUrlPipe})
				cb([]string{kasUrlPipe}) // same thing two times
				wg.Wait()
				cancel()
				<-ctx.Done()
			}),
	)
	gomock.InOrder(
		kasPool.EXPECT().
			Dial(gomock.Any(), kasUrlPipe).
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
	assert.Contains(t, tf.connections, kasUrlPipe)
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
	gomock.InOrder(
		querier.EXPECT().
			CachedKasUrlsByAgentId(testhelpers.AgentId),
		querier.EXPECT().
			PollKasUrlsByAgentId(gomock.Any(), testhelpers.AgentId, gomock.Any()).
			Do(func(ctx context.Context, agentId int64, cb tracker.PollKasUrlsByAgentIdCallback) {
				cb([]string{kasUrlPipe, "grpc://pipe2"})
				wg.Wait()
				cancel()
				<-ctx.Done()
			}),
	)
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

	gatewayKasVisitor, err := grpctool.NewStreamVisitor(&GatewayKasResponse{})
	require.NoError(t, err)

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
		gatewayKasVisitor,
	)
	return tf, querier, rpcApi, kasPool
}
