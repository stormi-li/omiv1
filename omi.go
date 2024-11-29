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

type Options struct {
	Addr     string
	Password string
	DB       int
}

func NewWeb(sourcePath, indexPath string, embeddedSource *embed.FS) *web.Web {
	return web.NewWeb(sourcePath, indexPath, embeddedSource)
}

func NewCache(cacheDir string, size int) *cache.Cache {
	return cache.NewCache(cacheDir, size)
}

func NewProxy(options *Options) *proxy.Proxy {
	return proxy.NewProxy(redis.NewClient(&redis.Options{
		Addr:     options.Addr,
		Password: options.Password,
		DB:       options.DB,
	}))
}

func NewRegister(options *Options) *register.Register {
	return register.NewRegister(redis.NewClient(&redis.Options{
		Addr:     options.Addr,
		Password: options.Password,
		DB:       options.DB,
	}))
}

func NewMonitor(options *Options) *monitor.Monitor {
	return monitor.NewMonitor(NewRegister(options))
}
