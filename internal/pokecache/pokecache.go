package pokecache

import (
	"fmt"
	"sync"
	"time"
)

const (
  DefaultCacheInterval = 5*time.Second
)

type Cache struct {
  cache     map[string]cacheEntry
  mu        *sync.Mutex
}

type cacheEntry struct {
  createdAt time.Time
  val       []byte
}

func NewCache(interval time.Duration) *Cache {
  c := &Cache{mu: &sync.Mutex{}, cache: make(map[string]cacheEntry)}
  go c.reapLoop(interval)

  return c
}

func (c *Cache) Add(key string, val []byte) {
  if key == "" {
    return
  }
  c.mu.Lock()
  defer c.mu.Unlock()

  c.cache[key] = cacheEntry{
    createdAt: time.Now(),
    val: val,
  }
}

func (c *Cache) Get(key string) (data []byte, found bool) {
  c.mu.Lock()
  defer c.mu.Unlock()

  entry, ok := c.cache[key]
  if !ok {
    return nil, false
  }
  
  fmt.Println("Retrieving cached data")
  return entry.val, true 
}

func (c *Cache) reapLoop(interval time.Duration) {
  ticker := time.NewTicker(interval)
  defer ticker.Stop() // This will run when reapLoop returns

  for range ticker.C {
    c.reap(interval)
  }
}

func (c *Cache) reap(interval time.Duration) {
  c.mu.Lock()
  defer c.mu.Unlock() // This will run when reap returns

  for k, v := range c.cache {
    if time.Since(v.createdAt) > interval {
        delete(c.cache, k)
    }
  }
}

