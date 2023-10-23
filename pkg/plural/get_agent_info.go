package plural

import (
	"context"

	"github.com/pluralsh/kuberentes-agent/pkg/api"
)

func GetAgentInfo(ctx context.Context, agentToken api.AgentToken, pluralURL string) (*api.AgentInfo, error) {
	_ = New(pluralURL, string(agentToken))
	return &api.AgentInfo{
		Id:            123456,
		ProjectId:     0,
		Name:          "fake-agent",
		DefaultBranch: "",
	}, nil
}
