package cache

import (
	"context"
	"time"
)

type GetItemDirectly func() (interface{}, error)

type ErrCacher interface {
	// GetError retrieves a cached error.
	// Returns nil if no cached error found or if there was a problem accessing the cache.
	GetError(ctx context.Context, key interface{}) error
	// CacheError puts error into the cache.
	CacheError(ctx context.Context, key interface{}, err error, errTtl time.Duration)
}

type CacheWithErr struct {
	cache     *Cache
	ttl       time.Duration
	errTtl    time.Duration
	errCacher ErrCacher
	// isCacheable determines whether an error is cacheable or not.
	// Returns true if cacheable and false otherwise.
	isCacheable func(error) bool
}

func NewWithError(ttl, errTtl time.Duration, errCacher ErrCacher, isCacheableFunc func(error) bool) *CacheWithErr {
	return &CacheWithErr{
		cache:       New(ttl),
		ttl:         ttl,
		errTtl:      errTtl,
		errCacher:   errCacher,
		isCacheable: isCacheableFunc,
	}
}

func (c *CacheWithErr) GetItem(ctx context.Context, key interface{}, f GetItemDirectly) (interface{}, error) {
	if c.ttl == 0 {
		return f()
	}
	c.cache.EvictExpiredEntries()
	entry := c.cache.GetOrCreateCacheEntry(key)
	if !entry.Lock(ctx) { // a concurrent caller may be refreshing the entry. Block until exclusive access is available.
		return nil, ctx.Err()
	}
	evictEntry := false
	defer func() {
		entry.Unlock()
		if evictEntry {
			// Currently, cache (e.g. in EvictExpiredEntries()) grabs the cache lock and then an entry's lock,
			// but only via TryLock(). We may need to use Lock() rather than TryLock() in the future in some
			// other method. That would lead to deadlocks if we grab an entry's lock and then such method is called
			// concurrently. Hence,	to future-proof the code, calling EvictEntry() after entry's lock has been
			// unlocked here.
			c.cache.EvictEntry(key, entry)
		}
	}()
	if entry.IsNeedRefreshLocked() {
		err := c.errCacher.GetError(ctx, key)
		if err != nil {
			evictEntry = true
			return nil, err
		}
		entry.Item, err = f()
		if err != nil {
			if c.isCacheable(err) {
				// cacheable error
				c.errCacher.CacheError(ctx, key, err, c.errTtl)
			}
			return nil, err
		}
		entry.Expires = time.Now().Add(c.ttl)
	}
	return entry.Item, nil
}
