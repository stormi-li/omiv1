package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

// HTTP代理
type HTTPProxy struct {
	Resolver  *Resolver
	Transport *http.Transport
}

func NewHTTPProxy(resolver *Resolver, transport *http.Transport) *HTTPProxy {

	return &HTTPProxy{
		Resolver:  resolver,
		Transport: transport,
	}
}

func (p *HTTPProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) *CapturedResponse {
	originalURL := *r.URL
	targetR, err := p.Resolver.Resolve(*r, false)
	if err != nil {
		return &CapturedResponse{
			Error: err,
		}
	}
	proxyURL := &url.URL{}
	proxyURL.Scheme = targetR.URL.Scheme
	proxyURL.Host = targetR.URL.Host

	reverseProxy := httputil.NewSingleHostReverseProxy(proxyURL)
	reverseProxy.Transport = p.Transport

	cw := CaptureResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

	reverseProxy.ServeHTTP(&cw, targetR)

	return &CapturedResponse{
		StatusCode:  cw.statusCode,
		Body:        cw.body,
		OriginalURL: &originalURL,
		TargetURL:   targetR.URL,
	}
}
