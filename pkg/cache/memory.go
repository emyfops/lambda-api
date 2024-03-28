package cache

import (
	"context"
	"github.com/Edouard127/lambda-rpc/pkg/background"
	"sync"
	"time"
)

// MemoryCache is an in-memory cache implementation.
type MemoryCache[T comparable, K any] struct {
	mu           sync.RWMutex    // Mutex for concurrent access
	persistent   bool            // Indicates if the cache is persistent
	data         map[T]K         // Cache data
	queue        map[T]time.Time // Tracks expiration time for cache entries
	timeoutEat   time.Duration   // Timeout duration for cache entries
	accumulation time.Duration   // Accumulation duration for cache entries
	ctx          context.Context // Context for the cache
}

const (
	eatInterval         = time.Millisecond * 1000
	defaultQueueSize    = 0
	defaultAccumulation = time.Millisecond * 1000
)

// NewPersistentMemoryCache creates a new persistent memory cache.
// The size parameter specifies the maximum number of entries that can be stored in the cache.
// To get a growable cache, set the size to 0.
func NewPersistentMemoryCache[T comparable, K any](size int64) *MemoryCache[T, K] {
	if size < 0 {
		size = defaultQueueSize
	}

	c := &MemoryCache[T, K]{
		persistent: true,
		data:       make(map[T]K, size),
		ctx:        context.Background(),
	}

	return c
}

// NewTempMemoryCache creates a new temporary memory cache.
// The timeout parameter specifies the duration after which a cache entry is considered expired.
// The size parameter specifies the maximum number of entries that can be stored in the cache.
// To get a growable cache, set the size to 0.
// The accumulation parameter specifies the duration added to the expiration time of a cache entry each time it is accessed.
func NewTempMemoryCache[T comparable, K any](timeout time.Duration, accumulation time.Duration, size int64) *MemoryCache[T, K] {
	if size < 0 {
		size = defaultQueueSize
	}

	if accumulation < 0 {
		accumulation = defaultAccumulation
	}

	c := &MemoryCache[T, K]{
		data:         make(map[T]K, size),
		queue:        make(map[T]time.Time, size),
		timeoutEat:   timeout,
		accumulation: accumulation,
		ctx:          context.Background(),
	}

	if !c.persistent {
		go background.Ticker(c.ctx, eatInterval, false, c.devour)
	}
	return c
}

// Get retrieves a value from the cache associated with the given key.
func (c *MemoryCache[T, K]) Get(key T) (K, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, ok := c.data[key]
	if ok && !c.persistent {
		c.queue[key] = c.queue[key].Add(c.accumulation)
	}
	return value, ok
}

// Set sets a value in the cache with the given key.
func (c *MemoryCache[T, K]) Set(key T, value K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
	if !c.persistent {
		c.queue[key] = time.Now().Add(c.timeoutEat)
	}
}

// Delete removes a value from the cache associated with the given key.
func (c *MemoryCache[T, K]) Delete(key T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
	delete(c.queue, key)
}

// Close stops the cache and clears its data.
func (c *MemoryCache[T, K]) Close() {
	c.ctx.Done()

	// Will be eaten by the GC when the goroutine exit
	c.data = nil
	c.queue = nil
}

// devour removes expired cache entries.
func (c *MemoryCache[T, K]) devour() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key := range c.queue {
		if c.queue[key].After(now) {
			c.Delete(key)
		}
	}
}
