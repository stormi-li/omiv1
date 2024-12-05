package cache

import (
	"fmt"
	"log"
	"time"

	"github.com/cockroachdb/pebble"
)

type Cache struct {
	State bool
	ttl   time.Duration
	DB    *pebble.DB
	Dir   string
	clean chan struct{}
}

type Item struct {
	Value   string
	CeateAt time.Time
}

func NewCache(dir string) *Cache {
	db, err := pebble.Open(dir, nil)
	if err != nil {
		log.Fatalln(err)
	}
	cache := &Cache{
		State: true,
		ttl:   time.Hour,
		Dir:   dir,
		DB:    db,
	}
	go cache.StartCleanupTask()
	return cache
}

func (c *Cache) Set(key string, data []byte) {
	if !c.State {
		return
	}
	item := fmt.Sprintf("%s|%s", data, time.Now())
	c.DB.Set([]byte(key), []byte(item), pebble.Sync)
}

func (c *Cache) Get(key string) []byte {
	if !c.State {
		return nil
	}
	value, closer, _ := c.DB.Get([]byte(key))
	closer.Close()
	return value
}

func (c *Cache) GetTTL() time.Duration {
	return c.ttl
}

func (c *Cache) SetTTL(ttl time.Duration) {
	c.ttl = ttl
	c.clean <- struct{}{}
}

func (c *Cache) StartCleanupTask() {
	ticker := time.NewTicker(c.ttl)
	for {
		select {
		case <-c.clean:
			c.CleanupExpiredEntries()
			ticker.Stop()
			ticker = time.NewTicker(c.ttl)
		case <-ticker.C:
			c.CleanupExpiredEntries()
		}
	}
}

func (c *Cache) CleanupExpiredEntries() {
	iter, err := c.DB.NewIter(nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer iter.Close()

	targetTime := time.Now().Add(-c.ttl)
	for iter.First(); iter.Valid(); iter.Next() {
		key := iter.Key()
		value := iter.Value()

		// 解析存储的值，提取过期时间
		parts := string(value)
		var item Item
		_, err := fmt.Sscanf(parts, "%s|%s", &item.Value, &item.CeateAt)
		if err != nil {
			continue
		}
		// 检查是否已过期
		if targetTime.After(item.CeateAt) {
			// 删除过期键
			c.DB.Delete(key, pebble.Sync)
		}
	}
}
