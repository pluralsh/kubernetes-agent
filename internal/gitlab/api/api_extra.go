package api

import (
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitaly/vendored/gitalypb"
)

func (g *GitalyInfo) ToApiGitalyInfo() api.GitalyInfo {
	return api.GitalyInfo{
		Address:  g.Address,
		Token:    g.Token,
		Features: g.Features,
	}
}

func (r *GitalyRepository) ToGitalyProtoRepository() *gitalypb.Repository {
	return &gitalypb.Repository{
		StorageName:   r.StorageName,
		RelativePath:  r.RelativePath,
		GlRepository:  r.GlRepository,
		GlProjectPath: r.GlProjectPath,
	}
}

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
