package memory

import (
	"sync"
	"time"
)

type item[K any] struct {
	value K
	exp   int64
}

func (i *item[K]) Expired() bool {
	return time.Now().UnixNano() > i.exp && i.exp != -1
}

// Cache is an in-memory cache implementation.
type Cache[T comparable, K any] struct {
	items     map[T]item[K] // Cache data
	mu        sync.RWMutex  // Mutex for concurrent access
	onEvicted func(T, K)    // Callback function when an item is evicted
}

const (
	NoExpiration time.Duration = -1
)

// NewCache creates a new persistent memory cache.
func NewCache[T comparable, K any]() *Cache[T, K] {
	return &Cache[T, K]{
		items: make(map[T]item[K]),
	}
}

func (c *Cache[T, K]) OnEvicted(f func(T, K)) {
	c.mu.Lock()
	c.onEvicted = f
	c.mu.Unlock()
}

// Get retrieves a value from the cache associated with the given key.
func (c *Cache[T, K]) Get(key T) (*K, bool) {
	c.mu.RLock()

	item, found := c.items[key]
	if !found {
		c.mu.RUnlock()
		return nil, false
	}

	if item.Expired() {
		c.mu.RUnlock()
		c.Delete(key)
		return nil, false // We can't return the pointer because the item is deleted
	}

	c.mu.RUnlock()
	return &item.value, true
}

// Set sets a value in the cache with the given key.
func (c *Cache[T, K]) Set(key T, value K, exp time.Duration) {
	c.mu.Lock()

	if exp > 0 {
		c.items[key] = item[K]{
			value: value,
			exp:   time.Now().Add(exp).UnixNano(),
		}
	} else {
		c.items[key] = item[K]{
			value: value,
			exp:   int64(NoExpiration),
		}
	}

	c.mu.Unlock()
}

// Delete removes a value from the cache associated with the given key.
func (c *Cache[T, K]) Delete(key T) {
	c.mu.Lock()

	if c.onEvicted != nil {
		c.onEvicted(key, c.items[key].value)
	}

	delete(c.items, key)
	c.mu.Unlock()
}

// Clear removes all values from the cache.
func (c *Cache[T, K]) Clear() {
	c.mu.Lock()
	c.items = make(map[T]item[K])
	c.mu.Unlock()
}
