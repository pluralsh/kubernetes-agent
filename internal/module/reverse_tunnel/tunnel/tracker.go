package tunnel

import (
	"context"
	"time"

	"github.com/redis/rueidis"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/redistool"
	"go.uber.org/zap"
)

type KasUrlsByAgentIdCallback func(kasUrl string) (bool /* done */, error)

type Querier interface {
	// KasUrlsByAgentId calls the callback with the list of kas URLs for a particular agent id.
	// Safe for concurrent use.
	KasUrlsByAgentId(ctx context.Context, agentId int64, cb KasUrlsByAgentIdCallback) error
}

// Registerer allows to register and unregister tunnels.
// Caller is responsible for periodically calling GC() and Refresh().
// Not safe for concurrent use.
type Registerer interface {
	// RegisterTunnel registers tunnel with the tracker.
	RegisterTunnel(ctx context.Context, agentId int64) error
	// UnregisterTunnel unregisters tunnel with the tracker.
	UnregisterTunnel(ctx context.Context, agentId int64) error
	// GC deletes expired tunnels from the underlying storage.
	GC() func(context.Context) (int /* keysDeleted */, error)
	// Refresh refreshes registered tunnels in the underlying storage.
	Refresh(ctx context.Context, nextRefresh time.Time) error
}

type Tracker interface {
	Registerer
	Querier
}

type RedisTracker struct {
	log                   *zap.Logger
	api                   modshared.Api
	ownPrivateApiUrl      string
	tunnelsByAgentIdCount map[int64]uint16
	tunnelsByAgentId      redistool.ExpiringHashInterface[int64, string] // agentId -> kas URL -> nil
}

func NewRedisTracker(log *zap.Logger, api modshared.Api, client rueidis.Client, agentKeyPrefix string,
	ttl time.Duration, ownPrivateApiUrl string) *RedisTracker {
	return &RedisTracker{
		log:                   log,
		api:                   api,
		ownPrivateApiUrl:      ownPrivateApiUrl,
		tunnelsByAgentIdCount: make(map[int64]uint16),
		tunnelsByAgentId:      redistool.NewExpiringHash(client, tunnelsByAgentIdHashKey(agentKeyPrefix), strToStr, ttl),
	}
}

func (t *RedisTracker) RegisterTunnel(ctx context.Context, agentId int64) error {
	cnt := t.tunnelsByAgentIdCount[agentId]
	cnt++
	t.tunnelsByAgentIdCount[agentId] = cnt
	if cnt == 1 {
		// First tunnel for this agentId
		return t.tunnelsByAgentId.Set(ctx, agentId, t.ownPrivateApiUrl, nil)
	}
	return nil
}

func (t *RedisTracker) UnregisterTunnel(ctx context.Context, agentId int64) error {
	cnt := t.tunnelsByAgentIdCount[agentId]
	cnt--
	if cnt == 0 {
		delete(t.tunnelsByAgentIdCount, agentId)
		return t.tunnelsByAgentId.Unset(ctx, agentId, t.ownPrivateApiUrl)
	} else {
		t.tunnelsByAgentIdCount[agentId] = cnt
		return nil
	}
}

func (t *RedisTracker) KasUrlsByAgentId(ctx context.Context, agentId int64, cb KasUrlsByAgentIdCallback) error {
	_, err := t.tunnelsByAgentId.Scan(ctx, agentId, func(rawHashKey string, value []byte, err error) (bool, error) {
		if err != nil {
			t.api.HandleProcessingError(ctx, t.log.With(logz.AgentId(agentId)), agentId, "Redis hash scan", err)
			return false, nil
		}
		return cb(rawHashKey)
	})
	return err
}

func (t *RedisTracker) Refresh(ctx context.Context, nextRefresh time.Time) error {
	return t.tunnelsByAgentId.Refresh(ctx, nextRefresh)
}

func (t *RedisTracker) GC() func(context.Context) (int /* keysDeleted */, error) {
	return t.tunnelsByAgentId.GC()
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
