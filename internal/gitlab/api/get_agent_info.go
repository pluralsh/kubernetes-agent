package api

import (
	"context"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/gitlab"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/prototool"
)

const (
	AgentInfoApiPath = "/api/v4/internal/kubernetes/agent_info"
)

func GetAgentInfo(ctx context.Context, client gitlab.ClientInterface, agentToken api.AgentToken, opts ...gitlab.DoOption) (*api.AgentInfo, error) {
	response := &GetAgentInfoResponse{}
	err := client.Do(ctx,
		joinOpts(opts,
			gitlab.WithPath(AgentInfoApiPath),
			gitlab.WithAgentToken(agentToken),
			gitlab.WithJWT(true),
			gitlab.WithResponseHandler(gitlab.JsonResponseHandler(&prototool.JsonBox{Message: response})),
		)...,
	)
	if err != nil {
		return nil, err
	}
	err = response.ValidateAll()
	if err != nil {
		return nil, err
	}
	return response.ToApiAgentInfo(), nil
}
