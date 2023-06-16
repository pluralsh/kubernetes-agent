package api

import (
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/api"
)

func (a *GetAgentInfoResponse) ToApiAgentInfo() *api.AgentInfo {
	return &api.AgentInfo{
		Id:            a.AgentId,
		ProjectId:     a.ProjectId,
		Name:          a.AgentName,
		GitalyInfo:    a.GitalyInfo.ToApiGitalyInfo(),
		Repository:    a.GitalyRepository.ToGitalyProtoRepository(),
		DefaultBranch: a.DefaultBranch,
	}
}

func (p *GetProjectInfoResponse) ToApiProjectInfo() *api.ProjectInfo {
	return &api.ProjectInfo{
		ProjectId:     p.ProjectId,
		GitalyInfo:    p.GitalyInfo.ToApiGitalyInfo(),
		Repository:    p.GitalyRepository.ToGitalyProtoRepository(),
		DefaultBranch: p.DefaultBranch,
	}
}
