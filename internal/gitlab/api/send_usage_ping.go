package api

import (
	"context"
	"net/http"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/gitlab"
)

const (
	UsagePingApiPath = "/api/v4/internal/kubernetes/usage_metrics"
)

func SendUsagePing(ctx context.Context, client gitlab.ClientInterface, counters map[string]int64, opts ...gitlab.DoOption) error {
	return client.Do(ctx,
		joinOpts(opts,
			gitlab.WithMethod(http.MethodPost),
			gitlab.WithPath(UsagePingApiPath),
			gitlab.WithJsonRequestBody(counters),
			gitlab.WithResponseHandler(gitlab.NoContentResponseHandler()),
			gitlab.WithJWT(true),
		)...,
	)
}
