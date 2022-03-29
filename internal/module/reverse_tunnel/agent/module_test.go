package agent

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/reverse_tunnel/info"
	"google.golang.org/grpc"
)

var (
	_ modagent.Module     = &module{}
	_ modagent.Factory    = &Factory{}
	_ connectionInterface = &mockConnection{}
)

func TestModule(t *testing.T) {
	mockConn := &mockConnection{}
	m := setupModule(mockConn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := m.Run(ctx, nil)
	require.NoError(t, err)
	assert.EqualValues(t, 1, mockConn.runCalled)
}

func setupModule(mockConn *mockConnection) module {
	return module{
		server:         grpc.NewServer(),
		numConnections: 1,
		connectionFactory: func(descriptor *info.AgentDescriptor) connectionInterface {
			return mockConn
		},
	}
}

type mockConnection struct {
	runCalled int32
}

func (m *mockConnection) Run(ctx context.Context) {
	atomic.AddInt32(&m.runCalled, 1)
	<-ctx.Done()
}
