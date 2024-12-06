package proxy

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/go-redis/redis/v8"
	cache "github.com/stormi-li/omiv1/omiproxy/omicache"
	"github.com/stormi-li/omiv1/omirpc"
	web "github.com/stormi-li/omiv1/omiweb"
)

type Proxy struct {
	HttpProxy      *HTTPProxy
	WebSocketProxy *WebSocketProxy
	Reslover       *Resolver
	Transport      *http.Transport
	Client         *http.Client
	Cache          *cache.Cache
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
		p.HttpProxy = NewHTTPProxy(p.Transport)
		p.WebSocketProxy = NewWebSocketProxy(p.Transport)
		p.Client = &http.Client{Transport: p.Transport}
	}
}

func (p *Proxy) ServeProxy(w http.ResponseWriter, r *http.Request, isDomain bool) *CapturedResponse {
	p.initProxy()
	if p.Cache != nil {
		data := p.Cache.Get(r.URL.String())
		if len(data) != 0 {
			web.WriterHeader(w, r)
			w.Write(data)
			return &CapturedResponse{TargetURL: r.URL}
		}
	}

	var targetURL *url.URL
	var err error
	if isDomain {
		targetURL, err = p.Reslover.ResolveDomain(r)
	} else {
		targetURL, err = p.Reslover.ResolvePath(r)
	}

	if err != nil {
		return &CapturedResponse{Error: err}
	}

	var cr *CapturedResponse
	if r.Header.Get("Upgrade") == "websocket" && strings.ToLower(r.Header.Get("Connection")) == "upgrade" {
		cr = p.WebSocketProxy.ServeWebSocket(w, r, targetURL)
	} else {
		cr = p.HttpProxy.ServeHTTP(w, r, targetURL)
	}

	if p.Cache != nil && r.Method == "GET" && cr.StatusCode == http.StatusOK && len(cr.Body) > 0 {
		p.Cache.Set(r.URL.String(), cr.Body)
	}
	return cr
}

func (p *Proxy) ServeDomainProxy(w http.ResponseWriter, r *http.Request) *CapturedResponse {
	return p.ServeProxy(w, r, true)
}

func (p *Proxy) ServePathProxy(w http.ResponseWriter, r *http.Request) *CapturedResponse {
	return p.ServeProxy(w, r, false)
}

func (p *Proxy) Call(serverName string, pattern string, v any) (*omirpc.Response, error) {
	p.initProxy()
	url := url.URL{Path: pattern}
	targetR, err := p.Reslover.Resolve(http.Request{URL: &url, Host: serverName, Header: http.Header{}}, false)
	if err != nil {
		return nil, err
	}

	return omirpc.Call(p.Client, targetR.URL.String(), v)
}
