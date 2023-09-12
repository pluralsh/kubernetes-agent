package redistool

import (
	"context"
	"time"

	"github.com/redis/rueidis"
)

// ExpiringHashApi represents a low-level API to work with a two-level hash: key K1 -> hashKey K2 -> value []byte.
// key identifies the hash; hashKey identifies the key in the hash; value is the value for the hashKey.
type ExpiringHashApi[K1 any, K2 any] interface {
	SetBuilder() SetBuilder[K1, K2]
	Unset(ctx context.Context, key K1, hashKey K2) error
}

type RedisExpiringHashApi[K1 any, K2 any] struct {
	Client         rueidis.Client
	Key1ToRedisKey KeyToRedisKey[K1]
	Key2ToRedisKey KeyToRedisKey[K2]
}

func (h *RedisExpiringHashApi[K1, K2]) SetBuilder() SetBuilder[K1, K2] {
	return &RedisSetBuilder[K1, K2]{
		client:         h.Client,
		key1ToRedisKey: h.Key1ToRedisKey,
		key2ToRedisKey: h.Key2ToRedisKey,
	}
}

func (h *RedisExpiringHashApi[K1, K2]) Unset(ctx context.Context, key K1, hashKey K2) error {
	hdelCmd := h.Client.B().Hdel().Key(h.Key1ToRedisKey(key)).Field(h.Key2ToRedisKey(hashKey)).Build()
	return h.Client.Do(ctx, hdelCmd).Error()
}
