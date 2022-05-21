package agent

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/kubernetes_api/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/httpz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/prototool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/tlstool"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/transport"
)

type server struct {
	rpc.UnimplementedKubernetesApiServer
	restConfig *rest.Config
	baseUrl    *url.URL
	userAgent  string
	via        string
}

func newServer(restConfig *rest.Config, baseUrl *url.URL, userAgent string) *server {
	return &server{
		restConfig: restConfig,
		baseUrl:    baseUrl,
		userAgent:  userAgent,
		via:        "gRPC/1.0 " + userAgent,
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
		HandleIoError: func(msg string, err error) error {
			return rpcApi.HandleIoError(log, msg, err)
		},
		HttpDo: s.httpDo,
	}
	return grpc2http.Pipe(server)
}

func (s *server) httpDo(ctx context.Context, h *grpctool.HttpRequest_Header, body io.Reader) (grpctool.DoResponse, error) {
	// 1. Construct request
	req, err := s.newRequest(ctx, h.Request, body)
	if err != nil {
		return grpctool.DoResponse{}, err
	}
	// 2. Construct rest config
	var headerExtra rpc.HeaderExtra
	if h.Extra != nil { // May not be there on older kas versions. Also, just be more robust.
		err = h.Extra.UnmarshalTo(&headerExtra)
		if err != nil {
			return grpctool.DoResponse{}, err
		}
	}
	restConfig, err := restImpersonationConfig(headerExtra.ImpConfig, s.restConfig, req)
	if err != nil {
		return grpctool.DoResponse{}, err
	}
	// 3. Construct round tripper
	var (
		rt        http.RoundTripper
		upgradeRT *httpz.UpgradeRoundTripper
	)
	transportCfg, err := restConfig.TransportConfig()
	if err != nil {
		return grpctool.DoResponse{}, err
	}
	isUpgrade := h.Request.IsUpgrade()
	if isUpgrade {
		var tlsConfig *tls.Config
		tlsConfig, err = transport.TLSConfigFor(transportCfg)
		if err != nil {
			return grpctool.DoResponse{}, err
		}
		if tlsConfig == nil { // transport.TLSConfigFor() can return nil
			tlsConfig = tlstool.DefaultClientTLSConfig()
		}
		tlsConfig.NextProtos = []string{"http/1.1"} // HTTP Upgrade doesn't work over HTTP/2, so enforce HTTP/1.1
		dialer := &net.Dialer{
			Timeout: 30 * time.Second,
		}
		upgradeRT = &httpz.UpgradeRoundTripper{
			Dialer: dialer,
			TlsDialer: &tls.Dialer{
				NetDialer: dialer,
				Config:    tlsConfig,
			},
		}
		rt, err = transport.HTTPWrappersForConfig(transportCfg, upgradeRT)
	} else {
		rt, err = transport.New(transportCfg) // returns pooled transports that reuse TCP connections
	}
	if err != nil {
		return grpctool.DoResponse{}, err
	}
	// 4. Make a request
	resp, err := rt.RoundTrip(req) // nolint: bodyclose
	if err != nil {
		ctxErr := ctx.Err()
		if ctxErr != nil {
			err = ctxErr // assume request errored out because of context
		}
		return grpctool.DoResponse{}, err
	}
	resp.Header[httpz.ViaHeader] = append(resp.Header[httpz.ViaHeader], fmt.Sprintf("%d.%d %s", resp.ProtoMajor, resp.ProtoMinor, s.userAgent))
	if isUpgrade {
		return grpctool.DoResponse{
			Resp:        resp,
			UpgradeConn: upgradeRT.Conn,
			ConnReader:  upgradeRT.ConnReader,
		}, nil
	} else {
		return grpctool.DoResponse{
			Resp: resp,
		}, nil
	}
}

func (s *server) newRequest(ctx context.Context, requestInfo *prototool.HttpRequest, body io.Reader) (*http.Request, error) {
	u := *s.baseUrl
	u.Path = requestInfo.UrlPath
	u.RawQuery = requestInfo.UrlQuery().Encode()

	req, err := http.NewRequestWithContext(ctx, requestInfo.Method, u.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header = requestInfo.HttpHeader()
	req.Header[httpz.ViaHeader] = append(req.Header[httpz.ViaHeader], s.via)
	return req, nil
}

func restImpersonationConfig(impConfig *rpc.ImpersonationConfig, restConfig *rest.Config, r *http.Request) (*rest.Config, error) {
	restImp := !isEmptyImpersonationConfig(restConfig.Impersonate)
	cfgImp := !impConfig.IsEmpty()
	reqImp := hasImpersonationHeaders(r)
	switch {
	case !restImp && !cfgImp && !reqImp:
		// No impersonation
	case restImp && !cfgImp && !reqImp:
		// Impersonation is configured in the rest config
	case !restImp && cfgImp && !reqImp:
		// Impersonation is configured in the agent config
		restConfig = rest.CopyConfig(restConfig) // copy to avoid mutating a potentially shared config object
		restConfig.Impersonate.UserName = impConfig.Username
		restConfig.Impersonate.UID = impConfig.Uid
		restConfig.Impersonate.Groups = impConfig.Groups
		restConfig.Impersonate.Extra = impConfig.GetExtraAsMap()
	case !restImp && !cfgImp && reqImp:
		// Impersonation is configured in the HTTP request
	default:
		// Nested impersonation support https://gitlab.com/gitlab-org/gitlab/-/issues/338664
		return nil, errors.New("nested impersonation is not supported - agent is already configured to impersonate an identity")
	}
	return restConfig, nil
}

func isEmptyImpersonationConfig(cfg rest.ImpersonationConfig) bool {
	return cfg.UserName == "" && len(cfg.Groups) == 0 && len(cfg.Extra) == 0
}

func hasImpersonationHeaders(r *http.Request) bool {
	for k := range r.Header {
		if isImpersonationHeader(k) {
			return true
		}
	}
	return false
}

func isImpersonationHeader(header string) bool {
	return header == transport.ImpersonateUserHeader ||
		header == transport.ImpersonateUIDHeader ||
		header == transport.ImpersonateGroupHeader ||
		strings.HasPrefix(header, transport.ImpersonateUserExtraHeaderPrefix)
}
