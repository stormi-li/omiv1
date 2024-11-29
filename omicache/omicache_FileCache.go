package cache

import (
	"net/http"
	"net/url"
	"sync"

	proxy "github.com/stormi-li/omiv1/omiproxy"
)

type Cache struct {
	lock           sync.RWMutex // 改为读写锁
	MaxSize        int
	FileMgr        *FileManager
	LRUManager     *LRUManager
	CacheHitCount  int // 缓存命中次数
	CacheMissCount int // 缓存未命中次数
	CacheClearNum  int // 缓存清除次数
}

// 初始化文件缓存
func NewCache(cacheDir string, maxSize int) *Cache {
	cache := &Cache{
		MaxSize:    maxSize,
		FileMgr:    NewFileManager(cacheDir),
		LRUManager: NewLRUManager(),
	}
	return cache
}

// 设置缓存文件
func (fc *Cache) Set(key string, data []byte) {
	size := len(data)
	fc.lock.Lock()
	defer fc.lock.Unlock()

	// 确保容量充足
	for fc.LRUManager.Size+size > fc.MaxSize {
		oldest := fc.LRUManager.RemoveOldest()
		if oldest != nil {
			_ = fc.FileMgr.DeleteFile(oldest.Key)
			fc.CacheClearNum++ // 清除计数
		}
	}

	if size > fc.MaxSize {
		return // 文件太大，直接丢弃
	}

	// 写入文件并更新 LRU
	if err := fc.FileMgr.WriteFile(key, data); err == nil {
		fc.LRUManager.Add(key, size)
	}
}

// 获取缓存文件
func (fc *Cache) Get(key string) ([]byte, bool) {
	fc.lock.RLock() // 使用读锁
	_, found := fc.LRUManager.Get(key)
	if !found {
		fc.lock.RUnlock() // 在获取写锁前释放读锁
		fc.lock.Lock()    // 获取写锁清理逻辑
		defer fc.lock.Unlock()

		fc.CacheMissCount++ // 未命中计数
		return nil, false
	}

	data, err := fc.FileMgr.ReadFile(key)
	fc.lock.RUnlock() // 提前释放读锁
	if err != nil {
		fc.lock.Lock()         // 升级为写锁，清除不存在的缓存项
		defer fc.lock.Unlock() // 确保锁的正确释放
		fc.LRUManager.Remove(key)
		fc.CacheMissCount++ // 未命中计数
		return nil, false
	}

	fc.lock.Lock()     // 写锁保护计数器更新
	fc.CacheHitCount++ // 命中计数
	fc.lock.Unlock()
	return data, true
}

// 删除缓存文件
func (fc *Cache) Del(key string) {
	fc.lock.Lock()
	defer fc.lock.Unlock()

	fc.LRUManager.Remove(key)
	_ = fc.FileMgr.DeleteFile(key)
}

// 当前缓存使用大小
func (fc *Cache) GetCurrentSize() int {
	fc.lock.RLock()
	defer fc.lock.RUnlock()
	return fc.LRUManager.Size
}

// 获取命中次数
func (fc *Cache) GetCacheHitCount() int {
	fc.lock.RLock()
	defer fc.lock.RUnlock()
	return fc.CacheHitCount
}

// 获取未命中次数
func (fc *Cache) GetCacheMissCount() int {
	fc.lock.RLock()
	defer fc.lock.RUnlock()
	return fc.CacheMissCount
}

// 获取清除次数
func (fc *Cache) GetCacheClearCount() int {
	fc.lock.RLock()
	defer fc.lock.RUnlock()
	return fc.CacheClearNum
}

// 获取缓存条数
func (fc *Cache) GetCacheNum() int {
	fc.lock.RLock()
	defer fc.lock.RUnlock()
	return fc.LRUManager.Count()
}

func (fc *Cache) ServeCache(w http.ResponseWriter, r *http.Request) bool {
	r.URL.Host = r.Host
	data, ok := fc.Get(r.URL.String())
	if ok {
		w.Write(data)
		return true
	}
	return false
}

func (fc *Cache) UpdateCache(url *url.URL, cr *proxy.CapturedResponse) {
	if cr.StatusCode == http.StatusOK && len(cr.Body.Bytes()) != 0 {
		fc.Set(url.String(), cr.Body.Bytes())
	}
}
