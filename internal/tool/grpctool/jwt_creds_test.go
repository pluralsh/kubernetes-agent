package grpctool

import (
	"context"
	"testing"

	"github.com/pluralsh/kuberentes-agent/internal/tool/grpctool/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/apimachinery/pkg/util/wait"
)

var (
	_ credentials.PerRPCCredentials = &JwtCredentials{}
)

const (
	secret   = "dfjnfkadskfadsnfkjadsgkasdbg"
	audience = "fasfadsf"
	issuer   = "cbcxvbvxbxb"
)

func TestJwtCredentialsProducesValidToken(t *testing.T) {
	c := &JwtCredentials{
		Secret:   []byte(secret),
		Audience: audience,
		Issuer:   issuer,
		Insecure: true,
	}
	auther := NewJWTAuther([]byte(secret), issuer, audience, func(ctx context.Context) *zap.Logger {
		return zaptest.NewLogger(t)
	})
	listener := NewDialListener()

	srv := grpc.NewServer(
		grpc.ChainStreamInterceptor(
			auther.StreamServerInterceptor,
		),
		grpc.ChainUnaryInterceptor(
			auther.UnaryServerInterceptor,
		),
	)
	test.RegisterTestingServer(srv, &test.GrpcTestingServer{
		UnaryFunc: func(ctx context.Context, request *test.Request) (*test.Response, error) {
			return &test.Response{
				Message: &test.Response_Scalar{Scalar: 123},
			}, nil
		},
		StreamingFunc: func(server test.Testing_StreamingRequestResponseServer) error {
			return server.Send(&test.Response{
				Message: &test.Response_Scalar{Scalar: 123},
			})
		},
	})
	var wg wait.Group
	defer wg.Wait()
	defer srv.GracefulStop()
	wg.Start(func() {
		assert.NoError(t, srv.Serve(listener))
	})
	conn, err := grpc.DialContext(context.Background(), "passthrough:pipe",
		grpc.WithContextDialer(listener.DialContext),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(c),
	)
	require.NoError(t, err)
	defer conn.Close()
	client := test.NewTestingClient(conn)
	_, err = client.RequestResponse(context.Background(), &test.Request{})
	require.NoError(t, err)
	stream, err := client.StreamingRequestResponse(context.Background())
	require.NoError(t, err)
	_, err = stream.Recv()
	require.NoError(t, err)
}
