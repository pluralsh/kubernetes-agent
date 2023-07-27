package grpctool_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_rpc"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	testcases = []struct {
		name             string
		handlerError     error
		expectReportCall bool
	}{
		{
			name:             "unknown error",
			handlerError:     status.Error(codes.Unknown, "some unknown error"),
			expectReportCall: true,
		},
		{
			name:             "canceled error",
			handlerError:     status.Error(codes.Canceled, "some canceled error"),
			expectReportCall: false,
		},
		{
			name:             "no error",
			handlerError:     nil,
			expectReportCall: false,
		},
	}
)

func TestServerErrorReporter_UnaryInterceptor(t *testing.T) {
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockServerErrorReporter := mock_rpc.NewMockServerErrorReporter(ctrl)

			if tc.expectReportCall {
				mockServerErrorReporter.EXPECT().Report(gomock.Any(), "some-method", tc.handlerError)
			}

			usHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
				return struct{}{}, tc.handlerError
			}

			usi := grpctool.UnaryServerErrorReporterInterceptor(mockServerErrorReporter)
			_, err := usi(context.Background(), struct{}{}, &grpc.UnaryServerInfo{FullMethod: "some-method"}, usHandler)

			require.ErrorIs(t, err, tc.handlerError)
		})
	}
}

func TestServerErrorReporter_StreamInterceptor(t *testing.T) {
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockServerErrorReporter := mock_rpc.NewMockServerErrorReporter(ctrl)
			mockServerStream := mock_rpc.NewMockServerStream(ctrl)

			if tc.expectReportCall {
				mockServerErrorReporter.EXPECT().Report(gomock.Any(), "some-method", tc.handlerError)
				mockServerStream.EXPECT().Context().Times(1)
			}

			ssHandler := func(interface{}, grpc.ServerStream) error {
				return tc.handlerError
			}

			ssi := grpctool.StreamServerErrorReporterInterceptor(mockServerErrorReporter)
			err := ssi(struct{}{}, mockServerStream, &grpc.StreamServerInfo{FullMethod: "some-method"}, ssHandler)

			require.ErrorIs(t, err, tc.handlerError)
		})
	}
}
