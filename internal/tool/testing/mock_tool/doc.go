package mock_tool

//go:generate mockgen.sh -destination "tool.go" -package "mock_tool" "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/errz" "ErrReporter"
