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

const (
	// https://datatracker.ietf.org/doc/html/rfc7230#section-5.7.1
	httpViaHeader = "Via"
)

type httpClient interface {
	// Do performs the request.
	// impConfig may be nil.
	Do(impConfig *rpc.ImpersonationConfig, r *http.Request) (*http.Response, error)
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
				var headerExtra rpc.HeaderExtra
				if h.Extra != nil { // May not be there on older kas versions. Also, just be more robust.
					err = h.Extra.UnmarshalTo(&headerExtra)
					if err != nil {
						return nil, err
					}
				}
				req.Header = h.Request.HttpHeader()
				req.Header.Add(httpViaHeader, via)
				resp, err := client.Do(headerExtra.ImpConfig, req)
				if err != nil {
					select {
					case <-ctx.Done(): // assume request errored out because of context
						return nil, ctx.Err()
					default:
						return nil, err
					}
				}
				resp.Header.Add(httpViaHeader, fmt.Sprintf("%d.%d %s", resp.ProtoMajor, resp.ProtoMinor, userAgent))
				return resp, nil
			},
		),
	}
}

func (m *server) MakeRequest(server rpc.KubernetesApi_MakeRequestServer) error {
	rpcApi := modagent.RpcApiFromContext(server.Context())
	return m.pipe.Pipe(rpcApi, server, modshared.NoAgentId)
}
