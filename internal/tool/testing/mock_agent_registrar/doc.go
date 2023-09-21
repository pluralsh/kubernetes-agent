package mock_agent_registrar

//go:generate mockgen.sh -destination "agent_registrar.go" -package "mock_agent_registrar" "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/agent_registrar/rpc" "AgentRegistrarClient"
