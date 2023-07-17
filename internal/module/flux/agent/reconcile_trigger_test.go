package agent

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/matcher"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_stdlib"
	"go.uber.org/mock/gomock"
)

func TestReconcileTrigger_WithKubeProxyApiUrl(t *testing.T) {
	// GIVEN
	ctrl := gomock.NewController(t)
	dummyKubeApiUrl := &url.URL{Scheme: "http", Host: "localhost", Path: "kubeapi"}
	mockRoundTripper := mock_stdlib.NewMockRoundTripper(ctrl)
	cfgUrl := "/api/v1/namespaces/flux-system/services/http:webhook-receiver:80/proxy"
	expectedUrl := &url.URL{Scheme: "http", Host: "localhost", Path: "/kubeapi/api/v1/namespaces/flux-system/services/http:webhook-receiver:80/proxy/some/webhook/path"}

	// setup mock expectations
	mockRoundTripper.EXPECT().RoundTrip(matcher.DoMatch(func(r *http.Request) bool {
		return r.URL.String() == expectedUrl.String()
	})).Return(&http.Response{StatusCode: http.StatusOK, Body: http.NoBody}, nil).Times(1)

	// WHEN
	rt, err := newGitRepositoryReconcileTrigger(cfgUrl, dummyKubeApiUrl, mockRoundTripper, nil)
	require.NoError(t, err)

	err = rt.reconcile(context.Background(), "/some/webhook/path")

	// THEN
	assert.NoError(t, err)
}

func TestReconcileTrigger_WithDefaultServiceUrl(t *testing.T) {
	// GIVEN
	ctrl := gomock.NewController(t)
	mockRoundTripper := mock_stdlib.NewMockRoundTripper(ctrl)
	expectedUrl := &url.URL{Scheme: "http", Host: "webhook-receiver.flux-system.svc.cluster.local", Path: "/some/webhook/path"}

	// setup mock expectations
	mockRoundTripper.EXPECT().RoundTrip(matcher.DoMatch(func(r *http.Request) bool {
		return r.URL.String() == expectedUrl.String()
	})).Return(&http.Response{StatusCode: http.StatusOK, Body: http.NoBody}, nil).Times(1)

	// WHEN
	rt, err := newGitRepositoryReconcileTrigger(defaultServiceApiBaseUrl, nil, nil, mockRoundTripper)
	require.NoError(t, err)

	err = rt.reconcile(context.Background(), "/some/webhook/path")

	// THEN
	assert.NoError(t, err)
}

func TestReconcileTrigger_WithCustomServiceUrl(t *testing.T) {
	// GIVEN
	ctrl := gomock.NewController(t)
	mockRoundTripper := mock_stdlib.NewMockRoundTripper(ctrl)
	expectedUrl := &url.URL{Scheme: "https", Host: "localhost", Path: "/some/webhook/path"}

	// setup mock expectations
	mockRoundTripper.EXPECT().RoundTrip(matcher.DoMatch(func(r *http.Request) bool {
		return r.URL.String() == expectedUrl.String()
	})).Return(&http.Response{StatusCode: http.StatusOK, Body: http.NoBody}, nil).Times(1)

	// WHEN
	rt, err := newGitRepositoryReconcileTrigger("https://localhost", nil, nil, mockRoundTripper)
	require.NoError(t, err)

	err = rt.reconcile(context.Background(), "/some/webhook/path")

	// THEN
	assert.NoError(t, err)
}

func TestReconcileTrigger_Failure(t *testing.T) {
	// GIVEN
	ctrl := gomock.NewController(t)
	mockRoundTripper := mock_stdlib.NewMockRoundTripper(ctrl)

	// setup mock expectations
	mockRoundTripper.EXPECT().RoundTrip(gomock.Any()).Return(&http.Response{StatusCode: http.StatusUnauthorized, Status: "401 Unauthorized", Body: http.NoBody}, nil).Times(1)

	// WHEN
	rt, err := newGitRepositoryReconcileTrigger("https://localhost", nil, nil, mockRoundTripper)
	require.NoError(t, err)

	err = rt.reconcile(context.Background(), "/some/webhook/path")

	// THEN
	assert.ErrorContains(t, err, "trigger to \"https://localhost/some/webhook/path\" returned status \"401 Unauthorized\"")
}
