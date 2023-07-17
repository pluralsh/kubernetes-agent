package grpctool_test

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/httpz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/prototool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/matcher"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_stdlib"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/testhelpers"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	requestBodyData  = "request_jkasdbfkadsbfkadbfkjasbfkasbdf"
	responseBodyData = "response_nlnflkwqnflkasdnflnasdlfnasldnflnl"

	requestUpgradeBodyData  = "upgrade_request_asdfjkasbfkasdf"
	responseUpgradeBodyData = "upgrade_response_asdfasdfadsf"
)

func TestGrpc2Http_HappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	server := mock_rpc.NewMockInboundGrpcToOutboundHttpStream(ctrl)
	server.EXPECT().
		Context().
		Return(context.Background()).
		MinTimes(1)
	sh := sendHeader()
	gomock.InOrder(mockRecvStream(server, true,
		&grpctool.HttpRequest{
			Message: &grpctool.HttpRequest_Header_{
				Header: sh,
			},
		},
		&grpctool.HttpRequest{
			Message: &grpctool.HttpRequest_Data_{
				Data: &grpctool.HttpRequest_Data{
					Data: []byte(requestBodyData[:1]),
				},
			},
		},
		&grpctool.HttpRequest{
			Message: &grpctool.HttpRequest_Data_{
				Data: &grpctool.HttpRequest_Data{
					Data: []byte(requestBodyData[1:]),
				},
			},
		},
		&grpctool.HttpRequest{
			Message: &grpctool.HttpRequest_Trailer_{
				Trailer: &grpctool.HttpRequest_Trailer{},
			},
		},
	)...)
	gomock.InOrder(mockSendStream(t, server,
		&grpctool.HttpResponse{
			Message: &grpctool.HttpResponse_Header_{
				Header: &grpctool.HttpResponse_Header{
					Response: &prototool.HttpResponse{
						StatusCode: http.StatusOK,
						Status:     "OK!",
						Header: map[string]*prototool.Values{
							"x1": {
								Value: []string{"a1", "a2"},
							},
						},
					},
				},
			},
		},
		&grpctool.HttpResponse{
			Message: &grpctool.HttpResponse_Data_{
				Data: &grpctool.HttpResponse_Data{
					Data: []byte(responseBodyData),
				},
			},
		},
		&grpctool.HttpResponse{
			Message: &grpctool.HttpResponse_Trailer_{
				Trailer: &grpctool.HttpResponse_Trailer{},
			},
		},
	)...)
	grpc2http := makeGrpc2http(t, func(ctx context.Context, header *grpctool.HttpRequest_Header, body io.Reader) (grpctool.DoResponse, error) {
		assert.Empty(t, cmp.Diff(header, sh, protocmp.Transform()))
		data, err := io.ReadAll(body)
		if !assert.NoError(t, err) {
			return grpctool.DoResponse{}, err
		}
		assert.Equal(t, requestBodyData, string(data))
		return grpctool.DoResponse{
			Resp: &http.Response{
				Status:     "OK!",
				StatusCode: http.StatusOK,
				Header: http.Header{
					"x1": []string{"a1", "a2"},
				},
				Body: io.NopCloser(strings.NewReader(responseBodyData)),
			},
		}, nil
	})
	err := grpc2http.Pipe(server)
	require.NoError(t, err)
}

