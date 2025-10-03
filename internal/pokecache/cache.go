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
	cacheMap map[string]cacheEntry
	mux      sync.Mutex
	interval time.Duration
}

func NewCache(interval time.Duration) Cache {
	newCache := Cache{}
	newCache.interval = time.Duration(interval.Seconds())
	newCache.cacheMap = map[string]cacheEntry{}
	go newCache.reapLoop()
	return newCache
}

func (c *Cache) Add(key string, passedVal []byte) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.cacheMap[key] = cacheEntry{
		createdAt: time.Now(),
		val:       passedVal,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mux.Lock()
	defer c.mux.Unlock()
	gotten, ok := c.cacheMap[key]
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
		for k, v := range c.cacheMap {
			if time.Since(v.createdAt) > c.interval {
				delete(c.cacheMap, k)
			}
		}
		c.mux.Unlock()
	}
}
