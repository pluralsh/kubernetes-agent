package plural

import (
	"context"

	"github.com/pluralsh/kuberentes-agent/internal/api"
	"github.com/pluralsh/kuberentes-agent/internal/gitlab"
	"github.com/pluralsh/kuberentes-agent/pkg/entity"
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
