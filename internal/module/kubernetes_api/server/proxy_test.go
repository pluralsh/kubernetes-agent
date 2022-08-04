package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gapi "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/gitlab/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/kubernetes_api/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/cache"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/httpz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/prototool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/matcher"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_gitlab"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_kubernetes_api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_usage_metrics"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/testhelpers"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/pkg/agentcfg"
	"gitlab.com/gitlab-org/labkit/correlation"
	"gitlab.com/gitlab-org/labkit/metrics"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	jobToken        = "asdfgasdfxadf"
	requestPath     = "/api/bla"
	requestPayload  = "asdfndaskjfadsbfjsadhvfjhavfjasvf"
	responsePayload = "jknkjnjkasdnfkjasdnfkasdnfjnkjn"
	queryParamValue = "query-param-value with a space"
	queryParamName  = "q with a space"
)

func TestProxy_JobTokenErrors(t *testing.T) {
	tests := []struct {
		name    string
		auth    []string
		message string
	}{
		{
			name:    "missing header",
			message: "GitLab Agent Server: Unauthorized: Authorization header: expecting token",
		},
		{
			name:    "multiple headers",
			auth:    []string{"a", "b"},
			message: "GitLab Agent Server: Unauthorized: Authorization header: expecting a single header, got 2",
		},
		{
			name:    "invalid format1",
			auth:    []string{"Token asdfadsf"},
			message: "GitLab Agent Server: Unauthorized: Authorization header: expecting Bearer token",
		},
		{
			name:    "invalid format2",
			auth:    []string{"Bearer asdfadsf"},
			message: "GitLab Agent Server: Unauthorized: Authorization header: invalid value",
		},
		{
			name:    "invalid agent id",
			auth:    []string{"Bearer ci:asdf:as"},
			message: `GitLab Agent Server: Unauthorized: Authorization header: failed to parse: strconv.ParseInt: parsing "asdf": invalid syntax`,
		},
		{
			name:    "empty token",
			auth:    []string{"Bearer ci:1:"},
			message: "GitLab Agent Server: Unauthorized: Authorization header: empty token",
		},
		{
			name:    "unknown token type",
			auth:    []string{"Bearer blabla:1:asd"},
			message: "GitLab Agent Server: Unauthorized: Authorization header: unknown token type",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, _, client, req, _, _ := setupProxyWithHandler(t, "/", func(w http.ResponseWriter, r *http.Request) {
				t.Fail() // unexpected invocation
			})
			delete(req.Header, httpz.AuthorizationHeader)
			if len(tc.auth) > 0 {
				req.Header[httpz.AuthorizationHeader] = tc.auth
			}
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()
			assert.EqualValues(t, http.StatusUnauthorized, resp.StatusCode)
			expected := metav1.Status{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Status",
					APIVersion: "v1",
				},
				Status:  metav1.StatusFailure,
				Message: tc.message + ". Correlation ID: " + resp.Header.Get(testhelpers.CorrelationIdHeader),
				Reason:  metav1.StatusReasonUnauthorized,
				Code:    http.StatusUnauthorized,
			}
			assert.Empty(t, cmp.Diff(expected, readStatus(t, resp)))
		})
	}
}

