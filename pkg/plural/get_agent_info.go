package plural

import (
	"context"

	"github.com/pluralsh/kuberentes-agent/pkg/api"
)

func GetAgentInfo(ctx context.Context, agentToken api.AgentToken, pluralURL string) (*api.AgentInfo, error) {
	client := New(pluralURL, string(agentToken))
	cluster, err := client.consoleClient.MyCluster(ctx)
	if err != nil {
		return nil, err
	}

	return &api.AgentInfo{
		Id:   123456,
		ClusterId: cluster.MyCluster.ID,
		Name: cluster.MyCluster.Name,
	}, nil
}
