package agentkapp

//go:generate mockgen.sh  -destination "mock_for_test.go" -package "agentkapp" "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/cmd/agentk/agentkapp" "Runner,LeaderElector"
