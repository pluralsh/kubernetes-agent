package tracker_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/reverse_tunnel/tracker"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_reverse_tunnel_tracker"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/testhelpers"
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
	aq := tracker.NewAggregatingQuerier(zaptest.NewLogger(t), q, api, testhelpers.NewPollConfig(time.Minute))
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
	aq := tracker.NewAggregatingQuerier(zaptest.NewLogger(t), q, api, testhelpers.NewPollConfig(time.Minute))
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
	aq := tracker.NewAggregatingQuerier(zaptest.NewLogger(t), q, api, testhelpers.NewPollConfig(time.Minute))
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
