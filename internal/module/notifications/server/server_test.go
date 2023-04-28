package server

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modserver/notifications"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/notifications/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/matcher"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_modserver"
	"go.uber.org/zap"
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
	publisher := NewMockPublisher(ctrl)
	publisher.EXPECT().Publish(
		gomock.Any(),
		notifications.GitPushEventsChannel,
		matcher.ProtoEq(t, notifications.Project{Id: 42, FullPath: "foo/bar"}))

	// setup server under test
	s := newServer(publisher)

	// WHEN
	_, err := s.GitPushEvent(ctx, &rpc.GitPushEventRequest{
		Project: &rpc.Project{Id: 42, FullPath: "foo/bar"},
	})

	// THEN
	require.NoError(t, err)
}

func TestServer_GitPushEvent_FailedPublish(t *testing.T) {
	// GIVEN
	// setup test fixtures
	ctrl := gomock.NewController(t)
	rpcApi := mock_modserver.NewMockRpcApi(ctrl)
	publisher := NewMockPublisher(ctrl)
	ctx := modserver.InjectRpcApi(context.Background(), rpcApi)
	rpcApi.EXPECT().Log().Return(zap.NewNop())
	rpcApi.EXPECT().HandleIoError(gomock.Any(), gomock.Any(), gomock.Any())
	publisher.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("some-error"))

	// setup server under test
	s := newServer(publisher)

	// WHEN
	_, err := s.GitPushEvent(ctx, &rpc.GitPushEventRequest{
		Project: &rpc.Project{Id: 42, FullPath: "foo/bar"},
	})

	// THEN
	assert.EqualError(t, err, "failed to handle git push event")
}
