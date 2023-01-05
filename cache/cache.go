package cache

import (
	"sync"
	"time"
)

type Cache struct {
	sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	Tasks             map[int]CTask
}

type ICache interface {
}

type CTask struct {
	Task       string
	Time       time.Time
	Expiration int64
}

func NewCache(defaultExpiration, cleanupInterval time.Duration) *Cache {
	tasks := make(map[int]CTask)

	cache := Cache{
		Tasks:             tasks,
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}

	if cleanupInterval > 0 {
		cache.StartGC()
	}
	return &cache
}

func (c *Cache) StartGC() {
	go c.GC()
}

func (c *Cache) GC() {
	for {
		<-time.After(c.cleanupInterval)

		if c.Tasks == nil {
			return
		}

		if keys := c.expiredKeys(); len(keys) != 0 {
			c.clearItems(keys)
		}
	}
}

func (c *Cache) expiredKeys() (keys []int) {

	c.RLock()

	defer c.RUnlock()

	for k, i := range c.Tasks {
		if time.Now().UnixNano() > i.Expiration && i.Expiration > 0 {
			keys = append(keys, k)
		}
	}

	return keys

}

func (c *Cache) clearItems(keys []int) {
	c.Lock()

	defer c.Unlock()

	for _, k := range keys {
		delete(c.Tasks, k)
	}
}

func (c *Cache) Set(key int, value string, created time.Time, duration time.Duration) {
	var expiration int64

	if duration == 0 {
		duration = c.defaultExpiration
	}

	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	c.Lock()

	defer c.Unlock()

	c.Tasks[key] = CTask{
		Task:       value,
		Time:       created,
		Expiration: expiration,
	}
}
