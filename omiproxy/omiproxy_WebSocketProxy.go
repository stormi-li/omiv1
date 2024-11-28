package proxy

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// WebSocket代理
type WebSocketProxy struct {
	Resolver *Resolver
	Dialer   websocket.Dialer
}

func NewWebSocketProxy(resolver *Resolver, transport *http.Transport) *WebSocketProxy {
	return &WebSocketProxy{
		Resolver: resolver,
		Dialer: websocket.Dialer{
			NetDialContext: transport.DialContext,
		},
	}
}

var upgrader = websocket.Upgrader{}

func (wp *WebSocketProxy) Forward(w http.ResponseWriter, r *http.Request) error {
	clientConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return fmt.Errorf("WebSocket升级失败: %v", err)
	}
	defer clientConn.Close()

	proxyURL, err := wp.Resolver.Resolve(*r.URL)
	if err != nil {
		return err
	}
	if proxyURL.Scheme == "https" {
		proxyURL.Scheme = "wss"
	} else {
		proxyURL.Scheme = "ws"
	}
	targetConn, _, err := wp.Dialer.Dial(proxyURL.String(), nil)
	if err != nil {
		return fmt.Errorf("无法连接到WebSocket服务器: %v", err)
	}
	defer targetConn.Close()

	errChan := make(chan error, 2)
	go wp.copyData(clientConn, targetConn, errChan)
	go wp.copyData(targetConn, clientConn, errChan)
	<-errChan
	return nil
}

func (wp *WebSocketProxy) copyData(src, dst *websocket.Conn, errChan chan error) {
	for {
		msgType, msg, err := src.ReadMessage()
		if err != nil {
			errChan <- err
			return
		}
		if err := dst.WriteMessage(msgType, msg); err != nil {
			errChan <- err
			return
		}
	}
}
