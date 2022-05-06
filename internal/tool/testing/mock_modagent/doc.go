package mock_modagent

//go:generate go run github.com/golang/mock/mockgen -destination "api.go" -package "mock_modagent" "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modagent" "Api,Factory,Module"
