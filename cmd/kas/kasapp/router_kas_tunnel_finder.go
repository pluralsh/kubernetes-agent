package kasapp

import (
	"context"
	"errors"
	"io"
	"sync"
	"time"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/reverse_tunnel/tracker"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/retry"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	tryNewKasInterval = 10 * time.Millisecond
)

var (
	proxyStreamDesc = grpc.StreamDesc{
		ServerStreams: true,
		ClientStreams: true,
	}

	// tunnelReadySentinelError is a sentinel error value to make stream visitor exit early.
	tunnelReadySentinelError = errors.New("")
)

type kasConnAttempt struct {
	cancel context.CancelFunc
}

type readyTunnel struct {
	kasUrl          string
	kasStream       grpc.ClientStream
	kasConn         grpctool.PoolConn
	kasStreamCancel context.CancelFunc
}

func (t readyTunnel) Done() {
	t.kasStreamCancel()
	t.kasConn.Done()
}

type tunnelFinder struct {
	log               *zap.Logger
	kasPool           grpctool.PoolInterface
	tunnelQuerier     tracker.PollingQuerier
	rpcApi            modserver.RpcApi
	fullMethod        string // /service/method
	ownPrivateApiUrl  string
	agentId           int64
	outgoingCtx       context.Context
	pollConfig        retry.PollConfigFactory
	gatewayKasVisitor *grpctool.StreamVisitor
	foundTunnel       chan readyTunnel
	noTunnel          chan struct{}
	wg                wait.Group
	pollCancel        context.CancelFunc

	mu          sync.Mutex                // protects the fields below
	connections map[string]kasConnAttempt // kas URL -> conn info
	kasUrls     []string                  // current known kas URLs for the agent id
	done        bool                      // successfully done searching
}

func newTunnelFinder(log *zap.Logger, kasPool grpctool.PoolInterface, tunnelQuerier tracker.PollingQuerier,
	rpcApi modserver.RpcApi, fullMethod string, ownPrivateApiUrl string, agentId int64, outgoingCtx context.Context,
	pollConfig retry.PollConfigFactory, gatewayKasVisitor *grpctool.StreamVisitor) *tunnelFinder {
	return &tunnelFinder{
		log:               log,
		kasPool:           kasPool,
		tunnelQuerier:     tunnelQuerier,
		rpcApi:            rpcApi,
		fullMethod:        fullMethod,
		ownPrivateApiUrl:  ownPrivateApiUrl,
		agentId:           agentId,
		outgoingCtx:       outgoingCtx,
		pollConfig:        pollConfig,
		gatewayKasVisitor: gatewayKasVisitor,
		foundTunnel:       make(chan readyTunnel),
		noTunnel:          make(chan struct{}),
		connections:       make(map[string]kasConnAttempt),
	}
}

func (f *tunnelFinder) Find(ctx context.Context) (readyTunnel, error) {
	defer f.wg.Wait()
	var pollCtx context.Context
	pollCtx, f.pollCancel = context.WithCancel(ctx)
	defer f.pollCancel()

	// Unconditionally connect to self ASAP.
	f.tryKasLocked(f.ownPrivateApiUrl) // nolint: contextcheck
	startedPolling := false
	// This flag is set when we've run out of kas URLs to try. When a new set of URLs is received, if this is set,
	// we try to connect to one of those URLs.
	needToTryNewKas := false

	// Timer is used to wake up the loop below after a certain amount of time has passed but there has been no activity,
	// in particular, a recently connected to kas didn't reply with noTunnel. If it's not replying, we need to try
	// another instance if it has been discovered.
	// If, for some reason, our own private API server doesn't respond with noTunnel/startStreaming in time, we
	// want to proceed with normal flow too.
	t := time.NewTimer(tryNewKasInterval)
	defer t.Stop()
	kasUrlsC := make(chan []string)
	f.kasUrls = f.tunnelQuerier.CachedKasUrlsByAgentId(f.agentId)
	done := ctx.Done()

	// Timer must have been stopped or has fired when this function is called
	tryNextKasWhenTimerNotRunning := func() {
		if f.tryNextKas() { // nolint: contextcheck
			// Connected to an instance.
			needToTryNewKas = false
			t.Reset(tryNewKasInterval)
		} else {
			// Couldn't find a kas instance we haven't connected to already.
			needToTryNewKas = true
			if !startedPolling {
				startedPolling = true
				// No more cached instances, start polling for kas instances.
				f.wg.Start(func() {
					pollDone := pollCtx.Done()
					f.tunnelQuerier.PollKasUrlsByAgentId(pollCtx, f.agentId, func(kasUrls []string) {
						select {
						case <-pollDone:
						case kasUrlsC <- kasUrls:
						}
					})
				})
			}
		}
	}

	for {
		select {
		case <-done:
			f.stopAllConnectionAttempts()
			return readyTunnel{}, ctx.Err()
		case <-f.noTunnel:
			stopAndDrain(t)
			tryNextKasWhenTimerNotRunning()
		case kasUrls := <-kasUrlsC:
			f.mu.Lock()
			f.kasUrls = kasUrls
			f.mu.Unlock()
			if !needToTryNewKas {
				continue
			}
			if f.tryNextKas() { // nolint: contextcheck
				// Connected to a new kas instance.
				needToTryNewKas = false
				stopAndDrain(t)
				t.Reset(tryNewKasInterval)
			}
		case <-t.C:
			tryNextKasWhenTimerNotRunning()
		case rt := <-f.foundTunnel:
			f.stopAllConnectionAttemptsExcept(rt.kasUrl)
			return rt, nil
		}
	}
}

