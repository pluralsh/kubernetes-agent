package agent_tracker

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/errz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/redistool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/syncz"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"
)

const (
	refreshOverlap           = 5 * time.Second
	connectedAgentsKey int64 = 0
)

type ConnectedAgentInfoCallback func(*ConnectedAgentInfo) (done bool, err error)

type Registerer interface {
	// RegisterConnection registers connection with the tracker.
	RegisterConnection(ctx context.Context, info *ConnectedAgentInfo) error
	// UnregisterConnection unregisters connection with the tracker.
	UnregisterConnection(ctx context.Context, info *ConnectedAgentInfo) error
}

type Querier interface {
	GetConnectionsByAgentId(ctx context.Context, agentId int64, cb ConnectedAgentInfoCallback) error
	GetConnectionsByProjectId(ctx context.Context, projectId int64, cb ConnectedAgentInfoCallback) error
	GetConnectedAgentsCount(ctx context.Context) (int64, error)
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

	// refreshMu is exclusively held during refresh process and non-exclusively held during de-registration.
	// This ensures refresh and de-registration never happen concurrently and hence just unregistered connections are
	// never written back into Redis by refresh process.
	refreshMu syncz.RWMutex
	// mu protects fields below
	mu                     sync.Mutex
	connectionsByAgentId   redistool.ExpiringHashInterface[int64, int64] // agentId -> connectionId -> info
	connectionsByProjectId redistool.ExpiringHashInterface[int64, int64] // projectId -> connectionId -> info
	connectedAgents        redistool.ExpiringHashInterface[int64, int64] // hash name -> agentId -> ""
}

func NewRedisTracker(log *zap.Logger, client redis.UniversalClient, agentKeyPrefix string, ttl, refreshPeriod, gcPeriod time.Duration) *RedisTracker {
	return &RedisTracker{
		log:                    log,
		refreshPeriod:          refreshPeriod,
		gcPeriod:               gcPeriod,
		refreshMu:              syncz.NewRWMutex(),
		connectionsByAgentId:   redistool.NewExpiringHash(client, connectionsByAgentIdHashKey(agentKeyPrefix), int64ToStr, ttl),
		connectionsByProjectId: redistool.NewExpiringHash(client, connectionsByProjectIdHashKey(agentKeyPrefix), int64ToStr, ttl),
		connectedAgents:        redistool.NewExpiringHash(client, connectedAgentsHashKey(agentKeyPrefix), int64ToStr, ttl),
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
			t.refreshRegistrations(ctx, time.Now().Add(t.refreshPeriod-refreshOverlap))
		case <-gcTicker.C:
			keysDeleted := t.runGC(ctx)
			if keysDeleted > 0 {
				t.log.Info("Deleted expired agent connections records", logz.RemovedHashKeys(keysDeleted))
			}
		}
	}
}

func (t *RedisTracker) RegisterConnection(ctx context.Context, info *ConnectedAgentInfo) error {
	infoBytes, err := proto.Marshal(info)
	if err != nil {
		// This should never happen
		return fmt.Errorf("failed to marshal object: %w", err)
	}
	return t.runIOFuncs(ctx, func() []redistool.IOFunc {
		return []redistool.IOFunc{
			t.connectionsByProjectId.Set(info.ProjectId, info.ConnectionId, infoBytes),
			t.connectionsByAgentId.Set(info.AgentId, info.ConnectionId, infoBytes),
			t.connectedAgents.Set(connectedAgentsKey, info.AgentId, nil),
		}
	})
}

func (t *RedisTracker) UnregisterConnection(ctx context.Context, info *ConnectedAgentInfo) error {
	if !t.refreshMu.RLock(ctx) {
		return ctx.Err()
	}
	defer t.refreshMu.RUnlock()
	return t.runIOFuncs(ctx, func() []redistool.IOFunc {
		t.connectedAgents.Forget(connectedAgentsKey, info.AgentId)
		return []redistool.IOFunc{
			t.connectionsByProjectId.Unset(info.ProjectId, info.ConnectionId),
			t.connectionsByAgentId.Unset(info.AgentId, info.ConnectionId),
		}
	})
}

func (t *RedisTracker) GetConnectionsByAgentId(ctx context.Context, agentId int64, cb ConnectedAgentInfoCallback) error {
	return getConnectionsByKey(ctx, t.log, t.connectionsByAgentId, agentId, cb)
}

func (t *RedisTracker) GetConnectionsByProjectId(ctx context.Context, projectId int64, cb ConnectedAgentInfoCallback) error {
	return getConnectionsByKey(ctx, t.log, t.connectionsByProjectId, projectId, cb)
}

