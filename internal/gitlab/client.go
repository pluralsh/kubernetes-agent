package gitlab

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/hashicorp/go-retryablehttp"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/httpz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/tracing"
	"gitlab.com/gitlab-org/labkit/correlation"
)

const (
	// This header carries the JWT token for gitlab-rails
	jwtRequestHeader  = "Gitlab-Kas-Api-Request"
	jwtValidFor       = 5 * time.Second
	jwtNotBefore      = 5 * time.Second
	jwtIssuer         = "gitlab-kas"
	jwtGitLabAudience = "gitlab"
)

type HTTPClient interface {
	Do(*retryablehttp.Request) (*http.Response, error)
}

type Client struct {
	Backend    *url.URL
	HTTPClient HTTPClient
	AuthSecret []byte
	UserAgent  string
}

func NewClient(backend *url.URL, authSecret []byte, opts ...ClientOption) *Client {
	o := applyClientOptions(opts)
	var transport http.RoundTripper = &http.Transport{
		Proxy:                 o.transportConfig.Proxy,
		DialContext:           o.transportConfig.DialContext,
		TLSClientConfig:       o.transportConfig.TLSClientConfig,
		TLSHandshakeTimeout:   o.transportConfig.ExpectContinueTimeout,
		MaxIdleConns:          o.transportConfig.MaxIdleConns,
		MaxIdleConnsPerHost:   o.transportConfig.MaxIdleConnsPerHost,
		MaxConnsPerHost:       o.transportConfig.MaxConnsPerHost,
		IdleConnTimeout:       o.transportConfig.IdleConnTimeout,
		ResponseHeaderTimeout: o.transportConfig.ResponseHeaderTimeout,
		ExpectContinueTimeout: o.transportConfig.ExpectContinueTimeout,
		ForceAttemptHTTP2:     o.transportConfig.ForceAttemptHTTP2,
	}
	if o.limiter != nil {
		transport = &httpz.RateLimitingRoundTripper{
			Delegate: transport,
			Limiter:  o.limiter,
		}
	}
	return &Client{
		Backend: backend,
		HTTPClient: &retryablehttp.Client{
			HTTPClient: &http.Client{
				Transport: tracing.NewRoundTripper(
					correlation.NewInstrumentedRoundTripper(
						transport,
						correlation.WithClientName(o.clientName),
					),
					tracing.WithRoundTripperTracer(o.tracer),
					tracing.WithLogger(o.log),
				),
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			},
			Logger:          o.retryConfig.Logger,
			RetryWaitMin:    o.retryConfig.RetryWaitMin,
			RetryWaitMax:    o.retryConfig.RetryWaitMax,
			RetryMax:        o.retryConfig.RetryMax,
			RequestLogHook:  o.retryConfig.RequestLogHook,
			ResponseLogHook: o.retryConfig.ResponseLogHook,
			CheckRetry:      o.retryConfig.CheckRetry,
			Backoff:         o.retryConfig.Backoff,
			ErrorHandler: func(resp *http.Response, err error, numTries int) (*http.Response, error) {
				// Just return the last response and error when ran out of retry attempts.
				return resp, err
			},
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
	u := *c.Backend
	u.Path = o.path
	u.RawQuery = o.query.Encode() // handles query == nil
	r, err := retryablehttp.NewRequest(o.method, u.String(), o.body)
	if err != nil {
		return fmt.Errorf("NewRequest: %w", err)
	}
	r = r.WithContext(ctx)
	if len(o.header) > 0 {
		r.Header = o.header
	}
	if o.withJWT {
		now := time.Now()
		claims := jwt.StandardClaims{
			Audience:  jwtGitLabAudience,
			ExpiresAt: now.Add(jwtValidFor).Unix(),
			IssuedAt:  now.Unix(),
			Issuer:    jwtIssuer,
			NotBefore: now.Add(-jwtNotBefore).Unix(),
		}
		signedClaims, claimsErr := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
			SignedString(c.AuthSecret)
		if claimsErr != nil {
			return fmt.Errorf("sign JWT: %w", claimsErr)
		}
		r.Header.Set(jwtRequestHeader, signedClaims)
	}
	if c.UserAgent != "" {
		r.Header.Set("User-Agent", c.UserAgent)
	}

	resp, err := c.HTTPClient.Do(r) // nolint: bodyclose
	if err != nil {
		select {
		case <-ctx.Done(): // assume request errored out because of context
			err = ctx.Err()
		default:
		}
	}
	return o.responseHandler.Handle(resp, err)
}