func (f *tunnelFinder) tryNextKas() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, kasUrl := range f.kasUrls {
		if _, ok := f.connections[kasUrl]; ok {
			continue // skip tunnel via kas that we have connected to already
		}
		f.tryKasLocked(kasUrl)
		return true
	}
	return false
}

func (f *tunnelFinder) tryKasLocked(kasUrl string) {
	connCtx, connCancel := context.WithCancel(f.outgoingCtx)
	f.connections[kasUrl] = kasConnAttempt{
		cancel: connCancel,
	}
	f.wg.Start(func() {
		f.tryKasAsync(connCtx, connCancel, kasUrl)
	})
}

func (f *tunnelFinder) tryKasAsync(ctx context.Context, cancel context.CancelFunc, kasUrl string) {
	log := f.log.With(logz.KasUrl(kasUrl)) // nolint:govet
	noTunnelSent := false
	_ = retry.PollWithBackoff(ctx, f.pollConfig(), func(ctx context.Context) (error, retry.AttemptResult) {
		success := false

		// 1. Dial another kas
		log.Debug("Trying tunnel")
		attemptCtx, attemptCancel := context.WithCancel(ctx)
		defer func() {
			if !success {
				attemptCancel()
				f.maybeStopTrying(kasUrl)
			}
		}()
		kasConn, err := f.kasPool.Dial(attemptCtx, kasUrl)
		if err != nil {
			f.rpcApi.HandleProcessingError(log, f.agentId, "Failed to dial kas", err)
			return nil, retry.Backoff
		}
		defer func() {
			if !success {
				kasConn.Done()
			}
		}()

		// 2. Open a stream to the desired service/method
		kasStream, err := kasConn.NewStream(
			attemptCtx,
			&proxyStreamDesc,
			f.fullMethod,
			grpc.ForceCodec(grpctool.RawCodecWithProtoFallback{}),
			grpc.WaitForReady(true),
		)
		if err != nil {
			f.rpcApi.HandleProcessingError(log, f.agentId, "Failed to open new stream to kas", err)
			return nil, retry.Backoff
		}

		// 3. Wait for the other kas to say it's ready to start streaming i.e. has a suitable tunnel to an agent
		err = f.gatewayKasVisitor.Visit(kasStream,
			grpctool.WithCallback(noTunnelFieldNumber, func(noTunnel *GatewayKasResponse_NoTunnel) error {
				trace.SpanFromContext(kasStream.Context()).AddEvent("No tunnel") // nolint: contextcheck
				if !noTunnelSent {                                               // send only once
					noTunnelSent = true
					// Let Find() know there is no tunnel available from that kas instantaneously.
					// A tunnel may still be found when a suitable agent connects later, but none available immediately.
					select {
					case <-attemptCtx.Done():
					case f.noTunnel <- struct{}{}:
					}
				}
				return nil
			}),
			grpctool.WithCallback(tunnelReadyFieldNumber, func(tunnelReady *GatewayKasResponse_TunnelReady) error {
				trace.SpanFromContext(kasStream.Context()).AddEvent("Ready")
				return tunnelReadySentinelError
			}),
			grpctool.WithNotExpectingToGet(codes.Internal, headerFieldNumber, messageFieldNumber, trailerFieldNumber, errorFieldNumber),
		)
		switch err { // nolint:errorlint
		case nil:
			// Gateway kas closed the connection cleanly, perhaps it's been open for too long
			return nil, retry.ContinueImmediately
		case tunnelReadySentinelError:
			// fallthrough
		default:
			f.rpcApi.HandleProcessingError(log, f.agentId, "RecvMsg(GatewayKasResponse)", err)
			return nil, retry.Backoff
		}

		// 4. Check if another goroutine has found a suitable tunnel already
		f.mu.Lock() // Ensure only one kas gets StartStreaming message
		if f.done {
			f.mu.Unlock()
			return nil, retry.Done
		}
		// 5. Tell the other kas we are starting streaming
		err = kasStream.SendMsg(&StartStreaming{})
		if err != nil {
			f.mu.Unlock()
			if err == io.EOF { // nolint:errorlint
				var frame grpctool.RawFrame
				err = kasStream.RecvMsg(&frame) // get the real error
			}
			_ = f.rpcApi.HandleIoError(log, "SendMsg(StartStreaming)", err)
			return nil, retry.Backoff
		}
		f.done = true
		f.mu.Unlock()
		f.pollCancel()
		rt := readyTunnel{
			kasUrl:          kasUrl,
			kasStream:       kasStream,
			kasConn:         kasConn,
			kasStreamCancel: cancel,
		}
		select {
		case <-attemptCtx.Done():
		case f.foundTunnel <- rt:
			success = true
		}
		return nil, retry.Done
	})
}

func (f *tunnelFinder) maybeStopTrying(tryingKasUrl string) {
	if tryingKasUrl == f.ownPrivateApiUrl {
		return // keep trying the own URL
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, kasUrl := range f.kasUrls {
		if kasUrl == tryingKasUrl {
			return // known URLs still contain this URL so keep trying it.
		}
	}
	attempt := f.connections[tryingKasUrl]
	delete(f.connections, tryingKasUrl)
	attempt.cancel()
}

func (f *tunnelFinder) stopAllConnectionAttemptsExcept(kasUrl string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for url, c := range f.connections {
		if url != kasUrl {
			c.cancel()
		}
	}
}

func (f *tunnelFinder) stopAllConnectionAttempts() {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, c := range f.connections {
		c.cancel()
	}
}

func stopAndDrain(t *time.Timer) {
	if !t.Stop() {
		select {
		case <-t.C:
		default:
		}
	}
}
