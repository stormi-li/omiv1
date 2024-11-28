package proxy

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/vmihailenco/msgpack"
)

// 主代理器
type Proxy struct {
	HttpProxy      *HTTPProxy
	WebSocketProxy *WebSocketProxy
	Reslover       *Resolver
	Transport      *http.Transport
	Client         *http.Client
}

func NewProxy(redisClient *redis.Client, transport *http.Transport) *Proxy {
	resolver := NewResolver(redisClient)
	return &Proxy{
		Transport:      transport,
		Client:         &http.Client{Transport: transport},
		Reslover:       resolver,
		HttpProxy:      NewHTTPProxy(resolver, transport),
		WebSocketProxy: NewWebSocketProxy(resolver, transport),
	}
}

func (p *Proxy) ServeHttp(w http.ResponseWriter, r *http.Request) error {
	if r.Header.Get("Upgrade") == "websocket" && strings.ToLower(r.Header.Get("Connection")) == "upgrade" {
		return p.WebSocketProxy.Forward(w, r)
	} else {
		return p.HttpProxy.Forward(w, r)
	}
}

// Response 组合 http.Response 并扩展方法
type Response struct {
	*http.Response
}

// OmiRead 读取响应的 Body 并解码到 v
func (response *Response) PRead(v any) error {
	if response.Body == nil {
		return fmt.Errorf("response body is nil")
	}

	defer response.Body.Close()

	// 读取 Body 内容
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := msgpack.Unmarshal(body, v); err != nil {
		return fmt.Errorf("failed to decode response body using msgpack: %w", err)
	}

	return nil
}

func (p *Proxy) Post(serverName string, pattern string, v any) (*Response, error) {
	url := url.URL{
		Host: serverName,
		Path: pattern,
	}
	targetURL, err := p.Reslover.Resolve(url)
	if err != nil {
		return nil, err
	}
	// 将 v 序列化为 JSON 数据
	jsonData, err := msgpack.Marshal(v)
	if err != nil {
		return nil, err
	}

	// 发起 POST 请求
	resp, err := p.Client.Post(targetURL.String(), "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}

	return &Response{Response: resp}, nil
}
