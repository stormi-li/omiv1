package main

import (
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	// WebSocket 服务器地址
	serverURL := "wss://example.com/ws"

	// HTTP 代理服务器地址
	proxyURL := "http://your-proxy-server:8080"

	// 创建代理 URL
	proxy, err := url.Parse(proxyURL)
	if err != nil {
		log.Fatalf("解析代理 URL 失败: %v", err)
	}

	// 自定义 HTTP Transport 并设置代理
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxy),
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	// 自定义 WebSocket Dialer
	dialer := websocket.Dialer{
		Proxy:            http.ProxyURL(proxy),             // 设置代理
		NetDialContext:   transport.DialContext,            // 使用自定义的 HTTP Transport
		TLSClientConfig:  transport.TLSClientConfig,        // 可选：TLS 配置
	}

	// 连接到 WebSocket 服务器
	conn, _, err := dialer.Dial(serverURL, nil)
	if err != nil {
		log.Fatalf("连接到 WebSocket 失败: %v", err)
	}
	defer conn.Close()

	log.Println("成功连接到 WebSocket 服务器")

	// 示例：发送消息到服务器
	err = conn.WriteMessage(websocket.TextMessage, []byte("Hello, WebSocket!"))
	if err != nil {
		log.Printf("发送消息失败: %v", err)
	}

	// 示例：接收服务器消息
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("读取消息失败: %v", err)
			break
		}
		log.Printf("收到消息: %s", message)
	}
}