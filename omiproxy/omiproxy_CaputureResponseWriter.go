package proxy

import (
	"bytes"
	"net/http"
)

type CaptureResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

// WriteHeader 捕获状态码
func (w *CaptureResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// Write 捕获响应体
func (w *CaptureResponseWriter) Write(data []byte) (int, error) {
	w.body.Write(data)                  // 将响应数据写入缓冲区
	return w.ResponseWriter.Write(data) // 继续写入原始响应
}
