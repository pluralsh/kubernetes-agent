package kasapp

import (
	"context"
	"io"
	"net"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/reverse_tunnel"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/reverse_tunnel/info"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/reverse_tunnel/tracker"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/grpctool/test"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_reverse_tunnel"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_reverse_tunnel_tracker"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/testhelpers"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/tlstool"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	"k8s.io/apimachinery/pkg/util/wait"
)

var (
	_ kasRouter          = (*router)(nil)
	_ grpc.StreamHandler = (*router)(nil).RouteToCorrectKasHandler
	_ grpc.StreamHandler = (*router)(nil).RouteToCorrectAgentHandler
)

func TestRouter_UnaryHappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	unaryResponse := &test.Response{Message: &test.Response_Scalar{Scalar: 123}}
	routingMetadata := modserver.RoutingMetadata(testhelpers.AgentId)
	payloadMD, responseMD, trailersMD := meta()
	payloadReq := &test.Request{S1: "123"}
	var (
		headerResp  metadata.MD
		trailerResp metadata.MD
	)
	tunnel := mock_reverse_tunnel.NewMockTunnel(ctrl)
	tunnelForwardStream := tunnel.EXPECT().
		ForwardStream(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Do(forwardStream(t, routingMetadata, payloadMD, payloadReq, unaryResponse, responseMD, trailersMD))
	runRouterTest(t, tunnel, tunnelForwardStream, func(client test.TestingClient) {
		ctx := metadata.NewOutgoingContext(context.Background(), metadata.Join(routingMetadata, payloadMD))
		// grpc.Header() and grpc.Trailer are ok here because it's a unary RPC.
		response, err := client.RequestResponse(ctx, payloadReq, grpc.Header(&headerResp), grpc.Trailer(&trailerResp)) // nolint: forbidigo
		require.NoError(t, err)
		assert.Empty(t, cmp.Diff(response, unaryResponse, protocmp.Transform()))
		mdContains(t, responseMD, headerResp)
		mdContains(t, trailersMD, trailerResp)
	})
}

func TestRouter_UnaryImmediateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	routingMetadata := modserver.RoutingMetadata(testhelpers.AgentId)
	statusWithDetails, err := status.New(codes.InvalidArgument, "Some expected error").
		WithDetails(&test.Request{S1: "some details of the error"})
	require.NoError(t, err)
	tunnel := mock_reverse_tunnel.NewMockTunnel(ctrl)
	tunnelForwardStream := tunnel.EXPECT().
		ForwardStream(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(statusWithDetails.Err())
	runRouterTest(t, tunnel, tunnelForwardStream, func(client test.TestingClient) {
		ctx := metadata.NewOutgoingContext(context.Background(), routingMetadata)
		_, err = client.RequestResponse(ctx, &test.Request{S1: "123"})
		require.Error(t, err)
		receivedStatus := status.Convert(err).Proto()
		assert.Empty(t, cmp.Diff(receivedStatus, statusWithDetails.Proto(), protocmp.Transform()))
	})
}

func TestRouter_UnaryErrorAfterHeader(t *testing.T) {
	ctrl := gomock.NewController(t)
	routingMetadata := modserver.RoutingMetadata(testhelpers.AgentId)
	payloadMD, responseMD, trailersMD := meta()
	statusWithDetails, err := status.New(codes.InvalidArgument, "Some expected error").
		WithDetails(&test.Request{S1: "some details of the error"})
	require.NoError(t, err)
	var (
		headerResp  metadata.MD
		trailerResp metadata.MD
	)
	tunnel := mock_reverse_tunnel.NewMockTunnel(ctrl)
	tunnelForwardStream := tunnel.EXPECT().
		ForwardStream(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(log *zap.Logger, rpcApi reverse_tunnel.RpcApi, incomingStream grpc.ServerStream, cb reverse_tunnel.TunnelDataCallback) error {
			verifyMeta(t, incomingStream, routingMetadata, payloadMD)
			assert.NoError(t, cb.Header(grpctool.MetaToValuesMap(responseMD)))
			assert.NoError(t, cb.Trailer(grpctool.MetaToValuesMap(trailersMD)))
			return statusWithDetails.Err()
		})
	runRouterTest(t, tunnel, tunnelForwardStream, func(client test.TestingClient) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		ctx = metadata.NewOutgoingContext(ctx, metadata.Join(routingMetadata, payloadMD))
		// grpc.Header() and grpc.Trailer are ok here because it's a unary RPC.
		_, err := client.RequestResponse(ctx, &test.Request{S1: "123"}, grpc.Header(&headerResp), grpc.Trailer(&trailerResp)) // nolint: forbidigo
		require.Error(t, err)
		receivedStatus := status.Convert(err).Proto()
		assert.Empty(t, cmp.Diff(receivedStatus, statusWithDetails.Proto(), protocmp.Transform()))
		mdContains(t, responseMD, headerResp)
		mdContains(t, trailersMD, trailerResp)
	})
}

