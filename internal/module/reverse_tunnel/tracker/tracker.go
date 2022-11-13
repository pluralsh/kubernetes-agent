package tracker

import (
	"context"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/redistool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/syncz"
	"go.uber.org/zap"
)

const refreshOverlap = 5 * time.Second

type KasUrlsByAgentIdCallback func(kasUrl string) (bool /* done */, error)

type Querier interface {
	KasUrlsByAgentId(ctx context.Context, agentId int64, cb KasUrlsByAgentIdCallback) error
}

type Registerer interface {
	// RegisterTunnel registers tunnel with the tracker.
	RegisterTunnel(ctx context.Context, agentId int64) error
	// UnregisterTunnel unregisters tunnel with the tracker.
	UnregisterTunnel(ctx context.Context, agentId int64) error
}

type Tracker interface {
	Registerer
	Querier
	Run(ctx context.Context) error
}

type RedisTracker struct {
	log              *zap.Logger
	refreshPeriod    time.Duration
	gcPeriod         time.Duration
	ownPrivateApiUrl string

	// mu protects fields below
	mu                    sync.Mutex
	tunnelsByAgentIdCount map[int64]uint16
	tunnelsByAgentId      redistool.ExpiringHashInterface[int64, string] // agentId -> kas URL -> nil
}

func NewRedisTracker(log *zap.Logger, client redis.UniversalClient, agentKeyPrefix string, ttl, refreshPeriod,
	gcPeriod time.Duration, ownPrivateApiUrl string) *RedisTracker {
	return &RedisTracker{
		log:                   log,
		refreshPeriod:         refreshPeriod,
		gcPeriod:              gcPeriod,
		ownPrivateApiUrl:      ownPrivateApiUrl,
		tunnelsByAgentIdCount: make(map[int64]uint16),
		tunnelsByAgentId:      redistool.NewExpiringHash(client, tunnelsByAgentIdHashKey(agentKeyPrefix), strToStr, ttl),
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

func (t *RedisTracker) RegisterTunnel(ctx context.Context, agentId int64) error {
	register := syncz.RunWithMutex(&t.mu, func() redistool.IOFunc {
		cnt := t.tunnelsByAgentIdCount[agentId]
		cnt++
		t.tunnelsByAgentIdCount[agentId] = cnt
		if cnt == 1 {
			// First tunnel for this agentId
			return t.tunnelsByAgentId.Set(agentId, t.ownPrivateApiUrl, nil)
		} else {
			return noopIO
		}
	})
	return register(ctx)
}

func (t *RedisTracker) UnregisterTunnel(ctx context.Context, agentId int64) error {
	unregister := syncz.RunWithMutex(&t.mu, func() redistool.IOFunc {
		cnt := t.tunnelsByAgentIdCount[agentId]
		cnt--
		if cnt == 0 {
			delete(t.tunnelsByAgentIdCount, agentId)
			return t.tunnelsByAgentId.Unset(agentId, t.ownPrivateApiUrl)
		} else {
			t.tunnelsByAgentIdCount[agentId] = cnt
			return noopIO
		}
	})
	return unregister(ctx)
}

func (t *RedisTracker) KasUrlsByAgentId(ctx context.Context, agentId int64, cb KasUrlsByAgentIdCallback) error {
	_, err := t.tunnelsByAgentId.Scan(ctx, agentId, func(rawHashKey string, value []byte, err error) (bool, error) {
		if err != nil {
			t.log.Error("Redis hash scan", logz.Error(err))
			return false, nil
		}
		return cb(rawHashKey)
	})
	return err
}

func (t *RedisTracker) refreshRegistrations(ctx context.Context, nextRefresh time.Time) error {
	refresh := syncz.RunWithMutex(&t.mu, func() redistool.IOFunc {
		return t.tunnelsByAgentId.Refresh(nextRefresh)
	})
	return refresh(ctx)
}

func (t *RedisTracker) runGC(ctx context.Context) (int /* keysDeleted */, error) {
	gc := syncz.RunWithMutex(&t.mu, func() func(ctx context.Context) (int /* keysDeleted */, error) {
		return t.tunnelsByAgentId.GC()
	})
	return gc(ctx)
}

// tunnelsByAgentIdHashKey returns a key for agentId -> (kasUrl -> nil).
func tunnelsByAgentIdHashKey(agentKeyPrefix string) redistool.KeyToRedisKey[int64] {
	prefix := agentKeyPrefix + ":kas_by_agent_id:"
	return func(agentId int64) string {
		return redistool.PrefixedInt64Key(prefix, agentId)
	}
}

func strToStr(key string) string {
	return key
}

func noopIO(ctx context.Context) error {
	return nil
}
