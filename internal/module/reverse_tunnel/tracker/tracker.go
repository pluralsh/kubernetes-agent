package tracker

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/redistool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/syncz"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const refreshOverlap = 5 * time.Second

type GetTunnelsByAgentIdCallback func(*TunnelInfo) (bool /* done */, error)

type Querier interface {
	GetTunnelsByAgentId(ctx context.Context, agentId int64, cb GetTunnelsByAgentIdCallback) error
}

type Registerer interface {
	// RegisterTunnel registers tunnel with the tracker.
	RegisterTunnel(ctx context.Context, info *TunnelInfo) error
	// UnregisterTunnel unregisters tunnel with the tracker.
	UnregisterTunnel(ctx context.Context, info *TunnelInfo) error
}

type Tracker interface {
	Registerer
	Querier
	Run(ctx context.Context) error
}

type RedisTracker struct {
	log           *zap.Logger
	refreshPeriod time.Duration
	gcPeriod      time.Duration

	// mu protects fields below
	mu               syncz.Mutex
	tunnelsByAgentId redistool.ExpiringHashInterface[int64, int64] // agentId -> connectionId -> TunnelInfo
}

func NewRedisTracker(log *zap.Logger, client redis.UniversalClient, agentKeyPrefix string, ttl, refreshPeriod, gcPeriod time.Duration) *RedisTracker {
	return &RedisTracker{
		log:              log,
		refreshPeriod:    refreshPeriod,
		gcPeriod:         gcPeriod,
		mu:               syncz.NewMutex(),
		tunnelsByAgentId: redistool.NewExpiringHash(client, tunnelsByAgentIdHashKey(agentKeyPrefix), int64ToStr, ttl),
	}
}

func (t *RedisTracker) Run(ctx context.Context) error {
	refreshTicker := time.NewTicker(t.refreshPeriod)
	defer refreshTicker.Stop()
	gcTicker := time.NewTicker(t.gcPeriod)
	defer gcTicker.Stop()
	done := ctx.Done()
	for {
		select {
		case <-done:
			return nil
		case <-refreshTicker.C:
			err := t.refreshRegistrations(ctx, time.Now().Add(t.refreshPeriod-refreshOverlap))
			if err != nil {
				t.log.Error("Failed to refresh data in Redis", logz.Error(err))
			}
		case <-gcTicker.C:
			deletedKeys, err := t.runGC(ctx)
			if err != nil {
				t.log.Error("Failed to GC data in Redis", logz.Error(err))
				// fallthrough
			}
			if deletedKeys > 0 {
				t.log.Info("Deleted expired agent tunnel records", logz.RemovedHashKeys(deletedKeys))
			}
		}
	}
}

func (t *RedisTracker) RegisterTunnel(ctx context.Context, info *TunnelInfo) error {
	infoBytes, err := proto.Marshal(info)
	if err != nil {
		// This should never happen
		return fmt.Errorf("failed to marshal tunnel info: %w", err)
	}
	var register redistool.IOFunc
	ok := t.mu.RunLocked(ctx, func() {
		register = t.tunnelsByAgentId.Set(info.AgentId, info.ConnectionId, infoBytes)
	})
	if !ok {
		return ctx.Err()
	}
	return register(ctx)
}

func (t *RedisTracker) UnregisterTunnel(ctx context.Context, info *TunnelInfo) error {
	var unregister redistool.IOFunc
	ok := t.mu.RunLocked(ctx, func() {
		unregister = t.tunnelsByAgentId.Unset(info.AgentId, info.ConnectionId)
	})
	if !ok {
		return ctx.Err()
	}
	return unregister(ctx)
}

func (t *RedisTracker) GetTunnelsByAgentId(ctx context.Context, agentId int64, cb GetTunnelsByAgentIdCallback) error {
	_, err := t.tunnelsByAgentId.Scan(ctx, agentId, func(value []byte, err error) (bool, error) {
		if err != nil {
			t.log.Error("Redis hash scan", logz.Error(err))
			return false, nil
		}
		var info TunnelInfo
		err = proto.Unmarshal(value, &info)
		if err != nil {
			t.log.Error("Redis proto.Unmarshal(TunnelInfo)", logz.Error(err))
			return false, nil
		}
		err = info.ValidateAll()
		if err != nil {
			t.log.Error("Redis proto.Unmarshal(TunnelInfo) validation", logz.Error(err))
			return false, nil
		}
		return cb(&info)
	})
	return err
}

func (t *RedisTracker) refreshRegistrations(ctx context.Context, nextRefresh time.Time) error {
	var refresh redistool.IOFunc
	ok := t.mu.RunLocked(ctx, func() {
		refresh = t.tunnelsByAgentId.Refresh(nextRefresh)
	})
	if !ok {
		return nil
	}
	return refresh(ctx)
}

func (t *RedisTracker) runGC(ctx context.Context) (int /* keysDeleted */, error) {
	var gc func(context.Context) (int /* keysDeleted */, error)
	ok := t.mu.RunLocked(ctx, func() {
		gc = t.tunnelsByAgentId.GC()
	})
	if !ok {
		return 0, nil
	}
	return gc(ctx)
}

// tunnelsByAgentIdHashKey returns a key for agentId -> (connectionId -> marshaled TunnelInfo).
func tunnelsByAgentIdHashKey(agentKeyPrefix string) redistool.KeyToRedisKey[int64] {
	prefix := agentKeyPrefix + ":conn_by_agent_id:"
	return func(agentId int64) string {
		return redistool.PrefixedInt64Key(prefix, agentId)
	}
}

func int64ToStr(key int64) string {
	return strconv.FormatInt(key, 10)
}
