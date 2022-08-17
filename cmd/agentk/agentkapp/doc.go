package agentkapp

//go:generate go run github.com/golang/mock/mockgen  -destination "mock_for_test.go" -package "agentkapp" "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/cmd/agentk/agentkapp" "Runner,LeaderElector"
