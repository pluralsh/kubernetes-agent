package reverse_tunnel

//go:generate mockgen.sh -self_package "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/reverse_tunnel" -destination "mock_for_test.go" -package "reverse_tunnel" "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/reverse_tunnel" "TunnelDataCallback"
