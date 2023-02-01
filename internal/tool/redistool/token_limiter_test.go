package redistool

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/testhelpers"
	"go.uber.org/zap/zaptest"
)

const (
	ctxKey = 23124
)

func TestTokenLimiterHappyPath(t *testing.T) {
	ctx, _, mock, limiter, key := setup(t)

	mock.ExpectGet(key).SetVal("0")
	mock.ExpectTxPipeline()
	mock.ExpectIncr(key).SetVal(1)
	mock.ExpectExpire(key, 59*time.Second).SetVal(true)
	mock.ExpectTxPipelineExec()

	require.True(t, limiter.Allow(ctx), "Allow when no token has been consumed")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTokenLimiterOverLimit(t *testing.T) {
	ctx, _, mock, limiter, key := setup(t)

	mock.ExpectGet(key).SetVal("1")

	require.False(t, limiter.Allow(ctx), "Do not allow when a token has been consumed")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTokenLimiterNotAllowedWhenGetError(t *testing.T) {
	ctx, rpcApi, mock, limiter, key := setup(t)
	err := errors.New("test connection error")
	mock.ExpectGet(key).SetErr(err)

	rpcApi.EXPECT().
		HandleProcessingError("redistool.TokenLimiter: error retrieving minute bucket count", err)

	require.False(t, limiter.Allow(ctx), "Do not allow when there is a connection error")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTokenLimiterNotAllowedWhenIncrError(t *testing.T) {
	ctx, rpcApi, mock, limiter, key := setup(t)

	mock.ExpectGet(key).SetVal("0")
	mock.ExpectTxPipeline()
	mock.ExpectIncr(key).SetVal(1)
	err := errors.New("test connection error")
	mock.ExpectExpire(key, 59*time.Second).SetErr(err)

	rpcApi.EXPECT().
		HandleProcessingError("redistool.TokenLimiter: error while incrementing token key count", err)

	require.False(t, limiter.Allow(ctx), "Do not allow when there is a connection error")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func setup(t *testing.T) (context.Context, *MockRpcApi, redismock.ClientMock, *TokenLimiter, string) {
	client, mock := redismock.NewClientMock()
	mock.MatchExpectationsInOrder(true)
	ctrl := gomock.NewController(t)
	rpcApi := NewMockRpcApi(ctrl)
	rpcApi.EXPECT().
		Log().
		Return(zaptest.NewLogger(t)).
		AnyTimes()
	limiter := NewTokenLimiter(client, "key_prefix", 1, func(ctx context.Context) RpcApi {
		rpcApi.EXPECT().
			RequestKey().
			Return(api.AgentToken2key(ctx.Value(ctxKey).(api.AgentToken)))
		return rpcApi
	})
	ctx := context.WithValue(context.Background(), ctxKey, testhelpers.AgentkToken) // nolint: staticcheck
	key := limiter.buildKey(api.AgentToken2key(testhelpers.AgentkToken))
	return ctx, rpcApi, mock, limiter, key
}
