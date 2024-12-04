package omi

import (
	"embed"
	"net/http"

	"github.com/go-redis/redis/v8"
	cert "github.com/stormi-li/omiv1/omicert"
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

func HandleFunc() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {}
}
