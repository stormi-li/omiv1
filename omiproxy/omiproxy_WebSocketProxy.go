package proxy

import (
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
			NetDialContext:  transport.DialContext,
			TLSClientConfig: transport.TLSClientConfig,
		},
	}
}

var upgrader = websocket.Upgrader{}

func (wp *WebSocketProxy) ServeWebSocket(w http.ResponseWriter, r *http.Request) *CapturedResponse {
	clientConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return &CapturedResponse{
			Error: err,
		}
	}
	defer clientConn.Close()

	targetR, err := wp.Resolver.Resolve(*r, true)
	if err != nil {
		return &CapturedResponse{
			Error: err,
		}
	}

	header := http.Header{}

	header.Set(KeyOriginalPath, targetR.Header.Get(KeyOriginalPath))
	header.Set(KeyProxyNodes, targetR.Header.Get(KeyProxyNodes))
	header.Set(KeyClientAddr, targetR.Header.Get(KeyClientAddr))

	targetConn, _, err := wp.Dialer.Dial(targetR.URL.String(), header)
	if err != nil {
		return &CapturedResponse{
			Error:     err,
			TargetURL: targetR.URL,
		}
	}
	defer targetConn.Close()

	wp.proxy(clientConn, targetConn)

	return &CapturedResponse{
		TargetURL: targetR.URL,
	}
}

func (wp *WebSocketProxy) proxy(src, dst *websocket.Conn) {
	errChan := make(chan error, 2)
	go wp.copyData(src, dst, errChan)
	go wp.copyData(dst, src, errChan)
	<-errChan
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
