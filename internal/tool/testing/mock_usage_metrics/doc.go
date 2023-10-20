package mock_usage_metrics

//go:generate mockgen.sh -destination "tool.go" "github.com/pluralsh/kuberentes-agent/internal/module/usage_metrics" "UsageTrackerInterface,Counter,UniqueCounter"
