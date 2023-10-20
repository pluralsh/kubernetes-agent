package mock_usage_metrics

//go:generate mockgen.sh -destination "api.go" "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/usage_metrics" "UsageTrackerInterface,Counter,UniqueCounter"
