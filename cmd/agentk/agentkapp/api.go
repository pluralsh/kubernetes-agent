package agentkapp

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	gitlab_access_rpc "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/gitlab_access/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/errz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/prototool"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	maxDataChunkSize = 32 * 1024

	headerFieldNumber  protoreflect.FieldNumber = 1
	dataFieldNumber    protoreflect.FieldNumber = 2
	trailerFieldNumber protoreflect.FieldNumber = 3
)

// agentAPI is an implementation of modagent.API.
type agentAPI struct {
	moduleName      string
	client          gitlab_access_rpc.GitlabAccessClient
	responseVisitor *grpctool.StreamVisitor
	featureTracker  *featureTracker
}

func (a *agentAPI) HandleProcessingError(ctx context.Context, log *zap.Logger, agentId int64, msg string, err error) {
	handleProcessingError(ctx, log, agentId, msg, err)
}

func (a *agentAPI) HandleSendError(log *zap.Logger, msg string, err error) error {
	return handleSendError(log, msg, err)
}

func (a *agentAPI) ToggleFeature(feature modagent.Feature, enabled bool) {
	a.featureTracker.ToggleFeature(feature, a.moduleName, enabled)
}

func (a *agentAPI) SubscribeToFeatureStatus(feature modagent.Feature, cb modagent.SubscribeCb) {
	a.featureTracker.Subscribe(feature, cb)
}

func (a *agentAPI) MakeGitLabRequest(ctx context.Context, path string, opts ...modagent.GitLabRequestOption) (*modagent.GitLabResponse, error) {
	config := modagent.ApplyRequestOptions(opts)
	ctx, cancel := context.WithCancel(ctx)
	client, errReq := a.client.MakeRequest(ctx)
	if errReq != nil {
		cancel()
		if config.Body != nil {
			_ = config.Body.Close()
		}
		return nil, errReq
	}
	pr, pw := io.Pipe()
	val := newValueOrError(func(err error) error {
		cancel()                                               // 1. Cancel the other goroutine and the client.
		err = grpctool.MaybeWrapWithCorrelationId(err, client) // 2. Get correlation id from header (can block if called before cancel)
		_ = pw.CloseWithError(err)                             // 3. Close the "write side" of the pipe
		return err
	})

	var wg wait.Group
	// Write request
	wg.Start(func() {
		err := a.makeRequest(client, path, config)
		if err != nil {
			val.SetError(err)
		}
	})
	// Read response
	wg.Start(func() {
		err := a.responseVisitor.Visit(client,
			grpctool.WithCallback(headerFieldNumber, func(header *grpctool.HttpResponse_Header) error {
				val.SetValue(&modagent.GitLabResponse{
					Status:     header.Response.Status,
					StatusCode: header.Response.StatusCode,
					Header:     header.Response.HttpHeader(),
					Body: cancelingReadCloser{
						ReadCloser: pr,
						cancel:     cancel,
					},
				})
				return nil
			}),
			grpctool.WithCallback(dataFieldNumber, func(data *grpctool.HttpResponse_Data) error {
				_, err := pw.Write(data.Data)
				return err
			}),
			grpctool.WithCallback(trailerFieldNumber, func(trailer *grpctool.HttpResponse_Trailer) error {
				return nil
			}),
			grpctool.WithEOFCallback(func() error {
				return pw.Close()
			}),
		)
		if err != nil {
			val.SetError(err)
		}
	})
	resp, err := val.Wait()
	if err != nil {
		wg.Wait() // Wait for both goroutines to finish before returning
		return nil, err
	}
	return resp.(*modagent.GitLabResponse), nil
}

