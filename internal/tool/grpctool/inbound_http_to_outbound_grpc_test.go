package grpctool_test

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/grpctool/test"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/httpz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/prototool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/matcher"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_kubernetes_api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_stdlib"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/testhelpers"
	"go.uber.org/zap/zaptest"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	requestPath = "/test"
)

func TestHttp2Grpc_HappyPath(t *testing.T) {
	mrClient, w, r, x := setupHttp2grpc(t, false)
	headerExtra := &test.Request{}
	send := mockSendHappy(t, mrClient, headerExtra, false)
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
						Data: []byte(responseBodyData),
					},
				},
			})),
		w.EXPECT().
			Write([]byte(responseBodyData)),
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

func TestHttp2Grpc_UpgradeHappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	mrClient, w, r, x := setupHttp2grpc(t, true)
	conn := mock_stdlib.NewMockConn(ctrl)
	headerExtra := &test.Request{}
	wh := make(http.Header)
	setReadDeadlineCall := conn.EXPECT().
		SetReadDeadline(time.Time{})
	send := mockSendHappy(t, mrClient, headerExtra, true)
	recv := []*gomock.Call{
		mrClient.EXPECT().
			RecvMsg(gomock.Any()).
			Do(testhelpers.RecvMsg(&grpctool.HttpResponse{
				Message: &grpctool.HttpResponse_Header_{
					Header: &grpctool.HttpResponse_Header{
						Response: &prototool.HttpResponse{
							StatusCode: http.StatusSwitchingProtocols,
							Status:     http.StatusText(http.StatusSwitchingProtocols),
							Header: map[string]*prototool.Values{
								"Resp-Header": {
									Value: []string{"x1", "x2"},
								},
								httpz.UpgradeHeader: {
									Value: []string{"http/x"},
								},
								httpz.ConnectionHeader: {
									Value: []string{"upgrade"},
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
			WriteHeader(http.StatusSwitchingProtocols).
			Do(func(status int) {
				// when WriteHeader is called, headers should have been set already
				assert.Equal(t, http.Header{
					"Resp-Header":          []string{"x1", "x2"},
					httpz.UpgradeHeader:    []string{"http/x"},
					httpz.ConnectionHeader: []string{"upgrade"},
				}, wh)
			}),
		mrClient.EXPECT().
			RecvMsg(gomock.Any()).
			Do(testhelpers.RecvMsg(&grpctool.HttpResponse{
				Message: &grpctool.HttpResponse_Data_{
					Data: &grpctool.HttpResponse_Data{
						Data: []byte(responseBodyData),
					},
				},
			})),
		w.EXPECT().
			Write([]byte(responseBodyData)),
		w.EXPECT().
			Flush(),
		mrClient.EXPECT().
			RecvMsg(gomock.Any()).
			Do(testhelpers.RecvMsg(&grpctool.HttpResponse{
				Message: &grpctool.HttpResponse_Trailer_{
					Trailer: &grpctool.HttpResponse_Trailer{},
				},
			})),
		w.EXPECT().
			Hijack().
			Return(conn, bufio.NewReadWriter(bufio.NewReader(conn), nil), nil),
		setReadDeadlineCall,
	}
	calls := send
	calls = append(calls, recv...)
	gomock.InOrder(calls...)
	connCloseCall := conn.EXPECT().Close()
	// pipeOutboundToInboundUpgraded
	gomock.InOrder(
		setReadDeadlineCall,
		mrClient.EXPECT().
			RecvMsg(gomock.Any()).
			Do(testhelpers.RecvMsg(&grpctool.HttpResponse{
				Message: &grpctool.HttpResponse_UpgradeData_{
					UpgradeData: &grpctool.HttpResponse_UpgradeData{
						Data: []byte(responseUpgradeBodyData),
					},
				},
			})),
		conn.EXPECT().
			SetWriteDeadline(gomock.Any()),
		conn.EXPECT().
			Write([]byte(responseUpgradeBodyData)),
		mrClient.EXPECT().
			RecvMsg(gomock.Any()).
			Return(io.EOF),
		connCloseCall,
	)
	// pipeInboundToOutboundUpgraded
	gomock.InOrder(
		setReadDeadlineCall,
		conn.EXPECT().
			Read(gomock.Any()).
			DoAndReturn(func(b []byte) (int, error) {
				return copy(b, requestUpgradeBodyData), io.EOF
			}),
		mrClient.EXPECT().
			Send(matcher.ProtoEq(t, &grpctool.HttpRequest{
				Message: &grpctool.HttpRequest_UpgradeData_{
					UpgradeData: &grpctool.HttpRequest_UpgradeData{
						Data: []byte(requestUpgradeBodyData),
					},
				},
			})),
		mrClient.EXPECT().CloseSend(),
		connCloseCall,
	)
	x.Pipe(mrClient, w, r, headerExtra)
}

func TestHttp2Grpc_ServerRefusesToUpgrade(t *testing.T) {
	mrClient, w, r, x := setupHttp2grpc(t, true)
	headerExtra := &test.Request{}
	wh := make(http.Header)
	extra, err := anypb.New(headerExtra)
	require.NoError(t, err)
	contentLength := int64(len(requestBodyData))
	send := mockSendHttp2grpcStream(t, mrClient, false,
		&grpctool.HttpRequest{
			Message: &grpctool.HttpRequest_Header_{
				Header: &grpctool.HttpRequest_Header{
					Request: &prototool.HttpRequest{
						Method: http.MethodGet,
						Header: map[string]*prototool.Values{
							"A": {
								Value: []string{"a1", "a2"},
							},
							httpz.UpgradeHeader: {
								Value: []string{"http/x"},
							},
							httpz.ConnectionHeader: {
								Value: []string{"upgrade"},
							},
						},
						UrlPath: requestPath,
						Query: map[string]*prototool.Values{
							"x": {
								Value: []string{"1"},
							},
						},
					},
					Extra:         extra,
					ContentLength: &contentLength,
				},
			},
		},
		&grpctool.HttpRequest{
			Message: &grpctool.HttpRequest_Data_{
				Data: &grpctool.HttpRequest_Data{
					Data: []byte(requestBodyData),
				},
			},
		},
		&grpctool.HttpRequest{
			Message: &grpctool.HttpRequest_Trailer_{
				Trailer: &grpctool.HttpRequest_Trailer{},
			},
		},
	)
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
				Message: &grpctool.HttpResponse_Trailer_{
					Trailer: &grpctool.HttpResponse_Trailer{},
				},
			})),
		mrClient.EXPECT().
			RecvMsg(gomock.Any()).
			Return(io.EOF),
		mrClient.EXPECT().CloseSend(),
	}
	calls := send
	calls = append(calls, recv...)
	gomock.InOrder(calls...)
	x.Pipe(mrClient, w, r, headerExtra)
}

func TestHttp2Grpc_HeaderRecvError(t *testing.T) {
	mrClient, w, r, x := setupHttp2grpc(t, false)
	headerExtra := &test.Request{}
	send := mockSendHappy(t, mrClient, headerExtra, false)
	recv := []*gomock.Call{
		mrClient.EXPECT().
			RecvMsg(gomock.Any()).
			Return(errors.New("no headers for you")),
		w.EXPECT().
			WriteHeader(http.StatusBadGateway),
	}
	calls := send
	calls = append(calls, recv...)
	gomock.InOrder(calls...)

	x.Pipe(mrClient, w, r, headerExtra)
}

func TestHttp2Grpc_ErrorAfterHeaderWritten(t *testing.T) {
	mrClient, w, r, x := setupHttp2grpc(t, false)
	headerExtra := &test.Request{}
	send := mockSendHappy(t, mrClient, headerExtra, false)
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

func TestHttp2Grpc_ErrorAfterBodyWritten(t *testing.T) {
	mrClient, w, r, x := setupHttp2grpc(t, false)
	headerExtra := &test.Request{}
	send := mockSendHappy(t, mrClient, headerExtra, false)
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
						Data: []byte(responseBodyData),
					},
				},
			})),
		w.EXPECT().
			Write([]byte(responseBodyData)),
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

