package mock_internalgitaly

//go:generate mockgen.sh -destination "internalgitaly.go" -package "mock_internalgitaly" "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitaly" "PoolInterface,FetchVisitor,PathEntryVisitor,FileVisitor,PathFetcherInterface,PollerInterface"
