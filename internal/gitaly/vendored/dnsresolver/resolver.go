package dnsresolver

import (
	"context"
	"net"
	"sync"
	"time"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitaly/vendored/backoff"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitaly/vendored/structerr"
	"google.golang.org/grpc/resolver"
)

type dnsResolver struct {
	retry backoff.Strategy

	ctx         context.Context
	cancel      context.CancelFunc
	cc          resolver.ClientConn
	host        string
	port        string
	refreshRate time.Duration
	lookup      dnsLookuper
	reqs        chan struct{}
	wg          sync.WaitGroup
}

const dnsLookupTimeout = 15 * time.Second

type dnsLookuper interface {
	LookupHost(context.Context, string) ([]string, error)
}

// ResolveNow signals the resolver to perform a DNS resolution immediately. This method returns
// without waiting for the result. The resolver treats this as a hint rather than a command. The
// client connection receives the resolution result asynchronously via `clientconn.UpdateState`
// This method also skip resolver caching because it's likely the client calls this method after
// encounter an error with recent subchannels.
func (d *dnsResolver) ResolveNow(resolver.ResolveNowOptions) {
	select {
	case d.reqs <- struct{}{}:
	default:
	}
}

// Close cancels all activities of this dns resolver. It waits until the watch goroutine exits.
func (d *dnsResolver) Close() {
	d.cancel()
	d.wg.Wait()
}

func (d *dnsResolver) watch() {
	defer d.wg.Done()

	// Exponential retry after failed to resolve or client connection failed to update its state
	var retries uint
	for {
		state, err := d.resolve()
		if err != nil {
			d.cc.ReportError(err)
		} else {
			err = d.updateState(state)
		}

		var timer *time.Timer
		if err == nil {
			timer = time.NewTimer(d.refreshRate)
			retries = 0
		} else {
			timer = time.NewTimer(d.retry.Backoff(retries))
			retries++
		}

		select {
		case <-d.ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
			// Refresh timer expires, issue another DNS lookup.
			continue
		case <-d.reqs:
			// If the resolver is requested to resolve now, force notify the client
			// connection. Typically, client connection contacts the resolver when any
			// of the subchannels change its connectivity state.
			timer.Stop()
		}
	}
}

func (d *dnsResolver) updateState(state *resolver.State) error {
	return d.cc.UpdateState(*state)
}

func (d *dnsResolver) resolve() (*resolver.State, error) {
	ctx, cancel := context.WithTimeout(d.ctx, dnsLookupTimeout)
	defer cancel()

	addrs, err := d.lookup.LookupHost(ctx, d.host)
	if err != nil {
		err = handleDNSError(err)
		return &resolver.State{Addresses: []resolver.Address{}}, err
	}
	newAddrs := make([]resolver.Address, 0, len(addrs))
	for _, a := range addrs {
		addr, ok := tryParseIP(a, d.port)
		if !ok {
			return nil, structerr.New("dns: error parsing dns record IP address %v", a)
		}
		newAddrs = append(newAddrs, resolver.Address{Addr: addr})
	}

	return &resolver.State{Addresses: newAddrs}, nil
}

// handleDNSError massages the error to fit into expectations of the gRPC model:
//   - Timeouts and temporary errors should be communicated to gRPC to attempt another DNS query (with
//     Backoff).
//   - Other errors should be suppressed (they may represent the absence of a TXT record).
func handleDNSError(err error) error {
	if dnsErr, ok := err.(*net.DNSError); ok && !dnsErr.IsTimeout && !dnsErr.IsTemporary {
		return nil
	}

	return structerr.New("dns: record resolve error: %w", err)
}
