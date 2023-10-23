package plural

import (
	"context"

	"github.com/pluralsh/kuberentes-agent/pkg/api"
)

func GetAgentInfo(ctx context.Context, agentToken string, pluralURL string) (*api.AgentInfo, error) {

	client := New(pluralURL, agentToken)
	client.
	return &api.AgentInfo{
		Id:            123456,
		ProjectId:     0,
		Name:          "fake-agent",
		DefaultBranch: "",
	}, nil
}
