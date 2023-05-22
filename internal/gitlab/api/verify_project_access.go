package api

import (
	"context"
	"net/url"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitlab"
)

const (
	VerifyProjectAccessApiPath = "/api/v4/internal/kubernetes/verify_project_access"
)

func VerifyProjectAccess(ctx context.Context, client gitlab.ClientInterface, agentToken api.AgentToken, projectId string, opts ...gitlab.DoOption) (bool, error) {
	err := client.Do(ctx,
		joinOpts(opts,
			gitlab.WithPath(VerifyProjectAccessApiPath),
			gitlab.WithQuery(url.Values{
				ProjectIdQueryParam: []string{projectId},
			}),
			gitlab.WithAgentToken(agentToken),
			gitlab.WithResponseHandler(gitlab.NoContentResponseHandler()),
			gitlab.WithJWT(true),
		)...,
	)
	if err == nil {
		return true, nil
	}

	if gitlab.IsForbidden(err) || gitlab.IsUnauthorized(err) || gitlab.IsNotFound(err) {
		return false, nil
	}

	return false, err
}
