package api

import (
	"context"
	"net/http"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/gitlab"
)

const (
	UsagePingApiPath = "/api/v4/internal/kubernetes/usage_metrics"
)

type UsagePingData struct {
	Counters       map[string]int64   `json:"counters,omitempty"`
	UniqueCounters map[string][]int64 `json:"unique_counters,omitempty"`
}

func SendUsagePing(ctx context.Context, client gitlab.ClientInterface, data UsagePingData, opts ...gitlab.DoOption) error {
	return client.Do(ctx,
		joinOpts(opts,
			gitlab.WithMethod(http.MethodPost),
			gitlab.WithPath(UsagePingApiPath),
			gitlab.WithJsonRequestBody(data),
			gitlab.WithResponseHandler(gitlab.NoContentResponseHandler()),
			gitlab.WithJWT(true),
		)...,
	)
}
