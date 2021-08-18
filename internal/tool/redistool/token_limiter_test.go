package redistool

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/testing/testhelpers"
	"go.uber.org/zap/zaptest"
)

func TestTokenLimiterHappyPath(t *testing.T) {
	ctx, mock, limiter, key := setup(t)

	mock.ExpectGet(key).SetVal("0")
	mock.ExpectTxPipeline()
	mock.ExpectIncr(key).SetVal(1)
	mock.ExpectExpire(key, 59*time.Second).SetVal(true)
	mock.ExpectTxPipelineExec()

	require.True(t, limiter.Allow(ctx), "Allow when no token has been consumed")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTokenLimiterOverLimit(t *testing.T) {
	ctx, mock, limiter, key := setup(t)

	mock.ExpectGet(key).SetVal("1")

	require.False(t, limiter.Allow(ctx), "Do not allow when a token has been consumed")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTokenLimiterNotAllowedWhenGetError(t *testing.T) {
	ctx, mock, limiter, key := setup(t)
	mock.ExpectGet(key).SetErr(errors.New("test connection error"))

	require.False(t, limiter.Allow(ctx), "Do not allow when there is a connection error")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTokenLimiterNotAllowedWhenIncrError(t *testing.T) {
	ctx, mock, limiter, key := setup(t)

	mock.ExpectGet(key).SetVal("0")
	mock.ExpectTxPipeline()
	mock.ExpectIncr(key).SetVal(1)
	mock.ExpectExpire(key, 59*time.Second).SetErr(errors.New("test connection error"))

	require.False(t, limiter.Allow(ctx), "Do not allow when there is a connection error")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func setup(t *testing.T) (context.Context, redismock.ClientMock, *TokenLimiter, string) {
	client, mock := redismock.NewClientMock()
	mock.MatchExpectationsInOrder(true)
	limiter := NewTokenLimiter(zaptest.NewLogger(t), client, "key_prefix", 1, tokenFromContext)
	ctx := context.WithValue(context.Background(), 123, testhelpers.AgentkToken) // nolint: staticcheck
	key := limiter.buildKey(string(testhelpers.AgentkToken))
	return ctx, mock, limiter, key
}

func tokenFromContext(ctx context.Context) string {
	return string(ctx.Value(123).(api.AgentToken))
}
