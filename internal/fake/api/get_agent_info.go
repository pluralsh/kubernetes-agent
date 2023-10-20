package api

import (
	"context"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/entity"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitlab"
)

const (
	AgentInfoApiPath = "/api/v4/internal/kubernetes/agent_info"
)

func GetAgentInfo(ctx context.Context, agentToken api.AgentToken, opts ...gitlab.DoOption) (*api.AgentInfo, error) {
	return &api.AgentInfo{
		Id:            123456,
		ProjectId:     0,
		Name:          "fake-agent",
		DefaultBranch: "",
		GitalyInfo: &entity.GitalyInfo{
			Address:  "127.0.0.1",
			Token:    "123",
			Features: nil,
		},
	}, nil
}
