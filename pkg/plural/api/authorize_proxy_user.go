package api

import (
	"context"
	"fmt"

	"github.com/pluralsh/kuberentes-agent/pkg/plural"
)

func AuthorizeProxyUser(ctx context.Context, token, clusterId, pluralURL string) (*AuthorizeProxyUserResponse, error) {
	client := plural.NewUnauthorized(pluralURL)
	_, err := client.Console.TokenExchange(ctx, fmt.Sprintf("plrl:%s:%s", clusterId, token))
	if err != nil {
		return nil, err
	}

	return &AuthorizeProxyUserResponse{
		AccessAs:      &AccessAsProxyAuthorization{
			AccessAs:      &AccessAsProxyAuthorization_Agent{
				Agent: &AccessAsAgentAuthorization{},
			},
		},
		User: &User{},
	}, nil
}