func TestRouter_StreamHappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	streamResponse := &test.Response{Message: &test.Response_Scalar{Scalar: 123}}
	routingMetadata := modserver.RoutingMetadata(testhelpers.AgentId)
	payloadMD, responseMD, trailersMD := meta()
	payloadReq := &test.Request{S1: "123"}
	tunnel := mock_reverse_tunnel.NewMockTunnel(ctrl)
	tunnelForwardStream := tunnel.EXPECT().
		ForwardStream(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Do(forwardStream(t, routingMetadata, payloadMD, payloadReq, streamResponse, responseMD, trailersMD))
	runRouterTest(t, tunnel, tunnelForwardStream, func(client test.TestingClient) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		ctx = metadata.NewOutgoingContext(ctx, metadata.Join(routingMetadata, payloadMD))
		stream, err := client.StreamingRequestResponse(ctx)
		require.NoError(t, err)
		err = stream.Send(payloadReq)
		require.NoError(t, err)
		err = stream.CloseSend()
		require.NoError(t, err)
		response, err := stream.Recv()
		require.NoError(t, err)
		assert.Empty(t, cmp.Diff(response, streamResponse, protocmp.Transform()))
		_, err = stream.Recv()
		assert.Equal(t, io.EOF, err)
		verifyHeaderAndTrailer(t, stream, responseMD, trailersMD)
	})
}

func TestRouter_StreamImmediateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	routingMetadata := modserver.RoutingMetadata(testhelpers.AgentId)
	statusWithDetails, err := status.New(codes.InvalidArgument, "Some expected error").
		WithDetails(&test.Request{S1: "some details of the error"})
	require.NoError(t, err)
	tunnel := mock_reverse_tunnel.NewMockTunnel(ctrl)
	tunnelForwardStream := tunnel.EXPECT().
		ForwardStream(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(statusWithDetails.Err())
	runRouterTest(t, tunnel, tunnelForwardStream, func(client test.TestingClient) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		ctx = metadata.NewOutgoingContext(ctx, routingMetadata)
		stream, err := client.StreamingRequestResponse(ctx)
		require.NoError(t, err)
		_, err = stream.Recv()
		require.Error(t, err)
		receivedStatus := status.Convert(err).Proto()
		assert.Empty(t, cmp.Diff(receivedStatus, statusWithDetails.Proto(), protocmp.Transform()))
	})
}

func TestRouter_StreamErrorAfterHeader(t *testing.T) {
	ctrl := gomock.NewController(t)
	routingMetadata := modserver.RoutingMetadata(testhelpers.AgentId)
	payloadMD, responseMD, trailersMD := meta()
	statusWithDetails, err := status.New(codes.InvalidArgument, "Some expected error").
		WithDetails(&test.Request{S1: "some details of the error"})
	require.NoError(t, err)
	tunnel := mock_reverse_tunnel.NewMockTunnel(ctrl)
	tunnelForwardStream := tunnel.EXPECT().
		ForwardStream(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(log *zap.Logger, rpcApi reverse_tunnel.RpcApi, incomingStream grpc.ServerStream, cb reverse_tunnel.TunnelDataCallback) error {
			verifyMeta(t, incomingStream, routingMetadata, payloadMD)
			assert.NoError(t, cb.Header(grpctool.MetaToValuesMap(responseMD)))
			assert.NoError(t, cb.Trailer(grpctool.MetaToValuesMap(trailersMD)))
			return statusWithDetails.Err()
		})
	runRouterTest(t, tunnel, tunnelForwardStream, func(client test.TestingClient) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		ctx = metadata.NewOutgoingContext(ctx, metadata.Join(routingMetadata, payloadMD))
		stream, err := client.StreamingRequestResponse(ctx)
		require.NoError(t, err)
		_, err = stream.Recv()
		require.Error(t, err)
		receivedStatus := status.Convert(err).Proto()
		assert.Empty(t, cmp.Diff(receivedStatus, statusWithDetails.Proto(), protocmp.Transform()))
		verifyHeaderAndTrailer(t, stream, responseMD, trailersMD)
	})
}

