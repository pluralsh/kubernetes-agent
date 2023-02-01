package redistool

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"
)

// KeyToRedisKey is used to convert typed key (key1 or key2) into a string.
// HSET key1 key2 value.
type KeyToRedisKey[K any] func(key K) string
type ScanCallback func(rawHashKey string, value []byte, err error) (bool /* done */, error)

// IOFunc is a function that should be called to perform the I/O of the requested operation.
// It is safe to call concurrently as it does not interfere with the hash's operation.
type IOFunc func(ctx context.Context) error

// ExpiringHashInterface represents a two-level hash: key K1 -> hashKey K2 -> value []byte.
// key identifies the hash; hashKey identifies the key in the hash; value is the value for the hashKey.
// It is not safe for concurrent use directly, but it allows to perform I/O with backing store concurrently by
// returning functions for doing that.
type ExpiringHashInterface[K1 any, K2 any] interface {
	Set(key K1, hashKey K2, value []byte) IOFunc
	Unset(key K1, hashKey K2) IOFunc
	// Forget only removes the item from the in-memory map.
	Forget(key K1, hashKey K2)
	Scan(ctx context.Context, key K1, cb ScanCallback) (int /* keysDeleted */, error)
	Len(ctx context.Context, key K1) (int64, error)
	// GC returns a function that iterates all relevant stored data and deletes expired entries.
	// The returned function can be called concurrently as it does not interfere with the hash's operation.
	// The function returns number of deleted Redis (hash) keys, including when an error occurs.
	// It only inspects/GCs hashes where it has entries. Other concurrent clients GC same and/or other corresponding hashes.
	// Hashes that don't have a corresponding client (e.g. because it crashed) will expire because of TTL on the hash key.
	GC() func(context.Context) (int /* keysDeleted */, error)
	// Clear clears all data in this hash and deletes it from the backing store.
	Clear(context.Context) (int, error)
	Refresh(nextRefresh time.Time) IOFunc
}

type ExpiringHash[K1 comparable, K2 comparable] struct {
	client         redis.UniversalClient
	key1ToRedisKey KeyToRedisKey[K1]
	key2ToRedisKey KeyToRedisKey[K2]
	ttl            time.Duration
	data           map[K1]map[K2]*ExpiringValue // key -> hash key -> value
}

func NewExpiringHash[K1 comparable, K2 comparable](client redis.UniversalClient, key1ToRedisKey KeyToRedisKey[K1],
	key2ToRedisKey KeyToRedisKey[K2], ttl time.Duration) *ExpiringHash[K1, K2] {
	return &ExpiringHash[K1, K2]{
		client:         client,
		key1ToRedisKey: key1ToRedisKey,
		key2ToRedisKey: key2ToRedisKey,
		ttl:            ttl,
		data:           make(map[K1]map[K2]*ExpiringValue),
	}
}

func (h *ExpiringHash[K1, K2]) Set(key K1, hashKey K2, value []byte) IOFunc {
	ev := &ExpiringValue{
		ExpiresAt: time.Now().Add(h.ttl).Unix(),
		Value:     value,
	}
	h.setData(key, hashKey, ev)
	return func(ctx context.Context) error {
		return h.refreshKey(ctx, key, []interface{}{hashKey, ev})
	}
}

func (h *ExpiringHash[K1, K2]) Unset(key K1, hashKey K2) IOFunc {
	h.unsetData(key, hashKey)
	return func(ctx context.Context) error {
		return h.client.HDel(ctx, h.key1ToRedisKey(key), h.key2ToRedisKey(hashKey)).Err()
	}
}

func (h *ExpiringHash[K1, K2]) Forget(key K1, hashKey K2) {
	h.unsetData(key, hashKey)
}

func (h *ExpiringHash[K1, K2]) Len(ctx context.Context, key K1) (size int64, retErr error) {
	redisKey := h.key1ToRedisKey(key)
	return h.client.HLen(ctx, redisKey).Result()
}

func (h *ExpiringHash[K1, K2]) scan(ctx context.Context, key K1, cb func(k, v string) (bool /*done*/, bool /*delete*/, error)) (keysDeleted int, retErr error) {
	redisKey := h.key1ToRedisKey(key)
	var keysToDelete []string
	defer func() {
		if len(keysToDelete) == 0 {
			return
		}
		err := h.client.HDel(ctx, redisKey, keysToDelete...).Err()
		if err != nil {
			if retErr == nil {
				retErr = err
			}
			return
		}
		keysDeleted = len(keysToDelete)
	}()
	// Scan keys of a hash. See https://redis.io/commands/scan
	iter := h.client.HScan(ctx, redisKey, 0, "", 0).Iterator()
	for iter.Next(ctx) {
		k := iter.Val()
		if !iter.Next(ctx) {
			err := iter.Err()
			if err != nil {
				return 0, err
			}
			// This shouldn't happen
			return 0, errors.New("invalid Redis reply")
		}
		v := iter.Val()
		done, del, err := cb(k, v)
		if del {
			keysToDelete = append(keysToDelete, k)
		}
		if err != nil || done {
			return 0, err
		}
	}
	return 0, iter.Err()
}

func (h *ExpiringHash[K1, K2]) Scan(ctx context.Context, key K1, cb ScanCallback) (keysDeleted int, retErr error) {
	now := time.Now().Unix()
	var msg ExpiringValue
	return h.scan(ctx, key, func(k, v string) (bool /*done*/, bool /*delete*/, error) {
		err := proto.Unmarshal([]byte(v), &msg)
		if err != nil {
			done, cbErr := cb(k, nil, fmt.Errorf("failed to unmarshal hash value from hashkey 0x%x: %w", k, err))
			return done, false, cbErr
		}
		if msg.ExpiresAt < now {
			return false, true, nil
		}
		done, cbErr := cb(k, msg.Value, nil)
		return done, false, cbErr
	})
}

