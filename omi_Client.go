package omi

import (
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	proxy "github.com/stormi-li/omiv1/omiproxy"
	cache "github.com/stormi-li/omiv1/omiproxy/omicache"
	register "github.com/stormi-li/omiv1/omiregister"
	"github.com/stormi-li/omiv1/omirpc"
)

type Options struct {
	Addr     string
	Password string
	DB       int
	CacheDir string
}

type Client struct {
	Options  *Options
	Proxy    *proxy.Proxy
	Register *register.Register
	Cache    *cache.Cache
}

func NewClient(options *Options) *Client {
	redisOptions := &redis.Options{
		Addr:     options.Addr,
		Password: options.Password,
		DB:       options.DB,
	}
	redisClient := redis.NewClient(redisOptions)

	omiRegister := register.NewRegister(redisClient)

	var omiCache *cache.Cache
	if options.CacheDir != "" {
		omiCache = cache.NewCache(options.CacheDir, omiRegister)
	}

	omiProxy := proxy.NewProxy(redisClient)
	omiProxy.Cache = omiCache

	return &Client{
		Options:  options,
		Proxy:    omiProxy,
		Register: omiRegister,
		Cache:    omiCache,
	}
}

func (c *Client) ServePathProxy(w http.ResponseWriter, r *http.Request) {
	c.Proxy.ServePathProxy(w, r)
}

func (c *Client) ServeDomainProxy(w http.ResponseWriter, r *http.Request) {
	c.Proxy.ServeDomainProxy(w, r)
}

func (c *Client) Call(serverName string, pattern string, v any) (*omirpc.Response, error) {
	return c.Proxy.Call(serverName, pattern, v)
}

func (c *Client) NewServeMux() *register.ServeMux {
	return register.NewServeMux(c.Register)
}

func (c *Client) RegisterAndServe(serverName, address string, handler http.Handler) {
	c.Register.RegisterAndServe(serverName, address, func(port string) {
		err := http.ListenAndServe(port, handler)
		log.Fatalln(err)
	})
}

func (c *Client) RegisterAndServeTLS(serverName, address, certFile, keyFile string, handler http.Handler) {
	c.Register.RegisterAndServeTLS(serverName, address, func(port string) {
		err := http.ListenAndServeTLS(port, certFile, keyFile, handler)
		log.Fatalln(err)
	})
}