func TestRouter_StreamVisitorErrorAfterErrorMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	routingMetadata := modserver.RoutingMetadata(testhelpers.AgentId)
	payloadMD, responseMD, trailersMD := meta()
	statusWithDetails, err := status.New(codes.InvalidArgument, "Some expected error").
		WithDetails(&test.Request{S1: "some details of the error"})
	require.NoError(t, err)
	tunnel := mock_reverse_tunnel.NewMockTunnel(ctrl)
	tunnelForwardStream := tunnel.EXPECT().
		ForwardStream(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(log *zap.Logger, rpcApi reverse_tunnel.RpcApi, incomingStream grpc.ServerStream, cb reverse_tunnel.TunnelDataCallback) error {
			verifyMeta(t, incomingStream, routingMetadata, payloadMD)
			assert.NoError(t, cb.Header(grpctool.MetaToValuesMap(responseMD)))
			assert.NoError(t, cb.Trailer(grpctool.MetaToValuesMap(trailersMD)))
			assert.NoError(t, cb.Error(statusWithDetails.Proto()))
			return status.Error(codes.Unavailable, "expected return error")
		})
	runRouterTest(t, tunnel, tunnelForwardStream, func(client test.TestingClient) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		ctx = metadata.NewOutgoingContext(ctx, metadata.Join(routingMetadata, payloadMD))
		stream, err := client.StreamingRequestResponse(ctx)
		require.NoError(t, err)
		_, err = stream.Recv()
		require.EqualError(t, err, "rpc error: code = Unavailable desc = expected return error")
		verifyHeaderAndTrailer(t, stream, responseMD, trailersMD)
	})
}

func tunnelInfo() *tracker.TunnelInfo {
	return &tracker.TunnelInfo{
		AgentDescriptor: &info.AgentDescriptor{
			Services: []*info.Service{
				{
					Name: "gitlab.agent.grpctool.test.Testing",
					Methods: []*info.Method{
						{
							Name: "RequestResponse",
						},
						{
							Name: "StreamingRequestResponse",
						},
					},
				},
			},
		},
		ConnectionId: 1312312313,
		AgentId:      testhelpers.AgentId,
		KasUrl:       "grpc://pipe",
	}
}

func meta() (metadata.MD, metadata.MD, metadata.MD) {
	payloadMD := metadata.Pairs("key1", "value1")
	responseMD := metadata.Pairs("key2", "value2")
	trailersMD := metadata.Pairs("key3", "value3")
	return payloadMD, responseMD, trailersMD
}

func verifyHeaderAndTrailer(t *testing.T, stream grpc.ClientStream, responseMD, trailersMD metadata.MD) {
	headerResp, err := stream.Header()
	require.NoError(t, err)
	mdContains(t, responseMD, headerResp)
	mdContains(t, trailersMD, stream.Trailer())
}

func forwardStream(t *testing.T, routingMetadata, payloadMD metadata.MD, payloadReq *test.Request, response *test.Response, responseMD, trailersMD metadata.MD) func(*zap.Logger, reverse_tunnel.RpcApi, grpc.ServerStream, reverse_tunnel.TunnelDataCallback) {
	return func(log *zap.Logger, rpcApi reverse_tunnel.RpcApi, incomingStream grpc.ServerStream, cb reverse_tunnel.TunnelDataCallback) {
		verifyMeta(t, incomingStream, routingMetadata, payloadMD)
		var req test.Request
		err := incomingStream.RecvMsg(&req)
		assert.NoError(t, err)
		assert.Empty(t, cmp.Diff(payloadReq, &req, protocmp.Transform()))
		data, err := proto.Marshal(response)
		assert.NoError(t, err)
		assert.NoError(t, cb.Header(grpctool.MetaToValuesMap(responseMD)))
		assert.NoError(t, cb.Message(data))
		assert.NoError(t, cb.Trailer(grpctool.MetaToValuesMap(trailersMD)))
	}
}

func verifyMeta(t *testing.T, incomingStream grpc.ServerStream, routingMetadata, payloadMD metadata.MD) {
	md, _ := metadata.FromIncomingContext(incomingStream.Context())
	for k := range routingMetadata { // no routing metadata is passed to the agent
		assert.NotContains(t, md, k)
	}
	mdContains(t, payloadMD, md)
}

func mdContains(t *testing.T, expectedMd metadata.MD, actualMd metadata.MD) {
	for k, v := range expectedMd {
		assert.Equalf(t, v, actualMd[k], "key: %s", k)
	}
}

