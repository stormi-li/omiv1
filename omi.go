package omi

import (
	"embed"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omiv1/omihttp"
	monitor "github.com/stormi-li/omiv1/ominitor"
	proxy "github.com/stormi-li/omiv1/omiproxy"
	cache "github.com/stormi-li/omiv1/omiproxy/omicache"
	register "github.com/stormi-li/omiv1/omiregister"
	cert "github.com/stormi-li/omiv1/omiregister/omicert"
	web "github.com/stormi-li/omiv1/omiweb"
)

type Options struct {
	// redis服务器地址，ip:port格式，比如：192.168.1.100:6379
	// 默认为 :6379
	Addr string
	// 默认为空，不进行认证。
	Password string
	// redis DB 数据库，默认为0
	DB int
	// 本地缓存路径，默认为不启用缓存
	CacheDir string
	// redis配置项，优先使用
	RedisOptions *redis.Options
}

type Client struct {
	Options  *Options
	Proxy    *proxy.Proxy
	Register *register.Register
	Cache    *cache.Cache
}

func NewClient(options *Options) *Client {
	var redisOptions *redis.Options
	if options.RedisOptions == nil {
		redisOptions = &redis.Options{
			Addr:     options.Addr,
			Password: options.Password,
			DB:       options.DB,
		}
	} else {
		redisOptions = options.RedisOptions
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

func (c *Client) SetTransport(transport *http.Transport) {
	c.Proxy.Transport = transport
}

func (c *Client) AddFilter(serverName, key string, handler func(value string) bool) {
	c.Proxy.AddFilter(serverName, key, handler)
}

func (c *Client) ServePathProxy(w http.ResponseWriter, r *http.Request) {
	c.Proxy.ServePathProxy(w, r)
}

func (c *Client) ServeDomainProxy(w http.ResponseWriter, r *http.Request) {
	c.Proxy.ServeDomainProxy(w, r)
}

func (c *Client) Post(serverName string, pattern string, v any, sType omihttp.SerializationType) (*omihttp.Response, error) {
	return c.Proxy.Post(serverName, pattern, v, sType)
}

func (c *Client) NewServeMux() *omihttp.ServeMux {
	return omihttp.NewServeMux()
}

func (c *Client) AddRegisterHandleFunc(key string, handler func() string) {
	c.Register.AddRegisterHandleFunc(key, handler)
}

func (c *Client) AddMessageHandleFunc(command string, handler func(message string)) {
	c.Register.AddMessageHandleFunc(command, handler)
}

func (c *Client) RegisterAndServe(serverName, address string, handler http.Handler) {
	c.Register.RegisterAndServe(serverName, address, handler)
}

func (c *Client) RegisterAndServeTLS(serverName, address, certFile, keyFile string, handler http.Handler) {
	c.Register.RegisterAndServeTLS(serverName, address, certFile, keyFile, handler)
}

func (c *Client) NewWebServer(embeddedSource *embed.FS) *web.Web {
	return web.NewWebServer(embeddedSource)
}

func (c *Client) NewMonitorMux() *omihttp.ServeMux {
	return monitor.NewMonitorMux(c.NewServeMux(), c.Register.RedisClient)
}

func (c *Client) GenerateTestCertAndKey() {
	cert.WriteDefaultCertAndKey()
}
