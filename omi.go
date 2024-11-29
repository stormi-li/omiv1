package omi

import (
	"embed"

	"github.com/go-redis/redis/v8"
	cache "github.com/stormi-li/omiv1/omicache"
	monitor "github.com/stormi-li/omiv1/ominitor"
	proxy "github.com/stormi-li/omiv1/omiproxy"
	register "github.com/stormi-li/omiv1/omiregister"
	web "github.com/stormi-li/omiv1/omiweb"
)

func NewWeb(sourcePath, indexPath string, embeddedSource *embed.FS) *web.Web {
	return web.NewWeb(sourcePath, indexPath, embeddedSource)
}

func NewCache(cacheDir string, size int) *cache.Cache {
	return cache.NewCache(cacheDir, size)
}

func NewProxy(redisClient *redis.Client) *proxy.Proxy {
	return proxy.NewProxy(redisClient)
}

func NewRegister(redisClient *redis.Client) *register.Register {
	return register.NewRegister(redisClient)
}

func NewMonitor(redisClient *redis.Client) *monitor.Monitor {
	return monitor.NewMonitor(NewRegister(redisClient))
}
