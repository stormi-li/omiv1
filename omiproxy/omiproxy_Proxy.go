package proxy

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omiv1/omirpc"
	rpc "github.com/stormi-li/omiv1/omirpc"
)

// 主代理器
type Proxy struct {
	HttpProxy      *HTTPProxy
	WebSocketProxy *WebSocketProxy
	Reslover       *Resolver
	Transport      *http.Transport
	Client         *http.Client
}

func NewProxy(redisClient *redis.Client) *Proxy {
	resolver := NewResolver(redisClient)
	return &Proxy{
		Transport: &http.Transport{},
		Reslover:  resolver,
	}
}

func (p *Proxy) AddFilter(serverName, key string, handle func(value string) bool) {
	p.Reslover.Router.Discover.AddFilter(serverName, key, handle)
}

func (p *Proxy) initProxy() {
	if p.HttpProxy == nil || p.WebSocketProxy == nil || p.Client == nil {
		p.HttpProxy = NewHTTPProxy(p.Reslover, p.Transport)
		p.WebSocketProxy = NewWebSocketProxy(p.Reslover, p.Transport)
		p.Client = &http.Client{Transport: p.Transport}
	}
}

func (p *Proxy) ServeProxy(w http.ResponseWriter, r *http.Request) *CapturedResponse {
	p.initProxy()
	if r.Header.Get("Upgrade") == "websocket" && strings.ToLower(r.Header.Get("Connection")) == "upgrade" {
		return p.WebSocketProxy.ServeWebSocket(w, r)
	} else {
		return p.HttpProxy.ServeHTTP(w, r)
	}
}

func (p *Proxy) Post(serverName string, pattern string, v any) (*rpc.Response, error) {
	p.initProxy()
	url := url.URL{Path: pattern}
	targetR, err := p.Reslover.Resolve(http.Request{URL: &url, Host: serverName, Header: http.Header{}}, false)
	if err != nil {
		return nil, err
	}

	return omirpc.Call(p.Client, targetR.URL.String(), v)
}
