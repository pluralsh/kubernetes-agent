package agent

//go:generate go run github.com/golang/mock/mockgen -self_package "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/gitops/agent" -destination "mock_for_test.go" -package "agent" "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/gitops/agent" "WorkerFactory,Worker"
