package tunnel

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/redistool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_redis"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/testhelpers"
	"go.uber.org/mock/gomock"
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
	r, hash := setupTracker(t)

	hash.EXPECT().
		Set(gomock.Any(), testhelpers.AgentId, selfUrl, gomock.Any())

	assert.NoError(t, r.RegisterTunnel(context.Background(), testhelpers.AgentId))
}

func TestUnregisterConnection(t *testing.T) {
	r, hash := setupTracker(t)

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
	r, hash := setupTracker(t)

	gomock.InOrder(
		hash.EXPECT().
			Set(gomock.Any(), testhelpers.AgentId, selfUrl, gomock.Any()),
		hash.EXPECT().
			Set(gomock.Any(), testhelpers.AgentId, selfUrl, gomock.Any()),
		hash.EXPECT().
			Unset(gomock.Any(), testhelpers.AgentId, selfUrl),
		hash.EXPECT().
			Unset(gomock.Any(), testhelpers.AgentId, selfUrl),
	)

	assert.NoError(t, r.RegisterTunnel(context.Background(), testhelpers.AgentId))
	assert.NoError(t, r.RegisterTunnel(context.Background(), testhelpers.AgentId))
	assert.NoError(t, r.UnregisterTunnel(context.Background(), testhelpers.AgentId))
	assert.NoError(t, r.UnregisterTunnel(context.Background(), testhelpers.AgentId))
}

func TestKasUrlsByAgentId_HappyPath(t *testing.T) {
	r, hash := setupTracker(t)
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
	kasUrls, err := r.KasUrlsByAgentId(context.Background(), testhelpers.AgentId)
	require.NoError(t, err)
	assert.Equal(t, []string{selfUrl}, kasUrls)
}

func TestKasUrlsByAgentId_ScanError(t *testing.T) {
	r, hash := setupTracker(t)
	hash.EXPECT().
		Scan(gomock.Any(), testhelpers.AgentId, gomock.Any()).
		Do(func(ctx context.Context, key interface{}, cb redistool.ScanCallback) (int, error) {
			done, err := cb("", nil, errors.New("intended error"))
			require.NoError(t, err)
			assert.False(t, done)
			return 0, nil
		})
	kasUrls, err := r.KasUrlsByAgentId(context.Background(), testhelpers.AgentId)
	assert.EqualError(t, err, "intended error")
	assert.Empty(t, kasUrls)
}

func setupTracker(t *testing.T) (*RedisTracker, *mock_redis.MockExpiringHashInterface[int64, string]) {
	ctrl := gomock.NewController(t)
	hash := mock_redis.NewMockExpiringHashInterface[int64, string](ctrl)
	return &RedisTracker{
		ownPrivateApiUrl: selfUrl,
		tunnelsByAgentId: hash,
	}, hash
}
