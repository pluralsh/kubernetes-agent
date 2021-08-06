package server

import (
	"context"
	"io"
	"net/http"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/gitlab"
	gapi "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/gitlab/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/gitlab_access/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/grpctool"
)

type server struct {
	rpc.UnimplementedGitlabAccessServer
	pipe *grpctool.InboundGrpcToOutboundHttp
}

func newServer(gitLabClient gitlab.ClientInterface) *server {
	return &server{
		pipe: grpctool.NewInboundGrpcToOutboundHttp(
			func(ctx context.Context, header *grpctool.HttpRequest_Header, body io.Reader) (*http.Response, error) {
				var extra rpc.HeaderExtra
				err := header.Extra.UnmarshalTo(&extra)
				if err != nil {
					return nil, err
				}
				return gapi.MakeModuleRequest(
					ctx,
					gitLabClient,
					modserver.RpcApiFromContext(ctx).AgentToken(),
					extra.ModuleName,
					header.Request.Method,
					header.Request.UrlPath,
					header.Request.UrlQuery(),
					header.Request.HttpHeader(),
					body,
				)
			},
		),
	}
}

func (s *server) MakeRequest(server rpc.GitlabAccess_MakeRequestServer) error {
	rpcApi := modserver.RpcApiFromContext(server.Context())
	return s.pipe.Pipe(rpcApi, server, modshared.NoAgentId)
}
