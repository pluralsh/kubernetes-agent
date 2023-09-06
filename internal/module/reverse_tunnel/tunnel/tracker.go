package tunnel

import (
	"context"
	"errors"
	"time"

	"github.com/redis/rueidis"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/redistool"
)

type Querier interface {
	// KasUrlsByAgentId returns the list of kas URLs for a particular agent id.
	// A partial list may be returned together with an error.
	// Safe for concurrent use.
	KasUrlsByAgentId(ctx context.Context, agentId int64) ([]string, error)
}

// Registerer allows to register and unregister tunnels.
// Caller is responsible for periodically calling GC() and Refresh().
// Not safe for concurrent use.
type Registerer interface {
	// RegisterTunnel registers tunnel with the tracker.
	RegisterTunnel(ctx context.Context, agentId int64) error
	// UnregisterTunnel unregisters tunnel with the tracker.
	UnregisterTunnel(ctx context.Context, agentId int64) error
	// Refresh refreshes registered tunnels in the underlying storage.
	Refresh(ctx context.Context, nextRefresh time.Time) error
}

type Tracker interface {
	Registerer
	Querier
}

type RedisTracker struct {
	ownPrivateApiUrl      string
	tunnelsByAgentIdCount map[int64]uint16
	tunnelsByAgentId      redistool.ExpiringHashInterface[int64, string] // agentId -> kas URL -> nil
}

func NewRedisTracker(client rueidis.Client, agentKeyPrefix string, ttl time.Duration, ownPrivateApiUrl string) *RedisTracker {
	return &RedisTracker{
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

func (t *RedisTracker) KasUrlsByAgentId(ctx context.Context, agentId int64) ([]string, error) {
	var urls []string
	var errs []error
	_, err := t.tunnelsByAgentId.Scan(ctx, agentId, func(rawHashKey string, value []byte, err error) (bool, error) {
		if err != nil {
			errs = append(errs, err)
			return false, nil
		}
		urls = append(urls, rawHashKey)
		return false, nil
	})
	if err != nil {
		errs = append(errs, err)
	}
	return urls, errors.Join(errs...)
}

func (t *RedisTracker) Refresh(ctx context.Context, nextRefresh time.Time) error {
	return t.tunnelsByAgentId.Refresh(ctx, nextRefresh)
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