func TestProxy_AllowedAgentsError(t *testing.T) {
	tests := []struct {
		allowedAgentsHttpStatus int
		expectedHttpStatus      int
		message                 string
		captureErr              bool
	}{
		{
			allowedAgentsHttpStatus: http.StatusUnauthorized, // token is invalid
			expectedHttpStatus:      http.StatusUnauthorized,
			message:                 "GitLab Agent Server: Unauthorized: CI job token: HTTP status code: 401 for path /api/v4/job/allowed_agents",
		},
		{
			allowedAgentsHttpStatus: http.StatusForbidden, // token is forbidden
			expectedHttpStatus:      http.StatusForbidden,
			message:                 "GitLab Agent Server: Forbidden: CI job token: HTTP status code: 403 for path /api/v4/job/allowed_agents",
		},
		{
			allowedAgentsHttpStatus: http.StatusNotFound, // agent is not found
			expectedHttpStatus:      http.StatusNotFound,
			message:                 "GitLab Agent Server: Not found: agents for CI job token: HTTP status code: 404 for path /api/v4/job/allowed_agents",
		},
		{
			allowedAgentsHttpStatus: http.StatusBadGateway, // some weird error
			expectedHttpStatus:      http.StatusInternalServerError,
			message:                 "GitLab Agent Server: Failed to get allowed agents for CI job token: HTTP status code: 502 for path /api/v4/job/allowed_agents",
			captureErr:              true,
		},
	}
	for _, tc := range tests {
		t.Run(strconv.Itoa(tc.allowedAgentsHttpStatus), func(t *testing.T) {
			api, _, client, req, _, _ := setupProxyWithHandler(t, "/", func(w http.ResponseWriter, r *http.Request) {
				assertToken(t, r)
				w.WriteHeader(tc.allowedAgentsHttpStatus)
			})
			if tc.captureErr {
				api.EXPECT().
					HandleProcessingError(gomock.Any(), gomock.Any(), testhelpers.AgentId, gomock.Any(),
						matcher.ErrorEq(fmt.Sprintf("HTTP status code: %d for path /api/v4/job/allowed_agents", tc.allowedAgentsHttpStatus)))
			}
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()
			assert.EqualValues(t, tc.expectedHttpStatus, resp.StatusCode)
			expected := metav1.Status{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Status",
					APIVersion: "v1",
				},
				Status:  metav1.StatusFailure,
				Message: tc.message + ". Correlation ID: " + resp.Header.Get(testhelpers.CorrelationIdHeader),
				Reason:  code2reason[int32(tc.expectedHttpStatus)],
				Code:    int32(tc.expectedHttpStatus),
			}
			assert.Empty(t, cmp.Diff(expected, readStatus(t, resp)))
		})
	}
}

func TestProxy_NoExpectedUrlPathPrefix(t *testing.T) {
	_, _, client, req, _, _ := setupProxyWithHandler(t, "/bla/", configGitLabHandler(t, nil, nil))
	req.URL.Path = requestPath
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.EqualValues(t, http.StatusBadRequest, resp.StatusCode)
	expected := metav1.Status{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Status",
			APIVersion: "v1",
		},
		Status:  metav1.StatusFailure,
		Message: "GitLab Agent Server: Bad request: URL does not start with expected prefix. Correlation ID: " + resp.Header.Get(testhelpers.CorrelationIdHeader),
		Reason:  metav1.StatusReasonBadRequest,
		Code:    http.StatusBadRequest,
	}
	assert.Empty(t, cmp.Diff(expected, readStatus(t, resp)))
}

func TestProxy_ForbiddenAgentId(t *testing.T) {
	_, _, client, req, _, _ := setupProxy(t)
	req.Header.Set(httpz.AuthorizationHeader, fmt.Sprintf("Bearer %s:%d:%s", tokenTypeCi, 15 /* disallowed id */, jobToken))
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.EqualValues(t, http.StatusForbidden, resp.StatusCode)
	expected := metav1.Status{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Status",
			APIVersion: "v1",
		},
		Status:  metav1.StatusFailure,
		Message: "GitLab Agent Server: Forbidden: agentId is not allowed. Correlation ID: " + resp.Header.Get(testhelpers.CorrelationIdHeader),
		Reason:  metav1.StatusReasonForbidden,
		Code:    http.StatusForbidden,
	}
	assert.Empty(t, cmp.Diff(expected, readStatus(t, resp)))
}

