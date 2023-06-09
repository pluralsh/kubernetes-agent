package redistool

import (
	"context"
	"time"

	"github.com/redis/rueidis"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/errz"
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
	ErrRep        errz.ErrReporter
	Client        rueidis.Client
	ErrMarshaler  ErrMarshaler
	KeyToRedisKey KeyToRedisKey[K]
}

func (c *ErrCacher[K]) GetError(ctx context.Context, key K) error {
	getCmd := c.Client.B().Get().Key(c.KeyToRedisKey(key)).Build()
	result, err := c.Client.Do(ctx, getCmd).ToString() // TODO use AsBytes() when have updated to Go 1.20
	if err != nil {
		if err != rueidis.Nil { // nolint:errorlint
			c.ErrRep.HandleProcessingError(ctx, c.Log, "Failed to get cached error from Redis", err)
		}
		return nil // Returns nil according to the interface contract.
	}
	if len(result) == 0 {
		return nil
	}
	e, err := c.ErrMarshaler.Unmarshal([]byte(result))
	if err != nil {
		c.ErrRep.HandleProcessingError(ctx, c.Log, "Failed to unmarshal cached error", err)
		return nil // Returns nil according to the interface contract.
	}
	return e
}

func (c *ErrCacher[K]) CacheError(ctx context.Context, key K, err error, errTtl time.Duration) {
	data, err := c.ErrMarshaler.Marshal(err)
	if err != nil {
		c.ErrRep.HandleProcessingError(ctx, c.Log, "Failed to marshal error for caching", err)
		return
	}
	setCmd := c.Client.B().Set().Key(c.KeyToRedisKey(key)).Value(rueidis.BinaryString(data)).Px(errTtl).Build()
	err = c.Client.Do(ctx, setCmd).Error()
	if err != nil {
		c.ErrRep.HandleProcessingError(ctx, c.Log, "Failed to cache error in Redis", err)
	}
}
