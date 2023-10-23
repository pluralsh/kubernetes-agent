package api

import (
	"context"

	"github.com/pluralsh/kuberentes-agent/pkg/api"
	"github.com/pluralsh/kuberentes-agent/pkg/gitlab"
)

func GetAgentInfo(ctx context.Context, agentToken api.AgentToken, opts ...gitlab.DoOption) (*api.AgentInfo, error) {
	return &api.AgentInfo{
		Id:            123456,
		ProjectId:     0,
		Name:          "fake-agent",
		DefaultBranch: "",
	}, nil
}
