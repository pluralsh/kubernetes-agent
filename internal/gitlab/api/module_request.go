package api

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/gitlab"
)

const (
	ModuleRequestApiPath = "/api/v4/internal/kubernetes/modules/"
)

func MakeModuleRequest(ctx context.Context, client gitlab.ClientInterface, agentToken api.AgentToken,
	moduleName, method, urlPath string, query url.Values, header http.Header, body io.Reader,
	opts ...gitlab.DoOption) (*http.Response, error) {
	var resp *http.Response
	err := client.Do(ctx,
		joinOpts(opts,
			gitlab.WithMethod(method),
			gitlab.WithPath(ModuleRequestApiPath+url.PathEscape(moduleName)+urlPath),
			gitlab.WithQuery(query),
			gitlab.WithHeader(header),
			gitlab.WithAgentToken(agentToken),
			gitlab.WithJWT(true),
			gitlab.WithRequestBody(body, ""),
			gitlab.WithResponseHandler(gitlab.NakedResponseHandler(&resp)),
		)...,
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
