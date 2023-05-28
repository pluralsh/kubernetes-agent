package grpctool_test

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/prototool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
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

func TestStatusErrorFromContext(t *testing.T) {
	canceled, cancel1 := context.WithCancel(context.Background())
	cancel1()
	expired, _ := context.WithDeadline(context.Background(), time.Now().Add(-time.Minute)) // nolint: govet
	tests := []struct {
		ctx         context.Context
		expectedMsg string
	}{
		{
			ctx:         context.Background(),
			expectedMsg: "rpc error: code = Unknown desc = 123: <nil>",
		},
		{
			ctx:         canceled,
			expectedMsg: "rpc error: code = Canceled desc = 123: context canceled",
		},
		{
			ctx:         expired,
			expectedMsg: "rpc error: code = DeadlineExceeded desc = 123: context deadline exceeded",
		},
	}
	for _, tc := range tests {
		t.Run(tc.expectedMsg, func(t *testing.T) {
			err := grpctool.StatusErrorFromContext(tc.ctx, "123")
			assert.EqualError(t, err, tc.expectedMsg)
		})
	}
}

func TestValuesMapToMeta(t *testing.T) {
	tests := []struct {
		input          map[string]*prototool.Values
		expectedOutput metadata.MD
	}{
		{
			input:          nil,
			expectedOutput: nil,
		},
		{
			input:          map[string]*prototool.Values{},
			expectedOutput: nil,
		},
		{
			input: map[string]*prototool.Values{
				"key1": {
					Value: []string{"s1", "s2"},
				},
				"key2": {
					Value: []string{},
				},
				"key3": {
					Value: []string{"s3", "s4"},
				},
			},
			expectedOutput: metadata.MD{
				"key1": []string{"s1", "s2"},
				"key2": []string{},
				"key3": []string{"s3", "s4"},
			},
		},
	}
	for i, tc := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			meta := grpctool.ValuesMapToMeta(tc.input)
			assert.Empty(t, cmp.Diff(tc.expectedOutput, meta))
		})
	}
}

func BenchmarkValuesMapToMeta(b *testing.B) {
	input := map[string]*prototool.Values{
		"key1": {
			Value: []string{"s1", "s2"},
		},
		"key2": {
			Value: []string{},
		},
		"key3": {
			Value: []string{"s3", "s4"},
		},
		"key4": {
			Value: []string{"s3", "s4"},
		},
		"key5": {
			Value: []string{"s3", "s4"},
		},
		"key6": {
			Value: []string{"s3", "s4"},
		},
	}
	var sink metadata.MD

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink = grpctool.ValuesMapToMeta(input)
	}
	_ = sink
}
