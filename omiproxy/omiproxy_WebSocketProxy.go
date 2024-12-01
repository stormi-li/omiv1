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

var omiproxyhead = "omi-proxy"

func (wp *WebSocketProxy) ServeWebSocket(w http.ResponseWriter, r *http.Request) *CapturedResponse {
	if r.Header.Get(omiproxyhead) == omiproxyhead {
		r.Host = ProxyHost
	}

	targetR, err := wp.Resolver.Resolve(*r)
	if err != nil {
		return &CapturedResponse{
			Error: err,
		}
	}

	clientConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		err := fmt.Errorf("WebSocket升级失败: %v", err)
		return &CapturedResponse{
			Error: err,
		}
	}
	defer clientConn.Close()

	if targetR.URL.Scheme == "https" {
		targetR.URL.Scheme = "wss"
	} else {
		targetR.URL.Scheme = "ws"
	}

	header := http.Header{}
	header.Set(omiproxyhead, omiproxyhead)

	targetConn, _, err := wp.Dialer.Dial(targetR.URL.String(), header)
	if err != nil {
		err := fmt.Errorf("无法连接到WebSocket服务器: %v", err)
		return &CapturedResponse{
			Error:     err,
			TargetURL: *targetR.URL,
		}
	}
	defer targetConn.Close()

	errChan := make(chan error, 2)
	go wp.copyData(clientConn, targetConn, errChan)
	go wp.copyData(targetConn, clientConn, errChan)
	<-errChan

	return &CapturedResponse{
		TargetURL: *targetR.URL,
	}
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
