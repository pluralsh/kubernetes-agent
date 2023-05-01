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
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestAuthorizeProxyUser(t *testing.T) {
	const (
		configProjectId int64  = 21
		agentId         int64  = 42
		accessType      string = "session_cookie"
		accessKey       string = "damndeliciouscookie"
		csrfToken       string = "token"
	)
	ctx, _ := testhelpers.CtxWithSpanContext(t)
	response := &AuthorizeProxyUserResponse{
		Agent: &AuthorizedAgentForUser{
			Id:            agentId,
			ConfigProject: &ConfigProject{Id: configProjectId},
		},
		User: &User{
			Id:       1234,
			Username: "any-user",
		},
		AccessAs: &AccessAsProxyAuthorization{
			AccessAs: &AccessAsProxyAuthorization_User{
				User: &AccessAsUserAuthorization{
					Projects: []*ProjectAccessCF{
						{
							Id:    configProjectId,
							Roles: []string{"Developer"},
						},
					},
				},
			},
		},
	}
	gitLabClient := mock_gitlab.SetupClient(t, AuthorizeProxyUserApiPath, func(w http.ResponseWriter, r *http.Request) {
		testhelpers.AssertRequestMethod(t, r, http.MethodPost)

		data, err := io.ReadAll(r.Body)
		if !assert.NoError(t, err) {
			return
		}
		actual := &AuthorizeProxyUserRequest{}
		err = protojson.Unmarshal(data, actual)
		if !assert.NoError(t, err) {
			return
		}
		expected := &AuthorizeProxyUserRequest{
			AgentId:    agentId,
			AccessType: accessType,
			AccessKey:  accessKey,
			CsrfToken:  csrfToken,
		}

		assert.Empty(t, cmp.Diff(expected, actual, protocmp.Transform()))
		testhelpers.RespondWithJSON(t, w, response)
	})

	auth, err := AuthorizeProxyUser(ctx, gitLabClient, agentId, accessType, accessKey, csrfToken)
	require.NoError(t, err)

	assert.Equal(t, response.Agent.Id, auth.Agent.Id)
}
