package api

import (
	"context"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/gitlab"
)

const (
	AgentInfoApiPath = "/api/v4/internal/kubernetes/agent_info"
)

type getAgentInfoResponse struct {
	ProjectId        int64                   `json:"project_id"`
	AgentId          int64                   `json:"agent_id"`
	AgentName        string                  `json:"agent_name"`
	GitalyInfo       gitlab.GitalyInfo       `json:"gitaly_info"`
	GitalyRepository gitlab.GitalyRepository `json:"gitaly_repository"`
	DefaultBranch    string                  `json:"default_branch"`
}

func GetAgentInfo(ctx context.Context, client gitlab.ClientInterface, agentToken api.AgentToken, opts ...gitlab.DoOption) (*api.AgentInfo, error) {
	response := getAgentInfoResponse{}
	err := client.Do(ctx,
		joinOpts(opts,
			gitlab.WithPath(AgentInfoApiPath),
			gitlab.WithAgentToken(agentToken),
			gitlab.WithJWT(true),
			gitlab.WithResponseHandler(gitlab.JsonResponseHandler(&response)),
		)...,
	)
	if err != nil {
		return nil, err
	}
	return &api.AgentInfo{
		Id:            response.AgentId,
		ProjectId:     response.ProjectId,
		Name:          response.AgentName,
		GitalyInfo:    response.GitalyInfo.ToGitalyInfo(),
		Repository:    response.GitalyRepository.ToProtoRepository(),
		DefaultBranch: response.DefaultBranch,
	}, nil
}
