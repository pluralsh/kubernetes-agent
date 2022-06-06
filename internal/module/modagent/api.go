package modagent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/textproto"
	"net/url"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/httpz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/pkg/agentcfg"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"k8s.io/kubectl/pkg/cmd/util"
)

// Feature describes a particular feature that can be enabled or disabled by a module.
// If module is configured in a way that requires use of a feature, it should indicate that it needs the feature.
// All features are disabled by default.
type Feature int

type SubscribeCb func(enabled bool)

const (
	// Invalid is an invalid sentinel value.
	// Treat default value as "invalid" to avoid accidental confusion.
	Invalid Feature = iota
)

var (
	KnownFeatures = map[Feature]string{}
)

// Config holds configuration for a Module.
type Config struct {
	// Log can be used for logging from the module.
	// It should not be used for logging from gRPC API methods. Use grpctool.LoggerFromContext(ctx) instead.
	Log       *zap.Logger
	AgentMeta *modshared.AgentMeta
	Api       Api
	// K8sUtilFactory provides means to interact with the Kubernetes cluster agentk is running in.
	K8sUtilFactory util.Factory
	// KasConn is the gRPC connection to gitlab-kas.
	KasConn grpc.ClientConnInterface
	// Server is a gRPC server that can be used to expose API endpoints to gitlab-kas and/or GitLab.
	// This can be used to add endpoints in Factory.New.
	// Request handlers can obtain the per-request logger using grpctool.LoggerFromContext(requestContext).
	Server *grpc.Server
	// AgentName is a string "gitlab-agent". Can be used as a user agent, server name, service name, etc.
	AgentName string
}

type GitLabResponse struct {
	Status     string // e.g. "200 OK"
	StatusCode int32  // e.g. 200
	Header     http.Header
	Body       io.ReadCloser
}

// Api provides the API for the module to use.
type Api interface {
	modshared.Api
	MakeGitLabRequest(ctx context.Context, path string, opts ...GitLabRequestOption) (*GitLabResponse, error)
	ToggleFeature(feature Feature, enabled bool)
	SubscribeToFeatureStatus(feature Feature, cb SubscribeCb)
}

// RpcApi provides the API for the module's gRPC handlers to use.
type RpcApi interface {
	modshared.RpcApi
}

type Factory interface {
	// New creates a new instance of a Module.
	New(*Config) (Module, error)
	// Name returns module's name.
	Name() string
	// UsesInternalServer returns true if the module makes requests to the internal API server.
	UsesInternalServer() bool
}

type Module interface {
	// Run starts the module.
	// Run can block until the context is canceled or exit with nil if there is nothing to do.
	// cfg is a channel that gets configuration updates sent to it. It's closed when the module should shut down.
	// cfg is a shared instance, must not be mutated. Module should make a copy if it needs to mutate the object.
	// Applying configuration may take time, the provided context may signal done if module should shut down.
	// cfg only provides the latest available configuration, intermediate configuration states are discarded.
	Run(ctx context.Context, cfg <-chan *agentcfg.AgentConfiguration) error
	// DefaultAndValidateConfiguration applies defaults and validates the passed configuration.
	// It is called each time on configuration update before sending it via the channel passed to Run().
	// cfg is a shared instance, module can mutate only the part of it that it owns and only inside of this method.
	DefaultAndValidateConfiguration(cfg *agentcfg.AgentConfiguration) error
	// Name returns module's name.
	Name() string
}

type GitLabRequestConfig struct {
	Method string
	Header http.Header
	Query  url.Values
	Body   io.ReadCloser
}

func defaultRequestConfig() *GitLabRequestConfig {
	return &GitLabRequestConfig{
		Method: http.MethodGet,
		Header: make(http.Header),
		Query:  make(url.Values),
	}
}

func ApplyRequestOptions(opts []GitLabRequestOption) (*GitLabRequestConfig, error) {
	c := defaultRequestConfig()
	var firstErr error
	for _, o := range opts {
		err := o(c)
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}
	if firstErr != nil { // return the first error but close the body first
		if c.Body != nil {
			_ = c.Body.Close()
		}
		return nil, firstErr
	}
	return c, nil
}

type GitLabRequestOption func(*GitLabRequestConfig) error

func WithRequestHeader(header string, values ...string) GitLabRequestOption {
	return func(c *GitLabRequestConfig) error {
		c.Header[textproto.CanonicalMIMEHeaderKey(header)] = values
		return nil
	}
}

func WithRequestQueryParam(key string, values ...string) GitLabRequestOption {
	return func(c *GitLabRequestConfig) error {
		c.Query[key] = values
		return nil
	}
}

// WithRequestBody specifies request body to send and HTTP Content-Type header if contentType is not empty.
// If body implements io.ReadCloser, its Close() method will be called once the data has been sent.
// If body is nil, no body or Content-Type header is sent.
func WithRequestBody(body io.Reader, contentType string) GitLabRequestOption {
	return func(c *GitLabRequestConfig) error {
		if body == nil {
			return nil
		}
		if rc, ok := body.(io.ReadCloser); ok {
			c.Body = rc
		} else {
			c.Body = io.NopCloser(body)
		}
		if contentType != "" {
			c.Header[httpz.ContentTypeHeader] = []string{contentType}
		}
		return nil
	}
}

func WithJsonRequestBody(body interface{}) GitLabRequestOption {
	return func(c *GitLabRequestConfig) error {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("WithJsonRequestBody: %w", err)
		}
		return WithRequestBody(bytes.NewReader(bodyBytes), "application/json")(c)
	}
}

// WithRequestMethod specifies request HTTP method.
func WithRequestMethod(method string) GitLabRequestOption {
	return func(c *GitLabRequestConfig) error {
		c.Method = method
		return nil
	}
}
