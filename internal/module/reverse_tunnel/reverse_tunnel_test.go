package reverse_tunnel_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/module/reverse_tunnel"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/tool/grpctool/test"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/tool/testing/testhelpers"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/testing/protocmp"
)

const (
	scalarNumber protoreflect.FieldNumber = 1
	x1Number     protoreflect.FieldNumber = 2
	dataNumber   protoreflect.FieldNumber = 3
	lastNumber   protoreflect.FieldNumber = 4

	metaKey    = "Cba"
	trailerKey = "Abc"
)

func TestStreamHappyPath(t *testing.T) {
	trailer := metadata.MD{}
	trailer.Set(trailerKey, "1", "2")
	ats := &agentTestingServer{
		streamingRequestResponse: func(server test.Testing_StreamingRequestResponseServer) error {
			recv, err := server.Recv()
			if err != nil {
				return status.Error(codes.Unavailable, "unavailable")
			}
			val, err := strconv.ParseInt(recv.S1, 10, 64)
			if err != nil {
				return status.Error(codes.Unavailable, "unavailable")
			}
			incomingContext, ok := metadata.FromIncomingContext(server.Context())
			if !ok {
				return status.Error(codes.Unavailable, "unavailable")
			}

			header := metadata.MD{}
			header.Set(metaKey, incomingContext.Get(metaKey)...)

			err = server.SetHeader(header)
			if err != nil {
				return status.Error(codes.Unavailable, "unavailable")
			}
			resps := []*test.Response{
				{
					Message: &test.Response_Scalar{
						Scalar: val,
					},
				},
				{
					Message: &test.Response_X1{
						X1: test.Enum1_v1,
					},
				},
				{
					Message: &test.Response_Data_{
						Data: &test.Response_Data{},
					},
				},
				{
					Message: &test.Response_Data_{
						Data: &test.Response_Data{},
					},
				},
				{
					Message: &test.Response_Last_{
						Last: &test.Response_Last{},
					},
				},
			}
			for _, resp := range resps {
				err = server.Send(resp)
				if err != nil {
					return status.Error(codes.Unavailable, "unavailable")
				}
			}
			server.SetTrailer(trailer)
			return nil
		},
	}
	runTest(t, ats, func(ctx context.Context, t *testing.T, client test.TestingClient) {
		for i := 0; i < 2; i++ { // test several sequential requests
			testStreamHappyPath(ctx, t, client, trailer)
		}
	})
}

func testStreamHappyPath(ctx context.Context, t *testing.T, client test.TestingClient, trailer metadata.MD) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	meta := metadata.MD{}
	meta.Set(metaKey, "3", "4")
	ctx = metadata.NewOutgoingContext(ctx, meta)
	stream, err := client.StreamingRequestResponse(ctx)
	require.NoError(t, err)
	err = stream.Send(&test.Request{
		S1: "123",
	})
	require.NoError(t, err)
	err = stream.CloseSend()
	require.NoError(t, err)
	var (
		scalarCalled int
		x1Called     int
		dataCalled   int
		lastCalled   int
		eofCalled    int
	)
	v, err := grpctool.NewStreamVisitor(&test.Response{})
	require.NoError(t, err)
	err = v.Visit(stream,
		grpctool.WithEOFCallback(func() error {
			eofCalled++
			return nil
		}),
		grpctool.WithCallback(scalarNumber, func(scalar int64) error {
			assert.EqualValues(t, 123, scalar)
			scalarCalled++
			return nil
		}),
		grpctool.WithCallback(x1Number, func(x1 test.Enum1) error {
			x1Called++
			return nil
		}),
		grpctool.WithCallback(dataNumber, func(data *test.Response_Data) error {
			dataCalled++
			return nil
		}),
		grpctool.WithCallback(lastNumber, func(last *test.Response_Last) error {
			lastCalled++
			return nil
		}),
	)
	require.NoError(t, err)
	assert.Equal(t, 1, scalarCalled)
	assert.Equal(t, 1, x1Called)
	assert.Equal(t, 2, dataCalled)
	assert.Equal(t, 1, lastCalled)
	assert.Equal(t, 1, eofCalled)
	assert.Equal(t, trailer, stream.Trailer())
	header, err := stream.Header()
	require.NoError(t, err)
	assert.Equal(t, meta.Get(metaKey), header.Get(metaKey))
}

func TestUnaryHappyPath(t *testing.T) {
	ats := &agentTestingServer{
		requestResponse: func(ctx context.Context, request *test.Request) (*test.Response, error) {
			val, err := strconv.ParseInt(request.S1, 10, 64)
			if err != nil {
				return nil, status.Error(codes.Unavailable, "unavailable")
			}
			incomingContext, _ := metadata.FromIncomingContext(ctx)
			meta := metadata.MD{}
			meta.Set(metaKey, incomingContext.Get(metaKey)...)
			err = grpc.SetHeader(ctx, meta)
			if err != nil {
				return nil, status.Error(codes.Unavailable, "unavailable")
			}
			trailer := metadata.MD{}
			trailer.Set(trailerKey, "1", "2")
			err = grpc.SetTrailer(ctx, trailer)
			if err != nil {
				return nil, status.Error(codes.Unavailable, "unavailable")
			}
			return &test.Response{
				Message: &test.Response_Scalar{
					Scalar: val,
				},
			}, nil
		},
	}
	runTest(t, ats, func(ctx context.Context, t *testing.T, client test.TestingClient) {
		for i := 0; i < 2; i++ { // test several sequential requests
			testUnaryHappyPath(ctx, t, client)
		}
	})
}

