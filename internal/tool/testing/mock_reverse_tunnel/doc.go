package mock_reverse_tunnel

//go:generate mockgen.sh -destination "api.go" -package "mock_reverse_tunnel" "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/reverse_tunnel" "TunnelHandler,TunnelFinder,Tunnel,FindHandle"
