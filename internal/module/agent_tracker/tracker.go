package agent_tracker

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/redistool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/retry"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"
)

type ConnectedAgentInfoCallback func(*ConnectedAgentInfo) (done bool, err error)

type Registerer interface {
	// RegisterConnection schedules the connection to be registered with the tracker.
	// Returns true on success and false if ctx signaled done.
	RegisterConnection(ctx context.Context, info *ConnectedAgentInfo) bool
	// UnregisterConnection schedules the connection to be unregistered with the tracker.
	// Returns true on success and false if ctx signaled done.
	UnregisterConnection(ctx context.Context, info *ConnectedAgentInfo) bool
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
	log                    *zap.Logger
	refreshPeriod          time.Duration
	gcPeriod               time.Duration
	connectionsByAgentId   redistool.ExpiringHashInterface // agentId -> connectionId -> info
	connectionsByProjectId redistool.ExpiringHashInterface // projectId -> connectionId -> info
	connectedAgents        redistool.ExpiringHashInterface // hash name -> agentId -> ""
	toRegister             chan *ConnectedAgentInfo
	toUnregister           chan *ConnectedAgentInfo
	gc                     retry.SingleRun
}

func NewRedisTracker(log *zap.Logger, client redis.UniversalClient, agentKeyPrefix string, ttl, refreshPeriod, gcPeriod time.Duration) *RedisTracker {
	return &RedisTracker{
		log:                    log,
		refreshPeriod:          refreshPeriod,
		gcPeriod:               gcPeriod,
		connectionsByAgentId:   redistool.NewExpiringHash(log, client, connectionsByAgentIdHashKey(agentKeyPrefix), ttl),
		connectionsByProjectId: redistool.NewExpiringHash(log, client, connectionsByProjectIdHashKey(agentKeyPrefix), ttl),
		connectedAgents:        redistool.NewExpiringHash(log, client, connectedAgentsHashKey(agentKeyPrefix), ttl),
		toRegister:             make(chan *ConnectedAgentInfo),
		toUnregister:           make(chan *ConnectedAgentInfo),
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
			err := t.refreshRegistrations(ctx)
			if err != nil {
				t.log.Error("Failed to refresh data in Redis", logz.Error(err))
			}
		case <-gcTicker.C:
			t.maybeRunGCAsync(ctx)
		case toReg := <-t.toRegister:
			err := t.registerConnection(ctx, toReg)
			if err != nil {
				t.log.Error("Failed to register connection", logz.Error(err))
			}
		case toUnreg := <-t.toUnregister:
			err := t.unregisterConnection(ctx, toUnreg)
			if err != nil {
				t.log.Error("Failed to unregister connection", logz.Error(err))
			}
		}
	}
}

func (t *RedisTracker) RegisterConnection(ctx context.Context, info *ConnectedAgentInfo) bool {
	select {
	case <-ctx.Done():
		return false
	case t.toRegister <- info:
		return true
	}
}

func (t *RedisTracker) UnregisterConnection(ctx context.Context, info *ConnectedAgentInfo) bool {
	select {
	case <-ctx.Done():
		return false
	case t.toUnregister <- info:
		return true
	}
}

func (t *RedisTracker) GetConnectionsByAgentId(ctx context.Context, agentId int64, cb ConnectedAgentInfoCallback) error {
	return t.getConnectionsByKey(ctx, t.connectionsByAgentId, agentId, cb)
}

func (t *RedisTracker) GetConnectionsByProjectId(ctx context.Context, projectId int64, cb ConnectedAgentInfoCallback) error {
	return t.getConnectionsByKey(ctx, t.connectionsByProjectId, projectId, cb)
}

func (t *RedisTracker) GetConnectedAgentsCount(ctx context.Context) (int64, error) {
	return t.connectedAgents.Len(ctx, nil)
}

