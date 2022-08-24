package gitlab

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/httpz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/tlstool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

const (
	// Default retry configuration
	defaultRetryWaitMin = 100 * time.Millisecond
	defaultRetryWaitMax = 30 * time.Second
	defaultRetryMax     = 4
)

type RetryConfig struct {
	// Logger instance. Can be either retryablehttp.Logger or retryablehttp.LeveledLogger
	Logger interface{}

	RetryWaitMin time.Duration // Minimum time to wait
	RetryWaitMax time.Duration // Maximum time to wait
	RetryMax     int           // Maximum number of retries

	// RequestLogHook allows a user-supplied function to be called
	// before each retry.
	RequestLogHook retryablehttp.RequestLogHook

	// ResponseLogHook allows a user-supplied function to be called
	// with the response from each HTTP request executed.
	ResponseLogHook retryablehttp.ResponseLogHook

	// CheckRetry specifies the policy for handling retries, and is called
	// after each request. The default policy is retryablehttp.DefaultRetryPolicy.
	CheckRetry retryablehttp.CheckRetry

	// Backoff specifies the policy for how long to wait between retries.
	// retryablehttp.DefaultBackoff is used by default.
	Backoff retryablehttp.Backoff
}

type transportConfig struct {
	Proxy                 func(*http.Request) (*url.URL, error)
	DialContext           func(ctx context.Context, network, address string) (net.Conn, error)
	TLSClientConfig       *tls.Config
	TLSHandshakeTimeout   time.Duration
	MaxIdleConns          int
	MaxIdleConnsPerHost   int
	MaxConnsPerHost       int
	IdleConnTimeout       time.Duration
	ResponseHeaderTimeout time.Duration
	ExpectContinueTimeout time.Duration
	ForceAttemptHTTP2     bool
}

// clientConfig holds configuration for the client.
type clientConfig struct {
	retryConfig     RetryConfig
	transportConfig transportConfig
	tracePropagator propagation.TextMapPropagator
	limiter         httpz.Limiter
	userAgent       string
}

// ClientOption to configure the client.
type ClientOption func(*clientConfig)

func applyClientOptions(opts []ClientOption) clientConfig {
	dialer := &net.Dialer{
		Timeout: 10 * time.Second,
	}
	config := clientConfig{
		retryConfig: RetryConfig{
			RetryWaitMin: defaultRetryWaitMin,
			RetryWaitMax: defaultRetryWaitMax,
			RetryMax:     defaultRetryMax,
			CheckRetry:   retryablehttp.DefaultRetryPolicy,
			Backoff:      retryablehttp.DefaultBackoff,
		},
		transportConfig: transportConfig{
			Proxy:                 http.ProxyFromEnvironment,
			DialContext:           dialer.DialContext,
			TLSClientConfig:       tlstool.DefaultClientTLSConfig(),
			TLSHandshakeTimeout:   10 * time.Second,
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   50,
			MaxConnsPerHost:       50,
			IdleConnTimeout:       90 * time.Second,
			ResponseHeaderTimeout: 20 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			ForceAttemptHTTP2:     true,
		},
		tracePropagator: otel.GetTextMapPropagator(),
		userAgent:       "",
	}
	for _, v := range opts {
		v(&config)
	}

	return config
}

// WithRetryConfig configures retry behavior.
func WithRetryConfig(retryConfig RetryConfig) ClientOption {
	return func(config *clientConfig) {
		config.retryConfig = retryConfig
	}
}

// WithTextMapPropagator sets a custom tracer to be used, otherwise the OTEL's global TextMapPropagator is used.
func WithTextMapPropagator(p propagation.TextMapPropagator) ClientOption {
	return func(config *clientConfig) {
		config.tracePropagator = p
	}
}

// WithUserAgent configures the User-Agent header on the http client.
func WithUserAgent(userAgent string) ClientOption {
	return func(config *clientConfig) {
		config.userAgent = userAgent
	}
}

// WithTLSConfig sets the TLS config to use.
func WithTLSConfig(tlsConfig *tls.Config) ClientOption {
	return func(config *clientConfig) {
		config.transportConfig.TLSClientConfig = tlsConfig
	}
}

// WithRateLimiter sets the rate limiter to use.
func WithRateLimiter(limiter httpz.Limiter) ClientOption {
	return func(config *clientConfig) {
		config.limiter = limiter
	}
}
