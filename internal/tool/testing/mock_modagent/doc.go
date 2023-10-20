package mock_modagent

//go:generate mockgen.sh -destination "api.go" -package "mock_modagent" "github.com/pluralsh/kuberentes-agent/internal/module/modagent" "Api,Factory,Module"
