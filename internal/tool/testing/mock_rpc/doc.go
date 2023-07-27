// Package mock_rpc contains mocks for gRPC interfaces.
package mock_rpc

//go:generate mockgen.sh -destination "grpc.go" -package "mock_rpc" "google.golang.org/grpc" "ServerStream,ClientStream,ClientConnInterface,ServerTransportStream"

//go:generate mockgen.sh -destination "gitops.go" -package "mock_rpc" "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/gitops/rpc" "GitopsClient,Gitops_GetObjectsToSynchronizeClient,Gitops_GetObjectsToSynchronizeServer,ObjectsToSynchronizeWatcherInterface"

//go:generate mockgen.sh -destination "agent_configuration.go" -package "mock_rpc" "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/agent_configuration/rpc" "AgentConfigurationClient,AgentConfiguration_GetConfigurationClient,AgentConfiguration_GetConfigurationServer,ConfigurationWatcherInterface"

//go:generate mockgen.sh -destination "gitlab_access.go" -package "mock_rpc" "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/gitlab_access/rpc" "GitlabAccess_MakeRequestServer"

//go:generate mockgen.sh -destination "grpctool.go" -package "mock_rpc" "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/grpctool" "InboundGrpcToOutboundHttpStream,PoolConn,PoolInterface,ServerErrorReporter"
