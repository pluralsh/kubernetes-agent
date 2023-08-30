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
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// refreshOverlap is the duration of an "overlap" between two refresh periods. It's a safety measure so that
	// a concurrent GC from another kas instance doesn't delete the data that is about to be refreshed.
	refreshOverlap = 5 * time.Second
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

type findTunnelRequest struct {
	agentId         int64
	service, method string
	retTun          chan<- *tunnelImpl
}

type findHandle struct {
	retTun    <-chan *tunnelImpl
	done      func()
	gotTunnel bool
}

func (h *findHandle) Get(ctx context.Context) (Tunnel, error) {
	select {
	case <-ctx.Done():
		return nil, grpctool.StatusErrorFromContext(ctx, "FindTunnel request aborted")
	case tun := <-h.retTun:
		h.gotTunnel = true
		if tun == nil {
			return nil, status.Error(codes.Unavailable, "kas is shutting down")
		}
		return tun, nil
	}
}

func (h *findHandle) Done() {
	if h.gotTunnel {
		// No cleanup needed if Get returned a tunnel.
		return
	}
	h.done()
}

type Registry struct {
	log                 *zap.Logger
	api                 modshared.Api
	tunnelRegisterer    Registerer
	refreshPeriod       time.Duration
	gcPeriod            time.Duration
	tunnelStreamVisitor *grpctool.StreamVisitor

	mu                    sync.Mutex
	tunsByAgentId         map[int64]map[*tunnelImpl]struct{}
	findRequestsByAgentId map[int64]map[*findTunnelRequest]struct{}
}

func NewRegistry(log *zap.Logger, api modshared.Api, tunnelRegisterer Registerer,
	refreshPeriod, gcPeriod time.Duration) (*Registry, error) {
	tunnelStreamVisitor, err := grpctool.NewStreamVisitor(&rpc.ConnectRequest{})
	if err != nil {
		return nil, err
	}
	return &Registry{
		log:                   log,
		api:                   api,
		tunnelRegisterer:      tunnelRegisterer,
		refreshPeriod:         refreshPeriod,
		gcPeriod:              gcPeriod,
		tunnelStreamVisitor:   tunnelStreamVisitor,
		tunsByAgentId:         make(map[int64]map[*tunnelImpl]struct{}),
		findRequestsByAgentId: make(map[int64]map[*findTunnelRequest]struct{}),
	}, nil
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
			err := r.refreshRegistrations(ctx, time.Now().Add(r.refreshPeriod-refreshOverlap))
			if err != nil {
				r.api.HandleProcessingError(ctx, r.log, modshared.NoAgentId, "Failed to refresh data in Redis", err)
			}
		case <-gcTicker.C:
			deletedKeys, err := r.runGC(ctx)
			if err != nil {
				r.api.HandleProcessingError(ctx, r.log, modshared.NoAgentId, "Failed to GC data in Redis", err)
				// fallthrough
			}
			if deletedKeys > 0 {
				r.log.Info("Deleted expired agent tunnel records", logz.RemovedHashKeys(deletedKeys))
			}
		}
	}
}

func (r *Registry) refreshRegistrations(ctx context.Context, nextRefresh time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.tunnelRegisterer.Refresh(ctx, nextRefresh)
}

func (r *Registry) runGC(ctx context.Context) (int /* keysDeleted */, error) {
	var gc func(ctx context.Context) (int /* keysDeleted */, error)
	func() {
		r.mu.Lock()
		defer r.mu.Unlock()
		gc = r.tunnelRegisterer.GC()
	}()
	return gc(ctx)
}

func (r *Registry) FindTunnel(agentId int64, service, method string) (bool /* tunnel found */, FindHandle) {
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
			return r.unregisterTunnelLocked(tun)
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
		r.api.HandleProcessingError(context.Background(), r.log.With(logz.AgentId(agentId)), agentId, "Failed to unregister tunnel", err)
	}
	return found, &findHandle{
		retTun: retTun,
		done: func() {
			err := func() error {
				r.mu.Lock()
				defer r.mu.Unlock()
				close(retTun)
				tun := <-retTun // will get nil if there was nothing in the channel or if registry is shutting down.
				if tun != nil {
					// Got the tunnel, but it's too late so return it to the registry.
					return r.onTunnelDoneLocked(tun)
				} else {
					r.deleteFindRequestLocked(ftr)
				}
				return nil
			}()
			if err != nil {
				r.api.HandleProcessingError(context.Background(), r.log.With(logz.AgentId(agentId)), agentId, "Failed to register tunnel", err)
			}
		},
	}
}

func (r *Registry) HandleTunnel(ctx context.Context, agentInfo *api.AgentInfo, server rpc.ReverseTunnel_ConnectServer) error {
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
		err = r.registerTunnelLocked(tun)
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
			err = r.unregisterTunnelLocked(tun) // nolint: contextcheck
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

func (r *Registry) registerTunnelLocked(toReg *tunnelImpl) error {
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
	return r.tunnelRegisterer.RegisterTunnel(context.Background(), agentId) // don't pass context to always register
}

func (r *Registry) unregisterTunnelLocked(toUnreg *tunnelImpl) error {
	agentId := toUnreg.agentId
	tunsByAgentId := r.tunsByAgentId[agentId]
	delete(tunsByAgentId, toUnreg)
	if len(tunsByAgentId) == 0 {
		delete(r.tunsByAgentId, agentId)
	}
	return r.tunnelRegisterer.UnregisterTunnel(context.Background(), agentId) // don't pass context to always unregister
}

func (r *Registry) onTunnelForward(tun *tunnelImpl) error {
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

func (r *Registry) onTunnelDone(tun *tunnelImpl) {
	var err error
	func() {
		r.mu.Lock()
		defer r.mu.Unlock()
		err = r.onTunnelDoneLocked(tun)
	}()
	if err != nil {
		r.api.HandleProcessingError(context.Background(), r.log.With(logz.AgentId(tun.agentId)), tun.agentId, "Failed to register tunnel", err)
	}
}

func (r *Registry) onTunnelDoneLocked(tun *tunnelImpl) error {
	switch tun.state {
	case stateReady:
		panic(errors.New("unreachable: ready -> done should never happen"))
	case stateFound:
		// Tunnel was found but was not used, Done() was called. Just put it back.
		return r.registerTunnelLocked(tun)
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

func (r *Registry) deleteFindRequestLocked(ftr *findTunnelRequest) {
	findRequestsForAgentId := r.findRequestsByAgentId[ftr.agentId]
	delete(findRequestsForAgentId, ftr)
	if len(findRequestsForAgentId) == 0 {
		delete(r.findRequestsByAgentId, ftr.agentId)
	}
}

// stopInternal aborts any open tunnels.
// It should not be necessary to abort tunnels when registry is used correctly i.e. this method is called after
// all tunnels have terminated gracefully.
func (r *Registry) stopInternal() (int, int) {
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
			err := r.unregisterTunnelLocked(tun)
			if err != nil {
				r.api.HandleProcessingError(context.Background(), r.log.With(logz.AgentId(agentId)), agentId, "Failed to unregister tunnel", err)
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
		r.api.HandleProcessingError(context.Background(), r.log, modshared.NoAgentId, "Stopped tunnels and aborted find requests", fmt.Errorf("num_tunnels=%d, num_find_requests=%d", stoppedTun, abortedFtr))
	}
	return stoppedTun, abortedFtr
}
