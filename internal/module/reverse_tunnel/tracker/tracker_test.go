package tracker

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/reverse_tunnel/info"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/redistool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/syncz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_redis"
	"go.uber.org/zap/zaptest"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
)

var (
	_ Registerer = &RedisTracker{}
	_ Tracker    = &RedisTracker{}
	_ Querier    = &RedisTracker{}
)

func TestRegisterConnection(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	r, hash, ti := setupTracker(t)

	hash.EXPECT().
		Set(ti.AgentId, ti.ConnectionId, gomock.Any()).
		Return(func(ctx context.Context) error {
			cancel()
			return nil
		})

	go func() {
		assert.NoError(t, r.RegisterTunnel(context.Background(), ti))
	}()

	require.NoError(t, r.Run(ctx))
}

func TestUnregisterConnection(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	r, hash, ti := setupTracker(t)

	gomock.InOrder(
		hash.EXPECT().
			Set(ti.AgentId, ti.ConnectionId, gomock.Any()).
			Return(nopIOFunc),
		hash.EXPECT().
			Unset(ti.AgentId, ti.ConnectionId).
			Return(func(ctx context.Context) error {
				cancel()
				return nil
			}),
	)

	go func() {
		assert.NoError(t, r.RegisterTunnel(context.Background(), ti))
		assert.NoError(t, r.UnregisterTunnel(context.Background(), ti))
	}()

	require.NoError(t, r.Run(ctx))
}

func TestGC(t *testing.T) {
	r, hash, _ := setupTracker(t)

	wasCalled := false

	hash.EXPECT().
		GC().
		Return(func(_ context.Context) (int, error) {
			wasCalled = true
			return 3, nil
		})

	deleted, err := r.runGC(context.Background())
	require.NoError(t, err)
	assert.EqualValues(t, 3, deleted)
	assert.True(t, wasCalled)
}

func TestRefreshRegistrations(t *testing.T) {
	r, hash, _ := setupTracker(t)

	hash.EXPECT().
		Refresh(gomock.Any()).
		Return(nopIOFunc)
	assert.NoError(t, r.refreshRegistrations(context.Background(), time.Now()))
}

func TestGetTunnelsByAgentId_HappyPath(t *testing.T) {
	r, hash, ti := setupTracker(t)
	tiBytes, err := proto.Marshal(ti)
	require.NoError(t, err)
	hash.EXPECT().
		Scan(gomock.Any(), ti.AgentId, gomock.Any()).
		Do(func(ctx context.Context, key interface{}, cb redistool.ScanCallback) (int, error) {
			var done bool
			done, err = cb(tiBytes, nil)
			if err != nil || done {
				return 0, err
			}
			return 0, nil
		})
	var cbCalled int
	err = r.GetTunnelsByAgentId(context.Background(), ti.AgentId, func(tunnelInfo *TunnelInfo) (bool, error) {
		cbCalled++
		assert.Empty(t, cmp.Diff(tunnelInfo, ti, protocmp.Transform()))
		return false, nil
	})
	require.NoError(t, err)
	assert.EqualValues(t, 1, cbCalled)
}

func TestGetTunnelsByAgentId_ScanError(t *testing.T) {
	r, hash, ti := setupTracker(t)
	hash.EXPECT().
		Scan(gomock.Any(), ti.AgentId, gomock.Any()).
		Do(func(ctx context.Context, key interface{}, cb redistool.ScanCallback) (int, error) {
			done, err := cb(nil, errors.New("intended error"))
			require.NoError(t, err)
			assert.False(t, done)
			return 0, nil
		})
	err := r.GetTunnelsByAgentId(context.Background(), ti.AgentId, func(tunnelInfo *TunnelInfo) (bool, error) {
		require.FailNow(t, "unexpected call")
		return false, nil
	})
	require.NoError(t, err)
}

func TestGetTunnelsByAgentId_UnmarshalError(t *testing.T) {
	r, hash, ti := setupTracker(t)
	hash.EXPECT().
		Scan(gomock.Any(), ti.AgentId, gomock.Any()).
		Do(func(ctx context.Context, key interface{}, cb redistool.ScanCallback) (int, error) {
			done, err := cb([]byte{1, 2, 3}, nil) // invalid data
			require.NoError(t, err)               // ignores error to keep going
			assert.False(t, done)
			return 0, nil
		})
	err := r.GetTunnelsByAgentId(context.Background(), ti.AgentId, func(tunnelInfo *TunnelInfo) (bool, error) {
		require.FailNow(t, "unexpected call")
		return false, nil
	})
	require.NoError(t, err)
}

func setupTracker(t *testing.T) (*RedisTracker, *mock_redis.MockExpiringHashInterface[int64], *TunnelInfo) {
	ctrl := gomock.NewController(t)
	hash := mock_redis.NewMockExpiringHashInterface[int64](ctrl)
	ti := &TunnelInfo{
		AgentDescriptor: &info.AgentDescriptor{
			Services: []*info.Service{
				{
					Name: "bla",
					Methods: []*info.Method{
						{
							Name: "bab",
						},
					},
				},
			},
		},
		ConnectionId: 123,
		AgentId:      543,
	}
	return &RedisTracker{
		log:              zaptest.NewLogger(t),
		refreshPeriod:    time.Minute,
		gcPeriod:         time.Minute,
		mu:               syncz.NewMutex(),
		tunnelsByAgentId: hash,
	}, hash, ti
}

func TestTunnelInfoSize(t *testing.T) {
	tiBytes, err := proto.Marshal(&TunnelInfo{
		AgentDescriptor: &info.AgentDescriptor{
			Services: []*info.Service{},
		},
		ConnectionId: 1231232,
		AgentId:      123123,
		KasUrl:       "grpcs://123.123.123.123:123",
	})
	require.NoError(t, err)
	data, err := proto.Marshal(&redistool.ExpiringValue{
		ExpiresAt: time.Now().Unix(),
		Value:     tiBytes,
	})
	require.NoError(t, err)
	t.Log(len(data))
}

func TestSupportsServiceAndMethod(t *testing.T) {
	ti := TunnelInfo{
		AgentDescriptor: &info.AgentDescriptor{
			Services: []*info.Service{
				{
					Name: "empire.fleet.DeathStar",
					Methods: []*info.Method{
						{
							Name: "BlastPlanet",
						},
					},
				},
			},
		},
	}
	assert.True(t, ti.SupportsServiceAndMethod("empire.fleet.DeathStar", "BlastPlanet"))
	assert.False(t, ti.SupportsServiceAndMethod("empire.fleet.DeathStar", "Explode"))
	assert.False(t, ti.SupportsServiceAndMethod("empire.fleet.hangar.DeathStar", "BlastPlanet"))
	assert.False(t, ti.SupportsServiceAndMethod("empire.fleet.hangar.DeathStar", "Debug"))
}

func nopIOFunc(ctx context.Context) error {
	return nil
}
