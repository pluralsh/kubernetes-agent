package api

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_gitlab"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/testhelpers"
)

func TestVerifyProjectAccess_Success(t *testing.T) {
	// GIVEN
	projectId := "bla/foo"
	ctx, traceId := testhelpers.CtxWithSpanContext(t)
	gitLabClient := mock_gitlab.SetupClient(t, VerifyProjectAccessApiPath, func(w http.ResponseWriter, r *http.Request) {
		testhelpers.AssertRequestMethod(t, r, http.MethodGet)
		testhelpers.AssertGetRequestIsCorrect(t, r, traceId)
		assert.Equal(t, projectId, r.URL.Query().Get(ProjectIdQueryParam))

		w.WriteHeader(http.StatusNoContent)
	})

	// WHEN
	hasAccess, err := VerifyProjectAccess(ctx, gitLabClient, testhelpers.AgentkToken, projectId)

	// THEN
	require.NoError(t, err)
	require.True(t, hasAccess)
}

func TestVerifyProjectAccess_NoAccessFailure(t *testing.T) {
	testcases := []int{http.StatusNotFound, http.StatusForbidden, http.StatusUnauthorized}

	for _, statusCode := range testcases {
		t.Run(fmt.Sprintf("response status code %d", statusCode), func(t *testing.T) {
			// GIVEN
			projectId := "bla/foo"
			ctx, traceId := testhelpers.CtxWithSpanContext(t)
			gitLabClient := mock_gitlab.SetupClient(t, VerifyProjectAccessApiPath, func(w http.ResponseWriter, r *http.Request) {
				testhelpers.AssertRequestMethod(t, r, http.MethodGet)
				testhelpers.AssertGetRequestIsCorrect(t, r, traceId)
				assert.Equal(t, projectId, r.URL.Query().Get(ProjectIdQueryParam))

				w.WriteHeader(statusCode)
			})

			// WHEN
			hasAccess, err := VerifyProjectAccess(ctx, gitLabClient, testhelpers.AgentkToken, projectId)

			// THEN
			require.NoError(t, err)
			require.False(t, hasAccess)
		})
	}
}

func TestVerifyProjectAccess_ApiFailure(t *testing.T) {
	testcases := []int{http.StatusBadRequest, http.StatusInternalServerError}

	for _, statusCode := range testcases {
		t.Run(fmt.Sprintf("response status code %d", statusCode), func(t *testing.T) {
			// GIVEN
			projectId := "bla/foo"
			ctx, traceId := testhelpers.CtxWithSpanContext(t)
			gitLabClient := mock_gitlab.SetupClient(t, VerifyProjectAccessApiPath, func(w http.ResponseWriter, r *http.Request) {
				testhelpers.AssertRequestMethod(t, r, http.MethodGet)
				testhelpers.AssertGetRequestIsCorrect(t, r, traceId)
				assert.Equal(t, projectId, r.URL.Query().Get(ProjectIdQueryParam))

				w.WriteHeader(statusCode)
			})

			// WHEN
			hasAccess, err := VerifyProjectAccess(ctx, gitLabClient, testhelpers.AgentkToken, projectId)

			// THEN
			require.Error(t, err)
			require.False(t, hasAccess)
		})
	}
}
