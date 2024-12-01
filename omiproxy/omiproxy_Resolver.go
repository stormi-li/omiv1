package proxy

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-redis/redis/v8"
)

// 定义接口用于服务地址解析
type Resolver struct {
	Router *Router
}

func NewResolver(redisClient *redis.Client) *Resolver {
	return &Resolver{
		Router: NewRouter(redisClient),
	}
}

const ProxyHost = "0.0.0.0:0"

func (resolver *Resolver) Resolve(r http.Request, iswebsocket bool) (*http.Request, error) {
	serverName := strings.Split(r.URL.Path, "/")[1]
	domainName := ""
	parts := strings.Split(r.Host, ":")
	r.Host = ProxyHost
	if len(parts) > 0 {
		domainName = parts[0]
	}
	if resolver.Router.Has(domainName) {
		r.URL.Host = resolver.Router.GetAddress(domainName)
		r.URL.Scheme = resolver.Router.addressMap[domainName][r.URL.Host]["Protocal"]
	} else if resolver.Router.Has(serverName) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/"+serverName)
		r.URL.Host = resolver.Router.GetAddress(serverName)
		r.URL.Scheme = resolver.Router.addressMap[serverName][r.URL.Host]["Protocal"]
	} else {
		return nil, fmt.Errorf("解析失败: %s", r.URL.String())
	}
	if r.URL.Scheme == "" {
		r.URL.Scheme = "http"
	}
	if iswebsocket {
		if r.URL.Scheme == "http" {
			r.URL.Scheme = "ws"
		} else {
			r.URL.Scheme = "wss"
		}
	}
	return &r, nil
}
