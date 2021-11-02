package grpctool

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sync"
	"time"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/logz"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"k8s.io/utils/clock"
)

const (
	evictIdleConnAfter = 1 * time.Hour
)

type PoolConn interface {
	grpc.ClientConnInterface
	Done()
}

var (
	_ PoolConn = &poolConn{}
)

type Pool struct {
	mx       sync.Mutex
	log      *zap.Logger
	dialOpts []grpc.DialOption
	conns    map[string]*connHolder // target -> conn
	clk      clock.PassiveClock
}

func (p *Pool) Close() error {
	p.mx.Lock()
	defer p.mx.Unlock()
	for target, conn := range p.conns {
		delete(p.conns, target)
		log := p.log.With(logz.PoolConnectionUrl(conn.targetUrl))
		if conn.numUsers > 0 {
			log.Sugar().Warnf("Closing pool connection that is being used by %d callers", conn.numUsers)
		}
		err := conn.Close()
		if err != nil {
			log.Error("Error closing pool connection", logz.Error(err))
		} else {
			log.Debug("Closed pool connection")
		}
	}
	return nil
}

func NewPool(log *zap.Logger, dialOpts ...grpc.DialOption) *Pool {
	return &Pool{
		log:      log,
		dialOpts: dialOpts,
		conns:    map[string]*connHolder{},
		clk:      clock.RealClock{},
	}
}

func (p *Pool) Dial(ctx context.Context, targetUrl string) (PoolConn, error) {
	u, err := url.Parse(targetUrl)
	if err != nil {
		return nil, err
	}
	var target string
	switch u.Scheme {
	case "grpc":
		target = u.Host
	//case "grpcs":
	// TODO support TLS
	default:
		return nil, fmt.Errorf("unsupported pool URL scheme in %s", targetUrl)
	}
	p.mx.Lock()
	defer p.mx.Unlock()
	conn := p.conns[target]
	if conn == nil {
		grpcConn, err := grpc.DialContext(ctx, target, p.dialOpts...)
		if err != nil {
			return nil, fmt.Errorf("pool gRPC dial: %w", err)
		}
		conn = &connHolder{
			ClientConn: grpcConn,
			targetUrl:  targetUrl,
		}
		p.conns[target] = conn
	}
	conn.numUsers++
	return &poolConn{
		connHolder: conn,
		done:       p.connDone,
	}, nil
}

func (p *Pool) connDone(conn *connHolder) {
	p.mx.Lock()
	defer p.mx.Unlock()
	conn.numUsers--
	conn.lastUsed = p.clk.Now()
	p.runGcLocked()
}

func (p *Pool) runGcLocked() {
	expireAt := p.clk.Now().Add(-evictIdleConnAfter)
	for target, conn := range p.conns {
		if conn.numUsers == 0 && conn.lastUsed.Before(expireAt) {
			delete(p.conns, target)
			err := conn.Close()
			if err != nil {
				p.log.Error("Error closing idle pool connection", logz.Error(err), logz.PoolConnectionUrl(conn.targetUrl))
			} else {
				p.log.Debug("Closed idle pool connection", logz.PoolConnectionUrl(conn.targetUrl))
			}
		}
	}
}

type connHolder struct {
	*grpc.ClientConn
	targetUrl string
	lastUsed  time.Time
	numUsers  int32 // protected by mutex
}

type poolConn struct {
	*connHolder
	done func(conn *connHolder)
}

func (c *poolConn) Done() {
	if c.done == nil {
		panic(errors.New("pool connection Done() called more than once"))
	}
	done := c.done
	c.done = nil
	done(c.connHolder)
}
