package proxy

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-redis/redis/v8"
	register "github.com/stormi-li/omiv1/omiregister"
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

const KeyOriginalPath = "Omi-OriginalPath"
const KeyProxyNodes = "Omi-ProxyNodes"
const KeyClientAddr = "Omi-ClientAddr"

func (resolver *Resolver) Resolve(r http.Request, iswebsocket bool) (*http.Request, error) {
	if r.URL.Path[0] != '/' {
		return nil, fmt.Errorf("路径必须以/开头")
	}
	serverName := strings.Split(r.URL.Path, "/")[1]
	domainName := ""
	parts := strings.Split(r.Host, ":")

	if len(parts) > 0 {
		domainName = parts[0]
	}

	originalPath := r.Header.Get(KeyOriginalPath)
	if originalPath == "" {
		r.Header.Set(KeyOriginalPath, r.URL.Path)
		r.Header.Set(KeyClientAddr, r.RemoteAddr)
	} else {
		domainName = ""
	}

	addresses := r.Header.Get(KeyProxyNodes)
	if addresses == "" {
		addresses = register.Address
	} else {
		addresses += "," + register.Address
	}
	r.Header.Set(KeyProxyNodes, addresses)

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

func (resolver *Resolver) ResolveDomain(r *http.Request) (*url.URL, error) {
	targetURL := *r.URL
	parts := strings.Split(r.Host, ":")

	if len(parts) == 0 {
		return nil, fmt.Errorf("域名解析失败: %s", r.Host)
	}
	domainName := parts[0]
	if resolver.Router.Has(domainName) {
		address := resolver.Router.GetAddress(domainName)
		targetURL.Host = address
		targetURL.Scheme = resolver.Router.addressMap[domainName][address]["Protocal"]
	} else {
		return nil, fmt.Errorf("域名解析失败: %s", r.Host)
	}
	return &targetURL, nil
}

func (resolver *Resolver) ResolvePath(r *http.Request) (*url.URL, error) {
	targetURL := *r.URL
	if r.URL.Path[0] != '/' {
		return nil, fmt.Errorf("路径解析失败: %s", r.URL.Path)
	}
	serverName := strings.Split(r.URL.Path, "/")[1]
	if resolver.Router.Has(serverName) {
		targetURL.Path = strings.TrimPrefix(r.URL.Path, "/"+serverName)
		address := resolver.Router.GetAddress(serverName)
		targetURL.Host = address
		targetURL.Scheme = resolver.Router.addressMap[serverName][address]["Protocal"]
	} else {
		return nil, fmt.Errorf("路径解析失败: %s", r.URL.Path)
	}
	return &targetURL, nil
}
