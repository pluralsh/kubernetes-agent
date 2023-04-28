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
	go aq.PollKasUrlsByAgentId(ctx, testhelpers.AgentId, func(newCycle bool, kasUrl string) bool {
		assert.Fail(t, "unexpected call")
		return false
	})
	aq.PollKasUrlsByAgentId(ctx, testhelpers.AgentId, func(newCycle bool, kasUrl string) bool {
		assert.Fail(t, "unexpected call")
		return false
	})
}

func TestPollKasUrlsByAgentId_StopsPollingOnDoneTrue(t *testing.T) {
	ctrl := gomock.NewController(t)
	q := mock_reverse_tunnel_tracker.NewMockQuerier(ctrl)
	q.EXPECT().
		KasUrlsByAgentId(gomock.Any(), testhelpers.AgentId, gomock.Any()).
		Do(func(ctx context.Context, agentId int64, cb tracker.KasUrlsByAgentIdCallback) {
			time.Sleep(time.Second) // give both goroutines some time to start and call PollKasUrlsByAgentId()
			done, err := cb("url1")
			assert.NoError(t, err)
			assert.False(t, done)
		})
	api := mock_modserver.NewMockApi(ctrl)
	aq := tracker.NewAggregatingQuerier(zaptest.NewLogger(t), q, api, testhelpers.NewPollConfig(time.Minute), time.Minute)
	var wg wait.Group
	defer wg.Wait()
	wg.Start(func() {
		aq.PollKasUrlsByAgentId(context.Background(), testhelpers.AgentId, func(newCycle bool, kasUrl string) bool {
			//time.Sleep(time.Second)
			return true
		})
	})
	wg.Start(func() {
		aq.PollKasUrlsByAgentId(context.Background(), testhelpers.AgentId, func(newCycle bool, kasUrl string) bool {
			//time.Sleep(time.Second)
			return true
		})
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
	aq.PollKasUrlsByAgentId(ctx, testhelpers.AgentId, func(newCycle bool, kasUrl string) bool {
		switch call {
		case 0:
			assert.True(t, newCycle)
			assert.Equal(t, "url1", kasUrl)
		default:
			assert.False(t, newCycle)
			assert.Equal(t, "url2", kasUrl)
		}
		call++
		return false
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
	aq.PollKasUrlsByAgentId(context.Background(), testhelpers.AgentId, func(newCycle bool, kasUrl string) bool {
		assert.Equal(t, "url1", kasUrl)
		return true
	})
	count := 0
	aq.PollKasUrlsByAgentId(context.Background(), testhelpers.AgentId, func(newCycle bool, kasUrl string) bool {
		switch count {
		case 0:
			assert.Equal(t, "url1", kasUrl) // from cache
			count++
			return false
		case 1:
			assert.Equal(t, "url2", kasUrl) // from redis
			count++
			return true
		default:
			assert.FailNow(t, "unexpected invocation")
		}
		return true
	})
	assert.EqualValues(t, 2, count)
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
				close(start1) // start concurrent query
				<-start2      // wait for the concurrent query to consume item from cache
				_, err := cb("url2")
				assert.NoError(t, err)
			}),
	)
	api := mock_modserver.NewMockApi(ctrl)
	aq := tracker.NewAggregatingQuerier(zaptest.NewLogger(t), q, api, testhelpers.NewPollConfig(time.Second), time.Minute)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	count1 := 0
	go aq.PollKasUrlsByAgentId(ctx, testhelpers.AgentId, func(newCycle bool, kasUrl string) bool {
		assert.True(t, newCycle)
		switch count1 {
		case 0:
			assert.Equal(t, "url1", kasUrl) // first call
			count1++
			return false
		case 1:
			assert.Equal(t, "url2", kasUrl) // second call
			count1++
			return true
		default:
			assert.FailNow(t, "unexpected invocation")
		}
		return true
	})
	<-start1
	count2 := 0
	aq.PollKasUrlsByAgentId(context.Background(), testhelpers.AgentId, func(newCycle bool, kasUrl string) bool {
		assert.True(t, newCycle)
		switch count2 {
		case 0:
			assert.Equal(t, "url1", kasUrl) // from cache
			count2++
			close(start2)
			return false
		case 1:
			assert.Equal(t, "url2", kasUrl) // from redis
			count2++
			return true
		default:
			assert.FailNow(t, "unexpected invocation")
		}
		return true
	})
	assert.EqualValues(t, 2, count2)
}

func TestPollKasUrlsByAgentId_GcRemovesExpiredCache(t *testing.T) {
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
	gcPeriod := time.Second
	aq := tracker.NewAggregatingQuerier(zaptest.NewLogger(t), q, api, testhelpers.NewPollConfig(time.Minute), gcPeriod)
	aq.PollKasUrlsByAgentId(context.Background(), testhelpers.AgentId, func(newCycle bool, kasUrl string) bool {
		return true
	})
	ctx, cancel := context.WithCancel(context.Background())
	var wg wait.Group
	defer wg.Wait()
	defer cancel()
	wg.Start(func() {
		_ = aq.Run(ctx)
	})
	time.Sleep(gcPeriod * 2)
	aq.PollKasUrlsByAgentId(context.Background(), testhelpers.AgentId, func(newCycle bool, kasUrl string) bool {
		assert.Equal(t, "url2", kasUrl) // not url1 from cache
		return true
	})
}
