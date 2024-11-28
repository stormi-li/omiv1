package omi

import (
	"embed"
	"net/http"

	"github.com/go-redis/redis/v8"
	monitor "github.com/stormi-li/omiv1/ominitor"
	proxy "github.com/stormi-li/omiv1/omiproxy"
	server "github.com/stormi-li/omiv1/omiserver"
	web "github.com/stormi-li/omiv1/omiweb"
)

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

func (c *Client) NewServer(serverName, address string) *server.Server {
	return server.NewServer(c.RedisClient, serverName, address)
}

func (c *Client) NewMonitor(serverName, address string) *monitor.Monitor {
	return monitor.NewMonitor(c.RedisClient, serverName, address)
}

func NewWeb(sourcePath, indexPath string, embeddedSource *embed.FS) *web.Web {
	return web.NewWeb(sourcePath, indexPath, embeddedSource)
}
