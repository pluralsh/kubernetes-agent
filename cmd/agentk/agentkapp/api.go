package agentkapp

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	gitlab_access_rpc "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/gitlab_access/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/errz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/memz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/prototool"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/anypb"
	"k8s.io/apimachinery/pkg/util/wait"
)

// agentAPI is an implementation of modagent.API.
type agentAPI struct {
	moduleName string
	agentId    *AgentIdHolder
	client     gitlab_access_rpc.GitlabAccessClient
}

func (a *agentAPI) GetAgentId(ctx context.Context) (int64, error) {
	return a.agentId.get(ctx)
}

func (a *agentAPI) TryGetAgentId() (int64, bool) {
	return a.agentId.tryGet()
}

func (a *agentAPI) HandleProcessingError(ctx context.Context, log *zap.Logger, agentId int64, msg string, err error) {
	handleProcessingError(ctx, log, agentId, msg, err)
}

func (a *agentAPI) MakeGitLabRequest(ctx context.Context, path string, opts ...modagent.GitLabRequestOption) (*modagent.GitLabResponse, error) {
	config, err := modagent.ApplyRequestOptions(opts)
	if err != nil {
		return nil, err
	}
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
	val := newValueOrError[*modagent.GitLabResponse](func(err error) error {
		cancel()                   // 1. Cancel the other goroutine and the client.
		_ = pw.CloseWithError(err) // 2. Close the "write side" of the pipe
		return err
	})

	var wg wait.Group
	// Write request
	wg.Start(func() {
		writeErr := a.makeRequest(client, path, config)
		if writeErr != nil {
			val.SetError(writeErr)
		}
	})
	// Read response
	wg.Start(func() {
		readErr := grpctool.HttpResponseStreamVisitor.Get().Visit(client,
			grpctool.WithCallback(grpctool.HttpResponseHeaderFieldNumber, func(header *grpctool.HttpResponse_Header) error {
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
			grpctool.WithCallback(grpctool.HttpResponseDataFieldNumber, func(data *grpctool.HttpResponse_Data) error {
				_, pwErr := pw.Write(data.Data)
				return pwErr
			}),
			grpctool.WithCallback(grpctool.HttpResponseTrailerFieldNumber, func(trailer *grpctool.HttpResponse_Trailer) error {
				return nil
			}),
			grpctool.WithEOFCallback(pw.Close),
			grpctool.WithNotExpectingToGet(codes.Internal, grpctool.HttpResponseUpgradeDataFieldNumber),
		)
		if readErr != nil {
			val.SetError(readErr)
		}
	})
	resp, err := val.Wait()
	if err != nil {
		wg.Wait() // Wait for both goroutines to finish before returning
		return nil, err
	}
	return resp, nil
}

func (a *agentAPI) makeRequest(client gitlab_access_rpc.GitlabAccess_MakeRequestClient, path string, config *modagent.GitLabRequestConfig) (retErr error) {
	var body io.ReadCloser
	if config.Body != nil {
		body = &onceReadCloser{
			ReadCloser: config.Body,
		}
		defer errz.SafeClose(body, &retErr)
	}
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
		if err == io.EOF { // nolint:errorlint
			return nil // the other goroutine will receive the error in RecvMsg()
		}
		return fmt.Errorf("send request header: %w", err) // wrap
	}
	if body != nil {
		err = a.sendRequestBody(client, body)
		if err != nil {
			return err
		}
	}
	err = client.Send(&grpctool.HttpRequest{
		Message: &grpctool.HttpRequest_Trailer_{
			Trailer: &grpctool.HttpRequest_Trailer{},
		},
	})
	if err != nil {
		if err == io.EOF { // nolint:errorlint
			return nil // the other goroutine will receive the error in RecvMsg()
		}
		return fmt.Errorf("send request trailer: %w", err) // wrap
	}
	err = client.CloseSend()
	if err != nil {
		return fmt.Errorf("close request stream: %w", err) // wrap
	}
	return nil
}

func (a *agentAPI) sendRequestBody(client gitlab_access_rpc.GitlabAccess_MakeRequestClient, body io.ReadCloser) (retErr error) {
	defer errz.SafeClose(body, &retErr) // close ASAP
	buffer := memz.Get32k()
	defer memz.Put32k(buffer)
	for {
		n, readErr := body.Read(buffer)
		if n > 0 { // handle n>0 before readErr != nil to ensure any consumed data gets forwarded
			sendErr := client.Send(&grpctool.HttpRequest{
				Message: &grpctool.HttpRequest_Data_{
					Data: &grpctool.HttpRequest_Data{
						Data: buffer[:n],
					}},
			})
			if sendErr != nil {
				if sendErr == io.EOF { // nolint:errorlint
					// the other goroutine will receive the error in RecvMsg()
					return nil
				}
				return fmt.Errorf("send request data: %w", sendErr) // wrap
			}
		}
		if readErr != nil {
			if readErr == io.EOF { // nolint:errorlint
				break
			}
			return fmt.Errorf("read request body: %w", readErr) // wrap
		}
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
		// don't add logz.TraceIdFromContext(ctx) here as it's been added to the logger already
		log.Error(msg, logz.Error(err))
	}
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
type valueOrError[T any] struct {
	mu      sync.Mutex
	locker  *sync.Cond
	onError func(error) error
	value   T
	err     error
	isSet   bool
}

func newValueOrError[T any](onError func(error) error) *valueOrError[T] {
	v := &valueOrError[T]{
		onError: onError,
	}
	v.locker = sync.NewCond(&v.mu)
	return v
}

func (v *valueOrError[T]) SetValue(value T) {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.isSet {
		return
	}
	v.isSet = true
	v.value = value
	v.locker.Broadcast()
}

func (v *valueOrError[T]) SetError(err error) {
	v.mu.Lock()
	defer v.mu.Unlock()
	err = v.onError(err)
	if v.isSet {
		return
	}
	v.isSet = true
	v.err = err
	v.locker.Broadcast()
}

// Wait returns a value or an error, blocking the caller until one of them is set.
func (v *valueOrError[T]) Wait() (T, error) {
	v.mu.Lock()
	defer v.mu.Unlock()
	if !v.isSet {
		v.locker.Wait()
	}
	return v.value, v.err
}

type onceReadCloser struct {
	io.ReadCloser
	once     sync.Once
	closeErr error
}

func (oc *onceReadCloser) Close() error {
	oc.once.Do(oc.close)
	return oc.closeErr
}

func (oc *onceReadCloser) close() {
	oc.closeErr = oc.ReadCloser.Close()
}