func (h *ExpiringHash[K1, K2]) GC() func(context.Context) (int, error) {
	// Copy keys for safe concurrent access.
	keys := make([]K1, 0, len(h.data))
	for key := range h.data {
		keys = append(keys, key)
	}
	return func(ctx context.Context) (int, error) {
		var deletedKeys int
		for _, key := range keys {
			deleted, err := h.gcHash(ctx, key)
			if err != nil {
				return deletedKeys, err
			}
			deletedKeys += deleted
		}
		return deletedKeys, nil
	}
}

// gcHash iterates a hash and removes all expired values.
// It assumes that values are marshaled ExpiringValue.
func (h *ExpiringHash[K1, K2]) gcHash(ctx context.Context, key K1) (int, error) {
	now := time.Now().Unix()
	var msg ExpiringValueTimestamp
	var firstErr error
	deleted, err := h.scan(ctx, key, func(k, v string) (bool /*done*/, bool /*delete*/, error) {
		err := proto.UnmarshalOptions{
			DiscardUnknown: true, // We know there is one more field, but we don't need it
		}.Unmarshal([]byte(v), &msg)
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			return false, false, nil
		}
		return false, msg.ExpiresAt < now, nil
	})
	if err != nil {
		return deleted, err
	}
	return deleted, firstErr
}

func (h *ExpiringHash[K1, K2]) Clear(ctx context.Context) (int, error) {
	var toDel []string
	keysDeleted := 0
	_, err := h.client.Pipelined(ctx, func(p redis.Pipeliner) error {
		// consider sending commands to Redis in batches to avoid accumulating too much in RAM.
		for k1, m := range h.data {
			toDel = toDel[:0] // reuse backing array, but reset length
			for k2 := range m {
				toDel = append(toDel, h.key2ToRedisKey(k2))
			}
			redisKey := h.key1ToRedisKey(k1)
			err := p.HDel(ctx, redisKey, toDel...).Err()
			if err != nil {
				return err
			}
			delete(h.data, k1)
			keysDeleted += len(toDel)
		}
		return nil
	})
	return keysDeleted, err
}

func (h *ExpiringHash[K1, K2]) Refresh(nextRefresh time.Time) IOFunc {
	argsMap := make(map[K1][]interface{}, len(h.data))
	for key, hashData := range h.data {
		args := h.prepareRefreshKey(hashData, nextRefresh)
		if len(args) == 0 {
			// Nothing to do for this key.
			continue
		}
		argsMap[key] = args
	}
	return func(ctx context.Context) error {
		var wg errgroup.Group
		for key, args := range argsMap {
			key := key
			args := args
			wg.Go(func() error {
				return h.refreshKey(ctx, key, args)
			})
		}
		return wg.Wait()
	}
}

func (h *ExpiringHash[K1, K2]) prepareRefreshKey(hashData map[K2]*ExpiringValue, nextRefresh time.Time) []interface{} {
	args := make([]interface{}, 0, len(hashData)*2)
	expiresAt := time.Now().Add(h.ttl).Unix()
	nextRefreshUnix := nextRefresh.Unix()
	for hashKey, value := range hashData {
		if value.ExpiresAt > nextRefreshUnix {
			// Expires after next refresh. Will be refreshed later, no need to refresh now.
			continue
		}
		value.ExpiresAt = expiresAt
		// Copy to decouple from the mutable instance in hashData. That way it's safe for concurrent access.
		valueCopy := &ExpiringValue{ExpiresAt: value.ExpiresAt, Value: value.Value}
		args = append(args, hashKey, valueCopy)
	}
	return args
}

func (h *ExpiringHash[K1, K2]) refreshKey(ctx context.Context, key K1, args []interface{}) error {
	var marshalErr error
	hsetArgs := make([]interface{}, 0, len(args))
	for i := 0; i < len(args); i += 2 {
		redisValue, err := proto.Marshal(args[i+1].(*ExpiringValue))
		if err != nil {
			// This should never happen
			if marshalErr == nil {
				marshalErr = fmt.Errorf("failed to marshal ExpiringValue: %w", err)
			}
			continue // skip this value
		}
		hsetArgs = append(hsetArgs, args[i], redisValue)
	}
	if len(hsetArgs) == 0 {
		return nil // nothing to do, all skipped.
	}
	redisKey := h.key1ToRedisKey(key)
	_, err := h.client.TxPipelined(ctx, func(p redis.Pipeliner) error {
		p.HSet(ctx, redisKey, hsetArgs) // nolint: asasalint
		p.PExpire(ctx, redisKey, h.ttl)
		return nil
	})
	if err != nil {
		return err
	}
	return marshalErr
}

func (h *ExpiringHash[K1, K2]) setData(key K1, hashKey K2, value *ExpiringValue) {
	nm := h.data[key]
	if nm == nil {
		nm = make(map[K2]*ExpiringValue, 1)
		h.data[key] = nm
	}
	nm[hashKey] = value
}

func (h *ExpiringHash[K1, K2]) unsetData(key K1, hashKey K2) {
	nm := h.data[key]
	delete(nm, hashKey)
	if len(nm) == 0 {
		delete(h.data, key)
	}
}

func PrefixedInt64Key(prefix string, key int64) string {
	var b strings.Builder
	b.WriteString(prefix)
	id := make([]byte, 8)
	binary.LittleEndian.PutUint64(id, uint64(key))
	b.Write(id)
	return b.String()
}
