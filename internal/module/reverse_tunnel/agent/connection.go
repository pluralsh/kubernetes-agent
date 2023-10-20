package agent

import (
	"context"
	"fmt"
	"io"

	"github.com/pluralsh/kuberentes-agent/internal/module/reverse_tunnel/info"
	"github.com/pluralsh/kuberentes-agent/internal/module/reverse_tunnel/rpc"
	"github.com/pluralsh/kuberentes-agent/internal/tool/grpctool"
	"github.com/pluralsh/kuberentes-agent/internal/tool/logz"
	"github.com/pluralsh/kuberentes-agent/internal/tool/retry"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	requestInfoNumber protoreflect.FieldNumber = 1
	messageNumber     protoreflect.FieldNumber = 2
	closeSendNumber   protoreflect.FieldNumber = 3
)

var (
	proxyStreamDesc = grpc.StreamDesc{
		ServerStreams: true,
		ClientStreams: true,
	}
)

type connectionInterface interface {
	Run(attemptCtx, pollCtx context.Context)
}

type connection struct {
	log                *zap.Logger
	descriptor         *info.AgentDescriptor
	client             rpc.ReverseTunnelClient
	internalServerConn grpc.ClientConnInterface
	streamVisitor      *grpctool.StreamVisitor
	pollConfig         retry.PollConfigFactory
	onActive           func(connectionInterface)
	onIdle             func(connectionInterface)
}

func (c *connection) Run(attemptCtx, pollCtx context.Context) {
	_ = retry.PollWithBackoff(pollCtx, c.pollConfig(), func(_ context.Context) (error, retry.AttemptResult) {
		err := c.attempt(attemptCtx)
		if err != nil {
			if grpctool.RequestCanceledOrTimedOut(err) {
				c.log.Debug("Canceled connection", logz.Error(err))
				return nil, retry.ContinueImmediately // handled a connection with a timeout/cancel, re-establish it immediately
			}
			c.log.Error("Error handling a connection", logz.Error(err))
			return nil, retry.Backoff
		}
		c.log.Debug("Handled a connection successfully")
		return nil, retry.ContinueImmediately // successfully handled a connection, re-establish it immediately
	})
	c.log.Debug("Connection done")
}

func (c *connection) attempt(ctx context.Context) (retErr error) {
	defer c.onIdle(c)
	ctx, cancel, stopPropagation := propagateUntil(ctx)
	defer cancel()

	tunnel, err := c.client.Connect(ctx, grpc.WaitForReady(true))
	if err != nil {
		return fmt.Errorf("Connect(): %w", err) // wrap
	}
	err = tunnel.Send(&rpc.ConnectRequest{
		Msg: &rpc.ConnectRequest_Descriptor_{
			Descriptor_: &rpc.Descriptor{
				AgentDescriptor: c.descriptor,
			},
		},
	})
	if err != nil {
		if err == io.EOF { // nolint:errorlint
			_, err = tunnel.Recv() // get the actual error
		}
		return fmt.Errorf("Send(descriptor): %w", err) // wrap
	}
	var (
		clientStream grpc.ClientStream
		g            errgroup.Group
	)
	// pipe tunnel -> internal client
	err1 := c.streamVisitor.Visit(tunnel,
		grpctool.WithCallback(requestInfoNumber, func(reqInfo *rpc.RequestInfo) error {
			c.onActive(c)
			outgoingCtx := metadata.NewOutgoingContext(ctx, reqInfo.Metadata())
			clientStream, err = c.internalServerConn.NewStream(outgoingCtx, &proxyStreamDesc, reqInfo.MethodName)
			if err != nil {
				return fmt.Errorf("NewStream(): %w", err)
			}
			// After this point we don't want context cancellation to interrupt request handling since this would
			// break a running request. We are calling stopPropagation() to ignore the passed context and let the
			// request finish properly.
			stopPropagation()
			g.Go(func() error {
				// pipe internal client -> tunnel
				return pipeInternalClientIntoTunnel(tunnel, clientStream)
			})
			return nil
		}),
		grpctool.WithCallback(messageNumber, func(message *rpc.Message) error {
			err = clientStream.SendMsg(&grpctool.RawFrame{
				Data: message.Data,
			})
			if err != nil {
				if err == io.EOF { // nolint:errorlint
					return nil // the other goroutine will receive the error in RecvMsg()
				}
				return fmt.Errorf("SendMsg(): %w", err)
			}
			return nil
		}),
		grpctool.WithCallback(closeSendNumber, func(closeSend *rpc.CloseSend) error {
			err = clientStream.CloseSend()
			if err != nil {
				return fmt.Errorf("CloseSend(): %w", err)
			}
			return nil
		}),
	)
	if err1 != nil {
		cancel()
	}
	err2 := g.Wait()
	if err1 == nil {
		err1 = err2
	}
	return err1
}

