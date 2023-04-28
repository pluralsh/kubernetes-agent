package gitlab

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/hashicorp/go-retryablehttp"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/httpz"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const (
	// This header carries the JWT token for gitlab-rails
	jwtRequestHeader  = "Gitlab-Kas-Api-Request"
	jwtValidFor       = 30 * time.Second
	jwtNotBefore      = 5 * time.Second
	jwtIssuer         = "gitlab-kas"
	jwtGitLabAudience = "gitlab"
)

type HTTPClient interface {
	Do(*retryablehttp.Request) (*http.Response, error)
}

type Client struct {
	Backend           *url.URL
	HTTPClient        HTTPClient
	HTTPClientNoRetry HTTPClient
	AuthSecret        []byte
	UserAgent         string
}

func NewClient(backend *url.URL, authSecret []byte, opts ...ClientOption) *Client {
	o := applyClientOptions(opts)
	var transport http.RoundTripper = &http.Transport{
		Proxy:                 o.transportConfig.Proxy,
		DialContext:           o.transportConfig.DialContext,
		TLSClientConfig:       o.transportConfig.TLSClientConfig,
		TLSHandshakeTimeout:   o.transportConfig.TLSHandshakeTimeout,
		MaxIdleConns:          o.transportConfig.MaxIdleConns,
		MaxIdleConnsPerHost:   o.transportConfig.MaxIdleConnsPerHost,
		MaxConnsPerHost:       o.transportConfig.MaxConnsPerHost,
		IdleConnTimeout:       o.transportConfig.IdleConnTimeout,
		ResponseHeaderTimeout: o.transportConfig.ResponseHeaderTimeout,
		ForceAttemptHTTP2:     o.transportConfig.ForceAttemptHTTP2,
	}
	if o.limiter != nil {
		transport = &httpz.RateLimitingRoundTripper{
			Delegate: transport,
			Limiter:  o.limiter,
		}
	}
	httpClient := &http.Client{
		Transport: otelhttp.NewTransport(
			transport,
			otelhttp.WithPropagators(o.tracePropagator),
		),
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	return &Client{
		Backend: backend,
		HTTPClient: &retryablehttp.Client{
			HTTPClient:      httpClient,
			Logger:          o.retryConfig.Logger,
			RetryWaitMin:    o.retryConfig.RetryWaitMin,
			RetryWaitMax:    o.retryConfig.RetryWaitMax,
			RetryMax:        o.retryConfig.RetryMax,
			RequestLogHook:  o.retryConfig.RequestLogHook,
			ResponseLogHook: o.retryConfig.ResponseLogHook,
			CheckRetry:      o.retryConfig.CheckRetry,
			Backoff:         o.retryConfig.Backoff,
			ErrorHandler:    errorHandler,
		},
		HTTPClientNoRetry: &retryablehttp.Client{
			HTTPClient:      httpClient,
			Logger:          o.retryConfig.Logger,
			RetryMax:        0,
			RequestLogHook:  o.retryConfig.RequestLogHook,
			ResponseLogHook: o.retryConfig.ResponseLogHook,
			CheckRetry:      o.retryConfig.CheckRetry,
			ErrorHandler:    errorHandler,
		},
		AuthSecret: authSecret,
		UserAgent:  o.userAgent,
	}
}

func (c *Client) Do(ctx context.Context, opts ...DoOption) error {
	o, err := applyDoOptions(opts)
	if err != nil {
		return err
	}
	r, err := retryablehttp.NewRequestWithContext(ctx, o.method, c.targetUrl(o.path, o.query), o.body)
	if err != nil {
		return fmt.Errorf("NewRequest: %w", err)
	}
	if len(o.header) > 0 {
		r.Header = o.header
	}
	if o.withJWT {
		signedClaims, signErr := c.jwtSignature()
		if signErr != nil {
			return signErr
		}
		r.Header[jwtRequestHeader] = []string{signedClaims}
	}
	if c.UserAgent != "" {
		r.Header[httpz.UserAgentHeader] = []string{c.UserAgent}
	}

	client := c.HTTPClient
	if o.noRetry {
		client = c.HTTPClientNoRetry
	}
	resp, err := client.Do(r) // nolint: bodyclose
	if err != nil {
		ctxErr := ctx.Err()
		if ctxErr != nil {
			err = ctxErr // assume request errored out because of context
		}
	}
	return o.responseHandler.Handle(resp, err)
}

func (c *Client) jwtSignature() (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Issuer:    jwtIssuer,
		Audience:  jwt.ClaimStrings{jwtGitLabAudience},
		ExpiresAt: jwt.NewNumericDate(now.Add(jwtValidFor)),
		NotBefore: jwt.NewNumericDate(now.Add(-jwtNotBefore)),
		IssuedAt:  jwt.NewNumericDate(now),
	}
	signedClaims, claimsErr := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString(c.AuthSecret)
	if claimsErr != nil {
		return "", fmt.Errorf("sign JWT: %w", claimsErr)
	}
	return signedClaims, nil
}

func (c *Client) targetUrl(path string, query url.Values) string {
	u := *c.Backend
	u.Path = joinUrlPaths(u.Path, path)
	q := query // may be nil
	if u.RawQuery != "" {
		if len(q) == 0 {
			// Nothing to do
		} else {
			// Merge queries
			uq := u.Query()
			for k, v := range q {
				uq[k] = v
			}
			q = uq
		}
	}
	u.RawQuery = q.Encode() // handles query == nil
	return u.String()
}

// errorHandler returns the last response and error when ran out of retry attempts.
// It masks values of URL query parameters.
func errorHandler(resp *http.Response, err error, numTries int) (*http.Response, error) {
	ue, ok := err.(*url.Error) // nolint: errorlint
	if !ok {
		return resp, err
	}
	u, parseErr := url.Parse(ue.URL)
	if parseErr != nil {
		return resp, err
	}
	if u.RawQuery != "" {
		maskURLQueryParams(u)
		ue.URL = u.String()
	}
	return resp, ue
}

func maskURLQueryParams(u *url.URL) {
	newVal := []string{"x"}
	q := u.Query()
	for k := range q {
		q[k] = newVal
	}
	u.RawQuery = q.Encode()
}

func joinUrlPaths(head, tail string) string {
	if head == "" {
		return tail
	}
	if tail == "" {
		return head
	}
	head = strings.TrimSuffix(head, "/")
	tail = strings.TrimPrefix(tail, "/")
	return head + "/" + tail
}
