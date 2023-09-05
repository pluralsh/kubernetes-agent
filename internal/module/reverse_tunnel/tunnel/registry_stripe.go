package tunnel

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/reverse_tunnel/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/logz"
	otelcodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type findTunnelRequest struct {
	agentId         int64
	service, method string
	retTun          chan<- *tunnelImpl
}

type findHandle struct {
	tracer    trace.Tracer
	retTun    <-chan *tunnelImpl
	done      func(context.Context)
	gotTunnel bool
}

func (h *findHandle) Get(ctx context.Context) (Tunnel, error) {
	ctx, span := h.tracer.Start(ctx, "findHandle.Get", trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	select {
	case <-ctx.Done():
		span.SetStatus(otelcodes.Error, "FindTunnel request aborted")
		span.RecordError(ctx.Err())
		return nil, grpctool.StatusErrorFromContext(ctx, "FindTunnel request aborted")
	case tun := <-h.retTun:
		h.gotTunnel = true
		if tun == nil {
			span.SetStatus(otelcodes.Error, "kas is shutting down")
			return nil, status.Error(codes.Unavailable, "kas is shutting down")
		}
		span.SetStatus(otelcodes.Ok, "")
		return tun, nil
	}
}

func (h *findHandle) Done(ctx context.Context) {
	ctx, span := h.tracer.Start(ctx, "findHandle.Done", trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	if h.gotTunnel {
		// No cleanup needed if Get returned a tunnel.
		return
	}
	h.done(ctx)
}

type registryStripe struct {
	log                   *zap.Logger
	api                   modshared.Api
	tracer                trace.Tracer
	tunnelStreamVisitor   *grpctool.StreamVisitor
	mu                    sync.Mutex
	tunnelTracker         Tracker
	tunsByAgentId         map[int64]map[*tunnelImpl]struct{}
	findRequestsByAgentId map[int64]map[*findTunnelRequest]struct{}
}

func (r *registryStripe) Refresh(ctx context.Context, nextRefresh time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.tunnelTracker.Refresh(ctx, nextRefresh)
}

func (r *registryStripe) FindTunnel(ctx context.Context, agentId int64, service, method string) (bool, FindHandle) {
	// Buffer 1 to not block on send when a tunnel is found before find request is registered.
	retTun := make(chan *tunnelImpl, 1) // can receive nil from it if Stop() is called
	ftr := &findTunnelRequest{
		agentId: agentId,
		service: service,
		method:  method,
		retTun:  retTun,
	}
	found := false
	err := func() error {
		r.mu.Lock()
		defer r.mu.Unlock()

		// 1. Check if we have a suitable tunnel
		for tun := range r.tunsByAgentId[agentId] {
			if !tun.agentDescriptor.SupportsServiceAndMethod(service, method) {
				continue
			}
			// Suitable tunnel found!
			tun.state = stateFound
			retTun <- tun // must not block because the reception is below
			found = true
			return r.unregisterTunnelLocked(ctx, tun)
		}
		// 2. No suitable tunnel found, add to the queue
		findRequestsForAgentId := r.findRequestsByAgentId[agentId]
		if findRequestsForAgentId == nil {
			findRequestsForAgentId = make(map[*findTunnelRequest]struct{}, 1)
			r.findRequestsByAgentId[agentId] = findRequestsForAgentId
		}
		findRequestsForAgentId[ftr] = struct{}{}
		return nil
	}()
	if err != nil {
		r.api.HandleProcessingError(ctx, r.log.With(logz.AgentId(agentId)), agentId, "Failed to unregister tunnel", err)
	}
	return found, &findHandle{
		tracer: r.tracer,
		retTun: retTun,
		done: func(ctx context.Context) {
			err := func() error {
				r.mu.Lock()
				defer r.mu.Unlock()
				close(retTun)
				tun := <-retTun // will get nil if there was nothing in the channel or if registry is shutting down.
				if tun != nil {
					// Got the tunnel, but it's too late so return it to the registry.
					return r.onTunnelDoneLocked(ctx, tun)
				} else {
					r.deleteFindRequestLocked(ftr)
				}
				return nil
			}()
			if err != nil {
				r.api.HandleProcessingError(ctx, r.log.With(logz.AgentId(agentId)), agentId, "Failed to register tunnel", err)
			}
		},
	}
}

func (r *registryStripe) HandleTunnel(ctx context.Context, agentInfo *api.AgentInfo, server rpc.ReverseTunnel_ConnectServer) error {
	recv, err := server.Recv()
	if err != nil {
		return err
	}
	descriptor, ok := recv.Msg.(*rpc.ConnectRequest_Descriptor_)
	if !ok {
		return status.Errorf(codes.InvalidArgument, "Invalid oneof value type: %T", recv.Msg)
	}
	retErr := make(chan error, 1)
	agentId := agentInfo.Id
	tun := &tunnelImpl{
		tunnel:              server,
		tunnelStreamVisitor: r.tunnelStreamVisitor,
		tunnelRetErr:        retErr,
		agentId:             agentId,
		agentDescriptor:     descriptor.Descriptor_.AgentDescriptor,
		state:               stateReady,
		onForward:           r.onTunnelForward,
		onDone:              r.onTunnelDone,
	}
	// Register
	func() {
		r.mu.Lock()
		defer r.mu.Unlock()
		err = r.registerTunnelLocked(ctx, tun)
	}()
	if err != nil {
		r.api.HandleProcessingError(ctx, r.log.With(logz.AgentId(agentId)), agentId, "Failed to register tunnel", err)
	}
	// Wait for return error or for cancellation
	select {
	case <-ctx.Done():
		// Context canceled
		r.mu.Lock()
		switch tun.state {
		case stateReady:
			tun.state = stateContextDone
			err = r.unregisterTunnelLocked(ctx, tun) // nolint: contextcheck
			r.mu.Unlock()
			if err != nil {
				r.api.HandleProcessingError(ctx, r.log.With(logz.AgentId(agentId)), agentId, "Failed to unregister tunnel", err)
			}
			return nil
		case stateFound:
			// Tunnel was found but hasn't been used yet, Done() hasn't been called.
			// Set state to stateContextDone so that ForwardStream() errors out without doing any I/O.
			tun.state = stateContextDone
			r.mu.Unlock()
			return nil
		case stateForwarding:
			// I/O on the stream will error out, just wait for the return value.
			r.mu.Unlock()
			return <-retErr
		case stateDone:
			// Forwarding has finished and then ctx signaled done. Return the result value from forwarding.
			r.mu.Unlock()
			return <-retErr
		case stateContextDone:
			// Cannot happen twice.
			r.mu.Unlock()
			panic(errors.New("unreachable"))
		default:
			// Should never happen
			r.mu.Unlock()
			panic(fmt.Errorf("invalid state: %d", tun.state))
		}
	case err = <-retErr:
		return err
	}
}

func (r *registryStripe) registerTunnelLocked(ctx context.Context, toReg *tunnelImpl) error {
	agentId := toReg.agentId
	// 1. Before registering the tunnel see if there is a find tunnel request waiting for it
	findRequestsForAgentId := r.findRequestsByAgentId[agentId]
	for ftr := range findRequestsForAgentId {
		if !toReg.agentDescriptor.SupportsServiceAndMethod(ftr.service, ftr.method) {
			continue
		}
		// Waiting request found!
		toReg.state = stateFound
		ftr.retTun <- toReg            // Satisfy the waiting request ASAP
		r.deleteFindRequestLocked(ftr) // Remove it from the queue
		return nil
	}

	// 2. Register the tunnel
	toReg.state = stateReady
	tunsByAgentId := r.tunsByAgentId[agentId]
	if tunsByAgentId == nil {
		tunsByAgentId = make(map[*tunnelImpl]struct{}, 1)
		r.tunsByAgentId[agentId] = tunsByAgentId
	}
	tunsByAgentId[toReg] = struct{}{}
	// don't pass the original context to always register
	return r.tunnelTracker.RegisterTunnel(contextWithoutCancel(ctx), agentId)
}

func (r *registryStripe) unregisterTunnelLocked(ctx context.Context, toUnreg *tunnelImpl) error {
	agentId := toUnreg.agentId
	tunsByAgentId := r.tunsByAgentId[agentId]
	delete(tunsByAgentId, toUnreg)
	if len(tunsByAgentId) == 0 {
		delete(r.tunsByAgentId, agentId)
	}
	// don't pass the original context to always unregister
	return r.tunnelTracker.UnregisterTunnel(contextWithoutCancel(ctx), agentId)
}

func (r *registryStripe) onTunnelForward(tun *tunnelImpl) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	switch tun.state {
	case stateReady:
		return status.Error(codes.Internal, "unreachable: ready -> forwarding should never happen")
	case stateFound:
		tun.state = stateForwarding
		return nil
	case stateForwarding:
		return status.Error(codes.Internal, "ForwardStream() called more than once")
	case stateDone:
		return status.Error(codes.Internal, "ForwardStream() called after Done()")
	case stateContextDone:
		return status.Error(codes.Canceled, "ForwardStream() called on done stream")
	default:
		return status.Errorf(codes.Internal, "unreachable: invalid state: %d", tun.state)
	}
}

func (r *registryStripe) onTunnelDone(ctx context.Context, tun *tunnelImpl) {
	ctx, span := r.tracer.Start(ctx, "registryStripe.onTunnelDone", trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	var err error
	func() {
		r.mu.Lock()
		defer r.mu.Unlock()
		err = r.onTunnelDoneLocked(ctx, tun)
	}()
	if err != nil {
		r.api.HandleProcessingError(ctx, r.log.With(logz.AgentId(tun.agentId)), tun.agentId, "Failed to register tunnel", err)
	}
}

func (r *registryStripe) onTunnelDoneLocked(ctx context.Context, tun *tunnelImpl) error {
	switch tun.state {
	case stateReady:
		panic(errors.New("unreachable: ready -> done should never happen"))
	case stateFound:
		// Tunnel was found but was not used, Done() was called. Just put it back.
		return r.registerTunnelLocked(ctx, tun)
	case stateForwarding:
		tun.state = stateDone
	case stateDone:
		panic(errors.New("Done() called more than once"))
	case stateContextDone:
	// Done() called after cancelled context in HandleTunnel(). Nothing to do.
	default:
		// Should never happen
		panic(fmt.Errorf("invalid state: %d", tun.state))
	}
	return nil
}

func (r *registryStripe) deleteFindRequestLocked(ftr *findTunnelRequest) {
	findRequestsForAgentId := r.findRequestsByAgentId[ftr.agentId]
	delete(findRequestsForAgentId, ftr)
	if len(findRequestsForAgentId) == 0 {
		delete(r.findRequestsByAgentId, ftr.agentId)
	}
}

// Stop aborts any open tunnels.
// It should not be necessary to abort tunnels when registry is used correctly i.e. this method is called after
// all tunnels have terminated gracefully.
func (r *registryStripe) Stop(ctx context.Context) (int /*stoppedTun*/, int /*abortedFtr*/) {
	stoppedTun := 0
	abortedFtr := 0

	r.mu.Lock()
	defer r.mu.Unlock()

	// 1. Abort all tunnels
	for agentId, tunSet := range r.tunsByAgentId {
		for tun := range tunSet {
			stoppedTun++
			tun.state = stateDone
			tun.tunnelRetErr <- nil // nil so that HandleTunnel() returns cleanly and agent immediately retries
			err := r.unregisterTunnelLocked(ctx, tun)
			if err != nil {
				r.api.HandleProcessingError(ctx, r.log.With(logz.AgentId(agentId)), agentId, "Failed to unregister tunnel", err)
			}
		}
	}

	// 2. Abort all waiting new stream requests
	for _, findRequestsForAgentId := range r.findRequestsByAgentId {
		for ftr := range findRequestsForAgentId {
			abortedFtr++
			ftr.retTun <- nil // respond ASAP, then do all the bookkeeping
			r.deleteFindRequestLocked(ftr)
		}
	}
	if stoppedTun > 0 || abortedFtr > 0 {
		r.api.HandleProcessingError(ctx, r.log, modshared.NoAgentId, "Stopped tunnels and aborted find requests", fmt.Errorf("num_tunnels=%d, num_find_requests=%d", stoppedTun, abortedFtr))
	}
	return stoppedTun, abortedFtr
}

// contextWithoutCancel returns an independent context with existing tracing span.
// TODO use https://pkg.go.dev/context#WithoutCancel after upgrading to Go 1.21.
func contextWithoutCancel(ctx context.Context) context.Context {
	return trace.ContextWithSpan(context.Background(), trace.SpanFromContext(ctx))
}
