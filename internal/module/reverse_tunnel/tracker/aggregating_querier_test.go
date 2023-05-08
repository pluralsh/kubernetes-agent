package tracker_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/reverse_tunnel/tracker"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_reverse_tunnel_tracker"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/testhelpers"
	"go.uber.org/zap/zaptest"
	"k8s.io/apimachinery/pkg/util/wait"
)

var (
	_ tracker.PollingQuerier = (*tracker.AggregatingQuerier)(nil)
)

func TestPollKasUrlsByAgentId_OnlyStartsSinglePoll(t *testing.T) {
	ctrl := gomock.NewController(t)
	q := mock_reverse_tunnel_tracker.NewMockQuerier(ctrl)
	q.EXPECT().
		KasUrlsByAgentId(gomock.Any(), testhelpers.AgentId, gomock.Any())
	api := mock_modserver.NewMockApi(ctrl)
	aq := tracker.NewAggregatingQuerier(zaptest.NewLogger(t), q, api, testhelpers.NewPollConfig(time.Minute), time.Minute)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	go aq.PollKasUrlsByAgentId(ctx, testhelpers.AgentId, func(kasUrls []string) {
		assert.Fail(t, "unexpected call")
	})
	aq.PollKasUrlsByAgentId(ctx, testhelpers.AgentId, func(kasUrls []string) {
		assert.Fail(t, "unexpected call")
	})
}

func TestPollKasUrlsByAgentId_PollingCycle(t *testing.T) {
	ctrl := gomock.NewController(t)
	q := mock_reverse_tunnel_tracker.NewMockQuerier(ctrl)
	q.EXPECT().
		KasUrlsByAgentId(gomock.Any(), testhelpers.AgentId, gomock.Any()).
		Do(func(ctx context.Context, agentId int64, cb tracker.KasUrlsByAgentIdCallback) {
			done, err := cb("url1")
			assert.NoError(t, err)
			assert.False(t, done)
			done, err = cb("url2")
			assert.NoError(t, err)
			assert.False(t, done)
		})
	api := mock_modserver.NewMockApi(ctrl)
	aq := tracker.NewAggregatingQuerier(zaptest.NewLogger(t), q, api, testhelpers.NewPollConfig(time.Minute), time.Minute)
	call := 0
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	aq.PollKasUrlsByAgentId(ctx, testhelpers.AgentId, func(kasUrls []string) {
		switch call {
		case 0:
			assert.Equal(t, []string{"url1", "url2"}, kasUrls)
		default:
			assert.FailNow(t, "unexpected invocation")
		}
		call++
	})
}

func TestPollKasUrlsByAgentId_CacheAfterStopped(t *testing.T) {
	ctrl := gomock.NewController(t)
	q := mock_reverse_tunnel_tracker.NewMockQuerier(ctrl)
	gomock.InOrder(
		q.EXPECT().
			KasUrlsByAgentId(gomock.Any(), testhelpers.AgentId, gomock.Any()).
			Do(func(ctx context.Context, agentId int64, cb tracker.KasUrlsByAgentIdCallback) {
				_, err := cb("url1")
				assert.NoError(t, err)
			}),
		q.EXPECT().
			KasUrlsByAgentId(gomock.Any(), testhelpers.AgentId, gomock.Any()).
			Do(func(ctx context.Context, agentId int64, cb tracker.KasUrlsByAgentIdCallback) {
				_, err := cb("url2")
				assert.NoError(t, err)
			}),
	)
	api := mock_modserver.NewMockApi(ctrl)
	aq := tracker.NewAggregatingQuerier(zaptest.NewLogger(t), q, api, testhelpers.NewPollConfig(time.Minute), time.Minute)
	ctx, cancel := context.WithCancel(context.Background())
	aq.PollKasUrlsByAgentId(ctx, testhelpers.AgentId, func(kasUrls []string) {
		assert.Equal(t, []string{"url1"}, kasUrls)
		cancel()
	})
	kasUrls := aq.CachedKasUrlsByAgentId(testhelpers.AgentId) // from cache
	assert.Equal(t, []string{"url1"}, kasUrls)
	ctx, cancel = context.WithCancel(context.Background())
	aq.PollKasUrlsByAgentId(ctx, testhelpers.AgentId, func(kasUrls []string) {
		assert.Equal(t, []string{"url2"}, kasUrls) // from redis
		cancel()
	})
}

