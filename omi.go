package omi

import (
	"embed"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omiv1/omihttp"
	monitor "github.com/stormi-li/omiv1/ominitor"
	proxy "github.com/stormi-li/omiv1/omiproxy"
	register "github.com/stormi-li/omiv1/omiregister"
	web "github.com/stormi-li/omiv1/omiweb"
)

func NewWeb(sourcePath, indexPath string, embeddedSource *embed.FS) *web.Web {
	return web.NewWeb(sourcePath, indexPath, embeddedSource)
}

func NewReadWriter() *omihttp.ReadWriter {
	return omihttp.NewReadWriter()
}

type Client struct {
	RedisClient *redis.Client
}

func NewClient(options *redis.Options) *Client {
	return &Client{
		RedisClient: redis.NewClient(options),
	}
}

func (c *Client) NewProxy(transport *http.Transport) *proxy.Proxy {
	return proxy.NewProxy(c.RedisClient, transport)
}

func (c *Client) NewRegister(serverName, address string) *register.Register {
	return register.NewRegister(c.RedisClient, serverName, address)
}

func (c *Client) NewMonitor(serverName, address string) *monitor.Monitor {
	return monitor.NewMonitor(c.NewRegister(serverName, address))
}
