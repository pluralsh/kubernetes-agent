// Package mock_rpc contains mocks for gRPC and HTTP interfaces.
package mock_rpc

import "net/http"

//go:generate go run github.com/golang/mock/mockgen -destination "grpc.go" -package "mock_rpc" "google.golang.org/grpc" "ServerStream,ClientStream,ClientConnInterface,ServerTransportStream"

//go:generate go run github.com/golang/mock/mockgen -destination "gitops.go" -package "mock_rpc" "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/gitops/rpc" "GitopsClient,Gitops_GetObjectsToSynchronizeClient,Gitops_GetObjectsToSynchronizeServer,ObjectsToSynchronizeWatcherInterface"

//go:generate go run github.com/golang/mock/mockgen -destination "agent_configuration.go" -package "mock_rpc" "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/agent_configuration/rpc" "AgentConfigurationClient,AgentConfiguration_GetConfigurationClient,AgentConfiguration_GetConfigurationServer,ConfigurationWatcherInterface"

//go:generate go run github.com/golang/mock/mockgen -destination "gitlab_access.go" -package "mock_rpc" "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/gitlab_access/rpc" "GitlabAccess_MakeRequestServer"

//go:generate go run github.com/golang/mock/mockgen -destination "grpctool.go" -package "mock_rpc" "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/grpctool" "InboundGrpcToOutboundHttpStream,PoolConn,PoolInterface"

//go:generate go run github.com/golang/mock/mockgen -destination "http.go" -package "mock_rpc" "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_rpc" "ResponseWriterFlusher"

//go:generate go run github.com/golang/mock/mockgen -destination "stdlib_net.go" -package "mock_rpc" "net" "Conn"

type ResponseWriterFlusher interface {
	http.ResponseWriter
	http.Flusher
	http.Hijacker
}