func TestPollKasUrlsByAgentId_CacheWhenRunning(t *testing.T) {
	ctrl := gomock.NewController(t)
	q := mock_reverse_tunnel_tracker.NewMockQuerier(ctrl)
	start1 := make(chan struct{})
	start2 := make(chan struct{})
	gomock.InOrder(
		q.EXPECT().
			KasUrlsByAgentId(gomock.Any(), testhelpers.AgentId, gomock.Any()).
			Do(func(ctx context.Context, agentId int64, cb tracker.KasUrlsByAgentIdCallback) {
				_, err := cb("url1")
				assert.NoError(t, err)
			}),
		q.EXPECT().
			KasUrlsByAgentId(gomock.Any(), testhelpers.AgentId, gomock.Any()).
			Do(func(ctx context.Context, agentId int64, cb tracker.KasUrlsByAgentIdCallback) {
				close(start1)                      // start concurrent query
				<-start2                           // wait for the concurrent query to consume item from cache
				time.Sleep(100 * time.Millisecond) // wait for aq.PollKasUrlsByAgentId() to register second callback
				_, err := cb("url2")
				assert.NoError(t, err)
			}),
	)
	api := mock_modserver.NewMockApi(ctrl)
	aq := tracker.NewAggregatingQuerier(zaptest.NewLogger(t), q, api, testhelpers.NewPollConfig(time.Second), time.Minute)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	count1 := 0
	go aq.PollKasUrlsByAgentId(ctx, testhelpers.AgentId, func(kasUrls []string) {
		switch count1 {
		case 0:
			assert.Equal(t, []string{"url1"}, kasUrls) // first call
		case 1:
			assert.Equal(t, []string{"url2"}, kasUrls) // second call
		default:
			assert.FailNow(t, "unexpected invocation")
		}
		count1++
	})
	<-start1
	kasUrls := aq.CachedKasUrlsByAgentId(testhelpers.AgentId)
	assert.Equal(t, []string{"url1"}, kasUrls) // from cache
	close(start2)
	count2 := 0
	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()
	aq.PollKasUrlsByAgentId(ctx2, testhelpers.AgentId, func(kasUrls []string) {
		switch count2 {
		case 0:
			assert.Equal(t, []string{"url2"}, kasUrls) // from redis
			cancel2()
		default:
			assert.FailNow(t, "unexpected invocation")
		}
		count2++
	})
	assert.EqualValues(t, 1, count2)
}

func TestPollKasUrlsByAgentId_GcRemovesExpiredCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	q := mock_reverse_tunnel_tracker.NewMockQuerier(ctrl)
	q.EXPECT().
		KasUrlsByAgentId(gomock.Any(), testhelpers.AgentId, gomock.Any()).
		Do(func(ctx context.Context, agentId int64, cb tracker.KasUrlsByAgentIdCallback) {
			_, err := cb("url1")
			assert.NoError(t, err)
		})
	api := mock_modserver.NewMockApi(ctrl)
	gcPeriod := time.Second
	aq := tracker.NewAggregatingQuerier(zaptest.NewLogger(t), q, api, testhelpers.NewPollConfig(time.Minute), gcPeriod)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	aq.PollKasUrlsByAgentId(ctx, testhelpers.AgentId, func(kasUrls []string) {
		cancel()
	})
	ctx, cancel = context.WithCancel(context.Background())
	var wg wait.Group
	defer wg.Wait()
	defer cancel()
	wg.Start(func() {
		_ = aq.Run(ctx)
	})
	time.Sleep(gcPeriod * 2)
	kasUrls := aq.CachedKasUrlsByAgentId(testhelpers.AgentId)
	assert.Empty(t, kasUrls)
}
