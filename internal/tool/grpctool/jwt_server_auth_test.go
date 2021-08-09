package grpctool_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/testing/mock_rpc"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	jwtAudience  = "valid_audience"
	jwtIssuer    = "valid_issuer"
	jwtValidFor  = 5 * time.Second
	jwtNotBefore = 1 * time.Second

	expectedReq    int32 = 345
	expectedResult int16 = 125

	expectedSrv int32 = 13123
)

var (
	secret = []byte("some random secret")
)

func TestJWTServerAuth(t *testing.T) {
	unaryInfo := &grpc.UnaryServerInfo{
		FullMethod: "bla",
	}
	streamInfo := &grpc.StreamServerInfo{
		FullMethod: "bla",
	}
	t.Run("happy path", func(t *testing.T) {
		jwtAuther := setupAuther()

		now := time.Now()
		claims := validClams(now)
		signedClaims, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
		require.NoError(t, err)

		ctx := incomingCtx(context.Background(), t, signedClaims)
		result, err := jwtAuther.UnaryServerInterceptor(ctx, expectedReq, unaryInfo, unaryHandler(ctx, t))
		require.NoError(t, err)
		assert.Equal(t, expectedResult, result)

		ctrl := gomock.NewController(t)
		stream := mock_rpc.NewMockServerStream(ctrl)
		stream.EXPECT().Context().Return(ctx)
		err = jwtAuther.StreamServerInterceptor(expectedSrv, stream, streamInfo, streamHandler(stream, t))
		require.NoError(t, err)
	})
	t.Run("missing header", func(t *testing.T) {
		jwtAuther := setupAuther()

		ctx := grpctool.InjectLogger(context.Background(), zaptest.NewLogger(t))
		_, err := jwtAuther.UnaryServerInterceptor(ctx, expectedReq, unaryInfo, unaryMustNotBeCalled(t))
		require.EqualError(t, err, "rpc error: code = Unauthenticated desc = Request unauthenticated with bearer")

		ctrl := gomock.NewController(t)
		stream := mock_rpc.NewMockServerStream(ctrl)
		stream.EXPECT().Context().Return(ctx)
		err = jwtAuther.StreamServerInterceptor(expectedSrv, stream, streamInfo, streamMustNotBeCalled(t))
		require.EqualError(t, err, "rpc error: code = Unauthenticated desc = Request unauthenticated with bearer")
	})
	t.Run("invalid token type", func(t *testing.T) {
		jwtAuther := setupAuther()

		now := time.Now()
		claims := validClams(now)
		signedClaims, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
		require.NoError(t, err)

		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "weird_type "+signedClaims))
		ctx = grpctool.InjectLogger(ctx, zaptest.NewLogger(t))
		_, err = jwtAuther.UnaryServerInterceptor(ctx, expectedReq, unaryInfo, unaryMustNotBeCalled(t))
		require.EqualError(t, err, "rpc error: code = Unauthenticated desc = Request unauthenticated with bearer")

		ctrl := gomock.NewController(t)
		stream := mock_rpc.NewMockServerStream(ctrl)
		stream.EXPECT().Context().Return(ctx)
		err = jwtAuther.StreamServerInterceptor(expectedSrv, stream, streamInfo, streamMustNotBeCalled(t))
		require.EqualError(t, err, "rpc error: code = Unauthenticated desc = Request unauthenticated with bearer")
	})
	t.Run("unexpected signing algorithm", func(t *testing.T) {
		jwtAuther := setupAuther()

		now := time.Now()
		claims := validClams(now)

		keyData, err := os.ReadFile("testdata/sample_key")
		require.NoError(t, err)
		rsaKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
		require.NoError(t, err)
		signedClaims, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(rsaKey)
		require.NoError(t, err)

		assertValidationFailed(t, signedClaims, jwtAuther, "JWT validation failed")
	})
	t.Run("none signing algorithm", func(t *testing.T) {
		jwtAuther := setupAuther()

		now := time.Now()
		claims := validClams(now)
		signedClaims, err := jwt.NewWithClaims(jwt.SigningMethodNone, claims).SignedString(jwt.UnsafeAllowNoneSignatureType)
		require.NoError(t, err)

		assertValidationFailed(t, signedClaims, jwtAuther, "JWT validation failed")
	})
	t.Run("unexpected audience", func(t *testing.T) {
		jwtAuther := setupAuther()

		now := time.Now()
		claims := validClams(now)
		claims.Audience = "blablabla"
		signedClaims, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
		require.NoError(t, err)

		assertValidationFailed(t, signedClaims, jwtAuther, "JWT audience validation failed")
	})
	t.Run("unexpected issuer", func(t *testing.T) {
		jwtAuther := setupAuther()

		now := time.Now()
		claims := validClams(now)
		claims.Issuer = "blablabla"
		signedClaims, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
		require.NoError(t, err)

		assertValidationFailed(t, signedClaims, jwtAuther, "JWT issuer validation failed")
	})
}

func setupAuther() *grpctool.JWTAuther {
	return grpctool.NewJWTAuther(secret, jwtIssuer, jwtAudience)
}

func assertValidationFailed(t *testing.T, signedClaims string, jwtAuther *grpctool.JWTAuther, errStr string) {
	unaryInfo := &grpc.UnaryServerInfo{
		FullMethod: "bla",
	}
	streamInfo := &grpc.StreamServerInfo{
		FullMethod: "bla",
	}
	ctx := incomingCtx(context.Background(), t, signedClaims)
	_, err := jwtAuther.UnaryServerInterceptor(ctx, expectedReq, unaryInfo, unaryMustNotBeCalled(t))
	require.EqualError(t, err, "rpc error: code = Unauthenticated desc = "+errStr)

	ctrl := gomock.NewController(t)
	stream := mock_rpc.NewMockServerStream(ctrl)
	stream.EXPECT().Context().Return(ctx)
	err = jwtAuther.StreamServerInterceptor(expectedSrv, stream, streamInfo, streamMustNotBeCalled(t))
	require.EqualError(t, err, "rpc error: code = Unauthenticated desc = "+errStr)
}

func validClams(now time.Time) jwt.StandardClaims {
	claims := jwt.StandardClaims{
		Audience:  jwtAudience,
		ExpiresAt: now.Add(jwtValidFor).Unix(),
		IssuedAt:  now.Unix(),
		Issuer:    jwtIssuer,
		NotBefore: now.Add(-jwtNotBefore).Unix(),
	}
	return claims
}

func incomingCtx(ctx context.Context, t *testing.T, token string) context.Context {
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("authorization", "bearer "+token))
	ctx = grpctool.InjectLogger(ctx, zaptest.NewLogger(t))
	return ctx
}

func unaryHandler(expectedCtx context.Context, t *testing.T) grpc.UnaryHandler {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		assert.Equal(t, expectedReq, req)
		assert.Equal(t, expectedCtx, ctx)

		return expectedResult, nil
	}
}

func streamHandler(expectedStream grpc.ServerStream, t *testing.T) grpc.StreamHandler {
	return func(srv interface{}, stream grpc.ServerStream) error {
		assert.Equal(t, expectedSrv, srv)
		assert.Equal(t, expectedStream, stream)

		return nil
	}
}

func unaryMustNotBeCalled(t *testing.T) grpc.UnaryHandler {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		require.FailNow(t, "handler must not be called")
		return nil, nil
	}
}

func streamMustNotBeCalled(t *testing.T) grpc.StreamHandler {
	return func(srv interface{}, stream grpc.ServerStream) error {
		require.FailNow(t, "handler must not be called")
		return nil
	}
}
