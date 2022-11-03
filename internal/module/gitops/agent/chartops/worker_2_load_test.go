package chartops

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/pkg/agentcfg"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestLoadValues(t *testing.T) {
	inline0, err := structpb.NewStruct(map[string]interface{}{
		"hello": "world",
		"some": map[string]interface{}{
			"nested": "value",
		},
		"leave": "me alone",
	})
	require.NoError(t, err)
	inline1, err := structpb.NewStruct(map[string]interface{}{
		"hello": "overworld",
		"some": map[string]interface{}{
			"nested": "overvalue",
		},
		"include": "me",
	})
	require.NoError(t, err)
	valuesCF := []*agentcfg.ChartValuesCF{
		{
			As: &agentcfg.ChartValuesCF_Inline{
				Inline: inline0,
			},
		},
		{
			As: &agentcfg.ChartValuesCF_Inline{
				Inline: inline1,
			},
		},
	}
	values, err := loadValues(valuesCF)
	require.NoError(t, err)

	require.Equal(t, values, map[string]interface{}{
		"hello": "overworld",
		"some": map[string]interface{}{
			"nested": "overvalue",
		},
		"leave":   "me alone",
		"include": "me",
	})
}
