package api

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/testing/protocmp"
)

var (
	_ json.Marshaler   = (*AllowedAgentsForJobAlias)(nil)
	_ json.Unmarshaler = (*AllowedAgentsForJobAlias)(nil)
)

func TestGetAllowedAgents_JSON(t *testing.T) {
	expected := &AllowedAgentsForJob{
		AllowedAgents: []*AllowedAgent{
			{
				Id: 123,
				ConfigProject: &ConfigProject{
					Id: 234,
				},
				Configuration: &Configuration{
					DefaultNamespace: "abc",
				},
			},
			{
				Id: 1,
				ConfigProject: &ConfigProject{
					Id: 2,
				},
				Configuration: &Configuration{
					DefaultNamespace: "", // empty
				},
			},
		},
		Job: &Job{
			Id: 32,
		},
		Pipeline: &Pipeline{
			Id: 3232,
		},
		Project: &Project{
			Id: 44,
		},
		User: &User{
			Id:       2323,
			Username: "abc",
		},
	}
	data, err := json.Marshal(expected)
	require.NoError(t, err)

	actual := &AllowedAgentsForJob{}
	err = json.Unmarshal(data, actual)
	require.NoError(t, err)

	assert.Empty(t, cmp.Diff(expected, actual, cmp.Transformer("AllowedAgentsForJob", transformAlias), protocmp.Transform()))
}

func transformAlias(val *AllowedAgentsForJobAlias) *AllowedAgentsForJob {
	return (*AllowedAgentsForJob)(val)
}
