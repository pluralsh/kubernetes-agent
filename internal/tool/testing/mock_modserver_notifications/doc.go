package mock_modserver_notifications

//go:generate go run github.com/golang/mock/mockgen -destination "api.go" -package "mock_modserver_notifications" "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modserver/notifications" "Subscriber"
