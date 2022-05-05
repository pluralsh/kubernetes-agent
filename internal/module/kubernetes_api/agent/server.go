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
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/httpz"
)

type httpClient interface {
	// Do performs the request.
	// impConfig may be nil.
	Do(impConfig *rpc.ImpersonationConfig, r *http.Request) (*http.Response, error)
}

type server struct {
	rpc.UnimplementedKubernetesApiServer
	userAgent string
	via       string
	client    httpClient
	baseUrl   *url.URL
}

func newServer(userAgent string, client httpClient, baseUrl *url.URL) *server {
	return &server{
		userAgent: userAgent,
		via:       "gRPC/1.0 " + userAgent,
		client:    client,
		baseUrl:   baseUrl,
	}
}

func (s *server) MakeRequest(server rpc.KubernetesApi_MakeRequestServer) error {
	rpcApi := modagent.RpcApiFromContext(server.Context())
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

func (s *server) httpDo(ctx context.Context, h *grpctool.HttpRequest_Header, body io.Reader) (*http.Response, error) {
	u := *s.baseUrl
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
	req.Header[httpz.ViaHeader] = append(req.Header[httpz.ViaHeader], s.via)
	resp, err := s.client.Do(headerExtra.ImpConfig, req)
	if err != nil {
		select {
		case <-ctx.Done(): // assume request errored out because of context
			return nil, ctx.Err()
		default:
			return nil, err
		}
	}
	resp.Header[httpz.ViaHeader] = append(resp.Header[httpz.ViaHeader], fmt.Sprintf("%d.%d %s", resp.ProtoMajor, resp.ProtoMinor, s.userAgent))
	return resp, nil
}
