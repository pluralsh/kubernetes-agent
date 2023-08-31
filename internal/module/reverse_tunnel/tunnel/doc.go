package tunnel

//go:generate mockgen.sh -self_package "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/reverse_tunnel/tunnel" -destination "mock_for_test.go" -package "tunnel" "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/reverse_tunnel/tunnel" "DataCallback,Querier,Tracker"
