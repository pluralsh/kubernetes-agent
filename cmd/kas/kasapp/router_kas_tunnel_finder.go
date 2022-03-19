package kasapp

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/reverse_tunnel/tracker"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/mathz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/retry"
	"go.uber.org/zap"
	"google.golang.org/grpc"
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

type tunnelFinder struct {
	log           *zap.Logger
	kasPool       KasPool
	tunnelQuerier tracker.Querier
	rpcApi        modserver.RpcApi
	fullMethod    string // /service/method
	agentId       int64
	outgoingCtx   context.Context
	foundTunnel   chan<- readyTunnel

	mu          sync.Mutex                // protects connections,done
	connections map[string]kasConnAttempt // kas URL -> conn info
	done        bool                      // successfully done searching
}

func (f *tunnelFinder) poll(ctx context.Context, pollConfig retry.PollConfig) {
	var tunnels []*tracker.TunnelInfo
	pollCtx, pollCancel := context.WithCancel(ctx)
	defer pollCancel()
	getTunnelsFunc := f.attemptToGetTunnels(pollCtx, &tunnels)
	_ = retry.PollImmediateUntil(pollCtx, pollConfig.Interval, func() ( /*done*/ bool, error) {
		err := retry.PollWithBackoff(pollCtx, pollConfig, getTunnelsFunc)
		if err != nil {
			return false, err // err can only be retry.ErrWaitTimeout
		}
		for _, tunnel := range tunnels {
			if f.handleTunnel(tunnel, pollCancel) { // nolint: contextcheck
				break
			}
		}
		return false, nil
	})
}

func (f *tunnelFinder) handleTunnel(tunnel *tracker.TunnelInfo, pollCancel context.CancelFunc) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.done {
		return true
	}
	if _, ok := f.connections[tunnel.KasUrl]; ok {
		return false // skip tunnel via kas that we have connected to already
	}
	connCtx, connCancel := context.WithCancel(f.outgoingCtx)
	f.connections[tunnel.KasUrl] = kasConnAttempt{
		cancel: connCancel,
	}
	go f.handleTunnelAsync(connCtx, connCancel, pollCancel, tunnel)
	return true
}

func (f *tunnelFinder) handleTunnelAsync(ctx context.Context, cancel, pollCancel context.CancelFunc, tunnel *tracker.TunnelInfo) {
	success := false
	defer func() {
		if !success {
			cancel()
		}
		f.mu.Lock()
		defer f.mu.Unlock()
		delete(f.connections, tunnel.KasUrl)
	}()

	// 1. Dial another kas
	log := f.log.With(logz.KasUrl(tunnel.KasUrl)) // nolint:govet
	log.Debug("Trying tunnel")
	kasConn, err := f.kasPool.Dial(ctx, tunnel.KasUrl)
	if err != nil {
		f.rpcApi.HandleProcessingError(log, f.agentId, "Failed to dial kas", err)
		return
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
		return
	}

	// 3. Wait for the other kas to say it's ready to start streaming i.e. has a suitable tunnel to an agent
	var kasResponse GatewayKasResponse
	err = kasStream.RecvMsg(&kasResponse) // Wait for the tunnel to be found
	if err != nil {
		if errors.Is(err, io.EOF) {
			// Gateway kas closed the connection cleanly, perhaps it's been open for too long
			return
		}
		f.rpcApi.HandleProcessingError(log, f.agentId, "RecvMsg(GatewayKasResponse)", err)
		return
	}
	if kasResponse.GetTunnelReady() == nil {
		f.rpcApi.HandleProcessingError(log, f.agentId, "GetTunnelReady()", fmt.Errorf("invalid oneof value type: %T", kasResponse.Msg))
		return
	}

	// 4. Check if another goroutine has found a suitable tunnel already
	f.mu.Lock()
	if f.done {
		f.mu.Unlock()
		return
	}
	f.done = true
	pollCancel()
	f.stopAllConnectionAttemptsExcept(tunnel.KasUrl)
	f.mu.Unlock()

	// 5. Tell the other kas we are starting streaming
	err = kasStream.SendMsg(&StartStreaming{})
	if err != nil {
		if errors.Is(err, io.EOF) {
			var frame grpctool.RawFrame
			err = kasStream.RecvMsg(&frame) // get the real error
		}
		_ = f.rpcApi.HandleSendError(log, "SendMsg(StartStreaming)", err)
		return
	}
	rt := readyTunnel{
		kasUrl:          tunnel.KasUrl,
		kasStream:       kasStream,
		kasConn:         kasConn,
		kasStreamCancel: cancel,
	}
	select {
	case <-ctx.Done():
	case f.foundTunnel <- rt:
		success = true
	}
}

// attemptToGetTunnels
// must return a gRPC status-compatible error or retry.ErrWaitTimeout.
func (f *tunnelFinder) attemptToGetTunnels(ctx context.Context, infosTarget *[]*tracker.TunnelInfo) retry.PollWithBackoffFunc {
	service, method := grpctool.SplitGrpcMethod(f.fullMethod)
	return func() (error, retry.AttemptResult) {
		var infos tunnelInfoCollector = (*infosTarget)[:0] // reuse target backing array
		err := f.tunnelQuerier.GetTunnelsByAgentId(ctx, f.agentId, infos.Collect(service, method))
		if err != nil {
			f.rpcApi.HandleProcessingError(f.log, f.agentId, "GetTunnelsByAgentId()", err)
			return nil, retry.Backoff
		}
		mathz.Shuffle(len(infos), func(i, j int) {
			infos[i], infos[j] = infos[j], infos[i]
		})
		*infosTarget = infos
		return nil, retry.Done
	}
}

func (f *tunnelFinder) stopAllConnectionAttemptsExcept(kasUrl string) {
	for url, c := range f.connections {
		if url != kasUrl {
			c.cancel()
		}
	}
}

type tunnelInfoCollector []*tracker.TunnelInfo

func (c *tunnelInfoCollector) Collect(service, method string) tracker.GetTunnelsByAgentIdCallback {
	return func(info *tracker.TunnelInfo) (bool /* done */, error) {
		if info.KasUrl == "" {
			// kas without a private API endpoint. Ignore it.
			// TODO this can be made mandatory if/when the env var the address is coming from is mandatory
			return false, nil
		}
		if !info.SupportsServiceAndMethod(service, method) {
			// This tunnel doesn't support required API. Ignore it.
			return false, nil
		}
		*c = append(*c, info)
		return false, nil
	}
}