func (t *RedisTracker) GetConnectedAgentsCount(ctx context.Context) (int64, error) {
	return t.connectedAgents.Len(ctx, connectedAgentsKey)
}

func (t *RedisTracker) refreshRegistrations(ctx context.Context, nextRefresh time.Time) {
	if !t.refreshMu.Lock(ctx) {
		return
	}
	defer t.refreshMu.Unlock()
	refreshFuncs := syncz.RunWithMutex(&t.mu, func() []redistool.IOFunc {
		return []redistool.IOFunc{
			t.connectionsByProjectId.Refresh(nextRefresh),
			t.connectionsByAgentId.Refresh(nextRefresh),
			t.connectedAgents.Refresh(nextRefresh),
		}
	})
	// No rush so run refresh sequentially to not stress RAM/CPU/Redis/network.
	// We have more important work to do that we shouldn't impact.
	for _, refresh := range refreshFuncs {
		err := refresh(ctx)
		if err != nil {
			if errz.ContextDone(err) {
				t.log.Debug("Redis hash data refresh interrupted", logz.Error(err))
				break
			}
			t.log.Error("Failed to refresh hash data in Redis", logz.Error(err))
			// continue anyway
		}
	}
}

func (t *RedisTracker) runGC(ctx context.Context) int {
	gcFuncs := syncz.RunWithMutex(&t.mu, func() []func(context.Context) (int, error) {
		return []func(context.Context) (int, error){
			t.connectionsByProjectId.GC(),
			t.connectionsByAgentId.GC(),
			t.connectedAgents.GC(),
		}
	})
	keysDeleted := 0
	// No rush so run GC sequentially to not stress RAM/CPU/Redis/network.
	// We have more important work to do that we shouldn't impact.
	for _, gc := range gcFuncs {
		deleted, err := gc(ctx)
		if err != nil {
			if errz.ContextDone(err) {
				t.log.Debug("Redis GC interrupted", logz.Error(err))
				break
			}
			t.log.Error("Failed to GC data in Redis", logz.Error(err))
			// continue anyway
		}
		keysDeleted += deleted
	}
	return keysDeleted
}

func (t *RedisTracker) runIOFuncs(ctx context.Context, f func() []redistool.IOFunc) error {
	ios := syncz.RunWithMutex(&t.mu, f)
	var g errgroup.Group
	for _, s := range ios {
		s := s
		g.Go(func() error {
			return s(ctx)
		})
	}
	return g.Wait()
}

func getConnectionsByKey[K1 any, K2 any](ctx context.Context, log *zap.Logger, hash redistool.ExpiringHashInterface[K1, K2], key K1, cb ConnectedAgentInfoCallback) error {
	_, err := hash.Scan(ctx, key, func(rawHashKey string, value []byte, err error) (bool, error) {
		if err != nil {
			log.Error("Redis hash scan", logz.Error(err))
			return false, nil
		}
		var info ConnectedAgentInfo
		err = proto.Unmarshal(value, &info)
		if err != nil {
			log.Error("Redis proto.Unmarshal(ConnectedAgentInfo)", logz.Error(err))
			return false, nil
		}
		return cb(&info)
	})
	return err
}

// connectionsByAgentIdHashKey returns a key for agentId -> (connectionId -> marshaled ConnectedAgentInfo).
func connectionsByAgentIdHashKey(agentKeyPrefix string) redistool.KeyToRedisKey[int64] {
	prefix := agentKeyPrefix + ":conn_by_agent_id:"
	return func(agentId int64) string {
		return redistool.PrefixedInt64Key(prefix, agentId)
	}
}

// connectionsByProjectIdHashKey returns a key for projectId -> (agentId ->marshaled ConnectedAgentInfo).
func connectionsByProjectIdHashKey(agentKeyPrefix string) redistool.KeyToRedisKey[int64] {
	prefix := agentKeyPrefix + ":conn_by_project_id:"
	return func(projectId int64) string {
		return redistool.PrefixedInt64Key(prefix, projectId)
	}
}

// connectedAgentsHashKey returns the key for the hash of connected agents.
func connectedAgentsHashKey(agentKeyPrefix string) redistool.KeyToRedisKey[int64] {
	prefix := agentKeyPrefix + ":connected_agents"
	return func(_ int64) string {
		return prefix
	}
}

type ConnectedAgentInfoCollector []*ConnectedAgentInfo

func (c *ConnectedAgentInfoCollector) Collect(info *ConnectedAgentInfo) (bool, error) {
	*c = append(*c, info)
	return false, nil
}

func int64ToStr(key int64) string {
	return strconv.FormatInt(key, 10)
}
