package grpctool_test

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/grpctool/test"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/prototool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/testing/matcher"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/testing/mock_kubernetes_api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/testing/mock_rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/testing/testhelpers"
	"go.uber.org/zap/zaptest"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	requestPath     = "/test"
	requestPayload  = "abcdefg"
	responsePayload = "jknkjnjkasdnfkjasdnfkasdnfjnkjn"
)

func TestHttp2grpc_HappyPath(t *testing.T) {
	mrClient, w, r, x := setupHttp2grpc(t)
	headerExtra := &test.Request{}
	send := mockSendHappy(t, mrClient, headerExtra)
	wh := make(http.Header)
	recv := []*gomock.Call{
		mrClient.EXPECT().
			RecvMsg(gomock.Any()).
			Do(testhelpers.RecvMsg(&grpctool.HttpResponse{
				Message: &grpctool.HttpResponse_Header_{
					Header: &grpctool.HttpResponse_Header{
						Response: &prototool.HttpResponse{
							StatusCode: http.StatusOK,
							Status:     http.StatusText(http.StatusOK),
							Header: map[string]*prototool.Values{
								"Resp-Header": {
									Value: []string{"x1", "x2"},
								},
							},
						},
					},
				},
			})),
		w.EXPECT().
			Header().
			Return(wh),
		w.EXPECT().
			WriteHeader(http.StatusOK).
			Do(func(status int) {
				// when WriteHeader is called, headers should have been set already
				assert.Equal(t, http.Header{
					"Resp-Header": []string{"x1", "x2"},
				}, wh)
			}),
		w.EXPECT().
			Flush(),
		mrClient.EXPECT().
			RecvMsg(gomock.Any()).
			Do(testhelpers.RecvMsg(&grpctool.HttpResponse{
				Message: &grpctool.HttpResponse_Data_{
					Data: &grpctool.HttpResponse_Data{
						Data: []byte(responsePayload),
					},
				},
			})),
		w.EXPECT().
			Write([]byte(responsePayload)),
		w.EXPECT().
			Flush(),
		mrClient.EXPECT().
			RecvMsg(gomock.Any()).
			Do(testhelpers.RecvMsg(&grpctool.HttpResponse{
				Message: &grpctool.HttpResponse_Trailer_{
					Trailer: &grpctool.HttpResponse_Trailer{},
				},
			})),
		mrClient.EXPECT().
			RecvMsg(gomock.Any()).
			Return(io.EOF),
	}
	calls := send
	calls = append(calls, recv...)
	gomock.InOrder(calls...)
	x.Pipe(mrClient, w, r, headerExtra)
}

func TestProxy_HeaderRecvError(t *testing.T) {
	mrClient, w, r, x := setupHttp2grpc(t)
	headerExtra := &test.Request{}
	send := mockSendHappy(t, mrClient, headerExtra)
	wh := make(http.Header)
	w.EXPECT().
		Header().
		Return(wh).
		MinTimes(1)
	recv := []*gomock.Call{
		mrClient.EXPECT().
			RecvMsg(gomock.Any()).
			Return(errors.New("no headers for you")),
		w.EXPECT().
			WriteHeader(http.StatusBadGateway).
			Do(func(status int) {
				// when WriteHeader is called, headers should have been set already
				assert.Equal(t, http.Header{
					"Content-Type":           []string{"text/plain; charset=utf-8"},
					"X-Content-Type-Options": []string{"nosniff"},
				}, wh)
			}),
		w.EXPECT().
			Write([]byte("HTTP->gRPC: failed to read gRPC response: no headers for you\n")),
	}
	calls := send
	calls = append(calls, recv...)
	gomock.InOrder(calls...)

	x.Pipe(mrClient, w, r, headerExtra)
}

func TestProxy_ErrorAfterHeaderWritten(t *testing.T) {
	mrClient, w, r, x := setupHttp2grpc(t)
	headerExtra := &test.Request{}
	send := mockSendHappy(t, mrClient, headerExtra)
	wh := make(http.Header)
	recv := []*gomock.Call{
		mrClient.EXPECT().
			RecvMsg(gomock.Any()).
			Do(testhelpers.RecvMsg(&grpctool.HttpResponse{
				Message: &grpctool.HttpResponse_Header_{
					Header: &grpctool.HttpResponse_Header{
						Response: &prototool.HttpResponse{
							StatusCode: http.StatusOK,
							Status:     http.StatusText(http.StatusOK),
							Header: map[string]*prototool.Values{
								"Resp-Header": {
									Value: []string{"x1", "x2"},
								},
							},
						},
					},
				},
			})),
		w.EXPECT().
			Header().
			Return(wh),
		w.EXPECT().
			WriteHeader(http.StatusOK).
			Do(func(status int) {
				// when WriteHeader is called, headers should have been set already
				assert.Equal(t, http.Header{
					"Resp-Header": []string{"x1", "x2"},
				}, wh)
			}),
		w.EXPECT().
			Flush(),
		mrClient.EXPECT().
			RecvMsg(gomock.Any()).
			Return(errors.New("no body for you")),
	}
	calls := send
	calls = append(calls, recv...)
	gomock.InOrder(calls...)

	assert.PanicsWithError(t, http.ErrAbortHandler.Error(), func() {
		x.Pipe(mrClient, w, r, headerExtra)
	})
}

