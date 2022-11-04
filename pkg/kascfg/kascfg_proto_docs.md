# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [pkg/kascfg/kascfg.proto](#pkg_kascfg_kascfg-proto)
    - [AgentCF](#gitlab-agent-kascfg-AgentCF)
    - [AgentConfigurationCF](#gitlab-agent-kascfg-AgentConfigurationCF)
    - [ApiCF](#gitlab-agent-kascfg-ApiCF)
    - [ConfigurationFile](#gitlab-agent-kascfg-ConfigurationFile)
    - [GitLabCF](#gitlab-agent-kascfg-GitLabCF)
    - [GitalyCF](#gitlab-agent-kascfg-GitalyCF)
    - [GitopsCF](#gitlab-agent-kascfg-GitopsCF)
    - [GoogleProfilerCF](#gitlab-agent-kascfg-GoogleProfilerCF)
    - [KubernetesApiCF](#gitlab-agent-kascfg-KubernetesApiCF)
    - [ListenAgentCF](#gitlab-agent-kascfg-ListenAgentCF)
    - [ListenApiCF](#gitlab-agent-kascfg-ListenApiCF)
    - [ListenKubernetesApiCF](#gitlab-agent-kascfg-ListenKubernetesApiCF)
    - [ListenPrivateApiCF](#gitlab-agent-kascfg-ListenPrivateApiCF)
    - [LivenessProbeCF](#gitlab-agent-kascfg-LivenessProbeCF)
    - [LoggingCF](#gitlab-agent-kascfg-LoggingCF)
    - [ObservabilityCF](#gitlab-agent-kascfg-ObservabilityCF)
    - [ObservabilityListenCF](#gitlab-agent-kascfg-ObservabilityListenCF)
    - [PrivateApiCF](#gitlab-agent-kascfg-PrivateApiCF)
    - [PrometheusCF](#gitlab-agent-kascfg-PrometheusCF)
    - [ReadinessProbeCF](#gitlab-agent-kascfg-ReadinessProbeCF)
    - [RedisCF](#gitlab-agent-kascfg-RedisCF)
    - [RedisSentinelCF](#gitlab-agent-kascfg-RedisSentinelCF)
    - [RedisServerCF](#gitlab-agent-kascfg-RedisServerCF)
    - [RedisTLSCF](#gitlab-agent-kascfg-RedisTLSCF)
    - [SentryCF](#gitlab-agent-kascfg-SentryCF)
    - [TokenBucketRateLimitCF](#gitlab-agent-kascfg-TokenBucketRateLimitCF)
    - [TracingCF](#gitlab-agent-kascfg-TracingCF)
  
    - [log_level_enum](#gitlab-agent-kascfg-log_level_enum)
  
- [Scalar Value Types](#scalar-value-types)



<a name="pkg_kascfg_kascfg-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## pkg/kascfg/kascfg.proto



<a name="gitlab-agent-kascfg-AgentCF"></a>

### AgentCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| listen | [ListenAgentCF](#gitlab-agent-kascfg-ListenAgentCF) |  | RPC listener configuration for agentk connections. |
| configuration | [AgentConfigurationCF](#gitlab-agent-kascfg-AgentConfigurationCF) |  | Configuration for agent&#39;s configuration repository. |
| gitops | [GitopsCF](#gitlab-agent-kascfg-GitopsCF) |  | Configuration for GitOps. |
| info_cache_ttl | [google.protobuf.Duration](#google-protobuf-Duration) |  | TTL for successful agent info lookups. /api/v4/internal/kubernetes/agent_info Set to zero to disable. |
| info_cache_error_ttl | [google.protobuf.Duration](#google-protobuf-Duration) |  | TTL for failed agent info lookups. /api/v4/internal/kubernetes/agent_info |
| redis_conn_info_ttl | [google.protobuf.Duration](#google-protobuf-Duration) |  | TTL for information about connected agents, stored in Redis. |
| redis_conn_info_refresh | [google.protobuf.Duration](#google-protobuf-Duration) |  | Refresh period for information about connected agents, stored in Redis. |
| redis_conn_info_gc | [google.protobuf.Duration](#google-protobuf-Duration) |  | Garbage collection period for information about connected agents, stored in Redis. If gitlab-kas crashes, another gitlab-kas instance will clean up stale data. This is how often this cleanup runs. |
| kubernetes_api | [KubernetesApiCF](#gitlab-agent-kascfg-KubernetesApiCF) |  | Configuration for exposing Kubernetes API. |






<a name="gitlab-agent-kascfg-AgentConfigurationCF"></a>

### AgentConfigurationCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| poll_period | [google.protobuf.Duration](#google-protobuf-Duration) |  | How often to poll agent&#39;s configuration repository for changes. |
| max_configuration_file_size | [uint32](#uint32) |  | Maximum file size of the agent configuration file. |






<a name="gitlab-agent-kascfg-ApiCF"></a>

### ApiCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| listen | [ListenApiCF](#gitlab-agent-kascfg-ListenApiCF) |  | RPC listener configuration for API connections. |






<a name="gitlab-agent-kascfg-ConfigurationFile"></a>

### ConfigurationFile
ConfigurationFile represents kas configuration file.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| gitlab | [GitLabCF](#gitlab-agent-kascfg-GitLabCF) |  | Configuration related to interaction with GitLab. |
| agent | [AgentCF](#gitlab-agent-kascfg-AgentCF) |  | Configuration related to the agent. Generally all configuration for user-facing features should be here. |
| observability | [ObservabilityCF](#gitlab-agent-kascfg-ObservabilityCF) |  | Configuration related to all things observability: metrics, tracing, monitoring, logging, usage metrics, profiling. |
| gitaly | [GitalyCF](#gitlab-agent-kascfg-GitalyCF) |  | Configuration related to interaction with Gitaly. |
| redis | [RedisCF](#gitlab-agent-kascfg-RedisCF) |  | Redis configurations available to kas. |
| api | [ApiCF](#gitlab-agent-kascfg-ApiCF) |  | Public API. |
| private_api | [PrivateApiCF](#gitlab-agent-kascfg-PrivateApiCF) |  | Private API for kas-&gt;kas communication. |






<a name="gitlab-agent-kascfg-GitLabCF"></a>

### GitLabCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| address | [string](#string) |  |  |
| authentication_secret_file | [string](#string) |  | Secret to generate JWT tokens to authenticate with GitLab. |
| ca_certificate_file | [string](#string) |  | Optional X.509 CA certificate for TLS in PEM format. Should be set for self-signed certificates. |
| api_rate_limit | [TokenBucketRateLimitCF](#gitlab-agent-kascfg-TokenBucketRateLimitCF) |  | Rate limiting configuration for talking to the GitLab API. |






<a name="gitlab-agent-kascfg-GitalyCF"></a>

### GitalyCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| global_api_rate_limit | [TokenBucketRateLimitCF](#gitlab-agent-kascfg-TokenBucketRateLimitCF) |  | Rate limit that is enforced across all Gitaly servers. |
| per_server_api_rate_limit | [TokenBucketRateLimitCF](#gitlab-agent-kascfg-TokenBucketRateLimitCF) |  | Rate limit that is enforced per each Gitaly server. |






<a name="gitlab-agent-kascfg-GitopsCF"></a>

### GitopsCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| poll_period | [google.protobuf.Duration](#google-protobuf-Duration) |  | How often to poll GitOps manifest repositories for changes. |
| project_info_cache_ttl | [google.protobuf.Duration](#google-protobuf-Duration) |  | TTL for successful project info lookups. /api/v4/internal/kubernetes/project_info Set to zero to disable. |
| project_info_cache_error_ttl | [google.protobuf.Duration](#google-protobuf-Duration) |  | TTL for failed project info lookups. /api/v4/internal/kubernetes/project_info |
| max_manifest_file_size | [uint32](#uint32) |  | Maximum size of a GitOps manifest file. |
| max_total_manifest_file_size | [uint32](#uint32) |  | Maximum total size of all GitOps manifest files per GitOps project. |
| max_number_of_paths | [uint32](#uint32) |  | Maximum number of scanned paths per GitOps project. |
| max_number_of_files | [uint32](#uint32) |  | Maximum number of scanned files across all paths per GitOps project. This limit ensures there are not too many files in the repository that we need to sift though to find *.yaml, *.yml, *.json files. All files and directories under a path are counted towards this limit. |






<a name="gitlab-agent-kascfg-GoogleProfilerCF"></a>

### GoogleProfilerCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| enabled | [bool](#bool) |  |  |
| project_id | [string](#string) |  |  |
| credentials_file | [string](#string) |  |  |
| debug_logging | [bool](#bool) |  |  |






<a name="gitlab-agent-kascfg-KubernetesApiCF"></a>

### KubernetesApiCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| listen | [ListenKubernetesApiCF](#gitlab-agent-kascfg-ListenKubernetesApiCF) |  | HTTP listener configuration for Kubernetes API connections. |
| url_path_prefix | [string](#string) |  | URL path prefix to remove from the incoming request URL. Should be `/` if no prefix trimming is needed. |
| allowed_agent_cache_ttl | [google.protobuf.Duration](#google-protobuf-Duration) |  | TTL for successful allowed agent lookups. /api/v4/job/allowed_agents Set to zero to disable. |
| allowed_agent_cache_error_ttl | [google.protobuf.Duration](#google-protobuf-Duration) |  | TTL for failed allowed agent lookups. /api/v4/job/allowed_agents |






<a name="gitlab-agent-kascfg-ListenAgentCF"></a>

### ListenAgentCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| network | [string](#string) | optional | Network type to listen on. Supported values: tcp, tcp4, tcp6, unix. |
| address | [string](#string) |  | Address to listen on. |
| websocket | [bool](#bool) |  | Enable &#34;gRPC through WebSocket&#34; listening mode. Rather than expecting gRPC directly, expect a WebSocket connection, from which a gRPC stream is then unpacked. |
| certificate_file | [string](#string) |  | X.509 certificate for TLS in PEM format. TLS is enabled iff both certificate_file and key_file are provided. |
| key_file | [string](#string) |  | X.509 key file for TLS in PEM format. TLS is enabled iff both certificate_file and key_file are provided. |
| connections_per_token_per_minute | [uint32](#uint32) |  | Maximum number of connections to allow per agent token per minute. |
| max_connection_age | [google.protobuf.Duration](#google-protobuf-Duration) |  | Max age of a connection. Connection is closed gracefully once it&#39;s too old and there is no streaming happening. |
| listen_grace_period | [google.protobuf.Duration](#google-protobuf-Duration) |  | How much time to wait before stopping accepting new connections on shutdown. |






<a name="gitlab-agent-kascfg-ListenApiCF"></a>

### ListenApiCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| network | [string](#string) | optional | Network type to listen on. Supported values: tcp, tcp4, tcp6, unix. |
| address | [string](#string) |  | Address to listen on. |
| authentication_secret_file | [string](#string) |  | Secret to verify JWT tokens. |
| certificate_file | [string](#string) |  | X.509 certificate for TLS in PEM format. TLS is enabled iff both certificate_file and key_file are provided. |
| key_file | [string](#string) |  | X.509 key file for TLS in PEM format. TLS is enabled iff both certificate_file and key_file are provided. |
| max_connection_age | [google.protobuf.Duration](#google-protobuf-Duration) |  | Max age of a connection. Connection is closed gracefully once it&#39;s too old and there is no streaming happening. |
| listen_grace_period | [google.protobuf.Duration](#google-protobuf-Duration) |  | How much time to wait before stopping accepting new connections on shutdown. |






<a name="gitlab-agent-kascfg-ListenKubernetesApiCF"></a>

### ListenKubernetesApiCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| network | [string](#string) | optional | Network type to listen on. Supported values: tcp, tcp4, tcp6, unix. |
| address | [string](#string) |  | Address to listen on. |
| certificate_file | [string](#string) |  | X.509 certificate for TLS in PEM format. TLS is enabled iff both certificate_file and key_file are provided. |
| key_file | [string](#string) |  | X.509 key file for TLS in PEM format. TLS is enabled iff both certificate_file and key_file are provided. |
| listen_grace_period | [google.protobuf.Duration](#google-protobuf-Duration) |  | How much time to wait before stopping accepting new connections on shutdown. |






<a name="gitlab-agent-kascfg-ListenPrivateApiCF"></a>

### ListenPrivateApiCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| network | [string](#string) | optional | Network type to listen on. Supported values: tcp, tcp4, tcp6, unix. |
| address | [string](#string) |  | Address to listen on. |
| authentication_secret_file | [string](#string) |  | Secret to verify JWT tokens. |
| certificate_file | [string](#string) |  | X.509 certificate for TLS in PEM format. TLS is enabled iff both certificate_file and key_file are provided. |
| key_file | [string](#string) |  | X.509 key file for TLS in PEM format. TLS is enabled iff both certificate_file and key_file are provided. |
| max_connection_age | [google.protobuf.Duration](#google-protobuf-Duration) |  | Max age of a connection. Connection is closed gracefully once it&#39;s too old and there is no streaming happening. |
| ca_certificate_file | [string](#string) |  | Optional X.509 CA certificate for TLS in PEM format. Should be set for self-signed certificates. |
| listen_grace_period | [google.protobuf.Duration](#google-protobuf-Duration) |  | How much time to wait before stopping accepting new connections on shutdown. |






<a name="gitlab-agent-kascfg-LivenessProbeCF"></a>

### LivenessProbeCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| url_path | [string](#string) |  | Expected URL path for requests. |






<a name="gitlab-agent-kascfg-LoggingCF"></a>

### LoggingCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| level | [log_level_enum](#gitlab-agent-kascfg-log_level_enum) |  |  |
| grpc_level | [log_level_enum](#gitlab-agent-kascfg-log_level_enum) | optional | optional to be able to tell when not set and use a different default value. |






<a name="gitlab-agent-kascfg-ObservabilityCF"></a>

### ObservabilityCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| usage_reporting_period | [google.protobuf.Duration](#google-protobuf-Duration) |  | How often to send usage metrics to the main application. /api/v4/internal/kubernetes/usage_ping Set to zero to disable. |
| listen | [ObservabilityListenCF](#gitlab-agent-kascfg-ObservabilityListenCF) |  | Listener configuration for HTTP endpoint that exposes Prometheus, pprof, liveness and readiness probes. |
| prometheus | [PrometheusCF](#gitlab-agent-kascfg-PrometheusCF) |  |  |
| tracing | [TracingCF](#gitlab-agent-kascfg-TracingCF) |  |  |
| sentry | [SentryCF](#gitlab-agent-kascfg-SentryCF) |  |  |
| logging | [LoggingCF](#gitlab-agent-kascfg-LoggingCF) |  |  |
| google_profiler | [GoogleProfilerCF](#gitlab-agent-kascfg-GoogleProfilerCF) |  | Configuration for the Google Cloud Profiler. See https://pkg.go.dev/cloud.google.com/go/profiler. |
| liveness_probe | [LivenessProbeCF](#gitlab-agent-kascfg-LivenessProbeCF) |  |  |
| readiness_probe | [ReadinessProbeCF](#gitlab-agent-kascfg-ReadinessProbeCF) |  |  |






<a name="gitlab-agent-kascfg-ObservabilityListenCF"></a>

### ObservabilityListenCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| network | [string](#string) | optional | Network type to listen on. Supported values: tcp, tcp4, tcp6, unix. |
| address | [string](#string) |  | Address to listen on. |
| certificate_file | [string](#string) | optional | X.509 certificate for TLS in PEM format. TLS is enabled iff both certificate_file and key_file are provided. |
| key_file | [string](#string) | optional | X.509 key file for TLS in PEM format. TLS is enabled iff both certificate_file and key_file are provided. |






<a name="gitlab-agent-kascfg-PrivateApiCF"></a>

### PrivateApiCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| listen | [ListenPrivateApiCF](#gitlab-agent-kascfg-ListenPrivateApiCF) |  | RPC listener configuration for API connections. |






<a name="gitlab-agent-kascfg-PrometheusCF"></a>

### PrometheusCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| url_path | [string](#string) |  | Expected URL path for requests. |






<a name="gitlab-agent-kascfg-ReadinessProbeCF"></a>

### ReadinessProbeCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| url_path | [string](#string) |  | Expected URL path for requests. |






<a name="gitlab-agent-kascfg-RedisCF"></a>

### RedisCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| server | [RedisServerCF](#gitlab-agent-kascfg-RedisServerCF) |  | Single-server Redis. |
| sentinel | [RedisSentinelCF](#gitlab-agent-kascfg-RedisSentinelCF) |  | Redis with Sentinel setup. See http://redis.io/topics/sentinel. |
| pool_size | [uint32](#uint32) |  | The max number of connections. |
| dial_timeout | [google.protobuf.Duration](#google-protobuf-Duration) |  | Dial timeout. |
| read_timeout | [google.protobuf.Duration](#google-protobuf-Duration) |  | Read timeout. |
| write_timeout | [google.protobuf.Duration](#google-protobuf-Duration) |  | Write timeout. |
| idle_timeout | [google.protobuf.Duration](#google-protobuf-Duration) |  | How long to keep TCP connections alive before closing. |
| key_prefix | [string](#string) |  | Key prefix for everything gitlab-kas stores in Redis. |
| username | [string](#string) |  | Use the specified Username to authenticate the current connection with one of the connections defined in the ACL list when connecting to a Redis 6.0 instance, or greater, that is using the Redis ACL system. |
| password_file | [string](#string) |  | Optional password. Must match the password specified in the requirepass server configuration option (if connecting to a Redis 5.0 instance, or lower), or the User Password when connecting to a Redis 6.0 instance, or greater, that is using the Redis ACL system. |
| network | [string](#string) |  | The network type, either tcp or unix. Default is tcp. |
| tls | [RedisTLSCF](#gitlab-agent-kascfg-RedisTLSCF) |  |  |






<a name="gitlab-agent-kascfg-RedisSentinelCF"></a>

### RedisSentinelCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| master_name | [string](#string) |  | The name of the sentinel master. |
| addresses | [string](#string) | repeated | The host:port addresses of the sentinels. |
| sentinel_password_file | [string](#string) |  | Sentinel password from &#34;requirepass &lt;password&gt;&#34; (if enabled) in Sentinel configuration |






<a name="gitlab-agent-kascfg-RedisServerCF"></a>

### RedisServerCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| address | [string](#string) |  | The host:port address of the node. |






<a name="gitlab-agent-kascfg-RedisTLSCF"></a>

### RedisTLSCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| enabled | [bool](#bool) |  | If true, uses TLS for the redis connection (only available if network is &#34;tcp&#34;) |
| certificate_file | [string](#string) |  | For mutual TLS, specify both certificate_file and key_file; otherwise, specify neither Optional custom X.509 certificate file for TLS in PEM format |
| key_file | [string](#string) |  | Optional custom X.509 key file for TLS in PEM format |
| ca_certificate_file | [string](#string) |  | Optional custom X.509 root CA file in PEM format, used to validate the Redis server&#39;s certificate (e.g. if the server has a self-signed certificate) |






<a name="gitlab-agent-kascfg-SentryCF"></a>

### SentryCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| dsn | [string](#string) |  | Sentry DSN https://docs.sentry.io/platforms/go/#configure |
| environment | [string](#string) |  | Sentry environment https://docs.sentry.io/product/sentry-basics/environments/ |






<a name="gitlab-agent-kascfg-TokenBucketRateLimitCF"></a>

### TokenBucketRateLimitCF
See https://pkg.go.dev/golang.org/x/time/rate#Limiter.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| refill_rate_per_second | [double](#double) |  | Number of events per second. A zero allows no events. How fast the &#34;token bucket&#34; is refilled. |
| bucket_size | [uint32](#uint32) |  | Maximum number of events that are allowed to happen in succession. Size of the &#34;token bucket&#34;. |






<a name="gitlab-agent-kascfg-TracingCF"></a>

### TracingCF



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| otlp_endpoint | [string](#string) | optional | URL to send traces to. Supported protocols are: http, https. Traces are protobuf encoded. Example: http://localhost:4317 |





 


<a name="gitlab-agent-kascfg-log_level_enum"></a>

### log_level_enum


| Name | Number | Description |
| ---- | ------ | ----------- |
| info | 0 | default value must be 0 |
| debug | 1 |  |
| warn | 2 |  |
| error | 3 |  |


 

 

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

