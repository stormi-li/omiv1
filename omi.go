package omi

import (
	"embed"

	"github.com/go-redis/redis/v8"
	monitor "github.com/stormi-li/omiv1/ominitor"
	proxy "github.com/stormi-li/omiv1/omiproxy"
	register "github.com/stormi-li/omiv1/omiregister"
	cert "github.com/stormi-li/omiv1/omiregister/omicert"
	web "github.com/stormi-li/omiv1/omiweb"
)

type Options struct {
	Addr     string
	Password string
	DB       int
	CacheDir string
}

func NewWeb(embeddedSource *embed.FS) *web.Web {
	return web.NewWeb(embeddedSource)
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

func WriteDefaultCertAndKey() {
	cert.WriteDefaultCertAndKey()
}
