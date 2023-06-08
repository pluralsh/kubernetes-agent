package mock_reverse_tunnel_tracker

//go:generate mockgen.sh -destination "tracker.go" -package "mock_reverse_tunnel_tracker" "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/reverse_tunnel/tracker" "Registerer,Querier,PollingQuerier"