func TestGrpc2Http_HappyPathNoBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	server := mock_rpc.NewMockInboundGrpcToOutboundHttpStream(ctrl)
	server.EXPECT().
		Context().
		Return(context.Background()).
		MinTimes(1)
	sh := sendHeader()
	contentLength := int64(0)
	sh.ContentLength = &contentLength
	gomock.InOrder(mockRecvStream(server, true,
		&grpctool.HttpRequest{
			Message: &grpctool.HttpRequest_Header_{
				Header: sh,
			},
		},
		&grpctool.HttpRequest{
			Message: &grpctool.HttpRequest_Trailer_{
				Trailer: &grpctool.HttpRequest_Trailer{},
			},
		},
	)...)
	gomock.InOrder(mockSendStream(t, server,
		&grpctool.HttpResponse{
			Message: &grpctool.HttpResponse_Header_{
				Header: &grpctool.HttpResponse_Header{
					Response: &prototool.HttpResponse{
						StatusCode: http.StatusOK,
						Status:     "OK!",
						Header: map[string]*prototool.Values{
							"x1": {
								Value: []string{"a1", "a2"},
							},
						},
					},
				},
			},
		},
		&grpctool.HttpResponse{
			Message: &grpctool.HttpResponse_Trailer_{
				Trailer: &grpctool.HttpResponse_Trailer{},
			},
		},
	)...)
	grpc2http := makeGrpc2http(t, func(ctx context.Context, header *grpctool.HttpRequest_Header, body io.Reader) (grpctool.DoResponse, error) {
		assert.IsType(t, http.NoBody, body)
		assert.Empty(t, cmp.Diff(header, sh, protocmp.Transform()))
		data, err := io.ReadAll(body)
		if !assert.NoError(t, err) {
			return grpctool.DoResponse{}, err
		}
		assert.Empty(t, data)
		return grpctool.DoResponse{
			Resp: &http.Response{
				Status:     "OK!",
				StatusCode: http.StatusOK,
				Header: http.Header{
					"x1": []string{"a1", "a2"},
				},
				Body: io.NopCloser(bytes.NewReader(nil)),
			},
		}, nil
	})
	err := grpc2http.Pipe(server)
	require.NoError(t, err)
}

func TestGrpc2Http_UpgradeHappyPathWithBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	server := mock_rpc.NewMockInboundGrpcToOutboundHttpStream(ctrl)
	server.EXPECT().
		Context().
		Return(context.Background()).
		MinTimes(1)
	conn := mock_stdlib.NewMockConn(ctrl)
	gomock.InOrder(
		conn.EXPECT().Write([]byte(requestUpgradeBodyData[:1])),
		conn.EXPECT().Write([]byte(requestUpgradeBodyData[1:])),
		conn.EXPECT().Close(),
	)
	sh := sendUpgradeHeader()
	gomock.InOrder(mockRecvStream(server, true,
		&grpctool.HttpRequest{
			Message: &grpctool.HttpRequest_Header_{
				Header: sh,
			},
		},
		&grpctool.HttpRequest{
			Message: &grpctool.HttpRequest_Data_{
				Data: &grpctool.HttpRequest_Data{
					Data: []byte(requestBodyData[:1]),
				},
			},
		},
		&grpctool.HttpRequest{
			Message: &grpctool.HttpRequest_Data_{
				Data: &grpctool.HttpRequest_Data{
					Data: []byte(requestBodyData[1:]),
				},
			},
		},
		&grpctool.HttpRequest{
			Message: &grpctool.HttpRequest_Trailer_{
				Trailer: &grpctool.HttpRequest_Trailer{},
			},
		},
		&grpctool.HttpRequest{
			Message: &grpctool.HttpRequest_UpgradeData_{
				UpgradeData: &grpctool.HttpRequest_UpgradeData{
					Data: []byte(requestUpgradeBodyData[:1]),
				},
			},
		},
		&grpctool.HttpRequest{
			Message: &grpctool.HttpRequest_UpgradeData_{
				UpgradeData: &grpctool.HttpRequest_UpgradeData{
					Data: []byte(requestUpgradeBodyData[1:]),
				},
			},
		},
	)...)
	gomock.InOrder(mockSendStream(t, server,
		&grpctool.HttpResponse{
			Message: &grpctool.HttpResponse_Header_{
				Header: &grpctool.HttpResponse_Header{
					Response: &prototool.HttpResponse{
						StatusCode: http.StatusSwitchingProtocols,
						Status:     "OK!",
						Header: map[string]*prototool.Values{
							"x1": {
								Value: []string{"a1", "a2"},
							},
						},
					},
				},
			},
		},
		&grpctool.HttpResponse{
			Message: &grpctool.HttpResponse_Data_{
				Data: &grpctool.HttpResponse_Data{
					Data: []byte(responseBodyData),
				},
			},
		},
		&grpctool.HttpResponse{
			Message: &grpctool.HttpResponse_Trailer_{
				Trailer: &grpctool.HttpResponse_Trailer{},
			},
		},
		&grpctool.HttpResponse{
			Message: &grpctool.HttpResponse_UpgradeData_{
				UpgradeData: &grpctool.HttpResponse_UpgradeData{
					Data: []byte(responseUpgradeBodyData),
				},
			},
		},
	)...)
	grpc2http := makeGrpc2http(t, func(ctx context.Context, header *grpctool.HttpRequest_Header, body io.Reader) (grpctool.DoResponse, error) {
		assert.Empty(t, cmp.Diff(header, sh, protocmp.Transform()))
		data, err := io.ReadAll(body)
		if !assert.NoError(t, err) {
			return grpctool.DoResponse{}, err
		}
		assert.Equal(t, requestBodyData, string(data))
		return grpctool.DoResponse{
			Resp: &http.Response{
				Status:     "OK!",
				StatusCode: http.StatusSwitchingProtocols,
				Header: http.Header{
					"x1": []string{"a1", "a2"},
				},
				Body: io.NopCloser(strings.NewReader(responseBodyData)),
			},
			UpgradeConn: conn,
			ConnReader:  bufio.NewReader(strings.NewReader(responseUpgradeBodyData)),
		}, nil
	})
	err := grpc2http.Pipe(server)
	require.NoError(t, err)
}

