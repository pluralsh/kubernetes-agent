package kasapp

import (
	"context"
	"errors"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/gitlab"
	gapi "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/gitlab/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/cache"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/logz"
	"gitlab.com/gitlab-org/labkit/errortracking"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	_ modserver.RpcApi = (*serverRpcApi)(nil)
)

type serverRpcApi struct {
	modshared.RpcApiStub
	ErrorTracker   errortracking.Tracker
	GitLabClient   gitlab.ClientInterface
	AgentInfoCache *cache.CacheWithErr
}

func (a *serverRpcApi) AgentToken() api.AgentToken {
	return api.AgentTokenFromContext(a.StreamCtx)
}

func (a *serverRpcApi) AgentInfo(ctx context.Context, log *zap.Logger) (*api.AgentInfo, error) {
	agentInfo, err := a.getAgentInfoCached(ctx)
	switch {
	case err == nil:
		return agentInfo, nil
	case errors.Is(err, context.Canceled):
		err = status.Error(codes.Canceled, err.Error())
	case errors.Is(err, context.DeadlineExceeded):
		err = status.Error(codes.DeadlineExceeded, err.Error())
	case gitlab.IsForbidden(err):
		err = status.Error(codes.PermissionDenied, "forbidden")
	case gitlab.IsUnauthorized(err):
		err = status.Error(codes.Unauthenticated, "unauthenticated")
	case gitlab.IsNotFound(err):
		err = status.Error(codes.NotFound, "agent not found")
	default:
		logAndCapture(ctx, a.ErrorTracker, log, modshared.NoAgentId, "AgentInfo()", err)
		err = status.Error(codes.Unavailable, "unavailable")
	}
	return nil, err
}

func (a *serverRpcApi) HandleProcessingError(log *zap.Logger, agentId int64, msg string, err error) {
	handleProcessingError(a.StreamCtx, a.ErrorTracker, log, agentId, msg, err)
}

func (a *serverRpcApi) HandleSendError(log *zap.Logger, msg string, err error) error {
	// The problem is almost certainly with the client's connection.
	// Still log it on Debug.
	if !grpctool.RequestCanceled(err) {
		log.Debug(msg, logz.Error(err))
	}
	return status.Error(codes.Unavailable, "gRPC send failed")
}

func (a *serverRpcApi) getAgentInfoCached(ctx context.Context) (*api.AgentInfo, error) {
	agentToken := a.AgentToken()
	agentInfo, err := a.AgentInfoCache.GetItem(ctx, agentToken, func() (interface{}, error) {
		return gapi.GetAgentInfo(ctx, a.GitLabClient, agentToken)
	})
	if err != nil {
		return nil, err
	}
	return agentInfo.(*api.AgentInfo), nil
}
