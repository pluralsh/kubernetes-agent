package grpctool_test

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/grpctool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestRequestCanceledOrTimedOut(t *testing.T) {
	t.Run("context errors", func(t *testing.T) {
		assert.True(t, grpctool.RequestCanceledOrTimedOut(context.Canceled))
		assert.True(t, grpctool.RequestCanceledOrTimedOut(context.DeadlineExceeded))
		assert.False(t, grpctool.RequestCanceledOrTimedOut(io.EOF))
	})
	t.Run("wrapped context errors", func(t *testing.T) {
		assert.True(t, grpctool.RequestCanceledOrTimedOut(fmt.Errorf("bla: %w", context.Canceled)))
		assert.True(t, grpctool.RequestCanceledOrTimedOut(fmt.Errorf("bla: %w", context.DeadlineExceeded)))
		assert.False(t, grpctool.RequestCanceledOrTimedOut(fmt.Errorf("bla: %w", io.EOF)))
	})
	t.Run("gRPC errors", func(t *testing.T) {
		assert.True(t, grpctool.RequestCanceledOrTimedOut(status.Error(codes.Canceled, "bla")))
		assert.True(t, grpctool.RequestCanceledOrTimedOut(status.Error(codes.DeadlineExceeded, "bla")))
		assert.False(t, grpctool.RequestCanceledOrTimedOut(status.Error(codes.Unavailable, "bla")))
	})
	t.Run("wrapped gRPC errors", func(t *testing.T) {
		assert.True(t, grpctool.RequestCanceledOrTimedOut(fmt.Errorf("bla: %w", status.Error(codes.Canceled, "bla"))))
		assert.True(t, grpctool.RequestCanceledOrTimedOut(fmt.Errorf("bla: %w", status.Error(codes.DeadlineExceeded, "bla"))))
		assert.False(t, grpctool.RequestCanceledOrTimedOut(fmt.Errorf("bla: %w", status.Error(codes.Unavailable, "bla"))))
	})
}

func TestRequestCanceled(t *testing.T) {
	t.Run("context errors", func(t *testing.T) {
		assert.True(t, grpctool.RequestCanceled(context.Canceled))
		assert.False(t, grpctool.RequestCanceled(context.DeadlineExceeded))
		assert.False(t, grpctool.RequestCanceled(io.EOF))
	})
	t.Run("wrapped context errors", func(t *testing.T) {
		assert.True(t, grpctool.RequestCanceled(fmt.Errorf("bla: %w", context.Canceled)))
		assert.False(t, grpctool.RequestCanceled(fmt.Errorf("bla: %w", context.DeadlineExceeded)))
		assert.False(t, grpctool.RequestCanceled(fmt.Errorf("bla: %w", io.EOF)))
	})
	t.Run("gRPC errors", func(t *testing.T) {
		assert.True(t, grpctool.RequestCanceled(status.Error(codes.Canceled, "bla")))
		assert.False(t, grpctool.RequestCanceled(status.Error(codes.DeadlineExceeded, "bla")))
		assert.False(t, grpctool.RequestCanceled(status.Error(codes.Unavailable, "bla")))
	})
	t.Run("wrapped gRPC errors", func(t *testing.T) {
		assert.True(t, grpctool.RequestCanceled(fmt.Errorf("bla: %w", status.Error(codes.Canceled, "bla"))))
		assert.False(t, grpctool.RequestCanceled(fmt.Errorf("bla: %w", status.Error(codes.DeadlineExceeded, "bla"))))
		assert.False(t, grpctool.RequestCanceled(fmt.Errorf("bla: %w", status.Error(codes.Unavailable, "bla"))))
	})
}

func TestRequestTimedOut(t *testing.T) {
	t.Run("context errors", func(t *testing.T) {
		assert.False(t, grpctool.RequestTimedOut(context.Canceled))
		assert.True(t, grpctool.RequestTimedOut(context.DeadlineExceeded))
		assert.False(t, grpctool.RequestTimedOut(io.EOF))
	})
	t.Run("wrapped context errors", func(t *testing.T) {
		assert.False(t, grpctool.RequestTimedOut(fmt.Errorf("bla: %w", context.Canceled)))
		assert.True(t, grpctool.RequestTimedOut(fmt.Errorf("bla: %w", context.DeadlineExceeded)))
		assert.False(t, grpctool.RequestTimedOut(fmt.Errorf("bla: %w", io.EOF)))
	})
	t.Run("gRPC errors", func(t *testing.T) {
		assert.False(t, grpctool.RequestTimedOut(status.Error(codes.Canceled, "bla")))
		assert.True(t, grpctool.RequestTimedOut(status.Error(codes.DeadlineExceeded, "bla")))
		assert.False(t, grpctool.RequestTimedOut(status.Error(codes.Unavailable, "bla")))
	})
	t.Run("wrapped gRPC errors", func(t *testing.T) {
		assert.False(t, grpctool.RequestTimedOut(fmt.Errorf("bla: %w", status.Error(codes.Canceled, "bla"))))
		assert.True(t, grpctool.RequestTimedOut(fmt.Errorf("bla: %w", status.Error(codes.DeadlineExceeded, "bla"))))
		assert.False(t, grpctool.RequestTimedOut(fmt.Errorf("bla: %w", status.Error(codes.Unavailable, "bla"))))
	})
}

func TestHandleIoError(t *testing.T) {
	tests := []struct {
		in       error
		expected error
	}{
		{
			in:       status.Error(codes.Canceled, "bla"),
			expected: status.Error(codes.Canceled, "msg: bla"),
		},
		{
			in:       status.Error(codes.DeadlineExceeded, "bla"),
			expected: status.Error(codes.DeadlineExceeded, "msg: bla"),
		},
		{
			in:       status.Error(codes.Internal, "bla"),
			expected: status.Error(codes.Internal, "msg: bla"),
		},
		{
			in:       status.Error(codes.Internal, "bla"),
			expected: status.Error(codes.Internal, "msg: bla"),
		},
		{
			in:       io.EOF,
			expected: status.Error(codes.Canceled, "msg: EOF"),
		},
		{
			in:       io.ErrUnexpectedEOF,
			expected: status.Error(codes.Canceled, "msg: unexpected EOF"),
		},
	}
	for i, tc := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := grpctool.HandleIoError("msg", tc.in)
			assert.Equal(t, tc.expected.Error(), actual.Error())
		})
	}
}
