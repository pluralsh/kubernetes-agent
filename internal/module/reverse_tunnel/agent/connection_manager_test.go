package agent

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/reverse_tunnel/info"
)

var (
	_ modagent.Module     = (*module)(nil)
	_ modagent.Factory    = (*Factory)(nil)
	_ connectionInterface = (*mockConnection)(nil)
)

func TestConnManager_StartsMinIdleConnectionsOnRun(t *testing.T) {
	cm, conns, mu := setupConnManager(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go cm.Run(ctx)
	require.Eventually(t, func() bool {
		cm.mu.Lock()
		defer cm.mu.Unlock()
		mu.Lock()
		defer mu.Unlock()
		return len(*conns) == int(cm.minIdleConnections)
	}, time.Minute, 10*time.Millisecond)
	cancel()
	cm.wg.Wait()
	require.Len(t, *conns, int(cm.minIdleConnections))
}

func TestConnManager_ScalesIdleConnectionsToMaxAndThenToMin(t *testing.T) {
	cm, conns, mu := setupConnManager(t)
	cm.maxIdleTime = 50 * time.Millisecond
	var activated int
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go cm.Run(ctx)
	// Scale to max
	require.Eventually(t, func() bool {
		mu.Lock()
		lenConns := len(*conns)
		assert.LessOrEqual(t, lenConns, int(cm.maxConnections))
		toActivate := make([]*mockConnection, lenConns-activated)
		copy(toActivate, (*conns)[activated:])
		mu.Unlock()
		for _, c := range toActivate { // activate any new connections, must not hold the mutex
			activated++
			c.onActive(c)
		}
		return lenConns == int(cm.maxConnections)
	}, time.Minute, 10*time.Millisecond)
	// Scale to min
	allConnectionsToIdle(cm)
	require.Eventually(t, func() bool { // poll until 50ms expire
		cm.mu.Lock()
		defer cm.mu.Unlock()
		return cm.idleConnections == cm.minIdleConnections && cm.activeConnections == 0 && len(cm.connections) == int(cm.minIdleConnections)
	}, time.Minute, 10*time.Millisecond)
	cancel()
	cm.wg.Wait()
	require.Len(t, *conns, int(cm.maxConnections))
}

func TestConnManager_ActiveRacesWithTimedOut(t *testing.T) {
	cm, conns, mu := setupConnManager(t)
	cf := cm.connectionFactory
	runBlock := make(chan struct{})
	cm.connectionFactory = func(agentDescriptor *info.AgentDescriptor, onActive, onIdle func(connectionInterface)) connectionInterface {
		c := cf(agentDescriptor, onActive, onIdle)
		c.(*mockConnection).runBlock = runBlock
		return c
	}
	cm.maxIdleTime = 10 * time.Millisecond
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go cm.Run(ctx)
	// Activate one, make it idle.
	require.Eventually(t, func() bool {
		mu.Lock()
		if len(*conns) == 0 {
			mu.Unlock()
			return false
		}
		first := (*conns)[0]
		mu.Unlock()
		first.onActive(first) // activate, this will start one more connection
		first.onIdle(first)   // make idle
		return true
	}, time.Minute, 10*time.Millisecond)
	// One of the connections should time out
	require.Eventually(t, func() bool {
		cm.mu.Lock()
		for c, i := range cm.connections {
			if i.state == timedOut {
				// the timed out and stopped connection
				cm.mu.Unlock()
				// Make sure calling onActive() on a timed out connection is fine
				c.(*mockConnection).onActive(c)
				// Make sure calling onIdle() on a timed out connection is fine
				c.(*mockConnection).onIdle(c)
				return true
			}
		}
		cm.mu.Unlock()
		return false
	}, time.Minute, 10*time.Millisecond)
	close(runBlock)
	cancel()
	cm.wg.Wait()
}

func allConnectionsToIdle(cm *connectionManager) {
	cm.mu.Lock()
	conns := make([]connectionInterface, 0, len(cm.connections))
	for c, i := range cm.connections {
		if i.state == active {
			conns = append(conns, c)
		}
	}
	cm.mu.Unlock()
	for _, c := range conns {
		c.(*mockConnection).onIdle(c)
	}
}

func setupConnManager(t *testing.T) (*connectionManager, *[]*mockConnection, *sync.Mutex) {
	t.Parallel()
	var conns []*mockConnection
	var mu sync.Mutex
	cm := &connectionManager{
		connections:        make(map[connectionInterface]connectionInfo),
		minIdleConnections: 1,
		maxConnections:     maxConnections,
		scaleUpStep:        2,
		maxIdleTime:        time.Minute,
		connectionFactory: func(agentDescriptor *info.AgentDescriptor, onActive, onIdle func(connectionInterface)) connectionInterface {
			c := &mockConnection{
				onActive: onActive,
				onIdle:   onIdle,
			}
			mu.Lock()
			defer mu.Unlock()
			conns = append(conns, c)
			return c
		},
	}
	t.Cleanup(func() {
		cm.wg.Wait()
		assert.Zero(t, cm.idleConnections)
		assert.Zero(t, cm.activeConnections)
		assert.Empty(t, cm.connections)
		for _, c := range conns {
			assert.EqualValues(t, 1, c.runCalled)
			assert.EqualValues(t, 1, c.stopped)
		}
	})
	return cm, &conns, &mu
}

type mockConnection struct {
	runCalled, stopped int32
	onActive, onIdle   func(connectionInterface)
	runBlock           chan struct{} // allows to block before running onStop()
}

func (m *mockConnection) Run(ctx context.Context) {
	defer atomic.AddInt32(&m.stopped, 1)
	atomic.AddInt32(&m.runCalled, 1)
	<-ctx.Done()
	if m.runBlock != nil {
		<-m.runBlock
	}
}
