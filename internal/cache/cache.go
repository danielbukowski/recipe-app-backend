// Package cache provides methods for interacting with Memcached.
package cache

import (
	"github.com/bradfitz/gomemcache/memcache"
)

// Cache interacts with a Memcached to retrieve or modify values.
type Cache struct {
	memcachedClient *memcache.Client
}

// New returns a new instance of Cache struct.
func New(memcachedClient *memcache.Client) *Cache {
	return &Cache{
		memcachedClient: memcachedClient,
	}
}

// DeleteItem removes an item by a key from Memcached.
func (c *Cache) DeleteItem(key string) error {
	return c.memcachedClient.Delete(key)
}

// GetItem fetches an item by a key from Memcached.
func (c *Cache) GetItem(key string) ([]byte, error) {
	item, err := c.memcachedClient.Get(key)
	if err != nil {
		return nil, err
	}
	return item.Value, nil
}

// InsertItem insert an item to Memcached or overwrites the already existing item.
func (c *Cache) InsertItem(key string, value []byte, expiration int32) error {
	item := memcache.Item{}
	item.Key = key
	item.Value = value
	item.Expiration = expiration

	return c.memcachedClient.Set(&item)
}
