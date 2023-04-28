package api

import (
	"io"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_gitlab"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/testhelpers"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/agentcfg"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestPostAgentConfiguration(t *testing.T) {
	id := "bla"
	config := &agentcfg.ConfigurationFile{
		Gitops: &agentcfg.GitopsCF{
			ManifestProjects: []*agentcfg.ManifestProjectCF{
				{
					Id: &id,
				},
			},
		},
		// don't need to test all fields, some is good enough
	}
	ctx, traceId := testhelpers.CtxWithSpanContext(t)
	c := mock_gitlab.SetupClient(t, AgentConfigurationApiPath, func(w http.ResponseWriter, r *http.Request) {
		testhelpers.AssertJWTSignature(t, r)
		testhelpers.AssertRequestMethod(t, r, http.MethodPost)
		testhelpers.AssertRequestContentTypeJson(t, r)
		testhelpers.AssertCommonRequestParams(t, r, traceId)
		data, err := io.ReadAll(r.Body)
		if !assert.NoError(t, err) {
			return
		}
		actual := &AgentConfigurationRequest{
			AgentConfig: &agentcfg.ConfigurationFile{},
		}
		err = protojson.Unmarshal(data, actual)
		if !assert.NoError(t, err) {
			return
		}
		expected := &AgentConfigurationRequest{
			AgentId:     testhelpers.AgentId,
			AgentConfig: config,
		}
		assert.Empty(t, cmp.Diff(expected, actual, protocmp.Transform()))
		w.WriteHeader(http.StatusNoContent)
	})
	err := PostAgentConfiguration(ctx, c, testhelpers.AgentId, config)
	require.NoError(t, err)
}