func setupHttp2grpc(t *testing.T, isUpgrade bool) (*mock_kubernetes_api.MockKubernetesApi_MakeRequestClient, *mock_stdlib.MockResponseWriterFlusher, *http.Request, grpctool.InboundHttpToOutboundGrpc) {
	ctrl := gomock.NewController(t)
	mrClient := mock_kubernetes_api.NewMockKubernetesApi_MakeRequestClient(ctrl)
	w := mock_stdlib.NewMockResponseWriterFlusher(ctrl)
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
		ContentLength: int64(len(requestBodyData)),
		Body:          io.NopCloser(strings.NewReader(requestBodyData)),
	}
	if isUpgrade {
		r.Header[httpz.ConnectionHeader] = []string{"upgrade"}
		r.Header[httpz.UpgradeHeader] = []string{"http/x"}
	}

	x := grpctool.InboundHttpToOutboundGrpc{
		Log: zaptest.NewLogger(t),
		HandleProcessingError: func(msg string, err error) {
			t.Error(msg, err)
		},
		WriteErrorResponse: func(w http.ResponseWriter, r *http.Request, eResp *grpctool.ErrResp) {
			w.WriteHeader(int(eResp.StatusCode))
		},
		MergeHeaders: func(outboundResponse, inboundResponse http.Header) {
			for k, v := range outboundResponse {
				inboundResponse[k] = append(inboundResponse[k], v...)
			}
		},
	}
	return mrClient, w, r, x
}

