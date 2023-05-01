package api

import (
	"context"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitlab"
)

const (
	AllowedAgentsApiPath = "/api/v4/job/allowed_agents"
)

func GetAllowedAgentsForJob(ctx context.Context, client gitlab.ClientInterface, jobToken string, opts ...gitlab.DoOption) (*AllowedAgentsForJob, error) {
	aa := &AllowedAgentsForJob{}
	err := client.Do(ctx,
		joinOpts(opts,
			gitlab.WithPath(AllowedAgentsApiPath),
			gitlab.WithJobToken(jobToken),
			gitlab.WithResponseHandler(gitlab.ProtoJsonResponseHandler(aa)),
		)...,
	)
	if err != nil {
		return nil, err
	}
	return aa, nil
}
