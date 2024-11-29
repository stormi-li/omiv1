package cache

import (
	"container/list"
	"sync"
)

type CacheItem struct {
	Key  string // 文件名
	Size int    // 文件大小
}

type LRUManager struct {
	lock     sync.Mutex
	itemMap  map[string]*list.Element
	itemList *list.List
	Size     int // 当前总大小
}

// 初始化 LRU 管理器
func NewLRUManager() *LRUManager {
	return &LRUManager{
		itemMap:  make(map[string]*list.Element),
		itemList: list.New(),
	}
}

// 添加缓存项
func (lru *LRUManager) Add(key string, size int) {
	lru.lock.Lock()
	defer lru.lock.Unlock()

	if elem, found := lru.itemMap[key]; found {
		lru.itemList.MoveToFront(elem)
		return
	}

	item := &CacheItem{Key: key, Size: size}
	elem := lru.itemList.PushFront(item)
	lru.itemMap[key] = elem
	lru.Size += size
}

// 获取并更新缓存项
func (lru *LRUManager) Get(key string) (*CacheItem, bool) {
	lru.lock.Lock()
	defer lru.lock.Unlock()

	elem, found := lru.itemMap[key]
	if !found {
		return nil, false
	}

	lru.itemList.MoveToFront(elem)
	return elem.Value.(*CacheItem), true
}

// 移除最旧的缓存项
func (lru *LRUManager) RemoveOldest() *CacheItem {
	lru.lock.Lock()
	defer lru.lock.Unlock()

	oldest := lru.itemList.Back()
	if oldest == nil {
		return nil
	}

	item := oldest.Value.(*CacheItem)
	lru.itemList.Remove(oldest)
	delete(lru.itemMap, item.Key)
	lru.Size -= item.Size
	return item
}

// 删除指定缓存项
func (lru *LRUManager) Remove(key string) {
	lru.lock.Lock()
	defer lru.lock.Unlock()

	elem, found := lru.itemMap[key]
	if !found {
		return
	}

	item := elem.Value.(*CacheItem)
	lru.itemList.Remove(elem)
	delete(lru.itemMap, item.Key)
	lru.Size -= item.Size
}

// 当前缓存项数目
func (lru *LRUManager) Count() int {
	lru.lock.Lock()
	defer lru.lock.Unlock()
	return lru.itemList.Len()
}