package mock_modshared

//go:generate mockgen.sh -destination "api.go" -package "mock_modshared" "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modshared" "RpcApi,Api"
