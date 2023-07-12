package redistool

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/rueidis"
	rmock "github.com/redis/rueidis/mock"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/testhelpers"
	"go.uber.org/zap/zaptest"
)

const (
	ctxKey = 23124
)

func BenchmarkBuildTokenLimiterKey(b *testing.B) {
	b.ReportAllocs()
	const prefix = "pref"
	var sink string
	requestKey := []byte{1, 2, 3, 4}
	for i := 0; i < b.N; i++ {
		sink = buildTokenLimiterKey(prefix, requestKey)
	}
	_ = sink
}

func TestTokenLimiterHappyPath(t *testing.T) {
	ctx, _, client, limiter, key := setup(t)

	client.EXPECT().
		Do(gomock.Any(), rmock.Match("GET", key)).
		Return(rmock.Result(rmock.RedisInt64(0)))
	client.EXPECT().
		DoMulti(gomock.Any(),
			rmock.Match("MULTI"),
			rmock.Match("INCR", key),
			rmock.Match("EXPIRE", key, "59"),
			rmock.Match("EXEC"),
		)

	require.True(t, limiter.Allow(ctx), "Allow when no token has been consumed")
}

func TestTokenLimiterOverLimit(t *testing.T) {
	ctx, _, client, limiter, key := setup(t)

	client.EXPECT().
		Do(gomock.Any(), rmock.Match("GET", key)).
		Return(rmock.Result(rmock.RedisInt64(1)))

	require.False(t, limiter.Allow(ctx), "Do not allow when a token has been consumed")
}

func TestTokenLimiterNotAllowedWhenGetError(t *testing.T) {
	ctx, rpcApi, client, limiter, key := setup(t)
	err := errors.New("test connection error")
	client.EXPECT().
		Do(gomock.Any(), rmock.Match("GET", key)).
		Return(rmock.ErrorResult(err))

	rpcApi.EXPECT().
		HandleProcessingError("redistool.TokenLimiter: error retrieving minute bucket count", err)

	require.False(t, limiter.Allow(ctx), "Do not allow when there is a connection error")
}

func TestTokenLimiterNotAllowedWhenIncrError(t *testing.T) {
	err := errors.New("test connection error")
	ctx, rpcApi, client, limiter, key := setup(t)

	client.EXPECT().
		Do(gomock.Any(), rmock.Match("GET", key)).
		Return(rmock.Result(rmock.RedisInt64(0)))
	client.EXPECT().
		DoMulti(gomock.Any(),
			rmock.Match("MULTI"),
			rmock.Match("INCR", key),
			rmock.Match("EXPIRE", key, "59"),
			rmock.Match("EXEC"),
		).
		Return([]rueidis.RedisResult{rmock.ErrorResult(err)})
	rpcApi.EXPECT().
		HandleProcessingError("redistool.TokenLimiter: error while incrementing token key count", err)

	require.False(t, limiter.Allow(ctx), "Do not allow when there is a connection error")
}

func setup(t *testing.T) (context.Context, *MockRpcApi, *rmock.Client, *TokenLimiter, string) {
	ctrl := gomock.NewController(t)
	client := rmock.NewClient(ctrl)
	rpcApi := NewMockRpcApi(ctrl)
	rpcApi.EXPECT().
		Log().
		Return(zaptest.NewLogger(t)).
		AnyTimes()

	limiter := NewTokenLimiter(client, "key_prefix", 1,
		prometheus.NewCounter(prometheus.CounterOpts{
			Name: "test",
		}),
		func(ctx context.Context) RpcApi {
			rpcApi.EXPECT().
				RequestKey().
				Return(api.AgentToken2key(ctx.Value(ctxKey).(api.AgentToken)))
			return rpcApi
		})
	ctx := context.WithValue(context.Background(), ctxKey, testhelpers.AgentkToken) // nolint: staticcheck
	key := buildTokenLimiterKey(limiter.keyPrefix, api.AgentToken2key(testhelpers.AgentkToken))
	return ctx, rpcApi, client, limiter, key
}
