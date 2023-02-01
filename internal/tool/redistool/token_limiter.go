package redistool

import (
	"context"
	"encoding/binary"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/logz"
	"go.uber.org/zap"
)

type RpcApi interface {
	Log() *zap.Logger
	HandleProcessingError(msg string, err error)
	RequestKey() []byte
}

// TokenLimiter is a redis-based rate limiter implementing the algorithm in https://redislabs.com/redis-best-practices/basic-rate-limiting/
type TokenLimiter struct {
	redisClient    redis.UniversalClient
	keyPrefix      string
	limitPerMinute uint64
	getApi         func(context.Context) RpcApi
}

// NewTokenLimiter returns a new TokenLimiter
func NewTokenLimiter(redisClient redis.UniversalClient, keyPrefix string,
	limitPerMinute uint64, getApi func(context.Context) RpcApi) *TokenLimiter {
	return &TokenLimiter{
		redisClient:    redisClient,
		keyPrefix:      keyPrefix,
		limitPerMinute: limitPerMinute,
		getApi:         getApi,
	}
}

// Allow consumes one limitable event from the token in the context
func (l *TokenLimiter) Allow(ctx context.Context) bool {
	api := l.getApi(ctx)
	key := l.buildKey(api.RequestKey())

	count, err := l.redisClient.Get(ctx, key).Uint64()
	if err != nil {
		if err != redis.Nil { // nolint:errorlint
			api.HandleProcessingError("redistool.TokenLimiter: error retrieving minute bucket count", err)
			return false
		}
		count = 0
	}
	if count >= l.limitPerMinute {
		api.Log().Debug("redistool.TokenLimiter: rate limit exceeded",
			logz.RedisKey([]byte(key)), logz.U64Count(count), logz.TokenLimit(l.limitPerMinute))
		return false
	}

	_, err = l.redisClient.TxPipelined(ctx, func(p redis.Pipeliner) error {
		p.Incr(ctx, key)
		p.Expire(ctx, key, 59*time.Second)
		return nil
	})
	if err != nil {
		api.HandleProcessingError("redistool.TokenLimiter: error while incrementing token key count", err)
		return false
	}

	return true
}

func (l *TokenLimiter) buildKey(requestKey []byte) string {
	currentMinute := time.Now().UTC().Minute()

	var result strings.Builder
	result.WriteString(l.keyPrefix)
	result.WriteByte(':')
	result.Write(requestKey)
	result.WriteByte(':')
	minuteBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(minuteBytes, uint16(currentMinute))
	result.Write(minuteBytes)

	return result.String()
}