func TestGrpc2Http_UpgradeHappyPathNoBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	server := mock_rpc.NewMockInboundGrpcToOutboundHttpStream(ctrl)
	server.EXPECT().
		Context().
		Return(context.Background()).
		MinTimes(1)
	conn := mock_stdlib.NewMockConn(ctrl)
	conn.EXPECT().Close()
	sh := sendUpgradeHeader()
	contentLength := int64(0)
	sh.ContentLength = &contentLength
	gomock.InOrder(mockRecvStream(server, true,
		&grpctool.HttpRequest{
			Message: &grpctool.HttpRequest_Header_{
				Header: sh,
			},
		},
		&grpctool.HttpRequest{
			Message: &grpctool.HttpRequest_Trailer_{
				Trailer: &grpctool.HttpRequest_Trailer{},
			},
		},
	)...)
	gomock.InOrder(mockSendStream(t, server,
		&grpctool.HttpResponse{
			Message: &grpctool.HttpResponse_Header_{
				Header: &grpctool.HttpResponse_Header{
					Response: &prototool.HttpResponse{
						StatusCode: http.StatusSwitchingProtocols,
						Status:     "OK!",
						Header: map[string]*prototool.Values{
							"x1": {
								Value: []string{"a1", "a2"},
							},
						},
					},
				},
			},
		},
		&grpctool.HttpResponse{
			Message: &grpctool.HttpResponse_Trailer_{
				Trailer: &grpctool.HttpResponse_Trailer{},
			},
		},
	)...)
	grpc2http := makeGrpc2http(t, func(ctx context.Context, header *grpctool.HttpRequest_Header, body io.Reader) (grpctool.DoResponse, error) {
		assert.IsType(t, http.NoBody, body)
		assert.Empty(t, cmp.Diff(header, sh, protocmp.Transform()))
		data, err := io.ReadAll(body)
		if !assert.NoError(t, err) {
			return grpctool.DoResponse{}, err
		}
		assert.Empty(t, data)
		return grpctool.DoResponse{
			Resp: &http.Response{
				Status:     "OK!",
				StatusCode: http.StatusSwitchingProtocols,
				Header: http.Header{
					"x1": []string{"a1", "a2"},
				},
				Body: io.NopCloser(bytes.NewReader(nil)),
			},
			UpgradeConn: conn,
			ConnReader:  bufio.NewReader(bytes.NewReader(nil)),
		}, nil
	})
	err := grpc2http.Pipe(server)
	require.NoError(t, err)
}