func TestProxy_ErrorAfterBodyWritten(t *testing.T) {
	mrClient, w, r, x := setupHttp2grpc(t)
	headerExtra := &test.Request{}
	send := mockSendHappy(t, mrClient, headerExtra)
	wh := make(http.Header)
	recv := []*gomock.Call{
		mrClient.EXPECT().
			RecvMsg(gomock.Any()).
			Do(testhelpers.RecvMsg(&grpctool.HttpResponse{
				Message: &grpctool.HttpResponse_Header_{
					Header: &grpctool.HttpResponse_Header{
						Response: &prototool.HttpResponse{
							StatusCode: http.StatusOK,
							Status:     http.StatusText(http.StatusOK),
							Header: map[string]*prototool.Values{
								"Resp-Header": {
									Value: []string{"x1", "x2"},
								},
							},
						},
					},
				},
			})),
		w.EXPECT().
			Header().
			Return(wh),
		w.EXPECT().
			WriteHeader(http.StatusOK).
			Do(func(status int) {
				// when WriteHeader is called, headers should have been set already
				assert.Equal(t, http.Header{
					"Resp-Header": []string{"x1", "x2"},
				}, wh)
			}),
		w.EXPECT().
			Flush(),
		mrClient.EXPECT().
			RecvMsg(gomock.Any()).
			Do(testhelpers.RecvMsg(&grpctool.HttpResponse{
				Message: &grpctool.HttpResponse_Data_{
					Data: &grpctool.HttpResponse_Data{
						Data: []byte(responsePayload),
					},
				},
			})),
		w.EXPECT().
			Write([]byte(responsePayload)),
		w.EXPECT().
			Flush(),
		mrClient.EXPECT().
			RecvMsg(gomock.Any()).
			Return(errors.New("no body for you")),
	}
	calls := send
	calls = append(calls, recv...)
	gomock.InOrder(calls...)

	assert.PanicsWithError(t, http.ErrAbortHandler.Error(), func() {
		x.Pipe(mrClient, w, r, headerExtra)
	})
}

func setupHttp2grpc(t *testing.T) (*mock_kubernetes_api.MockKubernetesApi_MakeRequestClient, *mock_rpc.MockResponseWriterFlusher, *http.Request, grpctool.InboundHttpToOutboundGrpc) {
	ctrl := gomock.NewController(t)
	mrClient := mock_kubernetes_api.NewMockKubernetesApi_MakeRequestClient(ctrl)
	w := mock_rpc.NewMockResponseWriterFlusher(ctrl)
	r := &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Scheme:   "http",
			Host:     "example.com",
			Path:     requestPath,
			RawQuery: "x=1",
		},
		Header: http.Header{
			"A": []string{"a1", "a2"},
		},
		Body: io.NopCloser(strings.NewReader(requestPayload)),
	}

	x := grpctool.InboundHttpToOutboundGrpc{
		Log: zaptest.NewLogger(t),
		HandleProcessingError: func(msg string, err error) {
		},
		MergeHeaders: func(fromOutbound, toInbound http.Header) {
			for k, v := range fromOutbound {
				toInbound[k] = append(toInbound[k], v...)
			}
		},
	}
	return mrClient, w, r, x
}

func mockSendHappy(t *testing.T, mrClient *mock_kubernetes_api.MockKubernetesApi_MakeRequestClient, headerExtra proto.Message) []*gomock.Call {
	extra, err := anypb.New(headerExtra)
	require.NoError(t, err)
	return mockSendHtto2grpcStream(t, mrClient,
		&grpctool.HttpRequest{
			Message: &grpctool.HttpRequest_Header_{
				Header: &grpctool.HttpRequest_Header{
					Request: &prototool.HttpRequest{
						Method: http.MethodGet,
						Header: map[string]*prototool.Values{
							"A": {
								Value: []string{"a1", "a2"},
							},
						},
						UrlPath: requestPath,
						Query: map[string]*prototool.Values{
							"x": {
								Value: []string{"1"},
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
}

func mockSendHtto2grpcStream(t *testing.T, client *mock_kubernetes_api.MockKubernetesApi_MakeRequestClient, msgs ...*grpctool.HttpRequest) []*gomock.Call {
	res := make([]*gomock.Call, 0, len(msgs)+1)
	for _, msg := range msgs {
		call := client.EXPECT().
			Send(matcher.ProtoEq(t, msg))
		res = append(res, call)
	}
	res = append(res, client.EXPECT().CloseSend())
	return res
}
