package api

import (
	"context"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitlab"
)

const (
	AgentInfoApiPath = "/api/v4/internal/kubernetes/agent_info"
)

func GetAgentInfo(ctx context.Context, agentToken api.AgentToken, opts ...gitlab.DoOption) (*api.AgentInfo, error) {
	return &api.AgentInfo{
		Id:            0,
		ProjectId:     0,
		Name:          "fake-agent",
		GitalyInfo:    nil,
		Repository:    nil,
		DefaultBranch: "",
	}, nil
}
