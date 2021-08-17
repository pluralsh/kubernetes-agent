package agent

import (
	"io"
	"net/http"
	"strings"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/kubernetes_api/rpc"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/transport"
)

type impersonatingClient struct {
	restConfig *rest.Config
}

func (c *impersonatingClient) Do(impConfig *rpc.ImpersonationConfig, r *http.Request) (*http.Response, error) {
	var config *rest.Config
	restImp := !isEmptyImpersonationConfig(c.restConfig.Impersonate)
	cfgImp := !impConfig.IsEmpty()
	reqImp := hasImpersonationHeaders(r)
	switch {
	case !restImp && !cfgImp && !reqImp:
		// No impersonation
		config = c.restConfig
	case restImp && !cfgImp && !reqImp:
		// Impersonation is configured in the rest config
		config = c.restConfig
	case !restImp && cfgImp && !reqImp:
		// Impersonation is configured in the agent config
		config = rest.CopyConfig(c.restConfig) // copy to avoid mutating a potentially shared config object
		config.Impersonate.UserName = impConfig.Username
		config.Impersonate.Groups = impConfig.Groups
		// TODO Add uid when we upgrade to Kubernetes 1.22 libraries
		config.Impersonate.Extra = impConfig.GetExtraAsMap()
	case !restImp && !cfgImp && reqImp:
		// Impersonation is configured in the HTTP request
		config = c.restConfig
	default:
		// Nested impersonation support https://gitlab.com/gitlab-org/gitlab/-/issues/338664
		return httpErrorResponse(http.StatusBadRequest, "Nested impersonation is not supported - agent is already configured to impersonate an identity"), nil
	}
	transportForConfig, err := rest.TransportFor(config)
	if err != nil {
		return nil, err
	}
	client := http.Client{
		Transport:     transportForConfig,
		CheckRedirect: useLastResponse,
	}
	return client.Do(r)
}

func httpErrorResponse(statusCode int, text string) *http.Response {
	return &http.Response{
		Status:     http.StatusText(statusCode),
		StatusCode: statusCode,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header: http.Header{ // What http.Error() returns
			"Content-Type":           {"text/plain; charset=utf-8"},
			"X-Content-Type-Options": {"nosniff"},
		},
		Body: io.NopCloser(strings.NewReader(text)),
	}
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
	// header==transport.ImpersonateUidHeader: TODO add when we upgrade to Kubernetes 1.22 libraries
	return header == transport.ImpersonateUserHeader || header == transport.ImpersonateGroupHeader || strings.HasPrefix(header, transport.ImpersonateUserExtraHeaderPrefix)
}

func useLastResponse(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}
