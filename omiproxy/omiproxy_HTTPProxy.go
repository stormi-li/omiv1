package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

// HTTP代理
type HTTPProxy struct {
	Transport *http.Transport
	Resolver  *Resolver
}

func NewHTTPProxy(resolver *Resolver, transport *http.Transport) *HTTPProxy {
	return &HTTPProxy{
		Transport: transport,
		Resolver:  resolver,
	}
}

func (p *HTTPProxy) Forward(w http.ResponseWriter, r *http.Request) error {
	r.URL.Host = r.Host
	targetURL, err := p.Resolver.Resolve(*r.URL)
	if err != nil {
		return err
	}

	if targetURL.Scheme == "" {
		targetURL.Scheme = "http"
	}

	proxyURL := &url.URL{
		Scheme: targetURL.Scheme,
		Host:   targetURL.Host,
	}

	proxy := httputil.NewSingleHostReverseProxy(proxyURL)
	r.URL.Path = targetURL.Path
	proxy.Transport = p.Transport
	proxy.ServeHTTP(w, r)
	return nil
}