func TestProxy_HappyPath(t *testing.T) {
	tests := []struct {
		name          string
		urlPathPrefix string
		config        *gapi.Configuration
		env           *gapi.Environment
		expectedExtra *rpc.HeaderExtra
	}{
		{
			name:          "no prefix",
			urlPathPrefix: "/",
			expectedExtra: &rpc.HeaderExtra{
				ImpConfig: &rpc.ImpersonationConfig{},
			},
		},
		{
			name:          "with prefix",
			urlPathPrefix: "/bla/",
			expectedExtra: &rpc.HeaderExtra{
				ImpConfig: &rpc.ImpersonationConfig{},
			},
		},
		{
			name:          "impersonate agent",
			urlPathPrefix: "/",
			config: &gapi.Configuration{
				AccessAs: &agentcfg.CiAccessAsCF{
					As: &agentcfg.CiAccessAsCF_Agent{
						Agent: &agentcfg.CiAccessAsAgentCF{},
					},
				},
			},
			expectedExtra: &rpc.HeaderExtra{
				ImpConfig: &rpc.ImpersonationConfig{},
			},
		},
		{
			name:          "impersonate",
			urlPathPrefix: "/",
			config: &gapi.Configuration{
				AccessAs: &agentcfg.CiAccessAsCF{
					As: &agentcfg.CiAccessAsCF_Impersonate{
						Impersonate: &agentcfg.CiAccessAsImpersonateCF{
							Username: "user1",
							Groups:   []string{"g1", "g2"},
							Uid:      "uid",
							Extra: []*agentcfg.ExtraKeyValCF{
								{
									Key: "k1",
									Val: []string{"v1", "v2"},
								},
							},
						},
					},
				},
			},
			expectedExtra: &rpc.HeaderExtra{
				ImpConfig: &rpc.ImpersonationConfig{
					Username: "user1",
					Groups:   []string{"g1", "g2"},
					Uid:      "uid",
					Extra: []*rpc.ExtraKeyVal{
						{
							Key: "k1",
							Val: []string{"v1", "v2"},
						},
					},
				},
			},
		},
		{
			name:          "impersonate ci job no env",
			urlPathPrefix: "/",
			config: &gapi.Configuration{
				AccessAs: &agentcfg.CiAccessAsCF{
					As: &agentcfg.CiAccessAsCF_CiJob{
						CiJob: &agentcfg.CiAccessAsCiJobCF{},
					},
				},
			},
			expectedExtra: &rpc.HeaderExtra{
				ImpConfig: &rpc.ImpersonationConfig{
					Username: "gitlab:ci_job:1",
					Groups:   []string{"gitlab:ci_job", "gitlab:group:6", "gitlab:project:3"},
					Extra: []*rpc.ExtraKeyVal{
						{
							Key: "agent.gitlab.com/id",
							Val: []string{"123"},
						},
						{
							Key: "agent.gitlab.com/config_project_id",
							Val: []string{"5"},
						},
						{
							Key: "agent.gitlab.com/project_id",
							Val: []string{"3"},
						},
						{
							Key: "agent.gitlab.com/ci_pipeline_id",
							Val: []string{"2"},
						},
						{
							Key: "agent.gitlab.com/ci_job_id",
							Val: []string{"1"},
						},
						{
							Key: "agent.gitlab.com/username",
							Val: []string{"testuser"},
						},
					},
				},
			},
		},
		{
			name:          "impersonate ci job prod env",
			urlPathPrefix: "/",
			config: &gapi.Configuration{
				AccessAs: &agentcfg.CiAccessAsCF{
					As: &agentcfg.CiAccessAsCF_CiJob{
						CiJob: &agentcfg.CiAccessAsCiJobCF{},
					},
				},
			},
			env: &gapi.Environment{
				Slug: "prod",
			},
			expectedExtra: &rpc.HeaderExtra{
				ImpConfig: &rpc.ImpersonationConfig{
					Username: "gitlab:ci_job:1",
					Groups:   []string{"gitlab:ci_job", "gitlab:group:6", "gitlab:project:3", "gitlab:project_env:3:prod"},
					Extra: []*rpc.ExtraKeyVal{
						{
							Key: "agent.gitlab.com/id",
							Val: []string{"123"},
						},
						{
							Key: "agent.gitlab.com/config_project_id",
							Val: []string{"5"},
						},
						{
							Key: "agent.gitlab.com/project_id",
							Val: []string{"3"},
						},
						{
							Key: "agent.gitlab.com/ci_pipeline_id",
							Val: []string{"2"},
						},
						{
							Key: "agent.gitlab.com/ci_job_id",
							Val: []string{"1"},
						},
						{
							Key: "agent.gitlab.com/username",
							Val: []string{"testuser"},
						},
						{
							Key: "agent.gitlab.com/environment_slug",
							Val: []string{"prod"},
						},
					},
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testProxyHappyPath(t, tc.urlPathPrefix, tc.expectedExtra, configGitLabHandler(t, tc.config, tc.env))
		})
	}
}

func testProxyHappyPath(t *testing.T, urlPathPrefix string, expectedExtra *rpc.HeaderExtra, handler func(http.ResponseWriter, *http.Request)) {
	_, k8sClient, client, req, requestCount, ciTunnelUsageSet := setupProxyWithHandler(t, urlPathPrefix, handler)
	requestCount.EXPECT().Inc()
	ciTunnelUsageSet.EXPECT().Add(testhelpers.AgentId)
	mrClient := mock_kubernetes_api.NewMockKubernetesApi_MakeRequestClient(gomock.NewController(t))
	mrCall := k8sClient.EXPECT().
		MakeRequest(gomock.Any()).
		DoAndReturn(func(ctx context.Context, opts ...grpc.CallOption) (rpc.KubernetesApi_MakeRequestClient, error) {
			requireCorrectOutgoingMeta(t, ctx)
			return mrClient, nil
		})
	extra, err := anypb.New(expectedExtra)
	require.NoError(t, err)
	send := mockSendStream(t, mrClient,
		&grpctool.HttpRequest{
			Message: &grpctool.HttpRequest_Header_{
				Header: &grpctool.HttpRequest_Header{
					Request: &prototool.HttpRequest{
						Method: http.MethodPost,
						Header: map[string]*prototool.Values{
							"Req-Header": {
								Value: []string{"x1", "x2"},
							},
							"Accept-Encoding": { // added by the Go client
								Value: []string{"gzip"},
							},
							httpz.UserAgentHeader: {
								Value: []string{"test-agent"},
							},
							"Content-Length": { // added by the Go client
								Value: []string{strconv.Itoa(len(requestPayload))},
							},
							httpz.ViaHeader: {
								Value: []string{"gRPC/1.0 sv1"},
							},
						},
						UrlPath: requestPath,
						Query: map[string]*prototool.Values{
							queryParamName: {
								Value: []string{queryParamValue},
							},
						},
					},
					Extra: extra,
				},
			},
		},
		&grpctool.HttpRequest{
			Message: &grpctool.HttpRequest_Data_{
				Data: &grpctool.HttpRequest_Data{
					Data: []byte(requestPayload),
				},
			},
		},
		&grpctool.HttpRequest{
			Message: &grpctool.HttpRequest_Trailer_{
				Trailer: &grpctool.HttpRequest_Trailer{},
			},
		},
	)
	recv := mockRecvStream(mrClient,
		&grpctool.HttpResponse{
			Message: &grpctool.HttpResponse_Header_{
				Header: &grpctool.HttpResponse_Header{
					Response: &prototool.HttpResponse{
						StatusCode: http.StatusOK,
						Status:     http.StatusText(http.StatusOK),
						Header: map[string]*prototool.Values{
							"Resp-Header": {
								Value: []string{"a1", "a2"},
							},
							"Content-Type": {
								Value: []string{"application/octet-stream"},
							},
							"Date": {
								Value: []string{"NOW!"},
							},
						},
					},
				},
			},
		},
		&grpctool.HttpResponse{
			Message: &grpctool.HttpResponse_Data_{
				Data: &grpctool.HttpResponse_Data{
					Data: []byte(responsePayload),
				},
			},
		},
		&grpctool.HttpResponse{
			Message: &grpctool.HttpResponse_Trailer_{
				Trailer: &grpctool.HttpResponse_Trailer{},
			},
		},
	)
	calls := []*gomock.Call{mrCall}
	calls = append(calls, send...)
	calls = append(calls, recv...)
	gomock.InOrder(calls...)

	req.Header.Set("Req-Header", "x1")
	req.Header.Add("Req-Header", "x2")
	req.Header.Set(httpz.UserAgentHeader, "test-agent") // added manually to override what is added by the Go client
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, resp.Body.Close())
	}()
	assert.EqualValues(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, responsePayload, string(readAll(t, resp.Body)))
	delete(resp.Header, "Date")
	assert.NotEmpty(t, resp.Header.Get("X-Request-Id"))
	delete(resp.Header, "X-Request-Id")
	assert.Empty(t, cmp.Diff(map[string][]string{
		"Resp-Header":  {"a1", "a2"},
		"Content-Type": {"application/octet-stream"},
		"Via":          {"gRPC/1.0 sv1"},
	}, (map[string][]string)(resp.Header)))
}

