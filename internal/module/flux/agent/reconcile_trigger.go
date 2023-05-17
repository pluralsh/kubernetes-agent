package agent

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/errz"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

const (
	defaultServiceApiBaseUrl = "http://webhook-receiver.flux-system.svc.cluster.local"
)

var (
	kubeProxyApiPathRegex = regexp.MustCompile("/api/v1/namespaces/[^/]+/services/[^/]+/proxy")
)

type reconcileTrigger interface {
	reconcile(ctx context.Context, webhookPath string) error
}

type gitrepositoryReconcileTrigger struct {
	baseUrl url.URL
	rt      http.RoundTripper
}

func newGitRepositoryReconcileTrigger(cfgUrl string, kubeApiUrl *url.URL, kubeApiRoundTripper http.RoundTripper, defaultRoundTripper http.RoundTripper) (*gitrepositoryReconcileTrigger, error) {
	if kubeProxyApiPathRegex.MatchString(cfgUrl) {
		u := *kubeApiUrl
		u.Path = path.Join(u.Path, cfgUrl)
		return &gitrepositoryReconcileTrigger{baseUrl: u, rt: kubeApiRoundTripper}, nil
	} else {
		u, err := url.Parse(cfgUrl)
		if err != nil {
			return nil, err
		}
		return &gitrepositoryReconcileTrigger{baseUrl: *u, rt: defaultRoundTripper}, nil
	}
}

func (t *gitrepositoryReconcileTrigger) reconcile(ctx context.Context, webhookPath string) (retErr error) {
	u := t.baseUrl
	u.Path = path.Join(u.Path, webhookPath)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), http.NoBody)
	if err != nil {
		return err
	}
	resp, err := t.rt.RoundTrip(req) // nolint:bodyclose
	if err != nil {
		return err
	}
	defer errz.SafeClose(resp.Body, &retErr)
	// draining response body so that the underlying transport can reuse the connection,
	// see https://pkg.go.dev/net/http#Response
	if _, err = io.Copy(io.Discard, io.LimitReader(resp.Body, 8*1024)); err != nil {
		return fmt.Errorf("failed to drain response body to reconciliation trigger request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("trigger to %q returned status %q", u.String(), resp.Status)
	}
	return nil
}

// This is a copy from k8s.io/client-go/rest/url_utils.go

// defaultServerUrlFor is shared between IsConfigTransportTLS and RESTClientFor. It
// requires Host and Version to be set prior to being called.
func defaultServerUrlFor(config *rest.Config) (*url.URL, string, error) {
	// TODO: move the default to secure when the apiserver supports TLS by default
	// config.Insecure is taken to mean "I want HTTPS but don't bother checking the certs against a CA."
	hasCA := len(config.CAFile) != 0 || len(config.CAData) != 0
	hasCert := len(config.CertFile) != 0 || len(config.CertData) != 0
	defaultTLS := hasCA || hasCert || config.Insecure
	host := config.Host
	if host == "" {
		host = "localhost"
	}

	if config.GroupVersion != nil {
		return rest.DefaultServerURL(host, config.APIPath, *config.GroupVersion, defaultTLS)
	}
	return rest.DefaultServerURL(host, config.APIPath, schema.GroupVersion{}, defaultTLS)
}
