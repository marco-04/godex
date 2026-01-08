package internal

import (
	"time"
	"sync"
)

type cacheEntry struct {
	createdAt time.Time
	val []byte
}

type Cache struct {
	cache map[string]cacheEntry
	interval time.Duration
	mux *sync.Mutex
}

func NewCache(interval time.Duration) Cache {
	cache := Cache {
		cache: make(map[string]cacheEntry),
		interval: interval,
		mux: &sync.Mutex{},
	}

	ticker := time.NewTicker(interval)

	go func(ticker *time.Ticker, cache Cache) {
		for range ticker.C {
			cache.reapLoop()
		}
	}(ticker, cache)

	return cache
}

func (c Cache) Add(key string, val []byte) {
	c.mux.Lock()
	c.cache[key] = cacheEntry {
		createdAt: time.Now(),
		val: val,
	}
	c.mux.Unlock()
}

func (c Cache) Get(key string) ([]byte, bool) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if entry, ok := c.cache[key]; ok {
		return entry.val, true
	} else {
		return nil, false
	}
}

func (c Cache) reapLoop() {
	for k, v := range c.cache {
		if !v.createdAt.Add(c.interval).After(time.Now()) {
			c.mux.Lock()
			delete(c.cache, k)
			c.mux.Unlock()
		}
	}
}
