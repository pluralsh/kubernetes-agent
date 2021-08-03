package agent

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/kubernetes_api/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/grpctool"
)

const (
	userAgentHeaderName = "User-Agent"
)

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type server struct {
	rpc.UnimplementedKubernetesApiServer
	api  modagent.API
	pipe *grpctool.InboundGrpcToOutboundHttp
}

func newServer(api modagent.API, userAgent string, client httpClient, baseUrl *url.URL) *server {
	return &server{
		api: api,
		pipe: grpctool.NewInboundGrpcToOutboundHttp(
			api,
			func(ctx context.Context, h *grpctool.HttpRequest_Header, body io.Reader) (*http.Response, error) {
				u := *baseUrl
				u.Path = h.Request.UrlPath
				u.RawQuery = h.Request.UrlQuery().Encode()

				req, err := http.NewRequestWithContext(ctx, h.Request.Method, u.String(), body)
				if err != nil {
					return nil, err
				}
				req.Header = h.Request.HttpHeader()
				ua := req.Header.Get(userAgentHeaderName)
				if ua == "" {
					ua = userAgent
				} else {
					ua = fmt.Sprintf("%s via %s", ua, userAgent)
				}
				req.Header.Set(userAgentHeaderName, ua)

				resp, err := client.Do(req)
				if err != nil {
					select {
					case <-ctx.Done(): // assume request errored out because of context
						return nil, ctx.Err()
					default:
						return nil, err
					}
				}
				return resp, nil
			},
		),
	}
}

func (m *server) MakeRequest(server rpc.KubernetesApi_MakeRequestServer) error {
	return m.pipe.Pipe(server)
}
