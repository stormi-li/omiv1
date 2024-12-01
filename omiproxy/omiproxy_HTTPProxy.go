package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

// HTTP代理
type HTTPProxy struct {
	Resolver     *Resolver
	ProxyURL     *url.URL
	ReverseProxy *httputil.ReverseProxy
}

func NewHTTPProxy(resolver *Resolver, transport *http.Transport) *HTTPProxy {
	proxyURL := &url.URL{}
	reverseProxy := httputil.NewSingleHostReverseProxy(proxyURL)
	reverseProxy.Transport = transport
	return &HTTPProxy{
		Resolver:     resolver,
		ProxyURL:     proxyURL,
		ReverseProxy: reverseProxy,
	}
}


func (p *HTTPProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) *CapturedResponse {
	targetR, err := p.Resolver.Resolve(*r, false)
	if err != nil {
		return &CapturedResponse{
			Error: err,
		}
	}

	p.ProxyURL.Scheme = targetR.URL.Scheme
	p.ProxyURL.Host = targetR.URL.Host

	cw := CaptureResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

	p.ReverseProxy.ServeHTTP(&cw, targetR)

	return &CapturedResponse{
		StatusCode: cw.statusCode,
		Body:       cw.body,
		TargetURL:  *targetR.URL,
	}
}
