package grpctool_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/errz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/testing/mock_rpc"
	"gitlab.com/gitlab-org/labkit/correlation"
	"go.uber.org/zap/zaptest"
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

func TestHandleSendError(t *testing.T) {
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
			in:       status.Error(codes.Internal, "bla: transport: the stream is done or WriteHeader was already called"),
			expected: status.Error(codes.Canceled, "msg: bla: transport: the stream is done or WriteHeader was already called"),
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
			actual := grpctool.HandleSendError(zaptest.NewLogger(t), "msg", tc.in)
			assert.Equal(t, tc.expected.Error(), actual.Error())
		})
	}
}

const metadataCorrelatorKey = "X-GitLab-Correlation-ID"

func TestMaybeWrapWithCorrelationId(t *testing.T) {
	t.Run("header error", func(t *testing.T) {
		stream := mock_rpc.NewMockClientStream(gomock.NewController(t))
		stream.EXPECT().Header().Return(nil, errors.New("header error"))
		err := errors.New("boom")
		wrappedErr := grpctool.MaybeWrapWithCorrelationId(err, stream)
		assert.Equal(t, err, wrappedErr)
	})
	t.Run("id present", func(t *testing.T) {
		id := correlation.SafeRandomID()
		stream := mock_rpc.NewMockClientStream(gomock.NewController(t))
		stream.EXPECT().Header().Return(metadata.Pairs(metadataCorrelatorKey, id), nil)
		err := errors.New("boom")
		wrappedErr := grpctool.MaybeWrapWithCorrelationId(err, stream)
		var errCorrelation errz.CorrelationError
		require.True(t, errors.As(wrappedErr, &errCorrelation))
		assert.Equal(t, id, errCorrelation.CorrelationId)
		assert.Equal(t, err, errCorrelation.Err)
	})
	t.Run("empty id", func(t *testing.T) {
		stream := mock_rpc.NewMockClientStream(gomock.NewController(t))
		stream.EXPECT().Header().Return(metadata.Pairs(metadataCorrelatorKey, ""), nil)
		err := errors.New("boom")
		wrappedErr := grpctool.MaybeWrapWithCorrelationId(err, stream)
		assert.Equal(t, err, wrappedErr)
	})
	t.Run("id missing", func(t *testing.T) {
		stream := mock_rpc.NewMockClientStream(gomock.NewController(t))
		stream.EXPECT().Header().Return(metadata.MD{}, nil)
		err := errors.New("boom")
		wrappedErr := grpctool.MaybeWrapWithCorrelationId(err, stream)
		assert.Equal(t, err, wrappedErr)
	})
}

func TestDeferMaybeWrapWithCorrelationId(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		stream := mock_rpc.NewMockClientStream(gomock.NewController(t))
		var wrappedErr error
		grpctool.DeferMaybeWrapWithCorrelationId(&wrappedErr, stream)
		assert.NoError(t, wrappedErr)
	})
	t.Run("header error", func(t *testing.T) {
		stream := mock_rpc.NewMockClientStream(gomock.NewController(t))
		stream.EXPECT().Header().Return(nil, errors.New("header error"))
		err := errors.New("boom")
		wrappedErr := err
		grpctool.DeferMaybeWrapWithCorrelationId(&wrappedErr, stream)
		assert.Equal(t, err, wrappedErr)
	})
	t.Run("id present", func(t *testing.T) {
		id := correlation.SafeRandomID()
		stream := mock_rpc.NewMockClientStream(gomock.NewController(t))
		stream.EXPECT().Header().Return(metadata.Pairs(metadataCorrelatorKey, id), nil)
		err := errors.New("boom")
		wrappedErr := err
		grpctool.DeferMaybeWrapWithCorrelationId(&wrappedErr, stream)
		var errCorrelation errz.CorrelationError
		require.True(t, errors.As(wrappedErr, &errCorrelation))
		assert.Equal(t, id, errCorrelation.CorrelationId)
		assert.Equal(t, err, errCorrelation.Err)
	})
	t.Run("empty id", func(t *testing.T) {
		stream := mock_rpc.NewMockClientStream(gomock.NewController(t))
		stream.EXPECT().Header().Return(metadata.Pairs(metadataCorrelatorKey, ""), nil)
		err := errors.New("boom")
		wrappedErr := err
		grpctool.DeferMaybeWrapWithCorrelationId(&wrappedErr, stream)
		assert.Equal(t, err, wrappedErr)
	})
	t.Run("id missing", func(t *testing.T) {
		stream := mock_rpc.NewMockClientStream(gomock.NewController(t))
		stream.EXPECT().Header().Return(metadata.MD{}, nil)
		err := errors.New("boom")
		wrappedErr := err
		grpctool.DeferMaybeWrapWithCorrelationId(&wrappedErr, stream)
		assert.Equal(t, err, wrappedErr)
	})
}
