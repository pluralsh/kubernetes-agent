package api

import (
	"context"
	"fmt"

	"github.com/pluralsh/polly/algorithms"

	"github.com/pluralsh/kuberentes-agent/pkg/plural"
)

func AuthorizeProxyUser(ctx context.Context, token, clusterId, pluralURL string) (*AuthorizeProxyUserResponse, error) {
	client := plural.NewUnauthorized(pluralURL)
	resp, err := client.Console.TokenExchange(ctx, fmt.Sprintf("plrl:%s:%s", clusterId, token))
	if err != nil {
		return nil, err
	}

	return &AuthorizeProxyUserResponse{
		AccessAs: &AccessAsProxyAuthorization{
			AccessAs: &AccessAsProxyAuthorization_User{
				User: &AccessAsUserAuthorization{
					Groups: algorithms.Map(resp.TokenExchange.Groups, func(g *struct {
						ID   string "json:\"id\" graphql:\"id\""
						Name string "json:\"name\" graphql:\"name\""
					}) string {
						return g.Name
					}),
					Roles: algorithms.Map(resp.TokenExchange.BoundRoles, func(r *struct {
						ID   string "json:\"id\" graphql:\"id\""
						Name string "json:\"name\" graphql:\"name\""
					}) string {
						return r.Name
					}),
				},
			},
		},
		User: &User{
			Id:       resp.TokenExchange.ID,
			Username: resp.TokenExchange.Email,
			Email:    resp.TokenExchange.Email,
		},
	}, nil
}
