package mock_modshared

//go:generate go run github.com/golang/mock/mockgen -destination "api.go" -package "mock_modshared" "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modshared" "RpcApi,Api"
