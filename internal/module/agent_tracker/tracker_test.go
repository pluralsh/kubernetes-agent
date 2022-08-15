package agent_tracker

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/redistool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_redis"
	"go.uber.org/zap/zaptest"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	_ Registerer                 = &RedisTracker{}
	_ Querier                    = &RedisTracker{}
	_ Tracker                    = &RedisTracker{}
	_ ConnectedAgentInfoCallback = (&ConnectedAgentInfoCollector{}).Collect
)

func TestRegisterConnection_HappyPath(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	r, connectedAgents, byAgentId, byProjectId, info := setupTracker(t)

	byProjectId.EXPECT().
		Set(info.ProjectId, info.ConnectionId, gomock.Any()).
		Return(nopIOFunc)
	byAgentId.EXPECT().
		Set(info.AgentId, info.ConnectionId, gomock.Any()).
		Return(nopIOFunc)
	connectedAgents.EXPECT().
		Set(nil, info.AgentId, gomock.Any()).
		Return(func(ctx context.Context) error {
			cancel()
			return nil
		})

	go func() {
		assert.True(t, r.RegisterConnection(context.Background(), info))
	}()

	require.NoError(t, r.Run(ctx))
}

func TestRegisterConnection_AllCalledOnError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	r, connectedAgents, byAgentId, byProjectId, info := setupTracker(t)

	byProjectId.EXPECT().
		Set(info.ProjectId, info.ConnectionId, gomock.Any()).
		Return(func(ctx context.Context) error { return errors.New("err1") })
	byAgentId.EXPECT().
		Set(info.AgentId, info.ConnectionId, gomock.Any()).
		Return(func(ctx context.Context) error { return errors.New("err1") })
	connectedAgents.EXPECT().
		Set(nil, info.AgentId, gomock.Any()).
		Return(func(ctx context.Context) error {
			cancel()
			return errors.New("err3")
		})

	go func() {
		assert.True(t, r.RegisterConnection(context.Background(), info))
	}()

	require.NoError(t, r.Run(ctx))
}

func TestUnregisterConnection_HappyPath(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	r, connectedAgents, byAgentId, byProjectId, info := setupTracker(t)

	gomock.InOrder(
		byProjectId.EXPECT().
			Set(info.ProjectId, info.ConnectionId, gomock.Any()).
			Return(nopIOFunc),
		byProjectId.EXPECT().
			Unset(info.ProjectId, info.ConnectionId).
			Return(nopIOFunc),
	)
	gomock.InOrder(
		byAgentId.EXPECT().
			Set(info.AgentId, info.ConnectionId, gomock.Any()).
			Return(nopIOFunc),
		byAgentId.EXPECT().
			Unset(info.AgentId, info.ConnectionId).
			Return(func(ctx context.Context) error {
				cancel()
				return nil
			}),
	)
	gomock.InOrder(
		connectedAgents.EXPECT().
			Set(nil, info.AgentId, gomock.Any()).
			Return(nopIOFunc),
		connectedAgents.EXPECT().
			Forget(nil, info.AgentId),
	)
	go func() {
		assert.True(t, r.RegisterConnection(context.Background(), info))
		assert.True(t, r.UnregisterConnection(context.Background(), info))
	}()

	require.NoError(t, r.Run(ctx))
}

func TestUnregisterConnection_AllCalledOnError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	r, connectedAgents, byAgentId, byProjectId, info := setupTracker(t)

	gomock.InOrder(
		byProjectId.EXPECT().
			Set(info.ProjectId, info.ConnectionId, gomock.Any()).
			Return(nopIOFunc),
		byProjectId.EXPECT().
			Unset(info.ProjectId, info.ConnectionId).
			Return(func(ctx context.Context) error { return errors.New("err1") }),
	)
	gomock.InOrder(
		byAgentId.EXPECT().
			Set(info.AgentId, info.ConnectionId, gomock.Any()).
			Return(nopIOFunc),
		byAgentId.EXPECT().
			Unset(info.AgentId, info.ConnectionId).
			Return(func(ctx context.Context) error {
				cancel()
				return errors.New("err1")
			}),
	)
	gomock.InOrder(
		connectedAgents.EXPECT().
			Set(nil, info.AgentId, gomock.Any()).
			Return(nopIOFunc),
		connectedAgents.EXPECT().
			Forget(nil, info.AgentId),
	)

	go func() {
		assert.True(t, r.RegisterConnection(context.Background(), info))
		assert.True(t, r.UnregisterConnection(context.Background(), info))
	}()

	require.NoError(t, r.Run(ctx))
}

func TestGC_HappyPath(t *testing.T) {
	r, connectedAgents, byAgentId, byProjectId, _ := setupTracker(t)

	wasCalled1 := false
	wasCalled2 := false
	wasCalled3 := false

	connectedAgents.EXPECT().
		GC().
		Return(func(_ context.Context) (int, error) {
			wasCalled3 = true
			return 3, nil
		})

	byAgentId.EXPECT().
		GC().
		Return(func(_ context.Context) (int, error) {
			wasCalled2 = true
			return 2, nil
		})

	byProjectId.EXPECT().
		GC().
		Return(func(_ context.Context) (int, error) {
			wasCalled1 = true
			return 1, nil
		})

	r.maybeRunGCAsync(context.Background())
	assert.Eventually(t, func() bool {
		return !r.gc.IsRunning()
	}, time.Second, 10*time.Millisecond)
	assert.True(t, wasCalled1)
	assert.True(t, wasCalled2)
	assert.True(t, wasCalled3)
}

