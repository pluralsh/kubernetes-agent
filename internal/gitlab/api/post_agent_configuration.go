package api

import (
	"context"
	"net/http"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/gitlab"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/prototool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/pkg/agentcfg"
)

const (
	AgentConfigurationApiPath = "/api/v4/internal/kubernetes/agent_configuration"
)

func PostAgentConfiguration(ctx context.Context, client gitlab.ClientInterface, agentId int64,
	config *agentcfg.ConfigurationFile, opts ...gitlab.DoOption) error {
	req := &AgentConfigurationRequest{
		AgentId:     agentId,
		AgentConfig: config,
	}
	err := req.ValidateAll()
	if err != nil {
		return err
	}
	err = client.Do(ctx,
		joinOpts(opts,
			gitlab.WithMethod(http.MethodPost),
			gitlab.WithPath(AgentConfigurationApiPath),
			gitlab.WithJWT(true),
			gitlab.WithJsonRequestBody(&prototool.JsonBox{Message: req}),
			gitlab.WithResponseHandler(gitlab.NoContentResponseHandler()),
		)...,
	)
	if err != nil {
		return err
	}
	return nil
}