func TestFormatStatusMessage(t *testing.T) {
	const cid = "abc"
	tests := []struct {
		name            string
		ctx             context.Context
		err             error
		expectedMessage string
	}{
		{
			name:            "no err, no correlation",
			ctx:             context.Background(),
			expectedMessage: "GitLab Agent Server: msg",
		},
		{
			name:            "no err, correlation",
			ctx:             correlation.ContextWithCorrelation(context.Background(), cid),
			expectedMessage: "GitLab Agent Server: msg. Correlation ID: abc",
		},
		{
			name:            "err, no correlation",
			ctx:             context.Background(),
			err:             errors.New("boom"),
			expectedMessage: "GitLab Agent Server: msg: boom",
		},
		{
			name:            "err, correlation",
			ctx:             correlation.ContextWithCorrelation(context.Background(), cid),
			err:             errors.New("boom"),
			expectedMessage: "GitLab Agent Server: msg: boom. Correlation ID: abc",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			msg := formatStatusMessage(test.ctx, "msg", test.err)
			assert.Equal(t, test.expectedMessage, msg)
		})
	}
}

func requireCorrectOutgoingMeta(t *testing.T, ctx context.Context) {
	md, _ := metadata.FromOutgoingContext(ctx)
	vals := md.Get(modserver.RoutingAgentIdMetadataKey)
	require.Len(t, vals, 1)
	agentId, err := strconv.ParseInt(vals[0], 10, 64)
	require.NoError(t, err)
	require.Equal(t, testhelpers.AgentId, agentId)
}