func (a *agentAPI) makeRequest(client gitlab_access_rpc.GitlabAccess_MakeRequestClient, path string, config *modagent.GitLabRequestConfig) (retErr error) {
	defer errz.SafeClose(config.Body, &retErr)
	extra, err := anypb.New(&gitlab_access_rpc.HeaderExtra{
		ModuleName: a.moduleName,
	})
	if err != nil {
		return err
	}
	err = client.Send(&grpctool.HttpRequest{
		Message: &grpctool.HttpRequest_Header_{
			Header: &grpctool.HttpRequest_Header{
				Request: &prototool.HttpRequest{
					Method:  config.Method,
					Header:  prototool.HttpHeaderToValuesMap(config.Header),
					UrlPath: path,
					Query:   prototool.UrlValuesToValuesMap(config.Query),
				},
				Extra: extra,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("send request header: %w", err) // wrap
	}
	if config.Body != nil {
		buffer := make([]byte, maxDataChunkSize)
		for {
			var n int
			n, err = config.Body.Read(buffer)
			if err != nil && !errors.Is(err, io.EOF) {
				return fmt.Errorf("send request body: %w", err) // wrap
			}
			if n > 0 { // handle n=0, err=io.EOF case
				sendErr := client.Send(&grpctool.HttpRequest{
					Message: &grpctool.HttpRequest_Data_{
						Data: &grpctool.HttpRequest_Data{
							Data: buffer[:n],
						}},
				})
				if sendErr != nil {
					return fmt.Errorf("send request data: %w", sendErr) // wrap
				}
			}
			if errors.Is(err, io.EOF) {
				break
			}
		}
	}
	err = client.Send(&grpctool.HttpRequest{
		Message: &grpctool.HttpRequest_Trailer_{
			Trailer: &grpctool.HttpRequest_Trailer{},
		},
	})
	if err != nil {
		return fmt.Errorf("send request trailer: %w", err) // wrap
	}
	err = client.CloseSend()
	if err != nil {
		return fmt.Errorf("close request stream: %w", err) // wrap
	}
	return nil
}

func handleProcessingError(ctx context.Context, log *zap.Logger, agentId int64, msg string, err error) { // nolint:unparam
	if grpctool.RequestCanceledOrTimedOut(err) {
		// An error caused by context signalling done
		return
	}
	var ue errz.UserError
	isUserError := errors.As(err, &ue)
	if isUserError {
		// TODO Don't log it, send it somewhere the user can see it https://gitlab.com/gitlab-org/gitlab/-/issues/277323
		// Log at Info for now.
		log.Info(msg, logz.Error(err))
	} else {
		// don't add logz.CorrelationIdFromContext(ctx) here as it's been added to the logger already
		log.Error(msg, logz.Error(err))
	}
}

func handleSendError(log *zap.Logger, msg string, err error) error {
	// The problem is almost certainly with the client's connection.
	// Still log it on Debug.
	log.Debug(msg, logz.Error(err))
	c := codes.Canceled
	if grpctool.RequestTimedOut(err) {
		c = codes.DeadlineExceeded
	}
	return status.Error(c, "gRPC send failed")
}

type cancelingReadCloser struct {
	io.ReadCloser
	cancel context.CancelFunc
}

func (c cancelingReadCloser) Close() error {
	// First close the pipe, then cancel the context.
	// Doing it the other way around can result in another goroutine calling CloseWithError() with the
	// cancellation error.
	err := c.ReadCloser.Close()
	c.cancel()
	return err
}

// valueOrError holds a value or an error. The first value/error is persisted, rest are discarded.
// onError is a callback that is called on each SetError() invocation, regardless of whether a value or an error
// has already been set.
// Thread safe.
type valueOrError struct {
	mx      sync.Mutex
	locker  *sync.Cond
	onError func(error) error
	value   interface{}
	err     error
	isSet   bool
}

func newValueOrError(onError func(error) error) *valueOrError {
	v := &valueOrError{
		onError: onError,
	}
	v.locker = sync.NewCond(&v.mx)
	return v
}

func (v *valueOrError) SetValue(value interface{}) {
	v.mx.Lock()
	defer v.mx.Unlock()
	if v.isSet {
		return
	}
	v.isSet = true
	v.value = value
	v.locker.Broadcast()
}

func (v *valueOrError) SetError(err error) {
	v.mx.Lock()
	defer v.mx.Unlock()
	err = v.onError(err)
	if v.isSet {
		return
	}
	v.isSet = true
	v.err = err
	v.locker.Broadcast()
}

// Wait returns a value or an error, blocking the caller until one of them is set.
func (v *valueOrError) Wait() (interface{}, error) {
	v.mx.Lock()
	defer v.mx.Unlock()
	if !v.isSet {
		v.locker.Wait()
	}
	return v.value, v.err
}
