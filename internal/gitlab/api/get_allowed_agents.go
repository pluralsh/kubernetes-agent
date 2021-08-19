package api

import (
	"context"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/gitlab"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	AllowedAgentsApiPath = "/api/v4/job/allowed_agents"
)

// AllowedAgentsForJobAlias ensures the protojson package is used for to/from JSON marshaling.
// See https://pkg.go.dev/google.golang.org/protobuf/encoding/protojson.
type AllowedAgentsForJobAlias AllowedAgentsForJob

func (a *AllowedAgentsForJobAlias) MarshalJSON() ([]byte, error) {
	typedA := (*AllowedAgentsForJob)(a)
	return protojson.Marshal(typedA)
}

func (a *AllowedAgentsForJobAlias) UnmarshalJSON(data []byte) error {
	typedA := (*AllowedAgentsForJob)(a)
	return protojson.Unmarshal(data, typedA)
}

func GetAllowedAgentsForJob(ctx context.Context, client gitlab.ClientInterface, jobToken string) (*AllowedAgentsForJob, error) {
	ji := &AllowedAgentsForJobAlias{}
	err := client.Do(ctx,
		gitlab.WithPath(AllowedAgentsApiPath),
		gitlab.WithJobToken(jobToken),
		gitlab.WithResponseHandler(gitlab.JsonResponseHandler(ji)),
	)
	if err != nil {
		return nil, err
	}
	return (*AllowedAgentsForJob)(ji), nil
}
