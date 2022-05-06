package api

import (
	"context"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/gitlab"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/prototool"
)

const (
	AllowedAgentsApiPath = "/api/v4/job/allowed_agents"
)

func GetAllowedAgentsForJob(ctx context.Context, client gitlab.ClientInterface, jobToken string, opts ...gitlab.DoOption) (*AllowedAgentsForJob, error) {
	aa := &AllowedAgentsForJob{}
	resp := &prototool.JsonBox{
		Message: aa,
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
	err = aa.Validate()
	if err != nil {
		return nil, err
	}
	return aa, nil
}
