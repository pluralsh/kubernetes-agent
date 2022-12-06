package kasapp

import (
	"context"
	"fmt"
	"io"
	"sync"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/reverse_tunnel/tracker"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/retry"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	speculationFactor = 2
)

var (
	proxyStreamDesc = grpc.StreamDesc{
		ServerStreams: true,
		ClientStreams: true,
	}
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
	log              *zap.Logger
	kasPool          grpctool.PoolInterface
	tunnelQuerier    tracker.PollingQuerier
	rpcApi           modserver.RpcApi
	fullMethod       string // /service/method
	ownPrivateApiUrl string
	agentId          int64
	outgoingCtx      context.Context
	pollConfig       retry.PollConfigFactory
	foundTunnel      chan<- readyTunnel

	wg          wait.Group
	mu          sync.Mutex                // protects connections,done
	connections map[string]kasConnAttempt // kas URL -> conn info
	done        bool                      // successfully done searching
}

func (f *tunnelFinder) Run(ctx context.Context) {
	defer f.wg.Wait()
	pollCtx, pollCancel := context.WithCancel(ctx)
	defer pollCancel()

	// Unconditionally connect to self.
	f.tryKasInternal(f.ownPrivateApiUrl, pollCancel) // nolint: contextcheck

	newKasConnections := 0
	f.tunnelQuerier.PollKasUrlsByAgentId(pollCtx, f.agentId, func(newCycle bool, kasUrl string) bool {
		if newCycle {
			newKasConnections = 0
		}
		if newKasConnections < speculationFactor {
			if f.tryKas(kasUrl, pollCancel) {
				newKasConnections++
			}
		}
		return false
	})
}

func (f *tunnelFinder) tryKas(kasUrl string, pollCancel context.CancelFunc) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.done {
		return true
	}
	if _, ok := f.connections[kasUrl]; ok {
		return false // skip tunnel via kas that we have connected to already
	}
	f.tryKasInternal(kasUrl, pollCancel)
	return true
}

func (f *tunnelFinder) tryKasInternal(kasUrl string, pollCancel context.CancelFunc) {
	connCtx, connCancel := context.WithCancel(f.outgoingCtx)
	f.connections[kasUrl] = kasConnAttempt{
		cancel: connCancel,
	}
	f.wg.Start(func() {
		f.tryKasAsync(connCtx, connCancel, pollCancel, kasUrl)
	})
}

func (f *tunnelFinder) tryKasAsync(ctx context.Context, cancel, pollCancel context.CancelFunc, kasUrl string) {
	log := f.log.With(logz.KasUrl(kasUrl)) // nolint:govet
	// err can only be retry.ErrWaitTimeout
	_ = retry.PollWithBackoff(ctx, f.pollConfig(), func(ctx context.Context) (error, retry.AttemptResult) {
		success := false

		// 1. Dial another kas
		log.Debug("Trying tunnel")
		kasConn, err := f.kasPool.Dial(ctx, kasUrl)
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
			ctx,
			&proxyStreamDesc,
			f.fullMethod,
			grpc.ForceCodec(grpctool.RawCodecWithProtoFallback{}),
		)
		if err != nil {
			f.rpcApi.HandleProcessingError(log, f.agentId, "Failed to open new stream to kas", err)
			return nil, retry.Backoff
		}

		// 3. Wait for the other kas to say it's ready to start streaming i.e. has a suitable tunnel to an agent
		var kasResponse GatewayKasResponse
		err = kasStream.RecvMsg(&kasResponse) // Wait for the tunnel to be found
		if err != nil {
			if err == io.EOF { // nolint:errorlint
				// Gateway kas closed the connection cleanly, perhaps it's been open for too long
				return nil, retry.ContinueImmediately
			}
			f.rpcApi.HandleProcessingError(log, f.agentId, "RecvMsg(GatewayKasResponse)", err)
			return nil, retry.Backoff
		}
		if kasResponse.GetTunnelReady() == nil {
			f.rpcApi.HandleProcessingError(log, f.agentId, "GetTunnelReady()", fmt.Errorf("invalid oneof value type: %T", kasResponse.Msg))
			return nil, retry.Backoff
		}

		// 4. Check if another goroutine has found a suitable tunnel already
		f.mu.Lock()
		if f.done {
			f.mu.Unlock()
			return nil, retry.Done
		}
		f.done = true
		pollCancel()
		f.stopAllConnectionAttemptsExcept(kasUrl)
		f.mu.Unlock()

		// 5. Tell the other kas we are starting streaming
		err = kasStream.SendMsg(&StartStreaming{})
		if err != nil {
			if err == io.EOF { // nolint:errorlint
				var frame grpctool.RawFrame
				err = kasStream.RecvMsg(&frame) // get the real error
			}
			_ = f.rpcApi.HandleIoError(log, "SendMsg(StartStreaming)", err)
			return nil, retry.Backoff
		}
		rt := readyTunnel{
			kasUrl:          kasUrl,
			kasStream:       kasStream,
			kasConn:         kasConn,
			kasStreamCancel: cancel,
		}
		select {
		case <-ctx.Done():
		case f.foundTunnel <- rt:
			success = true
		}
		return nil, retry.Done
	})
}

func (f *tunnelFinder) stopAllConnectionAttemptsExcept(kasUrl string) {
	for url, c := range f.connections {
		if url != kasUrl {
			c.cancel()
		}
	}
}
