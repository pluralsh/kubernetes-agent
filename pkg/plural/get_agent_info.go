package plural

import (
	"context"
	"encoding/binary"

	uuid "github.com/satori/go.uuid"

	"github.com/pluralsh/kuberentes-agent/pkg/api"
)

func GetAgentInfo(ctx context.Context, agentToken api.AgentToken, pluralURL string) (*api.AgentInfo, error) {
	client := New(pluralURL, string(agentToken))
	cluster, err := client.consoleClient.MyCluster(ctx)
	if err != nil {
		return nil, err
	}

	u, err := uuidToInt64(cluster.MyCluster.ID)
	if err != nil {
		return nil, err
	}

	return &api.AgentInfo{
		Id:   u,
		ClusterId: cluster.MyCluster.ID,
		Name: cluster.MyCluster.Name,
	}, nil
}

func uuidToInt64(id string) (int64, error) {
	u, err := uuid.FromString(id)
	// TODO: Figure out if we can do it better.
	// uint can end up overflowing int and end up with a negative value
	return int64(binary.BigEndian.Uint64(u[:8])), err
}
