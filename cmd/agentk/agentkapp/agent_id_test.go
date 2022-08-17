package agentkapp

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAgentId_GetReturnsAfterSet(t *testing.T) {
	h := newAgentIdHolder()
	require.NoError(t, h.set(1))
	id, err := h.get(context.Background())
	require.NoError(t, err)
	assert.EqualValues(t, 1, id)
}

func TestAgentId_GetReturnsAfterConcurrentSet(t *testing.T) {
	h := newAgentIdHolder()
	go func() {
		assert.NoError(t, h.set(1))
	}()
	id, err := h.get(context.Background())
	require.NoError(t, err)
	assert.EqualValues(t, 1, id)
}

func TestAgentId_GetTimesOut(t *testing.T) {
	h := newAgentIdHolder()
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	_, err := h.get(ctx)
	assert.Equal(t, context.DeadlineExceeded, err)
}

func TestAgentId_SetReturnsNoErrorOnSameId(t *testing.T) {
	h := newAgentIdHolder()
	require.NoError(t, h.set(1))
	assert.NoError(t, h.set(1))
}

func TestAgentId_SetReturnsErrorOnDifferentId(t *testing.T) {
	h := newAgentIdHolder()
	require.NoError(t, h.set(1))
	assert.EqualError(t, h.set(2), "agentId is already set to a different value: old 1, new 2")
}
