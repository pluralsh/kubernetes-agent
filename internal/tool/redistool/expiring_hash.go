package redistool

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"time"
	"unsafe"

	"github.com/redis/rueidis"
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
	client         rueidis.Client
	key1ToRedisKey KeyToRedisKey[K1]
	key2ToRedisKey KeyToRedisKey[K2]
	ttl            time.Duration
	data           map[K1]map[K2]*ExpiringValue // key -> hash key -> value
}

func NewExpiringHash[K1 comparable, K2 comparable](client rueidis.Client, key1ToRedisKey KeyToRedisKey[K1],
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
		return h.refreshKey(ctx, key, []refreshKey[K2]{
			{
				hashKey: hashKey,
				value: ExpiringValue{ // cannot copy ev directly
					ExpiresAt: ev.ExpiresAt,
					Value:     ev.Value,
				},
			},
		})
	}
}

func (h *ExpiringHash[K1, K2]) Unset(key K1, hashKey K2) IOFunc {
	h.unsetData(key, hashKey)
	return func(ctx context.Context) error {
		hdelCmd := h.client.B().Hdel().Key(h.key1ToRedisKey(key)).Field(h.key2ToRedisKey(hashKey)).Build()
		return h.client.Do(ctx, hdelCmd).Error()
	}
}

func (h *ExpiringHash[K1, K2]) Forget(key K1, hashKey K2) {
	h.unsetData(key, hashKey)
}

func (h *ExpiringHash[K1, K2]) Len(ctx context.Context, key K1) (size int64, retErr error) {
	hlenCmd := h.client.B().Hlen().Key(h.key1ToRedisKey(key)).Build()
	return h.client.Do(ctx, hlenCmd).AsInt64()
}

func (h *ExpiringHash[K1, K2]) scan(ctx context.Context, key K1, cb func(k, v string) (bool /*done*/, bool /*delete*/, error)) (keysDeleted int, retErr error) {
	redisKey := h.key1ToRedisKey(key)
	var keysToDelete []string
	defer func() {
		if len(keysToDelete) == 0 {
			return
		}
		hdelCmd := h.client.B().Hdel().Key(redisKey).Field(keysToDelete...).Build()
		err := h.client.Do(ctx, hdelCmd).Error()
		if err != nil {
			if retErr == nil {
				retErr = err
			}
			return
		}
		keysDeleted = len(keysToDelete)
	}()
	// Scan keys of a hash. See https://redis.io/commands/scan
	var se rueidis.ScanEntry
	var err error
	for more := true; more; more = se.Cursor != 0 {
		hscanCmd := h.client.B().Hscan().Key(redisKey).Cursor(se.Cursor).Build()
		se, err = h.client.Do(ctx, hscanCmd).AsScanEntry()
		if err != nil {
			return 0, err
		}
		if len(se.Elements)%2 != 0 {
			// This shouldn't happen
			return 0, errors.New("invalid Redis reply")
		}
		for i := 0; i < len(se.Elements); i += 2 {
			k := se.Elements[i]
			v := se.Elements[i+1]
			done, del, err := cb(k, v)
			if del {
				keysToDelete = append(keysToDelete, k)
			}
			if err != nil || done {
				return 0, err
			}
		}
	}
	return 0, nil
}

func (h *ExpiringHash[K1, K2]) Scan(ctx context.Context, key K1, cb ScanCallback) (keysDeleted int, retErr error) {
	now := time.Now().Unix()
	var msg ExpiringValue
	return h.scan(ctx, key, func(k, v string) (bool /*done*/, bool /*delete*/, error) {
		// Avoid creating a temporary copy
		vBytes := unsafe.Slice(unsafe.StringData(v), len(v))
		err := proto.Unmarshal(vBytes, &msg)
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
		// Avoid creating a temporary copy
		vBytes := unsafe.Slice(unsafe.StringData(v), len(v))
		err := proto.UnmarshalOptions{
			DiscardUnknown: true, // We know there is one more field, but we don't need it
		}.Unmarshal(vBytes, &msg)
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
	cmds := make([]rueidis.Completed, 0, len(h.data))
	for k1, m := range h.data {
		toDel = toDel[:0] // reuse backing array, but reset length
		for k2 := range m {
			toDel = append(toDel, h.key2ToRedisKey(k2))
		}
		cmds = append(cmds, h.client.B().Hdel().Key(h.key1ToRedisKey(k1)).Field(toDel...).Build())
		delete(h.data, k1)
		keysDeleted += len(toDel)
	}
	err := MultiFirstError(h.client.DoMulti(ctx, cmds...))
	return keysDeleted, err
}

func (h *ExpiringHash[K1, K2]) Refresh(nextRefresh time.Time) IOFunc {
	argsMap := make(map[K1][]refreshKey[K2], len(h.data))
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

func (h *ExpiringHash[K1, K2]) prepareRefreshKey(hashData map[K2]*ExpiringValue, nextRefresh time.Time) []refreshKey[K2] {
	args := make([]refreshKey[K2], 0, len(hashData))
	expiresAt := time.Now().Add(h.ttl).Unix()
	nextRefreshUnix := nextRefresh.Unix()
	for hashKey, value := range hashData {
		if value.ExpiresAt > nextRefreshUnix {
			// Expires after next refresh. Will be refreshed later, no need to refresh now.
			continue
		}
		value.ExpiresAt = expiresAt
		// Copy value to decouple from the mutable instance in hashData. That way it's safe for concurrent access.
		args = append(args, refreshKey[K2]{
			hashKey: hashKey,
			value:   ExpiringValue{ExpiresAt: value.ExpiresAt, Value: value.Value},
		})
	}
	return args
}

func (h *ExpiringHash[K1, K2]) refreshKey(ctx context.Context, key K1, args []refreshKey[K2]) error {
	var marshalErr error
	redisKey := h.key1ToRedisKey(key)
	hsetCmd := h.client.B().Hset().Key(redisKey).FieldValue()
	empty := true
	// Iterate indexes to avoid copying the value which has inlined proto message, which shouldn't be copied.
	for i := range args {
		redisValue, err := proto.Marshal(&args[i].value)
		if err != nil {
			// This should never happen
			if marshalErr == nil {
				marshalErr = fmt.Errorf("failed to marshal ExpiringValue: %w", err)
			}
			continue // skip this value
		}
		hsetCmd.FieldValue(h.key2ToRedisKey(args[i].hashKey), rueidis.BinaryString(redisValue))
		empty = false
	}
	if empty {
		return nil // nothing to do, all skipped.
	}
	resp := h.client.DoMulti(ctx,
		h.client.B().Multi().Build(),
		hsetCmd.Build(),
		h.client.B().Pexpire().Key(redisKey).Milliseconds(h.ttl.Milliseconds()).Build(),
		h.client.B().Exec().Build(),
	)
	err := MultiFirstError(resp)
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

type refreshKey[K2 any] struct {
	hashKey K2
	value   ExpiringValue
}

func PrefixedInt64Key(prefix string, key int64) string {
	b := make([]byte, 0, len(prefix)+8)
	b = append(b, prefix...)
	b = binary.LittleEndian.AppendUint64(b, uint64(key))

	return unsafe.String(unsafe.SliceData(b), len(b))
}
