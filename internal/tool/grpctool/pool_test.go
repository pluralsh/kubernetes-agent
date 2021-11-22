package grpctool

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc"
	clocktesting "k8s.io/utils/clock/testing"
)

const (
	t1 = "grpc://127.0.0.1:1"
	t2 = "grpc://127.0.0.1:2"
)

func TestKasPool_DialConnDifferentPort(t *testing.T) {
	p := NewPool(zaptest.NewLogger(t), grpc.WithInsecure())
	defer clz(t, p)
	c1, err := p.Dial(context.Background(), t1)
	require.NoError(t, err)
	c1.Done()
	c2, err := p.Dial(context.Background(), t2)
	require.NoError(t, err)
	assert.NotSame(t, c1, c2)
	c2.Done()
}

func TestKasPool_DialConnSequentialReuse(t *testing.T) {
	p := NewPool(zaptest.NewLogger(t), grpc.WithInsecure())
	defer clz(t, p)
	c1, err := p.Dial(context.Background(), t1)
	require.NoError(t, err)
	c1.Done()
	c2, err := p.Dial(context.Background(), t1)
	require.NoError(t, err)
	assert.Same(t, c1.(*poolConn).ClientConn, c2.(*poolConn).ClientConn)
	c2.Done()
}

func TestKasPool_DialConnConcurrentReuse(t *testing.T) {
	p := NewPool(zaptest.NewLogger(t), grpc.WithInsecure())
	defer clz(t, p)
	c1, err := p.Dial(context.Background(), t1)
	require.NoError(t, err)
	c2, err := p.Dial(context.Background(), t1)
	require.NoError(t, err)
	assert.Same(t, c1.(*poolConn).ClientConn, c2.(*poolConn).ClientConn)
	c1.Done()
	c2.Done()
}

func TestKasPool_CloseClosesAllConnections(t *testing.T) {
	p := NewPool(zaptest.NewLogger(t), grpc.WithInsecure())
	c, err := p.Dial(context.Background(), t1)
	require.NoError(t, err)
	c.Done()
	require.NoError(t, p.Close())
	assert.Empty(t, p.conns)
}

func TestKasPool_DonePanicsOnMultipleInvocations(t *testing.T) {
	p := NewPool(zaptest.NewLogger(t), grpc.WithInsecure())
	defer clz(t, p)
	c, err := p.Dial(context.Background(), t1)
	require.NoError(t, err)
	c.Done()
	assert.PanicsWithError(t, "pool connection Done() called more than once", func() {
		c.Done()
	})
}

func TestKasPool_DoneEvictsExpiredIdleConnections(t *testing.T) {
	start := time.Now()
	tClock := clocktesting.NewFakePassiveClock(start)
	p := &Pool{
		log:      zaptest.NewLogger(t),
		dialOpts: []grpc.DialOption{grpc.WithInsecure()},
		conns:    map[string]*connHolder{},
		clk:      tClock,
	}
	defer clz(t, p)
	c1, err := p.Dial(context.Background(), t1)
	require.NoError(t, err)
	c1.Done()
	tClock.SetTime(start.Add(2 * evictIdleConnAfter))
	p.runGcLocked()
	assert.Empty(t, p.conns)
}

func clz(t *testing.T, c io.Closer) {
	assert.NoError(t, c.Close())
}