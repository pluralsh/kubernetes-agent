package api

import (
	"context"
	"net/url"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/gitlab"
)

const (
	ProjectInfoApiPath  = "/api/v4/internal/kubernetes/project_info"
	ProjectIdQueryParam = "id"
)

func GetProjectInfo(ctx context.Context, client gitlab.ClientInterface, agentToken api.AgentToken, projectId string, opts ...gitlab.DoOption) (*api.ProjectInfo, error) {
	response := &GetProjectInfoResponse{}
	err := client.Do(ctx,
		joinOpts(opts,
			gitlab.WithPath(ProjectInfoApiPath),
			gitlab.WithQuery(url.Values{
				ProjectIdQueryParam: []string{projectId},
			}),
			gitlab.WithAgentToken(agentToken),
			gitlab.WithResponseHandler(gitlab.ProtoJsonResponseHandler(response)),
			gitlab.WithJWT(true),
		)...,
	)
	if err != nil {
		return nil, err
	}
	return response.ToApiProjectInfo(), nil
}
