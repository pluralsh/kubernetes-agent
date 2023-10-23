package api

import (
	"context"
	"net/http"

	gitlab2 "github.com/pluralsh/kuberentes-agent/pkg/gitlab"
)

const (
	UsagePingApiPath = "/api/v4/internal/kubernetes/usage_metrics"
)

type UsagePingData struct {
	Counters       map[string]int64   `json:"counters,omitempty"`
	UniqueCounters map[string][]int64 `json:"unique_counters,omitempty"`
}

func SendUsagePing(ctx context.Context, client gitlab2.ClientInterface, data UsagePingData, opts ...gitlab2.DoOption) error {
	return client.Do(ctx,
		joinOpts(opts,
			gitlab2.WithMethod(http.MethodPost),
			gitlab2.WithPath(UsagePingApiPath),
			gitlab2.WithJsonRequestBody(data),
			gitlab2.WithResponseHandler(gitlab2.NoContentResponseHandler()),
			gitlab2.WithJWT(true),
		)...,
	)
}
