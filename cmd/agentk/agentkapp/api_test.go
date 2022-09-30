package agentkapp

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gitlab_access_rpc "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/gitlab_access/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/prototool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/matcher"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_gitlab_access"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/testhelpers"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

var (
	_ modagent.Api = (*agentAPI)(nil)
)

const (
	httpMethod      = http.MethodPost
	urlPath         = "/bla"
	moduleName      = "mod1"
	requestPayload  = "asdfndaskjfadsbfjsadhvfjhavfjasvf"
	responsePayload = "jknkjnjkasdnfkjasdnfkasdnfjnkjn"
	queryParamValue = "query-param-value with a space"
	queryParamName  = "q with a space"
)

func TestMakeGitLabRequest_HappyPath(t *testing.T) {
	api, clientStream := setupApiWithStream(t)
	// Send goroutine
	extra, err := anypb.New(&gitlab_access_rpc.HeaderExtra{
		ModuleName: moduleName,
	})
	require.NoError(t, err)
	gomock.InOrder(mockSendStream(t, clientStream,
		&grpctool.HttpRequest{
			Message: &grpctool.HttpRequest_Header_{
				Header: &grpctool.HttpRequest_Header{
					Request: &prototool.HttpRequest{
						Method: http.MethodPost,
						Header: map[string]*prototool.Values{
							"Req-Header": {
								Value: []string{"x1", "x2"},
							},
							"Content-Type": {
								Value: []string{"text/plain"},
							},
						},
						UrlPath: urlPath,
						Query: map[string]*prototool.Values{
							queryParamName: {
								Value: []string{queryParamValue},
							},
						},
					},
					Extra: extra,
				},
			},
		},
		&grpctool.HttpRequest{
			Message: &grpctool.HttpRequest_Data_{
				Data: &grpctool.HttpRequest_Data{
					Data: []byte(requestPayload),
				},
			},
		},
		&grpctool.HttpRequest{
			Message: &grpctool.HttpRequest_Trailer_{
				Trailer: &grpctool.HttpRequest_Trailer{},
			},
		},
	)...)
	gomock.InOrder(mockRecvStream(clientStream,
		&grpctool.HttpResponse{
			Message: &grpctool.HttpResponse_Header_{
				Header: &grpctool.HttpResponse_Header{
					Response: &prototool.HttpResponse{
						StatusCode: http.StatusOK,
						Status:     http.StatusText(http.StatusOK),
						Header: map[string]*prototool.Values{
							"Resp-Header": {
								Value: []string{"a1", "a2"},
							},
						},
					},
				},
			},
		},
		&grpctool.HttpResponse{
			Message: &grpctool.HttpResponse_Data_{
				Data: &grpctool.HttpResponse_Data{
					Data: []byte(responsePayload),
				},
			},
		},
		&grpctool.HttpResponse{
			Message: &grpctool.HttpResponse_Trailer_{
				Trailer: &grpctool.HttpResponse_Trailer{},
			},
		})...)
	resp, err := api.MakeGitLabRequest(context.Background(), urlPath,
		modagent.WithRequestMethod(httpMethod),
		modagent.WithRequestQueryParam(queryParamName, queryParamValue),
		modagent.WithRequestHeader("Req-Header", "x1", "x2"),
		modagent.WithRequestBody(strings.NewReader(requestPayload), "text/plain"),
	)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, resp.Body.Close())
	}()
	assert.EqualValues(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, responsePayload, string(readAll(t, resp.Body)))
	assert.Empty(t, cmp.Diff(map[string][]string{
		"Resp-Header": {"a1", "a2"},
	}, (map[string][]string)(resp.Header)))
}

func TestMakeGitLabRequest_MakeRequestErrorClosesBody(t *testing.T) {
	api, client, _ := setupApi(t)
	body := newFailingReaderCloser()
	client.EXPECT().
		MakeRequest(gomock.Any()).
		Return(nil, errors.New("expected error"))
	_, err := api.MakeGitLabRequest(context.Background(), urlPath, modagent.WithRequestBody(body, "text/plain"))
	assert.EqualError(t, err, "expected error")
	assert.True(t, body.CloseCalled())
	assert.False(t, body.ReadCalled())
}

