package api

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/gitlab"
)

const (
	ModuleRequestApiPath = "/api/v4/internal/kubernetes/modules/"
)

func MakeModuleRequest(ctx context.Context, client gitlab.ClientInterface, agentToken api.AgentToken, moduleName, method, urlPath string, query url.Values, header http.Header, body io.Reader) (*http.Response, error) {
	var resp *http.Response
	err := client.Do(ctx,
		gitlab.WithMethod(method),
		gitlab.WithPath(ModuleRequestApiPath+url.PathEscape(moduleName)+urlPath),
		gitlab.WithQuery(query),
		gitlab.WithHeader(header),
		gitlab.WithAgentToken(agentToken),
		gitlab.WithJWT(true),
		gitlab.WithRequestBody(body, ""),
		gitlab.WithResponseHandler(gitlab.NakedResponseHandler(&resp)),
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
