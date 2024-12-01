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

func NewProxy(redisClient *redis.Client) *Proxy {
	resolver := NewResolver(redisClient)
	return &Proxy{
		Transport:     &http.Transport{},
		Reslover:      resolver,
		UnMarshalFunc: omihttp.UnMarshalFunc,
		MarshalFunc:   omihttp.MarshalFunc,
	}
}

type CapturedResponse struct {
	StatusCode int
	Body       bytes.Buffer
	Error      error
	TargetURL  url.URL
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

func (p *Proxy) Post(serverName string, pattern string, v any) (*omihttp.Response, error) {
	p.initProxy()
	url := url.URL{
		Host: serverName,
		Path: pattern,
	}
	targetR, err := p.Reslover.Resolve(http.Request{URL: &url}, false)
	if err != nil {
		return nil, err
	}
	// 将 v 序列化为 JSON 数据
	jsonData, err := p.MarshalFunc(v)
	if err != nil {
		return nil, err
	}

	// 发起 POST 请求
	resp, err := p.Client.Post(targetR.URL.String(), "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}

	return &omihttp.Response{Response: resp, UnMarshalFunc: p.UnMarshalFunc}, nil
}
