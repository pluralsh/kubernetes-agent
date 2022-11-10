package redistool

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/logz"
	"go.uber.org/zap"
)

type ErrMarshaler interface {
	// Marshal turns error into []byte.
	Marshal(error) ([]byte, error)
	// Unmarshal turns []byte into error.
	Unmarshal([]byte) (error, error)
}

type ErrCacher[K any] struct {
	Log           *zap.Logger
	Client        redis.UniversalClient
	ErrMarshaler  ErrMarshaler
	KeyToRedisKey KeyToRedisKey[K]
}

func (c *ErrCacher[K]) GetError(ctx context.Context, key K) error {
	result, err := c.Client.Get(ctx, c.KeyToRedisKey(key)).Bytes()
	if err != nil {
		if err != redis.Nil { // nolint:errorlint
			c.Log.Error("Failed to get cached error from Redis", logz.Error(err))
		}
		return nil // Returns nil according to the interface contract.
	}
	if len(result) == 0 {
		return nil
	}
	e, err := c.ErrMarshaler.Unmarshal(result)
	if err != nil {
		c.Log.Error("Failed to unmarshal cached error", logz.Error(err))
		return nil // Returns nil according to the interface contract.
	}
	return e
}

func (c *ErrCacher[K]) CacheError(ctx context.Context, key K, err error, errTtl time.Duration) {
	data, err := c.ErrMarshaler.Marshal(err)
	if err != nil {
		c.Log.Error("Failed to marshal error for caching", logz.Error(err))
		return
	}
	err = c.Client.Set(ctx, c.KeyToRedisKey(key), data, errTtl).Err()
	if err != nil {
		c.Log.Error("Failed to cache error in Redis", logz.Error(err))
	}
}