func TestGrpc2Http_ServerRefusesToUpgrade(t *testing.T) {
	ctrl := gomock.NewController(t)
	server := mock_rpc.NewMockInboundGrpcToOutboundHttpStream(ctrl)
	server.EXPECT().
		Context().
		Return(context.Background()).
		MinTimes(1)
	gomock.InOrder(mockRecvStream(server, true,
		&grpctool.HttpRequest{
			Message: &grpctool.HttpRequest_Header_{
				Header: sendUpgradeHeader(),
			},
		},
		&grpctool.HttpRequest{
			Message: &grpctool.HttpRequest_Trailer_{
				Trailer: &grpctool.HttpRequest_Trailer{},
			},
		},
	)...)
	gomock.InOrder(mockSendStream(t, server,
		&grpctool.HttpResponse{
			Message: &grpctool.HttpResponse_Header_{
				Header: &grpctool.HttpResponse_Header{
					Response: &prototool.HttpResponse{
						StatusCode: http.StatusOK,
						Status:     "OK!",
					},
				},
			},
		},
		&grpctool.HttpResponse{
			Message: &grpctool.HttpResponse_Trailer_{
				Trailer: &grpctool.HttpResponse_Trailer{},
			},
		},
	)...)
	grpc2http := makeGrpc2http(t, func(ctx context.Context, header *grpctool.HttpRequest_Header, body io.Reader) (grpctool.DoResponse, error) {
		return grpctool.DoResponse{
			Resp: &http.Response{
				Status:     "OK!",
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(nil)),
			},
		}, nil
	})
	err := grpc2http.Pipe(server)
	require.NoError(t, err)
}

func TestGrpc2Http_UpgradeMessageForNonUpgradeRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	server := mock_rpc.NewMockInboundGrpcToOutboundHttpStream(ctrl)
	server.EXPECT().
		Context().
		Return(context.Background()).
		MinTimes(1)
	gomock.InOrder(mockRecvStream(server, false,
		&grpctool.HttpRequest{
			Message: &grpctool.HttpRequest_Header_{
				Header: sendHeader(),
			},
		},
		&grpctool.HttpRequest{
			Message: &grpctool.HttpRequest_Trailer_{
				Trailer: &grpctool.HttpRequest_Trailer{},
			},
		},
		&grpctool.HttpRequest{
			Message: &grpctool.HttpRequest_UpgradeData_{
				UpgradeData: &grpctool.HttpRequest_UpgradeData{
					Data: []byte(requestUpgradeBodyData),
				},
			},
		},
	)...)
	gomock.InOrder(mockSendStream(t, server,
		&grpctool.HttpResponse{
			Message: &grpctool.HttpResponse_Header_{
				Header: &grpctool.HttpResponse_Header{
					Response: &prototool.HttpResponse{
						StatusCode: http.StatusOK,
						Status:     "OK!",
					},
				},
			},
		},
		&grpctool.HttpResponse{
			Message: &grpctool.HttpResponse_Trailer_{
				Trailer: &grpctool.HttpResponse_Trailer{},
			},
		},
	)...)
	grpc2http := makeGrpc2http(t, func(ctx context.Context, header *grpctool.HttpRequest_Header, body io.Reader) (grpctool.DoResponse, error) {
		return grpctool.DoResponse{
			Resp: &http.Response{
				Status:     "OK!",
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(nil)),
			},
		}, nil
	})
	err := grpc2http.Pipe(server)
	require.EqualError(t, err, "rpc error: code = Internal desc = unexpected HttpRequest_UpgradeData message for non-upgrade request")
}

func TestGrpc2Http_FailureWhenDataWasSentForRequestNotExpectingData(t *testing.T) {
	ctrl := gomock.NewController(t)
	server := mock_rpc.NewMockInboundGrpcToOutboundHttpStream(ctrl)
	server.EXPECT().
		Context().
		Return(context.Background()).
		MinTimes(1)
	sh := sendHeader()
	contentLength := int64(0)
	sh.ContentLength = &contentLength
	gomock.InOrder(mockRecvStream(server, false,
		&grpctool.HttpRequest{
			Message: &grpctool.HttpRequest_Header_{
				Header: sh,
			},
		},
		&grpctool.HttpRequest{
			Message: &grpctool.HttpRequest_Data_{
				Data: &grpctool.HttpRequest_Data{
					Data: []byte(requestBodyData),
				},
			},
		},
	)...)
	gomock.InOrder(mockSendStream(t, server,
		&grpctool.HttpResponse{
			Message: &grpctool.HttpResponse_Header_{
				Header: &grpctool.HttpResponse_Header{
					Response: &prototool.HttpResponse{
						StatusCode: http.StatusOK,
						Status:     "OK!",
					},
				},
			},
		},
		&grpctool.HttpResponse{
			Message: &grpctool.HttpResponse_Trailer_{
				Trailer: &grpctool.HttpResponse_Trailer{},
			},
		},
	)...)
	grpc2http := makeGrpc2http(t, func(ctx context.Context, header *grpctool.HttpRequest_Header, body io.Reader) (grpctool.DoResponse, error) {
		return grpctool.DoResponse{
			Resp: &http.Response{
				Status:     "OK!",
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(nil)),
			},
		}, nil
	})
	err := grpc2http.Pipe(server)
	require.EqualError(t, err, "rpc error: code = Internal desc = unexpected HttpRequest_Data message received")
}

