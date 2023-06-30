package mock_modagent

//go:generate mockgen.sh -destination "api.go" -package "mock_modagent" "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modagent" "Api,Factory,Module"