func testUnaryHappyPath(ctx context.Context, t *testing.T, client test.TestingClient) {
	meta := metadata.MD{}
	meta.Set(metaKey, "3", "4")
	ctx = metadata.NewOutgoingContext(ctx, meta)
	var (
		headerResp  metadata.MD
		trailerResp metadata.MD
	)
	resp, err := client.RequestResponse(ctx, &test.Request{
		S1: "123",
	}, grpc.Header(&headerResp), grpc.Trailer(&trailerResp))
	require.NoError(t, err)
	assert.EqualValues(t, 123, resp.Message.(*test.Response_Scalar).Scalar)
	assert.Equal(t, meta.Get(metaKey), headerResp.Get(metaKey))
	trailer := metadata.MD{}
	trailer.Set(trailerKey, "1", "2")
	assert.Equal(t, trailer, trailerResp)
}

func TestStreamError(t *testing.T) {
	statusWithDetails, err := status.New(codes.InvalidArgument, "Some expected error").
		WithDetails(&test.Request{S1: "some details of the error"})
	require.NoError(t, err)
	ats := &agentTestingServer{
		streamingRequestResponse: func(server test.Testing_StreamingRequestResponseServer) error {
			return statusWithDetails.Err()
		},
	}
	runTest(t, ats, func(ctx context.Context, t *testing.T, client test.TestingClient) {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		stream, err := client.StreamingRequestResponse(ctx)
		require.NoError(t, err)
		_, err = stream.Recv()
		require.Error(t, err)
		receivedStatus := status.Convert(err).Proto()
		assert.Empty(t, cmp.Diff(receivedStatus, statusWithDetails.Proto(), protocmp.Transform()))
	})
}

func TestUnaryError(t *testing.T) {
	statusWithDetails, err := status.New(codes.InvalidArgument, "Some expected error").
		WithDetails(&test.Request{S1: "some details of the error"})
	require.NoError(t, err)
	ats := &agentTestingServer{
		requestResponse: func(ctx context.Context, request *test.Request) (*test.Response, error) {
			return nil, statusWithDetails.Err()
		},
	}
	runTest(t, ats, func(ctx context.Context, t *testing.T, client test.TestingClient) {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		_, err := client.RequestResponse(ctx, &test.Request{
			S1: "123",
		})
		require.Error(t, err)
		receivedStatus := status.Convert(err).Proto()
		assert.Empty(t, cmp.Diff(receivedStatus, statusWithDetails.Proto(), protocmp.Transform()))
	})
}

func runTest(t *testing.T, ats test.TestingServer, f func(context.Context, *testing.T, test.TestingClient)) {
	// Start/stop
	g, ctx := errgroup.WithContext(context.Background())
	defer func() {
		assert.NoError(t, g.Wait())
	}()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Construct server and agent components
	runServer, kasConn, serverInternalServerConn, _, serverApi := serverConstructComponents(t)
	runAgent, agentInternalServer := agentConstructComponents(t, kasConn)
	agentInfo := testhelpers.AgentInfoObj()

	serverApi.EXPECT().
		GetAgentInfo(gomock.Any(), gomock.Any(), testhelpers.AgentkToken, false).
		Return(agentInfo, nil, false).
		MinTimes(1)

	test.RegisterTestingServer(agentInternalServer, ats)

	// Run all
	g.Go(func() error {
		return runServer(ctx)
	})
	g.Go(func() error {
		return runAgent(ctx)
	})

	// Test
	client := test.NewTestingClient(serverInternalServerConn)
	f(ctx, t, client)
}

type serverTestingServer struct {
	streamForwarder reverse_tunnel.IncomingConnectionHandler
}

func (s *serverTestingServer) RequestResponse(srv interface{}, server grpc.ServerStream) error {
	return s.streamForwarder.HandleIncomingConnection(testhelpers.AgentId, server)
}

func (s *serverTestingServer) StreamingRequestResponse(srv interface{}, server grpc.ServerStream) error {
	return s.streamForwarder.HandleIncomingConnection(testhelpers.AgentId, server)
}

// registerTestingServer is a test.RegisterTestingServer clone that's been modified to be compatible with
// reverse_tunnel.IncomingConnectionHandler.HandleIncomingConnection().
func registerTestingServer(s *grpc.Server, h *serverTestingServer) {
	// ServiceDesc must match test._Testing_serviceDesc
	s.RegisterService(&grpc.ServiceDesc{
		ServiceName: "gitlab.agent.grpctool.test.Testing",
		Streams: []grpc.StreamDesc{
			{
				StreamName:    "RequestResponse",
				Handler:       h.RequestResponse,
				ServerStreams: true,
				ClientStreams: true,
			},
			{
				StreamName:    "StreamingRequestResponse",
				Handler:       h.StreamingRequestResponse,
				ServerStreams: true,
				ClientStreams: true,
			},
		},
		Metadata: "internal/tool/grpctool/test/test.proto",
	}, nil)
}

var (
	_ test.TestingServer = &agentTestingServer{}
)

type agentTestingServer struct {
	requestResponse          func(context.Context, *test.Request) (*test.Response, error)
	streamingRequestResponse func(test.Testing_StreamingRequestResponseServer) error
}

func (s *agentTestingServer) RequestResponse(ctx context.Context, request *test.Request) (*test.Response, error) {
	return s.requestResponse(ctx, request)
}

func (s *agentTestingServer) StreamingRequestResponse(server test.Testing_StreamingRequestResponseServer) error {
	return s.streamingRequestResponse(server)
}