package omi

import (
	"embed"

	"github.com/go-redis/redis/v8"
	monitor "github.com/stormi-li/omiv1/ominitor"
	register "github.com/stormi-li/omiv1/omiregister"
	cert "github.com/stormi-li/omiv1/omiregister/omicert"
	web "github.com/stormi-li/omiv1/omiweb"
)

func NewWeb(embeddedSource *embed.FS) *web.Web {
	return web.NewWeb(embeddedSource)
}

func NewMonitor(options *Options) *monitor.Monitor {
	return monitor.NewMonitor(register.NewRegister(redis.NewClient(&redis.Options{
		Addr:     options.Addr,
		Password: options.Password,
		DB:       options.DB,
	})))
}

func WriteDefaultCertAndKey() {
	cert.WriteDefaultCertAndKey()
}
