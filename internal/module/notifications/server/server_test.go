package server

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/notifications/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/event"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
)

var (
	_ modserver.Module        = &module{}
	_ modserver.Factory       = &Factory{}
	_ rpc.NotificationsServer = &server{}
)

func TestServer_GitPushEvent_SuccessfulPublish(t *testing.T) {
	// GIVEN
	// setup test fixtures
	ctrl := gomock.NewController(t)
	rpcApi := mock_modserver.NewMockRpcApi(ctrl)
	ctx := modserver.InjectRpcApi(context.Background(), rpcApi)

	var gotEvent proto.Message
	// setup server under test
	s := newServer(func(ctx context.Context, e proto.Message) error {
		gotEvent = e
		return nil
	})

	// WHEN
	gitPushEvent := &event.GitPushEvent{
		Project: &event.Project{Id: 42, FullPath: "foo/bar"},
	}
	_, err := s.GitPushEvent(ctx, &rpc.GitPushEventRequest{
		Event: gitPushEvent,
	})

	// THEN
	require.NoError(t, err)
	assert.Empty(t, cmp.Diff(gotEvent, gitPushEvent, protocmp.Transform()))
}

func TestServer_GitPushEvent_FailedPublish(t *testing.T) {
	// GIVEN
	// setup test fixtures
	ctrl := gomock.NewController(t)
	rpcApi := mock_modserver.NewMockRpcApi(ctrl)
	ctx := modserver.InjectRpcApi(context.Background(), rpcApi)
	rpcApi.EXPECT().Log().Return(zap.NewNop())
	givenErr := errors.New("some error")
	rpcApi.EXPECT().
		HandleProcessingError(gomock.Any(), modshared.NoAgentId, gomock.Any(), givenErr)

	// setup server under test
	s := newServer(func(ctx context.Context, e proto.Message) error {
		return givenErr
	})

	// WHEN
	_, err := s.GitPushEvent(ctx, &rpc.GitPushEventRequest{
		Event: &event.GitPushEvent{
			Project: &event.Project{Id: 42, FullPath: "foo/bar"},
		},
	})

	// THEN
	assert.EqualError(t, err, "rpc error: code = Unavailable desc = Failed to publish received git push event: some error")
}
