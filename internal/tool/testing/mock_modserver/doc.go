package mock_modserver

//go:generate mockgen.sh -source "../../../module/modserver/api.go" -destination "api.go" -package "mock_modserver"

//go:generate mockgen.sh -destination "rpc_api.go" -package "mock_modserver" "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modserver" "RpcApi,AgentRpcApi"
