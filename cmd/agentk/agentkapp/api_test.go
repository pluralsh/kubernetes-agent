package agentkapp

import (
	"errors"
	"io"
	"net/http"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pluralsh/kuberentes-agent/internal/module/modagent"
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
