package api

import (
	"context"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/gitlab"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/prototool"
)

const (
	AllowedAgentsApiPath = "/api/v4/job/allowed_agents"
)

func GetAllowedAgentsForJob(ctx context.Context, client gitlab.ClientInterface, jobToken string, opts ...gitlab.DoOption) (*AllowedAgentsForJob, error) {
	resp := &prototool.JsonBox{
		Message: &AllowedAgentsForJob{},
	}
	err := client.Do(ctx,
		joinOpts(opts,
			gitlab.WithPath(AllowedAgentsApiPath),
			gitlab.WithJobToken(jobToken),
			gitlab.WithResponseHandler(gitlab.JsonResponseHandler(resp)),
		)...,
	)
	if err != nil {
		return nil, err
	}
	return resp.Message.(*AllowedAgentsForJob), nil
}