func pipeInternalClientIntoTunnel(tunnel rpc.ReverseTunnel_ConnectClient, clientStream grpc.ClientStream) error {
	header, err := clientStream.Header()
	if err != nil {
		return sendErrorToTunnel(err, tunnel)
	}
	err = tunnel.Send(&rpc.ConnectRequest{
		Msg: &rpc.ConnectRequest_Header{
			Header: &rpc.Header{
				Meta: grpctool.MetaToValuesMap(header),
			},
		},
	})
	if err != nil {
		if err == io.EOF { // nolint:errorlint
			return nil // the other goroutine will receive the error in RecvMsg()
		}
		return fmt.Errorf("Send(header): %w", err) // wrap
	}
	var frame grpctool.RawFrame
	for {
		recvErr := clientStream.RecvMsg(&frame)
		if recvErr != nil {
			// Trailer becomes available after RecvMsg() returns an error
			err = tunnel.Send(&rpc.ConnectRequest{
				Msg: &rpc.ConnectRequest_Trailer{
					Trailer: &rpc.Trailer{
						Meta: grpctool.MetaToValuesMap(clientStream.Trailer()),
					},
				},
			})
			if err != nil {
				if recvErr == io.EOF { // nolint:errorlint
					if err == io.EOF { // nolint:errorlint
						return nil // the other goroutine will receive the error in RecvMsg()
					}
					return fmt.Errorf("Send(trailer): %w", err) // wrap
				}
				return fmt.Errorf("Recv(RawFrame): %w", recvErr) // return recvErr as it happened first
			}
			if recvErr == io.EOF { // nolint:errorlint
				break
			}
			return sendErrorToTunnel(recvErr, tunnel)
		}
		err = tunnel.Send(&rpc.ConnectRequest{
			Msg: &rpc.ConnectRequest_Message{
				Message: &rpc.Message{
					Data: frame.Data,
				},
			},
		})
		if err != nil {
			if err == io.EOF { // nolint:errorlint
				return nil // the other goroutine will receive the error in RecvMsg()
			}
			return fmt.Errorf("Send(message): %w", err) // wrap
		}
	}
	err = tunnel.CloseSend()
	if err != nil {
		return fmt.Errorf("CloseSend(): %w", err) // wrap
	}
	return nil
}

func sendErrorToTunnel(errToSend error, tunnel rpc.ReverseTunnel_ConnectClient) error {
	err := tunnel.Send(&rpc.ConnectRequest{
		Msg: &rpc.ConnectRequest_Error{
			Error: &rpc.Error{
				Status: status.Convert(errToSend).Proto(),
			},
		},
	})
	if err != nil {
		if err == io.EOF { // nolint:errorlint
			return nil // the other goroutine will receive the error in RecvMsg()
		}
		return fmt.Errorf("Send(error): %w", err) // wrap
	}
	err = tunnel.CloseSend()
	if err != nil {
		return fmt.Errorf("CloseSend(): %w", err) // wrap
	}
	return nil
}

// propagateUntil propagates cancellation from ctx to the returned context.
// When stop is called, ctx cancellation will no longer be propagated to the returned context.
// To cancel the returned context use the returned context.CancelFunc.
func propagateUntil(ctx context.Context) (context.Context, context.CancelFunc, func() /* stop */) {
	ctxInternal, cancel := context.WithCancel(context.Background())
	stopPropagation := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			cancel()
		case <-ctxInternal.Done():
		case <-stopPropagation:
		}
	}()
	return ctxInternal, cancel, func() {
		close(stopPropagation)
	}
}