func TestMakeGitLabRequest_SendError(t *testing.T) {
	api, client, clientStream := setupApi(t)
	body := newFailingReaderCloser()
	var clientCtx context.Context
	client.EXPECT().
		MakeRequest(gomock.Any()).
		DoAndReturn(func(ctx context.Context, opts ...grpc.CallOption) (gitlab_access_rpc.GitlabAccess_MakeRequestClient, error) {
			clientCtx = ctx
			return clientStream, nil
		})
	clientStream.EXPECT().
		Send(gomock.Any()).
		Return(errors.New("expected error"))
	clientStream.EXPECT().
		RecvMsg(gomock.Any()).
		DoAndReturn(func(m interface{}) error {
			<-clientCtx.Done() // Blocks until context is canceled because of the send error.
			// Also return an error - this one must be ignored, the one from Send() should be used.
			return clientCtx.Err()
		})
	_, err := api.MakeGitLabRequest(context.Background(), urlPath, modagent.WithRequestBody(body, "text/plain"))
	assert.EqualError(t, err, "send request header: expected error")
	assert.True(t, body.CloseCalled())
	assert.False(t, body.ReadCalled())
	assert.EqualError(t, err, "send request header: expected error")
}

func TestMakeGitLabRequest_RecvError(t *testing.T) {
	api, client, clientStream := setupApi(t)
	body := newFailingReaderCloser()
	var clientCtx context.Context
	client.EXPECT().
		MakeRequest(gomock.Any()).
		DoAndReturn(func(ctx context.Context, opts ...grpc.CallOption) (gitlab_access_rpc.GitlabAccess_MakeRequestClient, error) {
			clientCtx = ctx
			return clientStream, nil
		})
	clientStream.EXPECT().
		Send(gomock.Any()).
		DoAndReturn(func(m interface{}) error {
			<-clientCtx.Done() // Blocks until context is canceled because of the send error.
			// Also return an error - this one must be ignored, the one from RecvMsg() should be used.
			return clientCtx.Err()
		})
	clientStream.EXPECT().
		RecvMsg(gomock.Any()).
		Return(errors.New("expected error"))
	_, err := api.MakeGitLabRequest(context.Background(), urlPath, modagent.WithRequestBody(body, "text/plain"))
	assert.EqualError(t, err, "expected error")
	assert.True(t, body.CloseCalled())
	assert.False(t, body.ReadCalled())
	assert.EqualError(t, err, "expected error")
}

func TestMakeGitLabRequest_LateRecvError(t *testing.T) {
	api, client, clientStream := setupApi(t)
	body := newFailingReaderCloser()
	var clientCtx context.Context
	client.EXPECT().
		MakeRequest(gomock.Any()).
		DoAndReturn(func(ctx context.Context, opts ...grpc.CallOption) (gitlab_access_rpc.GitlabAccess_MakeRequestClient, error) {
			clientCtx = ctx
			return clientStream, nil
		})
	clientStream.EXPECT().
		Send(gomock.Any()).
		DoAndReturn(func(m interface{}) error {
			<-clientCtx.Done() // Blocks until context is canceled because of the send error.
			// Also return an error - this one must be ignored, the one from RecvMsg() should be used.
			return clientCtx.Err()
		})
	gomock.InOrder(
		clientStream.EXPECT().
			RecvMsg(gomock.Any()).
			Do(testhelpers.RecvMsg(&grpctool.HttpResponse{
				Message: &grpctool.HttpResponse_Header_{
					Header: &grpctool.HttpResponse_Header{
						Response: &prototool.HttpResponse{
							StatusCode: http.StatusOK,
							Status:     http.StatusText(http.StatusOK),
						},
					},
				},
			})),
		clientStream.EXPECT().
			RecvMsg(gomock.Any()).
			Return(errors.New("expected error")),
	)
	resp, err := api.MakeGitLabRequest(context.Background(), urlPath, modagent.WithRequestBody(body, "text/plain"))
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, resp.Body.Close())
	}()
	assert.EqualValues(t, http.StatusOK, resp.StatusCode)
	_, err = io.ReadAll(resp.Body)
	assert.EqualError(t, err, "expected error")
	<-body.closeCalled // wait for async close
	assert.False(t, body.ReadCalled())
	assert.EqualError(t, err, "expected error")
}