func TestGC_AllCalledOnError(t *testing.T) {
	r, connectedAgents, byAgentId, byProjectId, _ := setupTracker(t)

	wasCalled1 := false
	wasCalled2 := false
	wasCalled3 := false

	connectedAgents.EXPECT().
		GC().
		Return(func(_ context.Context) (int, error) {
			wasCalled3 = true
			return 3, errors.New("err3")
		})

	byAgentId.EXPECT().
		GC().
		Return(func(_ context.Context) (int, error) {
			wasCalled2 = true
			return 2, errors.New("err2")
		})

	byProjectId.EXPECT().
		GC().
		Return(func(_ context.Context) (int, error) {
			wasCalled1 = true
			return 1, errors.New("err1")
		})

	r.maybeRunGCAsync(context.Background())
	assert.Eventually(t, func() bool {
		return !r.gc.IsRunning()
	}, time.Second, 10*time.Millisecond)
	assert.True(t, wasCalled1)
	assert.True(t, wasCalled2)
	assert.True(t, wasCalled3)
}

func TestRefresh_HappyPath(t *testing.T) {
	r, connectedAgents, byAgentId, byProjectId, _ := setupTracker(t)

	connectedAgents.EXPECT().
		Refresh(gomock.Any()).
		Return(nopIOFunc)
	byAgentId.EXPECT().
		Refresh(gomock.Any()).
		Return(nopIOFunc)
	byProjectId.EXPECT().
		Refresh(gomock.Any()).
		Return(nopIOFunc)
	assert.NoError(t, r.refreshRegistrations(context.Background(), time.Now()))
}

func TestRefresh_AllCalledOnError(t *testing.T) {
	r, connectedAgents, byAgentId, byProjectId, _ := setupTracker(t)

	connectedAgents.EXPECT().
		Refresh(gomock.Any()).
		Return(func(ctx context.Context) error { return errors.New("err3") })
	byAgentId.EXPECT().
		Refresh(gomock.Any()).
		Return(func(ctx context.Context) error { return errors.New("err1") })
	byProjectId.EXPECT().
		Refresh(gomock.Any()).
		Return(func(ctx context.Context) error { return errors.New("err2") })
	assert.Error(t, r.refreshRegistrations(context.Background(), time.Now()))
}

func TestGetConnectionsByProjectId_HappyPath(t *testing.T) {
	r, _, _, byProjectId, info := setupTracker(t)
	any, err := anypb.New(info)
	require.NoError(t, err)
	byProjectId.EXPECT().
		Scan(gomock.Any(), info.ProjectId, gomock.Any()).
		Do(func(ctx context.Context, key interface{}, cb redistool.ScanCallback) (int, error) {
			var done bool
			done, err = cb(any, nil)
			if err != nil || done {
				return 0, err
			}
			return 0, nil
		})
	var cbCalled int
	err = r.GetConnectionsByProjectId(context.Background(), info.ProjectId, func(i *ConnectedAgentInfo) (done bool, err error) {
		cbCalled++
		assert.Empty(t, cmp.Diff(i, info, protocmp.Transform()))
		return false, nil
	})
	require.NoError(t, err)
	assert.EqualValues(t, 1, cbCalled)
}

func TestGetConnectionsByProjectId_ScanError(t *testing.T) {
	r, _, _, byProjectId, info := setupTracker(t)
	byProjectId.EXPECT().
		Scan(gomock.Any(), info.ProjectId, gomock.Any()).
		Do(func(ctx context.Context, key interface{}, cb redistool.ScanCallback) (int, error) {
			done, err := cb(nil, errors.New("intended error"))
			require.NoError(t, err)
			assert.False(t, done)
			return 0, nil
		})
	err := r.GetConnectionsByProjectId(context.Background(), info.ProjectId, func(i *ConnectedAgentInfo) (done bool, err error) {
		require.FailNow(t, "unexpected call")
		return false, nil
	})
	require.NoError(t, err)
}

func TestGetConnectionsByProjectId_UnmarshalError(t *testing.T) {
	r, _, _, byProjectId, info := setupTracker(t)
	byProjectId.EXPECT().
		Scan(gomock.Any(), info.ProjectId, gomock.Any()).
		Do(func(ctx context.Context, key interface{}, cb redistool.ScanCallback) (int, error) {
			done, err := cb(&anypb.Any{
				TypeUrl: "gitlab.agent.agent_tracker.ConnectedAgentInfo", // valid
				Value:   []byte{1, 2, 3},                                 // invalid
			}, nil)
			require.NoError(t, err) // ignores error to keep going
			assert.False(t, done)
			return 0, nil
		})
	err := r.GetConnectionsByProjectId(context.Background(), info.ProjectId, func(i *ConnectedAgentInfo) (done bool, err error) {
		require.FailNow(t, "unexpected call")
		return false, nil
	})
	require.NoError(t, err)
}