func assertToken(t *testing.T, r *http.Request) bool {
	return assert.Equal(t, jobToken, r.Header.Get("Job-Token"))
}

func setupProxy(t *testing.T) (*mock_modserver.MockApi, *mock_kubernetes_api.MockKubernetesApiClient, *http.Client, *http.Request, *mock_usage_metrics.MockCounter, *mock_usage_metrics.MockUniqueCounter) {
	return setupProxyWithHandler(t, "/", configGitLabHandler(t, nil, nil))
}

func configGitLabHandler(t *testing.T, config *gapi.Configuration, env *gapi.Environment) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !assertToken(t, r) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		testhelpers.RespondWithJSON(t, w, &prototool.JsonBox{Message: &gapi.AllowedAgentsForJob{
			AllowedAgents: []*gapi.AllowedAgent{
				{
					Id: testhelpers.AgentId,
					ConfigProject: &gapi.ConfigProject{
						Id: 5,
					},
					Configuration: config,
				},
			},
			Job: &gapi.Job{
				Id: 1,
			},
			Pipeline: &gapi.Pipeline{
				Id: 2,
			},
			Project: &gapi.Project{
				Id: 3,
				Groups: []*gapi.Group{
					{
						Id: 6,
					},
				},
			},
			User: &gapi.User{
				Id:       testhelpers.AgentId,
				Username: "testuser",
			},
			Environment: env,
		}})
	}
}

