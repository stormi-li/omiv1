package omi

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	cache "github.com/stormi-li/omiv1/omiproxy/omicache"
	proxy "github.com/stormi-li/omiv1/omiproxy"
	register "github.com/stormi-li/omiv1/omiregister"
)

type Client struct {
	Options  *Options
	Proxy    *proxy.Proxy
	Register *register.Register
	cache    *cache.Cache
}

func NewClient(options *Options) *Client {
	var omiCache *cache.Cache
	if options.CacheDir != "" {
		omiCache = cache.NewCache(options.CacheDir)
	}

	redisOptions := &redis.Options{
		Addr:     options.Addr,
		Password: options.Password,
		DB:       options.DB,
	}
	redisClient := redis.NewClient(redisOptions)

	omiProxy := proxy.NewProxy(redisClient)
	omiProxy.Cache = omiCache
	omiRegister := register.NewRegister(redisClient)

	if omiCache != nil {
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
			return fmt.Sprintf("%.3fMB", float64(omiCache.DB.Metrics().DiskSpaceUsage())/1024/1024)
		})

		omiRegister.AddMessageHandleFunc("UpdateCacheState", func(message string) {
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
	}

	return &Client{
		Options:  options,
		Proxy:    omiProxy,
		Register: omiRegister,
		cache:    omiCache,
	}
}

