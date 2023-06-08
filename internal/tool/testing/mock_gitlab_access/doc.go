// Mocks for GitLab access RPC.
package mock_gitlab_access

//go:generate mockgen.sh -destination "gitlab_access.go" -package "mock_gitlab_access" "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/gitlab_access/rpc" "GitlabAccessClient,GitlabAccess_MakeRequestClient"
