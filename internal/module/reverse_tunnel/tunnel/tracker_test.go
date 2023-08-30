package tunnel

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/redistool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/matcher"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_redis"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/testhelpers"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"
)

var (
	_ Registerer = &RedisTracker{}
	_ Tracker    = &RedisTracker{}
	_ Querier    = &RedisTracker{}
)

const (
	selfUrl = "grpc://1.1.1.1:10"
)

func TestRegisterConnection(t *testing.T) {
	r, hash, _ := setupTracker(t)

	hash.EXPECT().
		Set(gomock.Any(), testhelpers.AgentId, selfUrl, gomock.Any())

	assert.NoError(t, r.RegisterTunnel(context.Background(), testhelpers.AgentId))
}

func TestRegisterConnection_TwoConnections(t *testing.T) {
	r, hash, _ := setupTracker(t)

	hash.EXPECT().
		Set(gomock.Any(), testhelpers.AgentId, selfUrl, gomock.Any())

	// Two registrations result in a single Set() call
	assert.NoError(t, r.RegisterTunnel(context.Background(), testhelpers.AgentId)) // first
	assert.NoError(t, r.RegisterTunnel(context.Background(), testhelpers.AgentId)) // second
}

func TestUnregisterConnection(t *testing.T) {
	r, hash, _ := setupTracker(t)

	gomock.InOrder(
		hash.EXPECT().
			Set(gomock.Any(), testhelpers.AgentId, selfUrl, gomock.Any()),
		hash.EXPECT().
			Unset(gomock.Any(), testhelpers.AgentId, selfUrl),
	)

	assert.NoError(t, r.RegisterTunnel(context.Background(), testhelpers.AgentId))
	assert.NoError(t, r.UnregisterTunnel(context.Background(), testhelpers.AgentId))
}

func TestUnregisterConnection_TwoConnections(t *testing.T) {
	r, hash, _ := setupTracker(t)

	gomock.InOrder(
		hash.EXPECT().
			Set(gomock.Any(), testhelpers.AgentId, selfUrl, gomock.Any()),
		hash.EXPECT().
			Unset(gomock.Any(), testhelpers.AgentId, selfUrl),
	)

	assert.NoError(t, r.RegisterTunnel(context.Background(), testhelpers.AgentId))
	assert.NoError(t, r.RegisterTunnel(context.Background(), testhelpers.AgentId))
	assert.NoError(t, r.UnregisterTunnel(context.Background(), testhelpers.AgentId))
	assert.NoError(t, r.UnregisterTunnel(context.Background(), testhelpers.AgentId))
}

// This test ensures Unset() is only called when there are no registered connections i.e. it is NOT called
// for two RegisterTunnel() and a single UnregisterTunnel().
func TestUnregisterConnection_TwoConnections_OneSet(t *testing.T) {
	r, hash, _ := setupTracker(t)

	hash.EXPECT().
		Set(gomock.Any(), testhelpers.AgentId, selfUrl, gomock.Any())

	assert.NoError(t, r.RegisterTunnel(context.Background(), testhelpers.AgentId))
	assert.NoError(t, r.RegisterTunnel(context.Background(), testhelpers.AgentId))
	assert.NoError(t, r.UnregisterTunnel(context.Background(), testhelpers.AgentId))
}

func TestKasUrlsByAgentId_HappyPath(t *testing.T) {
	r, hash, _ := setupTracker(t)
	hash.EXPECT().
		Scan(gomock.Any(), testhelpers.AgentId, gomock.Any()).
		Do(func(ctx context.Context, key interface{}, cb redistool.ScanCallback) (int, error) {
			var done bool
			done, err := cb(selfUrl, nil, nil)
			if err != nil || done {
				return 0, err
			}
			return 0, nil
		})
	var cbCalled int
	err := r.KasUrlsByAgentId(context.Background(), testhelpers.AgentId, func(kasUrl string) (bool, error) {
		cbCalled++
		assert.Equal(t, selfUrl, kasUrl)
		return false, nil
	})
	require.NoError(t, err)
	assert.EqualValues(t, 1, cbCalled)
}

func TestKasUrlsByAgentId_ScanError(t *testing.T) {
	r, hash, api := setupTracker(t)
	gomock.InOrder(
		hash.EXPECT().
			Scan(gomock.Any(), testhelpers.AgentId, gomock.Any()).
			Do(func(ctx context.Context, key interface{}, cb redistool.ScanCallback) (int, error) {
				done, err := cb("", nil, errors.New("intended error"))
				require.NoError(t, err)
				assert.False(t, done)
				return 0, nil
			}),
		api.EXPECT().
			HandleProcessingError(gomock.Any(), gomock.Any(), testhelpers.AgentId, "Redis hash scan", matcher.ErrorEq("intended error")),
	)
	err := r.KasUrlsByAgentId(context.Background(), testhelpers.AgentId, func(kasUrl string) (bool, error) {
		require.FailNow(t, "unexpected call")
		return false, nil
	})
	require.NoError(t, err)
}

func setupTracker(t *testing.T) (*RedisTracker, *mock_redis.MockExpiringHashInterface[int64, string], *mock_modshared.MockApi) {
	ctrl := gomock.NewController(t)
	api := mock_modshared.NewMockApi(ctrl)
	hash := mock_redis.NewMockExpiringHashInterface[int64, string](ctrl)
	return &RedisTracker{
		log:                   zaptest.NewLogger(t),
		api:                   api,
		ownPrivateApiUrl:      selfUrl,
		tunnelsByAgentIdCount: make(map[int64]uint16),
		tunnelsByAgentId:      hash,
	}, hash, api
}