func TestGrpc2Http_UpgradeMessageWhenServerRefusesToUpgrade(t *testing.T) {
	ctrl := gomock.NewController(t)
	server := mock_rpc.NewMockInboundGrpcToOutboundHttpStream(ctrl)
	server.EXPECT().
		Context().
		Return(context.Background()).
		MinTimes(1)
	gomock.InOrder(mockRecvStream(server, false,
		&grpctool.HttpRequest{
			Message: &grpctool.HttpRequest_Header_{
				Header: sendUpgradeHeader(),
			},
		},
		&grpctool.HttpRequest{
			Message: &grpctool.HttpRequest_Trailer_{
				Trailer: &grpctool.HttpRequest_Trailer{},
			},
		},
		&grpctool.HttpRequest{
			Message: &grpctool.HttpRequest_UpgradeData_{
				UpgradeData: &grpctool.HttpRequest_UpgradeData{
					Data: []byte(requestUpgradeBodyData),
				},
			},
		},
	)...)
	gomock.InOrder(mockSendStream(t, server,
		&grpctool.HttpResponse{
			Message: &grpctool.HttpResponse_Header_{
				Header: &grpctool.HttpResponse_Header{
					Response: &prototool.HttpResponse{
						StatusCode: http.StatusOK,
						Status:     "OK!",
					},
				},
			},
		},
		&grpctool.HttpResponse{
			Message: &grpctool.HttpResponse_Trailer_{
				Trailer: &grpctool.HttpResponse_Trailer{},
			},
		},
	)...)
	grpc2http := makeGrpc2http(t, func(ctx context.Context, header *grpctool.HttpRequest_Header, body io.Reader) (grpctool.DoResponse, error) {
		return grpctool.DoResponse{
			Resp: &http.Response{
				Status:     "OK!",
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(nil)),
			},
		}, nil
	})
	err := grpc2http.Pipe(server)
	require.EqualError(t, err, "rpc error: code = Internal desc = unexpected HttpRequest_UpgradeData message for HTTP status code 200")
}

// This test ensures PipeOutboundToInbound goroutine is unblocked on error in PipeInboundToOutbound.
func TestGrpc2Http_ErrorReceivingHeader(t *testing.T) {
	ctrl := gomock.NewController(t)
	server := mock_rpc.NewMockInboundGrpcToOutboundHttpStream(ctrl)
	server.EXPECT().
		Context().
		Return(context.Background()).
		MinTimes(1)
	server.EXPECT().
		RecvMsg(gomock.Any()).
		Return(status.Error(codes.DataLoss, "recv failed"))
	grpc2http := grpctool.InboundGrpcToOutboundHttp{
		Log: zaptest.NewLogger(t),
		HandleProcessingError: func(msg string, err error) {
			t.Error(msg, err)
		},
		HandleIoError: func(msg string, err error) error {
			t.Error(msg, err)
			return nil
		},
		HttpDo: func(ctx context.Context, header *grpctool.HttpRequest_Header, body io.Reader) (grpctool.DoResponse, error) {
			t.FailNow()
			return grpctool.DoResponse{}, nil
		},
	}
	err := grpc2http.Pipe(server)
	require.EqualError(t, err, "rpc error: code = DataLoss desc = recv failed")
}

