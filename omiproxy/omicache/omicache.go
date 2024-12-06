package cache

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/cockroachdb/pebble"
	register "github.com/stormi-li/omiv1/omiregister"
)

type Cache struct {
	State      bool
	ttl        time.Duration
	gcInterval time.Duration
	DB         *pebble.DB
	Dir        string
	clean      chan struct{}
}

type Item struct {
	Data    []byte
	CeateAt time.Time
}

func NewCache(dir string, omiRegister *register.Register) *Cache {
	db, err := pebble.Open(dir, nil)
	if err != nil {
		log.Fatalln(err)
	}
	omiCache := &Cache{
		State:      true,
		ttl:        time.Hour,
		gcInterval: time.Hour,
		Dir:        dir,
		DB:         db,
		clean:      make(chan struct{}),
	}
	go omiCache.StartCleanupTask()

	omiRegister.AddRegisterHandleFunc("CacheDir", func() string {
		return omiCache.Dir
	})
	omiRegister.AddRegisterHandleFunc("CacheState", func() string {
		if omiCache.State {
			return "open"
		}
		return "closed"
	})
	omiRegister.AddRegisterHandleFunc("CacheTTL", func() string {
		seconds := int(omiCache.GetTTL()) / int(time.Second)
		return strconv.Itoa(seconds) + "s"
	})
	omiRegister.AddRegisterHandleFunc("CacheSize", func() string {
		return fmt.Sprintf("%.3fMB", float64(omiCache.DiskSpaceUsage())/1024/1024)
	})

	omiRegister.AddMessageHandleFunc("SwitchCacheState", func(message string) {
		state, err := strconv.Atoi(message)
		if err == nil {
			if state == 1 {
				omiCache.State = true
			}
			if state == 0 {
				omiCache.State = false
			}
		}
	})
	omiRegister.AddMessageHandleFunc("UpdateCacheTTL", func(message string) {
		seconds, err := strconv.Atoi(message)
		if err == nil {
			omiCache.SetTTL(time.Duration(seconds) * time.Second)
		}
	})

	return omiCache
}

func (c *Cache) Set(key string, data []byte) {
	if !c.State {
		return
	}
	item := Item{Data: data, CeateAt: time.Now()}
	jsonBytes, err := json.Marshal(item)
	if err == nil {
		c.DB.Set([]byte(key), jsonBytes, pebble.Sync)
	}
}

func (c *Cache) Get(key string) []byte {
	if !c.State {
		return nil
	}
	value, closer, err := c.DB.Get([]byte(key))
	if err == nil {
		closer.Close()
	}
	var item Item
	json.Unmarshal(value, &item)
	return item.Data
}

func (c *Cache) DiskSpaceUsage() uint64 {
	return c.DB.Metrics().DiskSpaceUsage()
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
		var item Item
		err = json.Unmarshal(value, &item)
		if err != nil {
			continue
		}
		// 检查是否已过期
		if targetTime.After(item.CeateAt) {
			c.DB.Delete(key, pebble.Sync)
		}
	}
}
