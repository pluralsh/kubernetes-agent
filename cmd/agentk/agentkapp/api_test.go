package agentkapp

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/matcher"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_gitlab_access"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/testhelpers"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"
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
		moduleName:        moduleName,
		agentId:           NewValueHolder[int64](),
		gitLabExternalUrl: NewValueHolder[url.URL](),
	}, client, clientStream
}

func mockRecvStream(server *mock_gitlab_access.MockGitlabAccess_MakeRequestClient, msgs ...proto.Message) []any {
	res := make([]any, 0, len(msgs)+1)
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

func mockSendStream(t *testing.T, client *mock_gitlab_access.MockGitlabAccess_MakeRequestClient, msgs ...*grpctool.HttpRequest) []any {
	res := make([]any, 0, len(msgs)+1)
	for _, msg := range msgs {
		call := client.EXPECT().
			Send(matcher.ProtoEq(t, msg))
		res = append(res, call)
	}
	streamDone := make(chan struct{})
	res = append(res, client.EXPECT().
		CloseSend().
		Do(func() error {
			close(streamDone)
			return nil
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
