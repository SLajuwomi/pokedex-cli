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
	CacheMap map[string]cacheEntry
	mux      sync.Mutex
	interval time.Duration
}

func NewCache(interval time.Duration) *Cache {
	newCache := Cache{}
	newCache.interval = interval
	newCache.CacheMap = map[string]cacheEntry{}
	go newCache.reapLoop()
	return &newCache
}

func (c *Cache) Add(key string, passedVal []byte) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.CacheMap[key] = cacheEntry{
		createdAt: time.Now(),
		val:       passedVal,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mux.Lock()
	defer c.mux.Unlock()
	gotten, ok := c.CacheMap[key]
	if ok {
		return gotten.val, true
	} else {
		return nil, false
	}
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for {
		<-ticker.C
		c.mux.Lock()
		for k, v := range c.CacheMap {
			if time.Since(v.createdAt) > c.interval {
				delete(c.CacheMap, k)
			}
		}
		c.mux.Unlock()
	}
}
