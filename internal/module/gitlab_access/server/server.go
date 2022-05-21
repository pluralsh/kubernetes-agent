package server

import (
	"context"
	"errors"
	"io"

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
		HandleIoError: func(msg string, err error) error {
			return rpcApi.HandleIoError(log, msg, err)
		},
		HttpDo: s.httpDo,
	}
	return grpc2http.Pipe(server)
}

func (s *server) httpDo(ctx context.Context, header *grpctool.HttpRequest_Header, body io.Reader) (grpctool.DoResponse, error) {
	if header.Request.IsUpgrade() {
		return grpctool.DoResponse{}, errors.New("connection upgrade is not supported")
	}
	var extra rpc.HeaderExtra
	err := header.Extra.UnmarshalTo(&extra)
	if err != nil {
		return grpctool.DoResponse{}, err
	}
	resp, err := gapi.MakeModuleRequest( // nolint: bodyclose
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
	if err != nil {
		return grpctool.DoResponse{}, err
	}
	return grpctool.DoResponse{
		Resp: resp,
	}, nil
}
