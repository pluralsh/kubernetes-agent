package api

import (
	"context"

	"github.com/pluralsh/kuberentes-agent/pkg/api"
	"github.com/pluralsh/kuberentes-agent/pkg/entity"
	"github.com/pluralsh/kuberentes-agent/pkg/gitlab"
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
