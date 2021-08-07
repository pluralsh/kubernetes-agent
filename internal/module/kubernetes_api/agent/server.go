package agent

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/kubernetes_api/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/grpctool"
)

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type server struct {
	rpc.UnimplementedKubernetesApiServer
	pipe *grpctool.InboundGrpcToOutboundHttp
}

func newServer(userAgent string, client httpClient, baseUrl *url.URL) *server {
	via := "gRPC/1.0 " + userAgent
	return &server{
		pipe: grpctool.NewInboundGrpcToOutboundHttp(
			func(ctx context.Context, h *grpctool.HttpRequest_Header, body io.Reader) (*http.Response, error) {
				u := *baseUrl
				u.Path = h.Request.UrlPath
				u.RawQuery = h.Request.UrlQuery().Encode()

				req, err := http.NewRequestWithContext(ctx, h.Request.Method, u.String(), body)
				if err != nil {
					return nil, err
				}
				req.Header = h.Request.HttpHeader()
				req.Header.Add("Via", via)

				resp, err := client.Do(req)
				if err != nil {
					select {
					case <-ctx.Done(): // assume request errored out because of context
						return nil, ctx.Err()
					default:
						return nil, err
					}
				}
				resp.Header.Add("Via", fmt.Sprintf("%d.%d %s", resp.ProtoMajor, resp.ProtoMinor, userAgent))
				return resp, nil
			},
		),
	}
}

func (m *server) MakeRequest(server rpc.KubernetesApi_MakeRequestServer) error {
	rpcApi := modagent.RpcApiFromContext(server.Context())
	return m.pipe.Pipe(rpcApi, server, modshared.NoAgentId)
}
