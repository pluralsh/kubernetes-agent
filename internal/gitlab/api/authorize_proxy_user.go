package api

import (
	"context"
	"net/http"

	"github.com/pluralsh/kuberentes-agent/internal/gitlab"
)

const (
	AuthorizeProxyUserApiPath = "/api/v4/internal/kubernetes/authorize_proxy_user"
)

func AuthorizeProxyUser(ctx context.Context, client gitlab.ClientInterface, agentId int64, accessType, accessKey, csrfToken string, opts ...gitlab.DoOption) (*AuthorizeProxyUserResponse, error) {
	auth := &AuthorizeProxyUserResponse{}
	err := client.Do(ctx,
		joinOpts(opts,
			gitlab.WithMethod(http.MethodPost),
			gitlab.WithPath(AuthorizeProxyUserApiPath),
			gitlab.WithJWT(true),
			gitlab.WithProtoJsonRequestBody(&AuthorizeProxyUserRequest{
				AgentId:    agentId,
				AccessType: accessType,
				AccessKey:  accessKey,
				CsrfToken:  csrfToken,
			}),
			gitlab.WithResponseHandler(gitlab.ProtoJsonResponseHandlerWithStructuredErrReason(auth)),
		)...,
	)
	if err != nil {
		return nil, err
	}
	return auth, nil
}
