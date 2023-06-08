package mock_modserver

//go:generate mockgen.sh -destination "api.go" -package "mock_modserver" "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modserver" "Api,RpcApi,AgentRpcApi"