// This test ensures PipeInboundToOutbound goroutine is unblocked on error in PipeOutboundToInbound.
func TestGrpc2Http_ErrorReceivingHttpResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	server := mock_rpc.NewMockInboundGrpcToOutboundHttpStream(ctrl)
	server.EXPECT().
		Context().
		Return(context.Background()).
		MinTimes(1)
	block := make(chan struct{})
	gomock.InOrder(
		server.EXPECT().
			RecvMsg(gomock.Any()).
			Do(testhelpers.RecvMsg(&grpctool.HttpRequest{
				Message: &grpctool.HttpRequest_Header_{
					Header: sendUpgradeHeader(),
				},
			})),
		server.EXPECT().
			RecvMsg(gomock.Any()).
			Do(testhelpers.RecvMsg(&grpctool.HttpRequest{
				Message: &grpctool.HttpRequest_Trailer_{
					Trailer: &grpctool.HttpRequest_Trailer{},
				},
			})),
		server.EXPECT().
			RecvMsg(gomock.Any()).
			Do(func(msg interface{}) {
				testhelpers.SetValue(msg, &grpctool.HttpRequest{
					Message: &grpctool.HttpRequest_UpgradeData_{
						UpgradeData: &grpctool.HttpRequest_UpgradeData{
							Data: []byte(requestUpgradeBodyData),
						},
					},
				})
				close(block)
			}),
	)
	grpc2http := makeGrpc2http(t, func(ctx context.Context, header *grpctool.HttpRequest_Header, body io.Reader) (grpctool.DoResponse, error) {
		return grpctool.DoResponse{}, status.Error(codes.DataLoss, "DO failed")
	})
	err := grpc2http.Pipe(server)
	require.EqualError(t, err, "rpc error: code = DataLoss desc = DO failed")
	<-block // wait for the last recv to avoid a race
}

func makeGrpc2http(t *testing.T, f grpctool.HttpDo) *grpctool.InboundGrpcToOutboundHttp {
	return &grpctool.InboundGrpcToOutboundHttp{
		Log: zaptest.NewLogger(t),
		HandleProcessingError: func(msg string, err error) {
			t.Error(msg, err)
		},
		HandleIoError: func(msg string, err error) error {
			t.Error(msg, err)
			return nil
		},
		HttpDo: f,
	}
}

func sendHeader() *grpctool.HttpRequest_Header {
	return &grpctool.HttpRequest_Header{
		Request: &prototool.HttpRequest{
			Method: "BOOM",
			Header: map[string]*prototool.Values{
				"cxv": {
					Value: []string{"xx"},
				},
			},
			UrlPath: "/asd/asd/asd",
			Query: map[string]*prototool.Values{
				"adasd": {
					Value: []string{"a"},
				},
			},
		},
		Extra: &anypb.Any{
			TypeUrl: "sadfasdfasdfads",
			Value:   []byte{1, 2, 3},
		},
	}
}

func sendUpgradeHeader() *grpctool.HttpRequest_Header {
	sh := sendHeader()
	sh.Request.Header[httpz.UpgradeHeader] = &prototool.Values{
		Value: []string{"a"},
	}
	sh.Request.Header[httpz.ConnectionHeader] = &prototool.Values{
		Value: []string{"upgrade"},
	}
	return sh
}

func mockRecvStream(server *mock_rpc.MockInboundGrpcToOutboundHttpStream, eof bool, msgs ...proto.Message) []*gomock.Call {
	res := make([]*gomock.Call, 0, len(msgs)+1)
	for _, msg := range msgs {
		call := server.EXPECT().
			RecvMsg(gomock.Any()).
			Do(testhelpers.RecvMsg(msg))
		res = append(res, call)
	}
	if eof {
		call := server.EXPECT().
			RecvMsg(gomock.Any()).
			Return(io.EOF)
		res = append(res, call)
	}
	return res
}

func mockSendStream(t *testing.T, server *mock_rpc.MockInboundGrpcToOutboundHttpStream, msgs ...*grpctool.HttpResponse) []*gomock.Call {
	res := make([]*gomock.Call, 0, len(msgs))
	for _, msg := range msgs {
		call := server.EXPECT().
			Send(matcher.ProtoEq(t, msg))
		res = append(res, call)
	}
	return res
}
