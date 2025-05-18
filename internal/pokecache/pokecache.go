package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	mu              *sync.Mutex
	MapCacheEntries map[string]cacheEntry
}

func NewCache(interval time.Duration) Cache {
	c := Cache{
		mu:              &sync.Mutex{},
		MapCacheEntries: make(map[string]cacheEntry),
	}
	go c.reapLoop(interval)
	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.MapCacheEntries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	cacheE, ok := c.MapCacheEntries[key]
	if !ok {
		return nil, false
	}
	return cacheE.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		t := <-ticker.C
		c.mu.Lock()
		for key, cacheE := range c.MapCacheEntries {
			if t.Sub(cacheE.createdAt) >= interval {
				delete(c.MapCacheEntries, key)
			}
		}
		c.mu.Unlock()
	}
}
