package logz

// Do not add more dependencies to this package as it's depended upon by the whole codebase.

import (
	"context"
	"net"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// These constants are for type-safe zap field helpers that are not here to:
// - avoid adding a dependency or
// - because they are not generally useful.
// Field names are here to make it possible to see all field names that are in use.

const (
	AgentConfig = "agent_config"
)

func NetAddressFromAddr(addr net.Addr) zap.Field {
	return NetAddress(addr.String())
}

func NetNetworkFromAddr(addr net.Addr) zap.Field {
	return NetNetwork(addr.Network())
}

func NetAddress(listenAddress string) zap.Field {
	return zap.String("net_address", listenAddress)
}

func NetNetwork(listenNetwork string) zap.Field {
	return zap.String("net_network", listenNetwork)
}

func IsWebSocket(isWebSocket bool) zap.Field {
	return zap.Bool("is_websocket", isWebSocket)
}

func AgentId(agentId int64) zap.Field {
	return zap.Int64("agent_id", agentId)
}

func CommitId(commitId string) zap.Field {
	return zap.String("commit_id", commitId)
}

func NumberOfFilesVisited(n uint32) zap.Field {
	return zap.Uint32("files_visited", n)
}
func NumberOfFilesSent(n uint32) zap.Field {
	return zap.Uint32("files_sent", n)
}

// ProjectId is the human-readable GitLab project path (e.g. gitlab-org/gitlab).
func ProjectId(projectId string) zap.Field {
	return zap.String("project_id", projectId)
}

// WorkerId is an id of the work source such as project id or chart name. (e.g. gitlab-org/gitlab).
func WorkerId(workerId string) zap.Field {
	return zap.String("worker_id", workerId)
}

func TraceIdFromContext(ctx context.Context) zap.Field {
	return TraceId(trace.SpanContextFromContext(ctx).TraceID())
}

func TraceId(traceId trace.TraceID) zap.Field {
	if !traceId.IsValid() {
		return zap.Skip()
	}
	return zap.String("trace_id", traceId.String())
}

// Use for any keys in Redis.
func RedisKey(key []byte) zap.Field {
	return zap.Binary("redis_key", key)
}

// Use for any integer counters.
func U64Count(count uint64) zap.Field {
	return zap.Uint64("count", count)
}

// Use for any integer counters.
func TokenLimit(limit uint64) zap.Field {
	return zap.Uint64("token_limit", limit)
}

func RemovedHashKeys(n int) zap.Field {
	return zap.Int("removed_hash_keys", n)
}

// GitLab-kas or agentk module name.
func ModuleName(name string) zap.Field {
	return zap.String("mod_name", name)
}

func KasUrl(kasUrl string) zap.Field {
	return zap.String("kas_url", kasUrl)
}

func PoolConnectionUrl(poolConnUrl string) zap.Field {
	return zap.String("pool_conn_url", poolConnUrl)
}

func UrlPathPrefix(urlPrefix string) zap.Field {
	return zap.String("url_path_prefix", urlPrefix)
}

func UrlPath(url string) zap.Field {
	return zap.String("url_path", url)
}

func GrpcService(service string) zap.Field {
	return zap.String("grpc_service", service)
}

func GrpcMethod(method string) zap.Field {
	return zap.String("grpc_method", method)
}

func VulnerabilitiesCount(n int) zap.Field {
	return zap.Int("vulnerabilities_count", n)
}

func ReportName(name string) zap.Field {
	return zap.String("report_name", name)
}

func Kind(kind string) zap.Field {
	return zap.String("k8s_kind", kind)
}

func Error(err error) zap.Field {
	return zap.Error(err) // nolint:forbidigo
}

func NumberOfTunnels(n int) zap.Field {
	return zap.Int("num_tunnels", n)
}

func NumberOfTunnelFindRequests(n int) zap.Field {
	return zap.Int("num_find_requests", n)
}
