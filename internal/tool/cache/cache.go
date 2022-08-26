package cache

import (
	"sync"
	"time"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/syncz"
)

type Entry struct {
	// protects state in this object.
	syncz.Mutex
	// Expires holds the time when this entry should be removed from the cache.
	Expires time.Time
	// Item is the cached item.
	Item interface{}
}

func (e *Entry) IsNeedRefreshLocked() bool {
	return e.IsEmptyLocked() || e.IsExpiredLocked(time.Now())
}

func (e *Entry) IsEmptyLocked() bool {
	return e.Item == nil
}

func (e *Entry) IsExpiredLocked(t time.Time) bool {
	return e.Expires.Before(t)
}

type Cache struct {
	mu                    sync.Mutex
	data                  map[interface{}]*Entry
	expirationCheckPeriod time.Duration
	nextExpirationCheck   time.Time
}

func New(expirationCheckPeriod time.Duration) *Cache {
	return &Cache{
		data:                  make(map[interface{}]*Entry),
		expirationCheckPeriod: expirationCheckPeriod,
	}
}

func (c *Cache) EvictExpiredEntries() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	if now.Before(c.nextExpirationCheck) {
		return
	}
	c.nextExpirationCheck = now.Add(c.expirationCheckPeriod)
	for key, entry := range c.data {
		func() {
			if !entry.TryLock() {
				// entry is busy, skip
				return
			}
			defer entry.Unlock()
			if entry.IsExpiredLocked(now) {
				delete(c.data, key)
			}
		}()
	}
}

func (c *Cache) GetOrCreateCacheEntry(key interface{}) *Entry {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry := c.data[key]
	if entry != nil {
		return entry
	}
	entry = &Entry{
		Mutex: syncz.NewMutex(),
	}
	c.data[key] = entry
	return entry
}
