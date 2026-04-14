package cache

import (
	"sync"
	"time"
)

// CacheItem represents a single cached item with its expiration time
type CacheItem struct {
	Data      interface{}
	ExpiresAt time.Time
}

// Cache is a thread-safe in-memory key-value store
type Cache struct {
	store map[string]CacheItem
	mu    sync.RWMutex
}

// NewCache initializes and returns a new Cache instance
func NewCache() *Cache {
	return &Cache{
		store: make(map[string]CacheItem),
	}
}

// Get retrieves an item from the cache.
// It returns the data and true if the key exists and hasn't expired.
// If the key doesn't exist or has expired, it returns nil and false.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	item, exists := c.store[key]
	c.mu.RUnlock()

	if !exists {
		return nil, false
	}

	// Check for expiration
	if time.Now().After(item.ExpiresAt) {
		// Lazy deletion of expired items
		c.Delete(key)
		return nil, false
	}

	return item.Data, true
}

// Set adds or updates an item in the cache with the given TTL (Time To Live).
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.store[key] = CacheItem{
		Data:      value,
		ExpiresAt: time.Now().Add(ttl),
	}
}

// Delete removes an item from the cache by its key.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key)
}