func TestGetConnectionsByAgentId_HappyPath(t *testing.T) {
	r, _, byAgentId, _, info := setupTracker(t)
	any, err := anypb.New(info)
	require.NoError(t, err)
	byAgentId.EXPECT().
		Scan(gomock.Any(), info.AgentId, gomock.Any()).
		Do(func(ctx context.Context, key interface{}, cb redistool.ScanCallback) (int, error) {
			var done bool
			done, err = cb(any, nil)
			if err != nil || done {
				return 0, err
			}
			return 0, nil
		})
	var cbCalled int
	err = r.GetConnectionsByAgentId(context.Background(), info.AgentId, func(i *ConnectedAgentInfo) (done bool, err error) {
		cbCalled++
		assert.Empty(t, cmp.Diff(i, info, protocmp.Transform()))
		return false, nil
	})
	require.NoError(t, err)
	assert.EqualValues(t, 1, cbCalled)
}

func TestGetConnectionsByAgentId_ScanError(t *testing.T) {
	r, _, byAgentId, _, info := setupTracker(t)
	byAgentId.EXPECT().
		Scan(gomock.Any(), info.AgentId, gomock.Any()).
		Do(func(ctx context.Context, key interface{}, cb redistool.ScanCallback) (int, error) {
			done, err := cb(nil, errors.New("intended error"))
			require.NoError(t, err)
			assert.False(t, done)
			return 0, nil
		})
	err := r.GetConnectionsByAgentId(context.Background(), info.AgentId, func(i *ConnectedAgentInfo) (done bool, err error) {
		require.FailNow(t, "unexpected call")
		return false, nil
	})
	require.NoError(t, err)
}

func TestGetConnectionsByAgentId_UnmarshalError(t *testing.T) {
	r, _, byAgentId, _, info := setupTracker(t)
	byAgentId.EXPECT().
		Scan(gomock.Any(), info.AgentId, gomock.Any()).
		Do(func(ctx context.Context, key interface{}, cb redistool.ScanCallback) (int, error) {
			done, err := cb(&anypb.Any{
				TypeUrl: "gitlab.agent.agent_tracker.ConnectedAgentInfo", // valid
				Value:   []byte{1, 2, 3},                                 // invalid
			}, nil)
			require.NoError(t, err) // ignores error to keep going
			assert.False(t, done)
			return 0, nil
		})
	err := r.GetConnectionsByAgentId(context.Background(), info.AgentId, func(i *ConnectedAgentInfo) (done bool, err error) {
		require.FailNow(t, "unexpected call")
		return false, nil
	})
	require.NoError(t, err)
}

func TestGetConnectedAgentsCount_HappyPath(t *testing.T) {
	r, connectedAgents, _, _, _ := setupTracker(t)
	connectedAgents.EXPECT().
		Len(gomock.Any(), nil).
		Return(int64(1), nil)
	size, err := r.GetConnectedAgentsCount(context.Background())
	require.NoError(t, err)
	assert.EqualValues(t, 1, size)
}

func TestGetConnectedAgentsCount_LenError(t *testing.T) {
	r, connectedAgents, _, _, _ := setupTracker(t)
	connectedAgents.EXPECT().
		Len(gomock.Any(), gomock.Any()).
		Return(int64(0), errors.New("intended error"))
	size, err := r.GetConnectedAgentsCount(context.Background())
	require.Error(t, err)
	assert.Zero(t, size)
}

func nopIOFunc(ctx context.Context) error {
	return nil
}

func setupTracker(t *testing.T) (*RedisTracker, *mock_redis.MockExpiringHashInterface, *mock_redis.MockExpiringHashInterface, *mock_redis.MockExpiringHashInterface, *ConnectedAgentInfo) {
	ctrl := gomock.NewController(t)
	connectedAgents := mock_redis.NewMockExpiringHashInterface(ctrl)
	byAgentId := mock_redis.NewMockExpiringHashInterface(ctrl)
	byProjectId := mock_redis.NewMockExpiringHashInterface(ctrl)
	tr := &RedisTracker{
		log:                    zaptest.NewLogger(t),
		refreshPeriod:          time.Minute,
		gcPeriod:               time.Minute,
		connectionsByAgentId:   byAgentId,
		connectionsByProjectId: byProjectId,
		connectedAgents:        connectedAgents,
		toRegister:             make(chan *ConnectedAgentInfo),
		toUnregister:           make(chan *ConnectedAgentInfo),
	}
	return tr, connectedAgents, byAgentId, byProjectId, connInfo()
}

func connInfo() *ConnectedAgentInfo {
	return &ConnectedAgentInfo{
		AgentMeta: &modshared.AgentMeta{
			Version:      "v1.2.3",
			CommitId:     "123123",
			PodNamespace: "ns",
			PodName:      "name",
		},
		ConnectedAt:  timestamppb.Now(),
		ConnectionId: 123,
		AgentId:      345,
		ProjectId:    456,
	}
}
