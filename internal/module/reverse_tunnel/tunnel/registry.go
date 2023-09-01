package tunnel

import (
	"context"
	"sync/atomic"
	"time"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/reverse_tunnel/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/syncz"
	"go.opentelemetry.io/otel/attribute"
	otelcodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	// refreshOverlap is the duration of an "overlap" between two refresh periods. It's a safety measure so that
	// a concurrent GC from another kas instance doesn't delete the data that is about to be refreshed.
	refreshOverlap = 5 * time.Second
	stopTimeout    = 5 * time.Second
	stripeBits     = 8
)

const (
	traceTunnelFoundAttr    attribute.Key = "found"
	traceDeletedKeysAttr    attribute.Key = "deletedKeys"
	traceStoppedTunnelsAttr attribute.Key = "stoppedTunnels"
	traceAbortedFTRAttr     attribute.Key = "abortedFTR"
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
	// ctx is used for tracing only.
	Done(ctx context.Context)
}

type Finder interface {
	// FindTunnel starts searching for a tunnel to a matching agentk.
	// Found tunnel is:
	// - to an agent with provided id.
	// - supports handling provided gRPC service and method.
	// Tunnel found boolean indicates whether a suitable tunnel is immediately available from the
	// returned FindHandle object.
	FindTunnel(ctx context.Context, agentId int64, service, method string) (bool, FindHandle)
}

type Registry struct {
	log           *zap.Logger
	api           modshared.Api
	tracer        trace.Tracer
	refreshPeriod time.Duration
	gcPeriod      time.Duration
	stripes       syncz.StripedValue[registryStripe]
}

func NewRegistry(log *zap.Logger, api modshared.Api, tracer trace.Tracer, refreshPeriod, gcPeriod time.Duration,
	newTunnelTracker func() Tracker) (*Registry, error) {
	tunnelStreamVisitor, err := grpctool.NewStreamVisitor(&rpc.ConnectRequest{})
	if err != nil {
		return nil, err
	}
	return &Registry{
		log:           log,
		api:           api,
		tracer:        tracer,
		refreshPeriod: refreshPeriod,
		gcPeriod:      gcPeriod,
		stripes: syncz.NewStripedValueInit(stripeBits, func() registryStripe {
			return registryStripe{
				log:                   log,
				api:                   api,
				tracer:                tracer,
				tunnelStreamVisitor:   tunnelStreamVisitor,
				tunnelTracker:         newTunnelTracker(),
				tunsByAgentId:         make(map[int64]map[*tunnelImpl]struct{}),
				findRequestsByAgentId: make(map[int64]map[*findTunnelRequest]struct{}),
			}
		}),
	}, nil
}

func (r *Registry) FindTunnel(ctx context.Context, agentId int64, service, method string) (bool, FindHandle) {
	ctx, span := r.tracer.Start(ctx, "Registry.FindTunnel", trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	// Use GetPointer() to avoid copying the embedded mutex.
	found, th := r.stripes.GetPointer(agentId).FindTunnel(ctx, agentId, service, method)
	span.SetAttributes(traceTunnelFoundAttr.Bool(found))
	return found, th
}

func (r *Registry) HandleTunnel(ctx context.Context, agentInfo *api.AgentInfo, server rpc.ReverseTunnel_ConnectServer) error {
	ctx, span := r.tracer.Start(ctx, "Registry.HandleTunnel", trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End() // we don't add the returned error to the span as it's added by the gRPC OTEL stats handler already.

	// Use GetPointer() to avoid copying the embedded mutex.
	return r.stripes.GetPointer(agentInfo.Id).HandleTunnel(ctx, agentInfo, server)
}

func (r *Registry) KasUrlsByAgentId(ctx context.Context, agentId int64, cb KasUrlsByAgentIdCallback) error {
	ctx, span := r.tracer.Start(ctx, "Registry.KasUrlsByAgentId", trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	// Use GetPointer() to avoid copying the embedded mutex.
	return r.stripes.GetPointer(agentId).tunnelTracker.KasUrlsByAgentId(ctx, agentId, cb)
}

func (r *Registry) Run(ctx context.Context) error {
	defer r.stopInternal(ctx) // nolint: contextcheck
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
func (r *Registry) stopInternal(ctx context.Context) (int /*stoppedTun*/, int /*abortedFtr*/) {
	ctx = contextWithoutCancel(ctx)
	ctx, cancel := context.WithTimeout(ctx, stopTimeout)
	defer cancel()
	ctx, span := r.tracer.Start(ctx, "Registry.stopInternal", trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	var wg wait.Group
	var stoppedTun, abortedFtr atomic.Int32

	for s := range r.stripes.Stripes { // use index var to avoid copying embedded mutex
		s := s
		wg.Start(func() {
			stopCtx, stopSpan := r.tracer.Start(ctx, "registryStripe.Stop", trace.WithSpanKind(trace.SpanKindInternal))
			defer stopSpan.End()

			st, aftr := r.stripes.Stripes[s].Stop(stopCtx)
			stoppedTun.Add(int32(st))
			abortedFtr.Add(int32(aftr))
			stopSpan.SetAttributes(traceStoppedTunnelsAttr.Int(st), traceAbortedFTRAttr.Int(aftr))
		})
	}
	wg.Wait()

	v1 := int(stoppedTun.Load())
	v2 := int(abortedFtr.Load())
	span.SetAttributes(traceStoppedTunnelsAttr.Int(v1), traceAbortedFTRAttr.Int(v2))
	return v1, v2
}

func (r *Registry) refreshRegistrations(ctx context.Context, nextRefresh time.Time) {
	ctx, span := r.tracer.Start(ctx, "Registry.refreshRegistrations", trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	for s := range r.stripes.Stripes { // use index var to avoid copying embedded mutex
		func() {
			refreshCtx, refreshSpan := r.tracer.Start(ctx, "registryStripe.Refresh", trace.WithSpanKind(trace.SpanKindInternal))
			defer refreshSpan.End()

			err := r.stripes.Stripes[s].Refresh(refreshCtx, nextRefresh)
			if err != nil {
				r.api.HandleProcessingError(refreshCtx, r.log, modshared.NoAgentId, "Failed to refresh data", err)
				refreshSpan.SetStatus(otelcodes.Error, err.Error())
				// fallthrough
			} else {
				refreshSpan.SetStatus(otelcodes.Ok, "")
			}
		}()
	}
}

func (r *Registry) runGC(ctx context.Context) int {
	ctx, span := r.tracer.Start(ctx, "Registry.runGC", trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	total := 0
	for s := range r.stripes.Stripes { // use index var to avoid copying embedded mutex
		func() {
			gcCtx, gcSpan := r.tracer.Start(ctx, "registryStripe.GC", trace.WithSpanKind(trace.SpanKindInternal))
			defer gcSpan.End()

			deletedKeys, err := r.stripes.Stripes[s].GC(gcCtx)
			if err != nil {
				r.api.HandleProcessingError(gcCtx, r.log, modshared.NoAgentId, "Failed to GC data", err)
				gcSpan.SetStatus(otelcodes.Error, err.Error())
				// fallthrough
			} else {
				gcSpan.SetStatus(otelcodes.Ok, "")
			}
			total += deletedKeys
			gcSpan.SetAttributes(traceDeletedKeysAttr.Int(deletedKeys))
		}()
	}
	span.SetAttributes(traceDeletedKeysAttr.Int(total))
	if total > 0 {
		r.log.Info("Deleted expired agent tunnel records", logz.RemovedHashKeys(total))
	}
	return total
}
