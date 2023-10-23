package plural

import (
	"context"

	"github.com/pluralsh/kuberentes-agent/internal/api"
	"github.com/pluralsh/kuberentes-agent/pkg/plural"
)

func GetAgentInfo(ctx context.Context, agentToken string, pluralURL string) (*api.AgentInfo, error) {

	client := plural.New(pluralURL, agentToken)
	client.
	return &api.AgentInfo{
		Id:            123456,
		ProjectId:     0,
		Name:          "fake-agent",
		DefaultBranch: "",
	}, nil
}