func setupApiWithStream(t *testing.T) (*agentAPI, *mock_gitlab_access.MockGitlabAccess_MakeRequestClient) {
	api, client, clientStream := setupApi(t)
	client.EXPECT().
		MakeRequest(gomock.Any()).
		Return(clientStream, nil)
	return api, clientStream
}

func setupApi(t *testing.T) (*agentAPI, *mock_gitlab_access.MockGitlabAccessClient, *mock_gitlab_access.MockGitlabAccess_MakeRequestClient) {
	ctrl := gomock.NewController(t)
	client := mock_gitlab_access.NewMockGitlabAccessClient(ctrl)
	clientStream := mock_gitlab_access.NewMockGitlabAccess_MakeRequestClient(ctrl)
	return &agentAPI{
		moduleName: moduleName,
		agentId:    NewAgentIdHolder(),
		client:     client,
	}, client, clientStream
}

func mockRecvStream(server *mock_gitlab_access.MockGitlabAccess_MakeRequestClient, msgs ...proto.Message) []*gomock.Call {
	res := make([]*gomock.Call, 0, len(msgs)+1)
	for _, msg := range msgs {
		call := server.EXPECT().
			RecvMsg(gomock.Any()).
			Do(testhelpers.RecvMsg(msg))
		res = append(res, call)
	}
	call := server.EXPECT().
		RecvMsg(gomock.Any()).
		Return(io.EOF)
	res = append(res, call)
	return res
}

func mockSendStream(t *testing.T, client *mock_gitlab_access.MockGitlabAccess_MakeRequestClient, msgs ...*grpctool.HttpRequest) []*gomock.Call {
	res := make([]*gomock.Call, 0, len(msgs)+1)
	for _, msg := range msgs {
		call := client.EXPECT().
			Send(matcher.ProtoEq(t, msg))
		res = append(res, call)
	}
	streamDone := make(chan struct{})
	res = append(res, client.EXPECT().
		CloseSend().
		Do(func() {
			close(streamDone)
		}))
	t.Cleanup(func() {
		// The sending is done concurrently and test can finish earlier than the sending goroutine is done sending.
		// In that case there will be a missing expected invocation. Wait for it to finish before proceeding.
		// t.Cleanup() processes added functions in LIFO order, so this one should be executed before the validation
		// function (added by gomock.NewController()).
		<-streamDone
	})
	return res
}

func readAll(t *testing.T, r io.Reader) []byte {
	data, err := io.ReadAll(r)
	require.NoError(t, err)
	return data
}

type failingReaderCloser struct {
	readCalled  chan struct{}
	closeCalled chan struct{}
	readOnce    sync.Once
	closeOnce   sync.Once
}

func newFailingReaderCloser() *failingReaderCloser {
	return &failingReaderCloser{
		readCalled:  make(chan struct{}),
		closeCalled: make(chan struct{}),
	}
}

func (c *failingReaderCloser) Read(p []byte) (n int, err error) {
	c.readOnce.Do(func() {
		close(c.readCalled)
	})
	return 0, errors.New("expected read error")
}

func (c *failingReaderCloser) Close() error {
	c.closeOnce.Do(func() {
		close(c.closeCalled)
	})
	return errors.New("expected close error")
}

func (c *failingReaderCloser) ReadCalled() bool {
	select {
	case <-c.readCalled:
		return true
	default:
		return false
	}
}

func (c *failingReaderCloser) CloseCalled() bool {
	select {
	case <-c.closeCalled:
		return true
	default:
		return false
	}
}
