package mock_cache

//go:generate go run github.com/golang/mock/mockgen -destination "cache.go" -package "mock_cache" "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/cache" "ErrCacher"
