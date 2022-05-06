package agent

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/reverse_tunnel/info"
	"k8s.io/apimachinery/pkg/util/wait"
)

type state int8

const (
	invalid state = iota // nolint: deadcode,varcheck
	idle
	active
	timedOut
)

type connectionInfo struct {
	cancel    context.CancelFunc
	idleTimer *time.Timer
	state     state
}

// connectionManager manages a pool of connections and their lifecycles.
type connectionManager struct {
	mu          sync.Mutex // protects connections,idleConnections,activeConnections
	connections map[connectionInterface]connectionInfo
	// Counters to track connections in those states. There may be timedOut connections in the map too.
	idleConnections   int32
	activeConnections int32

	wg wait.Group

	// minIdleConnections is the minimum number of connections that are not streaming a request.
	minIdleConnections int32
	// maxConnections is the maximum number of connections (idle and active).
	maxConnections int32
	// scaleUpStep is the number of new connections to start when below minIdleConnections.
	scaleUpStep int32
	// maxIdleTime is the maximum duration of time a connection can stay in an idle state.
	maxIdleTime       time.Duration
	connectionFactory connectionFactory
	agentDescriptor   *info.AgentDescriptor
}

func (m *connectionManager) Run(ctx context.Context) {
	defer m.wg.Wait() // blocks here until ctx is done and all connections exit
	m.mu.Lock()
	defer m.mu.Unlock()
	for m.idleConnections < m.minIdleConnections {
		m.startConnectionLocked(ctx)
	}
}

func (m *connectionManager) startConnectionLocked(rootCtx context.Context) {
	m.idleConnections++
	c := m.connectionFactory(m.agentDescriptor,
		func(c connectionInterface) {
			m.onActive(rootCtx, c)
		},
		m.onIdle)
	ctx, cancel := context.WithCancel(rootCtx)
	m.connections[c] = connectionInfo{
		cancel: cancel,
		idleTimer: time.AfterFunc(m.maxIdleTime, func() {
			m.onIdleTimeout(c)
		}),
		state: idle,
	}
	m.wg.StartWithContext(ctx, func(ctx context.Context) {
		defer m.onStop(c)
		c.Run(ctx)
	})
}

func (m *connectionManager) onActive(rootCtx context.Context, c connectionInterface) {
	m.mu.Lock()
	defer m.mu.Unlock()
	i := m.connections[c]
	switch i.state { // nolint: exhaustive
	case idle: // idle -> active transition
		i.state = active
		i.idleTimer.Stop()
		m.connections[c] = i
		m.idleConnections--
		m.activeConnections++
		if m.idleConnections < m.minIdleConnections {
			// Not enough idle connections. Must scale up the number of connections.
			// Ensure we don't go above the limit.
			scaleBy := m.scaleUpStep
			haveConnections := m.idleConnections + m.activeConnections
			canSpawnConnections := m.maxConnections - haveConnections
			if scaleBy > canSpawnConnections {
				scaleBy = canSpawnConnections
			}
			for ; scaleBy > 0; scaleBy-- {
				m.startConnectionLocked(rootCtx)
			}
		}
	case timedOut:
	// The connection has timed out already, nothing to do here.
	// Timeout handler cancels the context so the connection will quickly stop.
	// It'll be removed from the map soon.
	case active:
		panic(errors.New("connection is already active"))
	default:
		panic(fmt.Errorf("unknown state: %d", i.state))
	}
}

func (m *connectionManager) onIdle(c connectionInterface) {
	m.mu.Lock()
	defer m.mu.Unlock()
	i := m.connections[c]
	switch i.state { // nolint: exhaustive
	case idle:
		// nothing to do, already in the idle state
	case active: // active -> idle transition
		i.state = idle
		i.idleTimer.Reset(m.maxIdleTime)
		m.connections[c] = i
		m.idleConnections++
		m.activeConnections--
	case timedOut:
		// The connection has timed out already, nothing to do here. It'll be removed from the map soon.
	default:
		panic(fmt.Errorf("unknown state: %d", i.state))
	}
}

func (m *connectionManager) onIdleTimeout(c connectionInterface) {
	m.mu.Lock()
	defer m.mu.Unlock()
	i, ok := m.connections[c]
	if !ok {
		// The connection has been removed from the set before this method acquired the lock and got here.
		return
	}
	switch i.state { // nolint: exhaustive
	case idle: // idle -> timed out transition
		if m.idleConnections > m.minIdleConnections {
			// Enough idle connections, can close this one since it timed out
			i.state = timedOut
			i.cancel()
			m.connections[c] = i
			m.idleConnections--
		} else {
			// We are at the minimum of idle connections. Keep this connection by resetting its timeout timer.
			i.idleTimer.Reset(m.maxIdleTime)
		}
	case active:
		// Timer fired concurrently with connection transitioning from idle to active.
		// The connection is active now and there is nothing to do, just ignore this invocation.
	case timedOut:
		// Timer fired multiple times (high system load? clock skew?), ignore.
	default:
		panic(fmt.Errorf("unknown state: %d", i.state))
	}
}

func (m *connectionManager) onStop(c connectionInterface) {
	m.mu.Lock()
	defer m.mu.Unlock()
	i := m.connections[c]
	i.idleTimer.Stop()
	delete(m.connections, c)
	if i.state != timedOut {
		// onIdleTimeout() decrements this field on timeout.
		// It's decremented here too to handle context done situation.
		m.idleConnections--
	}
}
