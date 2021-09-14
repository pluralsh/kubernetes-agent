package redistool

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/logz"
	"go.uber.org/zap"
)

type RpcApi interface {
	Log() *zap.Logger
	HandleProcessingError(msg string, err error)
	AgentToken() string
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
	key := l.buildKey(api.AgentToken())

	count, err := l.redisClient.Get(ctx, key).Uint64()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
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

func (l *TokenLimiter) buildKey(token string) string {
	// We use only the first half of the token as a key. Under the assumption of
	// a randomly generated token of length at least 50, with an alphabet of at least
	//
	// - upper-case characters (26)
	// - lower-case characters (26),
	// - numbers (10),
	//
	// (see https://gitlab.com/gitlab-org/gitlab/blob/master/app/models/clusters/agent_token.rb)
	//
	// we have at least 62^25 different possible token prefixes. Since the token is
	// randomly generated, to obtain the token from this hash, one would have to
	// also guess the second half, and validate it by attempting to log in (kas
	// cannot validate tokens on its own)
	n := len(token) / 2
	tokenHash := sha256.Sum256([]byte(token[:n]))

	currentMinute := time.Now().UTC().Minute()

	var result strings.Builder
	result.WriteString(l.keyPrefix)
	result.WriteByte(':')
	result.Write(tokenHash[:])
	result.WriteByte(':')
	minuteBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(minuteBytes, uint16(currentMinute))
	result.Write(minuteBytes)

	return result.String()
}
