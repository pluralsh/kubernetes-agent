package api

import (
	"context"
	"net/http"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/gitlab"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/prototool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/pkg/agentcfg"
)

const (
	AgentConfigurationApiPath = "/api/v4/internal/kubernetes/agent_configuration"
)

type agentConfigurationRequest struct {
	AgentId     int64             `json:"agent_id"`
	AgentConfig prototool.JsonBox `json:"agent_config"`
}

func PostAgentConfiguration(ctx context.Context, client gitlab.ClientInterface, agentId int64, config *agentcfg.ConfigurationFile) error {
	return client.Do(ctx,
		gitlab.WithMethod(http.MethodPost),
		gitlab.WithPath(AgentConfigurationApiPath),
		gitlab.WithJWT(true),
		gitlab.WithJsonRequestBody(&agentConfigurationRequest{
			AgentId:     agentId,
			AgentConfig: prototool.JsonBox{Message: config},
		}),
		gitlab.WithResponseHandler(gitlab.NoContentResponseHandler()),
	)
}