// test:client(default codec) --> kas1:internal server(raw codec) --> router_kas handler -->
// client from kas_pool(raw wih fallback codec) --> kas2:private server(raw wih fallback codec) -->
// router_agent handler --> tunnel finder --> tunnel.ForwardStream()
func runRouterTest(t *testing.T, tunnel *mock_reverse_tunnel.MockTunnel, tunnelForwardStream *gomock.Call, runTest func(client test.TestingClient)) {
	ctrl := gomock.NewController(t)
	log := zaptest.NewLogger(t)
	querier := mock_reverse_tunnel_tracker.NewMockQuerier(ctrl)
	finder := mock_reverse_tunnel.NewMockTunnelFinder(ctrl)
	internalServerListener := grpctool.NewDialListener()
	defer internalServerListener.Close()
	privateApiServerListener := grpctool.NewDialListener()
	defer privateApiServerListener.Close()

	gomock.InOrder(
		querier.EXPECT().
			GetTunnelsByAgentId(gomock.Any(), testhelpers.AgentId, gomock.Any()).
			DoAndReturn(func(ctx context.Context, agentId int64, cb tracker.GetTunnelsByAgentIdCallback) error {
				done, err := cb(tunnelInfo())
				assert.False(t, done)
				return err
			}),
		finder.EXPECT().
			FindTunnel(gomock.Any(), testhelpers.AgentId, gomock.Any(), gomock.Any()).
			Return(tunnel, nil),
		tunnelForwardStream,
		tunnel.EXPECT().Done(),
	)
	factory := func(ctx context.Context, fullMethodName string) modserver.RpcApi {
		return &serverRpcApi{
			RpcApiStub: modshared.RpcApiStub{
				StreamCtx: ctx,
				Logger:    log,
			},
			sentryHubRoot: sentry.NewHub(nil, sentry.NewScope()),
		}
	}

	internalServer := grpc.NewServer(
		grpc.StatsHandler(grpctool.NewServerMaxConnAgeStatsHandler(context.Background(), 0)),
		grpc.ChainStreamInterceptor(
			modserver.StreamRpcApiInterceptor(factory),
		),
		grpc.ChainUnaryInterceptor(
			modserver.UnaryRpcApiInterceptor(factory),
		),
		grpc.ForceServerCodec(grpctool.RawCodec{}),
	)
	privateApiServer := grpc.NewServer(
		grpc.StatsHandler(grpctool.NewServerMaxConnAgeStatsHandler(context.Background(), 0)),
		grpc.ChainStreamInterceptor(
			modserver.StreamRpcApiInterceptor(factory),
		),
		grpc.ChainUnaryInterceptor(
			modserver.UnaryRpcApiInterceptor(factory),
		),
		grpc.ForceServerCodec(grpctool.RawCodecWithProtoFallback{}),
	)
	gatewayKasVisitor, err := grpctool.NewStreamVisitor(&GatewayKasResponse{})
	require.NoError(t, err)
	r := &router{
		kasPool: grpctool.NewPool(log,
			credentials.NewTLS(tlstool.DefaultClientTLSConfig()),
			grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
				if addr == "self" {
					<-ctx.Done()
					return nil, ctx.Err()
				} else {
					return privateApiServerListener.DialContext(ctx, addr)
				}
			}),
		),
		tunnelQuerier:             querier,
		tunnelFinder:              finder,
		ownPrivateApiUrl:          selfAddr,
		pollConfig:                testhelpers.NewPollConfig(time.Minute),
		internalServer:            internalServer,
		privateApiServer:          privateApiServer,
		gatewayKasVisitor:         gatewayKasVisitor,
		kasRoutingDurationSuccess: prometheus.ObserverFunc(func(f float64) {}),
		kasRoutingDurationError:   prometheus.ObserverFunc(func(f float64) {}),
	}
	r.RegisterAgentApi(&test.Testing_ServiceDesc)
	var wg wait.Group
	defer wg.Wait()
	defer internalServer.GracefulStop()
	defer privateApiServer.GracefulStop()
	wg.Start(func() {
		assert.NoError(t, internalServer.Serve(internalServerListener))
	})
	wg.Start(func() {
		assert.NoError(t, privateApiServer.Serve(privateApiServerListener))
	})
	internalServerConn, err := grpc.DialContext(context.Background(), "pipe",
		grpc.WithContextDialer(internalServerListener.DialContext),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainStreamInterceptor(
			grpctool.StreamClientValidatingInterceptor,
		),
		grpc.WithChainUnaryInterceptor(
			grpctool.UnaryClientValidatingInterceptor,
		),
	)
	require.NoError(t, err)
	defer internalServerConn.Close()
	client := test.NewTestingClient(internalServerConn)
	runTest(client)
}