func setupProxyWithHandler(t *testing.T, urlPathPrefix string, handler func(http.ResponseWriter, *http.Request)) (*mock_modserver.MockApi, *mock_kubernetes_api.MockKubernetesApiClient, *http.Client, *http.Request, *mock_usage_metrics.MockCounter, *mock_usage_metrics.MockUniqueCounter) {
	ctrl := gomock.NewController(t)
	mockApi := mock_modserver.NewMockApi(ctrl)
	k8sClient := mock_kubernetes_api.NewMockKubernetesApiClient(ctrl)
	requestCount := mock_usage_metrics.NewMockCounter(ctrl)
	ciTunnelUsageSet := mock_usage_metrics.NewMockUniqueCounter(ctrl)

	p := kubernetesApiProxy{
		log:                  zaptest.NewLogger(t),
		api:                  mockApi,
		kubernetesApiClient:  k8sClient,
		gitLabClient:         mock_gitlab.SetupClient(t, gapi.AllowedAgentsApiPath, handler),
		allowedAgentsCache:   cache.NewWithError(0, 0, func(err error) bool { return false }),
		requestCounter:       requestCount,
		ciTunnelUsersCounter: ciTunnelUsageSet,
		metricsHttpHandlerFactory: func(next http.Handler, opts ...metrics.HandlerOption) http.Handler {
			return next
		},
		responseSerializer: serializer.NewCodecFactory(runtime.NewScheme()),
		serverName:         "sv1",
		serverVia:          "gRPC/1.0 sv1",
		urlPathPrefix:      urlPathPrefix,
	}
	listener := grpctool.NewDialListener()
	var wg wait.Group
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(func() {
		cancel()
		wg.Wait()
		listener.Close()
	})
	wg.Start(func() {
		assert.NoError(t, p.Run(ctx, listener))
	})
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return listener.DialContext(ctx, addr)
			},
		},
	}
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"http://any_host_will_do.local"+path.Join(urlPathPrefix, requestPath)+"?"+url.QueryEscape(queryParamName)+"="+url.QueryEscape(queryParamValue),
		strings.NewReader(requestPayload),
	)
	require.NoError(t, err)
	req.Header.Set(httpz.AuthorizationHeader, fmt.Sprintf("Bearer %s:%d:%s", tokenTypeCi, testhelpers.AgentId, jobToken))
	return mockApi, k8sClient, client, req, requestCount, ciTunnelUsageSet
}

func mockRecvStream(server *mock_kubernetes_api.MockKubernetesApi_MakeRequestClient, msgs ...proto.Message) []*gomock.Call {
	res := make([]*gomock.Call, 0, len(msgs)+1)
	for _, msg := range msgs {
		call := server.EXPECT().
			RecvMsg(gomock.Any()).
			Do(testhelpers.RecvMsg(msg))
		res = append(res, call)
	}
	call := server.EXPECT().
		RecvMsg(gomock.Any()).
		Return(io.EOF)
	res = append(res, call)
	return res
}

func mockSendStream(t *testing.T, client *mock_kubernetes_api.MockKubernetesApi_MakeRequestClient, msgs ...*grpctool.HttpRequest) []*gomock.Call {
	res := make([]*gomock.Call, 0, len(msgs)+1)
	for _, msg := range msgs {
		call := client.EXPECT().
			Send(matcher.ProtoEq(t, msg))
		res = append(res, call)
	}
	res = append(res, client.EXPECT().CloseSend())
	return res
}

func readAll(t *testing.T, r io.Reader) []byte {
	data, err := io.ReadAll(r)
	require.NoError(t, err)
	return data
}

func readStatus(t *testing.T, resp *http.Response) metav1.Status {
	data := readAll(t, resp.Body)
	var s metav1.Status
	negotiator := runtime.NewClientNegotiator(serializer.NewCodecFactory(runtime.NewScheme()), schema.GroupVersion{})
	decoder, err := negotiator.Decoder(resp.Header.Get(httpz.ContentTypeHeader), nil)
	require.NoError(t, err)
	obj, _, err := decoder.Decode(data, nil, &s)
	require.NoError(t, err)
	return *obj.(*metav1.Status)
}
