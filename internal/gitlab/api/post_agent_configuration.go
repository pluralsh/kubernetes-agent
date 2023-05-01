package api

import (
	"context"
	"net/http"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitlab"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/agentcfg"
)

const (
	AgentConfigurationApiPath = "/api/v4/internal/kubernetes/agent_configuration"
)

func PostAgentConfiguration(ctx context.Context, client gitlab.ClientInterface, agentId int64,
	config *agentcfg.ConfigurationFile, opts ...gitlab.DoOption) error {
	return client.Do(ctx,
		joinOpts(opts,
			gitlab.WithMethod(http.MethodPost),
			gitlab.WithPath(AgentConfigurationApiPath),
			gitlab.WithJWT(true),
			gitlab.WithProtoJsonRequestBody(&AgentConfigurationRequest{
				AgentId:     agentId,
				AgentConfig: config,
			}),
			gitlab.WithResponseHandler(gitlab.NoContentResponseHandler()),
		)...,
	)
}