func (t *RedisTracker) getConnectionsByKey(ctx context.Context, hash redistool.ExpiringHashInterface, key interface{}, cb ConnectedAgentInfoCallback) error {
	_, err := hash.Scan(ctx, key, func(value *anypb.Any, err error) (bool, error) {
		if err != nil {
			t.log.Error("Redis hash scan", logz.Error(err))
			return false, nil
		}
		var info ConnectedAgentInfo
		err = value.UnmarshalTo(&info)
		if err != nil {
			t.log.Error("Redis proto.UnmarshalTo(ConnectedAgentInfo)", logz.Error(err))
			return false, nil
		}
		return cb(&info)
	})
	return err
}

func (t *RedisTracker) registerConnection(ctx context.Context, info *ConnectedAgentInfo) error {
	infoAny, err := anypb.New(info)
	if err != nil {
		// This should never happen
		return err
	}
	// Ensure data is put into all sets, even if there was an error
	err1 := t.connectionsByProjectId.Set(ctx, info.ProjectId, info.ConnectionId, infoAny)
	err2 := t.connectionsByAgentId.Set(ctx, info.AgentId, info.ConnectionId, infoAny)
	err3 := t.connectedAgents.Set(ctx, nil, info.AgentId, nil)
	if err1 == nil {
		err1 = err2
	}
	if err1 == nil {
		err1 = err3
	}
	return err1
}

func (t *RedisTracker) unregisterConnection(ctx context.Context, unreg *ConnectedAgentInfo) error {
	err1 := t.connectionsByProjectId.Unset(ctx, unreg.ProjectId, unreg.ConnectionId)
	err2 := t.connectionsByAgentId.Unset(ctx, unreg.AgentId, unreg.ConnectionId)
	if err1 == nil {
		err1 = err2
	}
	return err1
}

func (t *RedisTracker) refreshRegistrations(ctx context.Context) error {
	err1 := t.connectionsByProjectId.Refresh(ctx)
	err2 := t.connectionsByAgentId.Refresh(ctx)
	err3 := t.connectedAgents.Refresh(ctx)

	if err1 == nil {
		err1 = err2
	}
	if err1 == nil {
		err1 = err3
	}
	return err1
}

func (t *RedisTracker) maybeRunGCAsync(ctx context.Context) {
	gc1 := t.connectionsByProjectId.GC()
	gc2 := t.connectionsByAgentId.GC()
	gc3 := t.connectedAgents.GC()
	t.gc.Run(func() {
		keysDeleted := 0
		for _, gc := range []func(context.Context) (int, error){gc1, gc2, gc3} {
			deleted, err := gc(ctx)
			if err != nil {
				t.log.Error("Failed to GC data in Redis", logz.Error(err))
				// continue anyway
			}
			keysDeleted += deleted
		}
		if keysDeleted > 0 {
			t.log.Info("Deleted expired agent connections records", logz.RemovedHashKeys(keysDeleted))
		}
	})
}

// connectionsByAgentIdHashKey returns a key for agentId -> (connectionId -> marshaled ConnectedAgentInfo).
func connectionsByAgentIdHashKey(agentKeyPrefix string) redistool.KeyToRedisKey {
	prefix := agentKeyPrefix + ":conn_by_agent_id:"
	return func(agentId interface{}) string {
		return redistool.PrefixedInt64Key(prefix, agentId.(int64))
	}
}

// connectionsByProjectIdHashKey returns a key for projectId -> (agentId ->marshaled ConnectedAgentInfo).
func connectionsByProjectIdHashKey(agentKeyPrefix string) redistool.KeyToRedisKey {
	prefix := agentKeyPrefix + ":conn_by_project_id:"
	return func(projectId interface{}) string {
		return redistool.PrefixedInt64Key(prefix, projectId.(int64))
	}
}

// connectedAgentsHashKey returns the key for the hash of connected agents.
func connectedAgentsHashKey(agentKeyPrefix string) redistool.KeyToRedisKey {
	prefix := agentKeyPrefix + ":connected_agents"
	return func(_ interface{}) string {
		return prefix
	}
}

type ConnectedAgentInfoCollector []*ConnectedAgentInfo

func (c *ConnectedAgentInfoCollector) Collect(info *ConnectedAgentInfo) (bool, error) {
	*c = append(*c, info)
	return false, nil
}
