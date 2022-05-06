package server

import (
	"context"
	"io"
	"net/http"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/gitlab"
	gapi "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/gitlab/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/gitlab_access/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/grpctool"
)

type server struct {
	rpc.UnimplementedGitlabAccessServer
	gitLabClient gitlab.ClientInterface
}

func newServer(gitLabClient gitlab.ClientInterface) *server {
	return &server{
		gitLabClient: gitLabClient,
	}
}

func (s *server) MakeRequest(server rpc.GitlabAccess_MakeRequestServer) error {
	rpcApi := modserver.AgentRpcApiFromContext(server.Context())
	log := rpcApi.Log()
	grpc2http := grpctool.InboundGrpcToOutboundHttp{
		Log: log,
		HandleProcessingError: func(msg string, err error) {
			rpcApi.HandleProcessingError(log, modshared.NoAgentId, msg, err)
		},
		HandleSendError: func(msg string, err error) error {
			return rpcApi.HandleSendError(log, msg, err)
		},
		HttpDo: s.httpDo,
	}
	return grpc2http.Pipe(server)
}

func (s *server) httpDo(ctx context.Context, header *grpctool.HttpRequest_Header, body io.Reader) (*http.Response, error) {
	var extra rpc.HeaderExtra
	err := header.Extra.UnmarshalTo(&extra)
	if err != nil {
		return nil, err
	}
	return gapi.MakeModuleRequest(
		ctx,
		s.gitLabClient,
		modserver.AgentRpcApiFromContext(ctx).AgentToken(),
		extra.ModuleName,
		header.Request.Method,
		header.Request.UrlPath,
		header.Request.UrlQuery(),
		header.Request.HttpHeader(),
		body,
	)
}
