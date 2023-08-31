package tunnel

import (
	"context"
	"time"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/reverse_tunnel/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/syncz"
	"go.uber.org/zap"
)

const (
	// refreshOverlap is the duration of an "overlap" between two refresh periods. It's a safety measure so that
	// a concurrent GC from another kas instance doesn't delete the data that is about to be refreshed.
	refreshOverlap = 5 * time.Second
	stripeBits     = 8
)

type Handler interface {
	// HandleTunnel is called with server-side interface of the reverse tunnel.
	// It registers the tunnel and blocks, waiting for a request to proxy through the tunnel.
	// The method returns the error value to return to gRPC framework.
	// ctx can be used to unblock the method if the tunnel is not being used already.
	// ctx should be a child of the server's context.
	HandleTunnel(ctx context.Context, agentInfo *api.AgentInfo, server rpc.ReverseTunnel_ConnectServer) error
}

type FindHandle interface {
	// Get finds a tunnel to an agentk.
	// It waits for a matching tunnel to proxy a connection through. When a matching tunnel is found, it is returned.
	// It returns gRPC status errors only, ready to return from RPC handler.
	Get(ctx context.Context) (Tunnel, error)
	// Done must be called to free resources of this FindHandle instance.
	Done()
}

type Finder interface {
	// FindTunnel starts searching for a tunnel to a matching agentk.
	// Found tunnel is:
	// - to an agent with provided id.
	// - supports handling provided gRPC service and method.
	// Tunnel found boolean indicates whether a suitable tunnel is immediately available from the
	// returned FindHandle object.
	FindTunnel(agentId int64, service, method string) (bool /* tunnel found */, FindHandle)
}

type Registry struct {
	log           *zap.Logger
	api           modshared.Api
	refreshPeriod time.Duration
	gcPeriod      time.Duration
	stripes       syncz.StripedValue[registryStripe]
}

func NewRegistry(log *zap.Logger, api modshared.Api, refreshPeriod, gcPeriod time.Duration,
	newTunnelTracker func() Tracker) (*Registry, error) {
	tunnelStreamVisitor, err := grpctool.NewStreamVisitor(&rpc.ConnectRequest{})
	if err != nil {
		return nil, err
	}
	return &Registry{
		log:           log,
		api:           api,
		refreshPeriod: refreshPeriod,
		gcPeriod:      gcPeriod,
		stripes: syncz.NewStripedValueInit(stripeBits, func() registryStripe {
			return registryStripe{
				log:                   log,
				api:                   api,
				tunnelStreamVisitor:   tunnelStreamVisitor,
				tunnelTracker:         newTunnelTracker(),
				tunsByAgentId:         make(map[int64]map[*tunnelImpl]struct{}),
				findRequestsByAgentId: make(map[int64]map[*findTunnelRequest]struct{}),
			}
		}),
	}, nil
}

func (r *Registry) FindTunnel(agentId int64, service, method string) (bool, FindHandle) {
	// Use GetPointer() to avoid copying the embedded mutex.
	return r.stripes.GetPointer(agentId).FindTunnel(agentId, service, method)
}

func (r *Registry) HandleTunnel(ctx context.Context, agentInfo *api.AgentInfo, server rpc.ReverseTunnel_ConnectServer) error {
	// Use GetPointer() to avoid copying the embedded mutex.
	return r.stripes.GetPointer(agentInfo.Id).HandleTunnel(ctx, agentInfo, server)
}

func (r *Registry) KasUrlsByAgentId(ctx context.Context, agentId int64, cb KasUrlsByAgentIdCallback) error {
	// Use GetPointer() to avoid copying the embedded mutex.
	return r.stripes.GetPointer(agentId).tunnelTracker.KasUrlsByAgentId(ctx, agentId, cb)
}

func (r *Registry) Run(ctx context.Context) error {
	defer r.stopInternal() // nolint: contextcheck
	refreshTicker := time.NewTicker(r.refreshPeriod)
	defer refreshTicker.Stop()
	gcTicker := time.NewTicker(r.gcPeriod)
	defer gcTicker.Stop()
	done := ctx.Done()
	for {
		select {
		case <-done:
			return nil
		case <-refreshTicker.C:
			r.refreshRegistrations(ctx, time.Now().Add(r.refreshPeriod-refreshOverlap))
		case <-gcTicker.C:
			r.runGC(ctx)
		}
	}
}

// stopInternal aborts any open tunnels.
// It should not be necessary to abort tunnels when registry is used correctly i.e. this method is called after
// all tunnels have terminated gracefully.
func (r *Registry) stopInternal() (int /*stoppedTun*/, int /*abortedFtr*/) {
	var stoppedTun, abortedFtr int
	for s := range r.stripes.Stripes { // use index var to avoid copying embedded mutex
		st, aftr := r.stripes.Stripes[s].Stop()
		stoppedTun += st
		abortedFtr += aftr
	}
	return stoppedTun, abortedFtr
}

func (r *Registry) refreshRegistrations(ctx context.Context, nextRefresh time.Time) {
	for s := range r.stripes.Stripes { // use index var to avoid copying embedded mutex
		err := r.stripes.Stripes[s].Refresh(ctx, nextRefresh)
		if err != nil {
			r.api.HandleProcessingError(ctx, r.log, modshared.NoAgentId, "Failed to refresh data in Redis", err)
		}
	}
}

func (r *Registry) runGC(ctx context.Context) int {
	total := 0
	for s := range r.stripes.Stripes { // use index var to avoid copying embedded mutex
		deletedKeys, err := r.stripes.Stripes[s].GC(ctx)
		if err != nil {
			r.api.HandleProcessingError(ctx, r.log, modshared.NoAgentId, "Failed to GC data in Redis", err)
			// fallthrough
		}
		total += deletedKeys
	}
	if total > 0 {
		r.log.Info("Deleted expired agent tunnel records", logz.RemovedHashKeys(total))
	}
	return total
}
