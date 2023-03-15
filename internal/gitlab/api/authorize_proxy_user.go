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

func AuthorizeProxyUser(ctx context.Context, client gitlab.ClientInterface, agentId int64, accessType, accessKey, csrfToken string, opts ...gitlab.DoOption) (*AuthorizeProxyUserResponse, error) {
	req := &AuthorizeProxyUserRequest{
		AgentId:    agentId,
		AccessType: accessType,
		AccessKey:  accessKey,
		CsrfToken:  csrfToken,
	}
	err := req.ValidateAll()
	if err != nil {
		return nil, err
	}

	auth := &AuthorizeProxyUserResponse{}
	err = client.Do(ctx,
		joinOpts(opts,
			gitlab.WithMethod(http.MethodPost),
			gitlab.WithPath(AuthorizeProxyUserApiPath),
			gitlab.WithJWT(true),
			gitlab.WithJsonRequestBody(&prototool.JsonBox{Message: req}),
			gitlab.WithResponseHandler(gitlab.JsonResponseHandler(&prototool.JsonBox{Message: auth})),
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
