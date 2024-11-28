package proxy

import (
	"bytes"
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

type captureResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

// WriteHeader 捕获状态码
func (w *captureResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// Write 捕获响应体
func (w *captureResponseWriter) Write(data []byte) (int, error) {
	w.body.Write(data)                  // 将响应数据写入缓冲区
	return w.ResponseWriter.Write(data) // 继续写入原始响应
}

func (p *HTTPProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) *CapturedResponse {
	r.URL.Host = r.Host
	targetURL, err := p.Resolver.Resolve(*r.URL)
	if err != nil {
		return &CapturedResponse{
			Error:     err,
			Schema:    r.URL.Scheme,
			OriginURL: *r.URL,
		}
	}

	if targetURL.Scheme == "" {
		targetURL.Scheme = "http"
	}

	proxyURL := &url.URL{
		Scheme: targetURL.Scheme,
		Host:   targetURL.Host,
	}
	cw := captureResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
	proxy := httputil.NewSingleHostReverseProxy(proxyURL)
	r.URL.Path = targetURL.Path
	proxy.Transport = p.Transport
	proxy.ServeHTTP(&cw, r)
	return &CapturedResponse{
		StatusCode: cw.statusCode,
		Body:       cw.body,
		Schema:     r.URL.Scheme,
		OriginURL:  *r.URL,
		TargetURL:  *proxyURL,
	}
}
