package cache

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrKeyNotFound = errors.New("key not found")
	ErrKeyExpired  = errors.New("key expired")
)

// Item represents a cached item with expiration
type Item struct {
	Value      string
	Expiration int64
}

// Cache is a simple in-memory cache with expiration
type Cache struct {
	items map[string]Item
	mu    sync.RWMutex
}

// NewCache creates a new cache
func NewCache() *Cache {
	cache := &Cache{
		items: make(map[string]Item),
	}

	// Start a goroutine to clean expired items
	go cache.janitor()

	return cache
}

// Set adds a key-value pair to the cache with a TTL
func (c *Cache) Set(_ context.Context, key string, value string, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	expiration := time.Now().Add(ttl).UnixNano()
	c.items[key] = Item{
		Value:      value,
		Expiration: expiration,
	}

	return nil
}

// Get retrieves a value by key
func (c *Cache) Get(_ context.Context, key string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return "", ErrKeyNotFound
	}

	// Check if the item has expired
	if item.Expiration < time.Now().UnixNano() {
		return "", ErrKeyExpired
	}

	return item.Value, nil
}

// Delete removes a key from the cache
func (c *Cache) Delete(_ context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
	return nil
}

// janitor cleans expired items from the cache
func (c *Cache) janitor() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.deleteExpired()
	}
}

// deleteExpired removes expired items from the cache
func (c *Cache) deleteExpired() {
	now := time.Now().UnixNano()

	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range c.items {
		if v.Expiration < now {
			delete(c.items, k)
		}
	}
}

// Global cache instance
var (
	globalCache *Cache
	once        sync.Once
)

// GetCache returns the global cache instance
func GetCache() *Cache {
	once.Do(func() {
		globalCache = NewCache()
	})
	return globalCache
}

// InitCache initializes the cache
func InitCache() error {
	GetCache() // Initialize the singleton
	return nil
}

// SetWithTTL sets a key-value pair with a TTL
func SetWithTTL(ctx context.Context, key string, value string, ttl time.Duration) error {
	return GetCache().Set(ctx, key, value, ttl)
}

// Get retrieves a value by key
func Get(ctx context.Context, key string) (string, error) {
	return GetCache().Get(ctx, key)
}

// Delete removes a key
func Delete(ctx context.Context, key string) error {
	return GetCache().Delete(ctx, key)
}
