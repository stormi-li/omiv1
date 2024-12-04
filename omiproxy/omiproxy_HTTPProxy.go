package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

// HTTP代理
type HTTPProxy struct {
	Transport *http.Transport
}

func NewHTTPProxy(transport *http.Transport) *HTTPProxy {

	return &HTTPProxy{
		Transport: transport,
	}
}

func (p *HTTPProxy) ServeHTTP(w http.ResponseWriter, r *http.Request, targetURL *url.URL) *CapturedResponse {
	reverseProxy := httputil.NewSingleHostReverseProxy(targetURL)
	reverseProxy.Transport = p.Transport

	cw := CaptureResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

	originalPath := r.URL.Path
	r.URL.Path = targetURL.Path
	targetURL.Path = ""
	reverseProxy.ServeHTTP(&cw, r)
	targetURL.Path = r.URL.Path
	r.URL.Path = originalPath
	return &CapturedResponse{
		StatusCode: cw.statusCode,
		Body:       cw.body,
		TargetURL:  targetURL,
	}
}
