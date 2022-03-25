package kasapp

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/reverse_tunnel/tracker"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/testing/mock_modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/testing/mock_reverse_tunnel_tracker"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/testing/testhelpers"
	"go.uber.org/zap/zaptest"
)

func TestTunnelFinder_NoDialsForNonMatchingService(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tChan, tf, querier, _, _ := setupTunnelFinder(ctx, t)
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
	go tf.poll(ctx, testhelpers.NewPollConfig(100*time.Millisecond)()) // 10 times
	select {
	case <-ctx.Done():
	case <-tChan:
		t.FailNow()
	}
	assert.Eventually(t, func() bool {
		tf.mu.Lock()
		defer tf.mu.Unlock()
		return len(tf.connections) == 0
	}, time.Second, 10*time.Millisecond) // wait for goroutines to stop
}

func TestTunnelFinder_PollStartsSingleGoroutineForUrl(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		wg.Wait()
		time.Sleep(200 * time.Millisecond)
		cancel()
	}()
	tChan, tf, querier, rpcApi, kasPool := setupTunnelFinder(ctx, t)

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
		}).
		MinTimes(2)
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

	go tf.poll(ctx, testhelpers.NewPollConfig(100*time.Millisecond)())
	select {
	case <-ctx.Done():
	case <-tChan:
		t.FailNow()
	}
	assert.Eventually(t, func() bool {
		tf.mu.Lock()
		defer tf.mu.Unlock()
		return len(tf.connections) == 0
	}, time.Second, 10*time.Millisecond) // wait for goroutines to stop
}

func TestTunnelFinder_PollStartsGoroutineForEachUrl(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var dialStartWg, dialStopWg sync.WaitGroup
	dialStartWg.Add(2)
	go func() {
		dialStartWg.Wait()
		time.Sleep(200 * time.Millisecond)
		cancel()
	}()
	tChan, tf, querier, rpcApi, kasPool := setupTunnelFinder(ctx, t)

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
		}).
		MinTimes(2)
	kasPool.EXPECT().
		Dial(gomock.Any(), "grpc://pipe").
		DoAndReturn(func(ctx context.Context, targetUrl string) (grpctool.PoolConn, error) {
			dialStartWg.Done()
			<-ctx.Done() // block to simulate a long running dial
			dialStopWg.Wait()
			return nil, ctx.Err()
		})
	kasPool.EXPECT().
		Dial(gomock.Any(), "grpc://pipe2").
		DoAndReturn(func(ctx context.Context, targetUrl string) (grpctool.PoolConn, error) {
			dialStartWg.Done()
			<-ctx.Done() // block to simulate a long running dial
			dialStopWg.Wait()
			return nil, ctx.Err()
		})
	rpcApi.EXPECT().
		HandleProcessingError(gomock.Any(), testhelpers.AgentId, gomock.Any(), gomock.Any()).
		Times(2)

	dialStopWg.Add(1)
	go func() {
		defer dialStopWg.Done() // unblock dials once polling is done. This is to avoid flakes.
		tf.poll(ctx, testhelpers.NewPollConfig(100*time.Millisecond)())
	}()
	select {
	case <-ctx.Done():
	case <-tChan:
		t.FailNow()
	}
	assert.Eventually(t, func() bool {
		tf.mu.Lock()
		defer tf.mu.Unlock()
		return len(tf.connections) == 0
	}, time.Second, 10*time.Millisecond) // wait for goroutines to stop
}

func setupTunnelFinder(ctx context.Context, t *testing.T) (chan readyTunnel, *tunnelFinder, *mock_reverse_tunnel_tracker.MockQuerier, *mock_modserver.MockRpcApi, *MockKasPool) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	querier := mock_reverse_tunnel_tracker.NewMockQuerier(ctrl)
	rpcApi := mock_modserver.NewMockRpcApi(ctrl)
	kasPool := NewMockKasPool(ctrl)

	tChan := make(chan readyTunnel)
	tf := &tunnelFinder{
		log:           zaptest.NewLogger(t),
		kasPool:       kasPool,
		tunnelQuerier: querier,
		rpcApi:        rpcApi,
		fullMethod:    "/gitlab.agent.grpctool.test.Testing/RequestResponse",
		agentId:       testhelpers.AgentId,
		outgoingCtx:   ctx,
		foundTunnel:   tChan,
		connections:   make(map[string]kasConnAttempt),
	}
	return tChan, tf, querier, rpcApi, kasPool
}
