package server

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modserver"
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
	ctrl := gomock.NewController(t)
	rpcApi := mock_modserver.NewMockRpcApi(ctrl)
	publisher := NewMockPublisher(ctrl)
	ctx := modserver.InjectRpcApi(context.Background(), rpcApi)

	// setup mock expectations
	rpcApi.EXPECT().Log().Return(zap.NewNop())
	publisher.EXPECT().Publish(
		gomock.Any(),
		git_push_events_channel,
		matcher.ProtoEq(t, rpc.Project{Id: 42, PathWithNamespace: "foo/bar"}))

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
	ctrl := gomock.NewController(t)
	rpcApi := mock_modserver.NewMockRpcApi(ctrl)
	publisher := NewMockPublisher(ctrl)
	ctx := modserver.InjectRpcApi(context.Background(), rpcApi)

	// setup mock expectations
	rpcApi.EXPECT().Log().Return(zap.NewNop())
	rpcApi.EXPECT().HandleIoError(gomock.Any(), gomock.Any(), gomock.Any())
	publisher.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("some-error"))

	s := newServer(publisher)

	// WHEN
	_, err := s.GitPushEvent(ctx, &rpc.GitPushEventRequest{
		Project: &rpc.Project{Id: 42, FullPath: "foo/bar"},
	})

	// THEN
	assert.EqualError(t, err, "failed to handle git push event")
}
