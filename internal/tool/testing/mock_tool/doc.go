package mock_tool

//go:generate go run github.com/golang/mock/mockgen -destination "tool.go" -package "mock_tool" "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/errz" "ErrReporter"
