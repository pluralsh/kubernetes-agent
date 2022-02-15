package agent

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/kubernetes_api/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modagent"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/transport"
	cmdtesting "k8s.io/kubectl/pkg/cmd/testing"
)

var (
	_ httpClient       = (*impersonatingClient)(nil)
	_ modagent.Module  = (*module)(nil)
	_ modagent.Factory = (*Factory)(nil)
)

func TestClientImpersonation(t *testing.T) {
	restImpConfig := rest.ImpersonationConfig{
		UserName: "ruser1",
		UID:      "ruid1",
		Groups:   []string{"rg1", "rg2"},
		Extra: map[string][]string{
			"rx": {"rx1", "rx2"},
		},
	}
	impConfig := &rpc.ImpersonationConfig{
		Username: "iuser1",
		Groups:   []string{"ig1", "ig2"},
		Uid:      "iuid",
		Extra: []*rpc.ExtraKeyVal{
			{
				Key: "ix",
				Val: []string{"ix1", "ix2"},
			},
		},
	}
	requestHeader := http.Header{}
	requestHeader.Set(transport.ImpersonateUserHeader, "huser1")
	requestHeader.Set(transport.ImpersonateUIDHeader, "huid1")
	requestHeader.Set(transport.ImpersonateGroupHeader, "hg1")
	requestHeader.Add(transport.ImpersonateGroupHeader, "hg2")
	requestHeader.Set(transport.ImpersonateUserExtraHeaderPrefix+"Hx", "hx1")
	requestHeader.Add(transport.ImpersonateUserExtraHeaderPrefix+"Hx", "hx2")

	tests := []struct {
		name                  string
		restImpConfig         rest.ImpersonationConfig
		impConfig             *rpc.ImpersonationConfig
		requestHeader         http.Header
		expectedRequestHeader http.Header
		expectedStatus        int
	}{
		{
			name:           "no impersonation",
			impConfig:      &rpc.ImpersonationConfig{},
			expectedStatus: http.StatusOK,
		},
		{
			name:          "rest config",
			restImpConfig: restImpConfig,
			impConfig:     &rpc.ImpersonationConfig{},
			expectedRequestHeader: http.Header{
				transport.ImpersonateUserHeader:                   {"ruser1"},
				transport.ImpersonateUIDHeader:                    {"ruid1"},
				transport.ImpersonateGroupHeader:                  {"rg1", "rg2"},
				transport.ImpersonateUserExtraHeaderPrefix + "Rx": {"rx1", "rx2"},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "rest config and impConfig",
			restImpConfig:  restImpConfig,
			impConfig:      impConfig,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "rest config and requestHeader",
			restImpConfig:  restImpConfig,
			requestHeader:  requestHeader,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "rest config and impConfig and requestHeader",
			restImpConfig:  restImpConfig,
			impConfig:      impConfig,
			requestHeader:  requestHeader,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "impConfig and requestHeader",
			impConfig:      impConfig,
			requestHeader:  requestHeader,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:          "requestHeader",
			requestHeader: requestHeader,
			expectedRequestHeader: http.Header{
				transport.ImpersonateUserHeader:                   {"huser1"},
				transport.ImpersonateUIDHeader:                    {"huid1"},
				transport.ImpersonateGroupHeader:                  {"hg1", "hg2"},
				transport.ImpersonateUserExtraHeaderPrefix + "Hx": {"hx1", "hx2"},
			},
			expectedStatus: http.StatusOK,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tf := cmdtesting.NewTestFactory()
			defer tf.Cleanup()
			config, err := tf.ToRESTConfig()
			require.NoError(t, err)
			config.Impersonate = tc.restImpConfig
			rt := &testRoundTripper{
				Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("")),
				},
			}
			config.Transport = rt
			c := impersonatingClient{
				restConfig: config,
			}
			r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
			require.NoError(t, err)
			if tc.requestHeader != nil {
				r.Header = tc.requestHeader
			}

			resp, err := c.Do(tc.impConfig, r)
			require.NoError(t, err)
			defer resp.Body.Close()

			for k, v := range tc.expectedRequestHeader {
				assert.Equal(t, v, rt.Request.Header[k], k)
			}
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
		})
	}
}

type testRoundTripper struct {
	Request  *http.Request
	Response *http.Response
	Err      error
}

func (rt *testRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	rt.Request = req
	return rt.Response, rt.Err
}