func mockSendHappy(t *testing.T, mrClient *mock_kubernetes_api.MockKubernetesApi_MakeRequestClient, headerExtra proto.Message, isUpgrade bool) []*gomock.Call {
	extra, err := anypb.New(headerExtra)
	require.NoError(t, err)
	header := map[string]*prototool.Values{
		"A": {
			Value: []string{"a1", "a2"},
		},
	}
	if isUpgrade {
		header[httpz.UpgradeHeader] = &prototool.Values{
			Value: []string{"http/x"},
		}
		header[httpz.ConnectionHeader] = &prototool.Values{
			Value: []string{"upgrade"},
		}
	}
	contentLength := int64(len(requestBodyData))
	return mockSendHttp2grpcStream(t, mrClient, !isUpgrade,
		&grpctool.HttpRequest{
			Message: &grpctool.HttpRequest_Header_{
				Header: &grpctool.HttpRequest_Header{
					Request: &prototool.HttpRequest{
						Method:  http.MethodGet,
						Header:  header,
						UrlPath: requestPath,
						Query: map[string]*prototool.Values{
							"x": {
								Value: []string{"1"},
							},
						},
					},
					Extra:         extra,
					ContentLength: &contentLength,
				},
			},
		},
		&grpctool.HttpRequest{
			Message: &grpctool.HttpRequest_Data_{
				Data: &grpctool.HttpRequest_Data{
					Data: []byte(requestBodyData),
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

func mockSendHttp2grpcStream(t *testing.T, client *mock_kubernetes_api.MockKubernetesApi_MakeRequestClient, close bool, msgs ...*grpctool.HttpRequest) []*gomock.Call {
	res := make([]*gomock.Call, 0, len(msgs)+1)
	for _, msg := range msgs {
		call := client.EXPECT().
			Send(matcher.ProtoEq(t, msg))
		res = append(res, call)
	}
	if close {
		res = append(res, client.EXPECT().CloseSend())
	}
	return res
}
