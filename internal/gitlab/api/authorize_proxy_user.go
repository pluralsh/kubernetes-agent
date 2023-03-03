package api

import (
	"context"
	"net/http"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/gitlab"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/prototool"
)

const (
	AuthorizeProxyUserApiPath = "/api/v4/internal/kubernetes/authorize_proxy_user"
)

type AuthorizeProxyUserRequest struct {
	AgentId    int64  `json:"agent_id"`
	AccessType string `json:"access_type"`
	AccessKey  string `json:"access_key"`
	CsrfToken  string `json:"csrf_token"`
}

func AuthorizeProxyUser(ctx context.Context, client gitlab.ClientInterface, agentId int64, accessType, accessKey, csrfToken string, opts ...gitlab.DoOption) (*AuthorizeProxyUserResponse, error) {
	auth := &AuthorizeProxyUserResponse{}
	jsonBox := &prototool.JsonBox{
		Message: auth,
	}
	err := client.Do(ctx,
		joinOpts(opts,
			gitlab.WithMethod(http.MethodPost),
			gitlab.WithPath(AuthorizeProxyUserApiPath),
			gitlab.WithJWT(true),
			gitlab.WithJsonRequestBody(&AuthorizeProxyUserRequest{
				AgentId:    agentId,
				AccessType: accessType,
				AccessKey:  accessKey,
				CsrfToken:  csrfToken,
			}),
			gitlab.WithResponseHandler(gitlab.JsonResponseHandler(jsonBox)),
		)...,
	)
	if err != nil {
		return nil, err
	}
	err = auth.ValidateAll()
	if err != nil {
		return nil, err
	}
	return auth, nil
}
