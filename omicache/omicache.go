package cache

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/dgraph-io/badger/v3"
)

var TTL = time.Hour
var GCInterval = time.Hour

type BadgerCache struct {
	State      bool
	Size       int64
	TTL        time.Duration
	GCInterval time.Duration
	DB         *badger.DB
	Dir        string
}

func NewCache(dir string) *BadgerCache {
	db, err := badger.Open(badger.DefaultOptions(dir))
	if err != nil {
		log.Fatalln(err)
	}
	cache := &BadgerCache{
		State:      false,
		Size:       0,
		TTL:        TTL,
		GCInterval: GCInterval,
		Dir:        dir,
		DB:         db,
	}
	cache.Refresh()
	return cache
}

func (c *BadgerCache) Refresh() {
	go func() {
		for {
			if c.State {
				c.Size = getDirSize(c.Dir)
			}
			time.Sleep(10 * time.Second)
		}
	}()
	go func() {
		for {
			if c.State {
				c.DB.RunValueLogGC(0.5)
			}
			time.Sleep(GCInterval)
		}
	}()
}

func (c *BadgerCache) Set(key string, data []byte) {
	if !c.State {
		return
	}
	c.DB.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(key), data).WithTTL(TTL)
		return txn.SetEntry(e)
	})
}

func (c *BadgerCache) Get(key string) []byte {
	if !c.State {
		return nil
	}
	data := []byte{}
	c.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("name"))
		if err != nil {
			return err
		}
		// 读取值
		return item.Value(func(val []byte) error {
			data = val
			return nil
		})
	})
	return data
}

func getDirSize(path string) int64 {
	var size int64
	filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size
}
