package modserver

import (
	"context"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitaly"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitlab"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/observability"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/usage_metrics"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/kascfg"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	// RoutingHopPrefix is a metadata key prefix that is used for metadata keys that should be consumed by
	// the gateway kas instances and not passed along to agentk.
	RoutingHopPrefix = "kas-hop-"
	// RoutingAgentIdMetadataKey is used to pass destination agent id in request metadata
	// from the routing kas instance, that is handling the incoming request, to the gateway kas instance,
	// that is forwarding the request to an agentk.
	RoutingAgentIdMetadataKey = RoutingHopPrefix + "routing-agent-id"

	// TraceIdSentryField is the name of the Sentry field for trace ID.
	TraceIdSentryField     = "trace_id"
	GrpcServiceSentryField = "grpc.service"
	GrpcMethodSentryField  = "grpc.method"
)

// ApplyDefaults is a signature of a public function, exposed by modules to perform defaulting.
// The function should be called ApplyDefaults.
type ApplyDefaults func(*kascfg.ConfigurationFile)

// Config holds configuration for a Module.
type Config struct {
	// Log can be used for logging from the module.
	// It should not be used for logging from gRPC Api methods. Use grpctool.LoggerFromContext(ctx) instead.
	Log          *zap.Logger
	Api          Api
	Config       *kascfg.ConfigurationFile
	GitLabClient gitlab.ClientInterface
	// Registerer allows to register metrics.
	// Metrics should be registered in Run and unregistered before Run returns.
	Registerer   prometheus.Registerer
	UsageTracker usage_metrics.UsageTrackerRegisterer
	// AgentServer is the gRPC server agentk is talking to.
	// This can be used to add endpoints in Factory.New.
	// Request handlers can obtain the per-request logger using grpctool.LoggerFromContext(requestContext).
	AgentServer *grpc.Server
	// ApiServer is the gRPC server GitLab is talking to.
	// This can be used to add endpoints in Factory.New.
	// Request handlers can obtain the per-request logger using grpctool.LoggerFromContext(requestContext).
	ApiServer *grpc.Server
	// RegisterAgentApi allows to register a gRPC Api endpoint that kas proxies to agentk.
	RegisterAgentApi func(*grpc.ServiceDesc)
	// AgentConn is a gRPC connection that can be used to send requests to an agentk instance.
	// Agent Id must be specified in the request metadata in RoutingAgentIdMetadataKey field.
	// Make sure factory returns modshared.ModuleStartAfterServers if module uses this connection.
	AgentConn       grpc.ClientConnInterface
	Gitaly          gitaly.PoolInterface
	TraceProvider   trace.TracerProvider
	TracePropagator propagation.TextMapPropagator
	MeterProvider   metric.MeterProvider
	RedisClient     redis.UniversalClient
	// KasName is a string "gitlab-kas". Can be used as a user agent, server name, service name, etc.
	KasName string
	// Version is gitlab-kas version.
	Version string
	// CommitId is gitlab-kas commit sha.
	CommitId string
	// ProbeRegistry is for registering liveness probes and readiness probes
	ProbeRegistry *observability.ProbeRegistry
}

type GitPushEventCallback func(ctx context.Context, project *Project)

// Api provides the API for the module to use.
type Api interface {
	modshared.Api
	OnGitPushEvent(ctx context.Context, callback GitPushEventCallback)
}

type Factory interface {
	modshared.Factory
	// New creates a new instance of a Module.
	New(*Config) (Module, error)
}

type Module interface {
	// Run starts the module.
	// Run can block until the context is canceled or exit with nil if there is nothing to do.
	Run(context.Context) error
	// Name returns module's name.
	Name() string
}

func RoutingMetadata(agentId int64) metadata.MD {
	return metadata.MD{
		RoutingAgentIdMetadataKey: []string{strconv.FormatInt(agentId, 10)},
	}
}
