package proxy

import (
	"bytes"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omiv1/omihttp"
)

// 主代理器
type Proxy struct {
	HttpProxy      *HTTPProxy
	WebSocketProxy *WebSocketProxy
	Reslover       *Resolver
	Transport      *http.Transport
	Client         *http.Client
	UnMarshalFunc  func(data []byte, v any) error
	MarshalFunc    func(v any) ([]byte, error)
}

func NewProxy(redisClient *redis.Client, transport *http.Transport) *Proxy {
	resolver := NewResolver(redisClient)
	return &Proxy{
		Transport:      transport,
		Client:         &http.Client{Transport: transport},
		Reslover:       resolver,
		HttpProxy:      NewHTTPProxy(resolver, transport),
		WebSocketProxy: NewWebSocketProxy(resolver, transport),
		UnMarshalFunc:  omihttp.UnMarshalFunc,
		MarshalFunc:    omihttp.MarshalFunc,
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	if r.Header.Get("Upgrade") == "websocket" && strings.ToLower(r.Header.Get("Connection")) == "upgrade" {
		return p.WebSocketProxy.ServeWebSocket(w, r)
	} else {
		return p.HttpProxy.ServeHTTP(w, r)
	}
}

// Response 组合 http.Response 并扩展方法

func (p *Proxy) Post(serverName string, pattern string, v any) (*omihttp.Response, error) {
	url := url.URL{
		Host: serverName,
		Path: pattern,
	}
	targetURL, err := p.Reslover.Resolve(url)
	if err != nil {
		return nil, err
	}
	// 将 v 序列化为 JSON 数据
	jsonData, err := p.MarshalFunc(v)
	if err != nil {
		return nil, err
	}

	// 发起 POST 请求
	resp, err := p.Client.Post(targetURL.String(), "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}

	return &omihttp.Response{Response: resp, UnMarshalFunc: p.UnMarshalFunc}, nil
}
