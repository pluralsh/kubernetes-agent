package agent_tracker

import (
	"testing"
	"time"

	"github.com/pluralsh/kuberentes-agent/internal/tool/redistool"
	"github.com/pluralsh/kuberentes-agent/pkg/entity"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestConnectedAgentInfoSize(t *testing.T) {
	infoBytes, err := proto.Marshal(&ConnectedAgentInfo{
		AgentMeta: &entity.AgentMeta{
			Version:      "v1.0.0",
			CommitId:     "f500e3e",
			PodNamespace: "gitlab-agent",
			PodName:      "agentk-g7x6j",
		},
		ConnectedAt:  timestamppb.Now(),
		ConnectionId: 1231232,
		AgentId:      123123,
		ProjectId:    3232323,
	})
	require.NoError(t, err)
	data, err := proto.Marshal(&redistool.ExpiringValue{
		ExpiresAt: time.Now().Unix(),
		Value:     infoBytes,
	})
	require.NoError(t, err)
	t.Log(len(data))
}
