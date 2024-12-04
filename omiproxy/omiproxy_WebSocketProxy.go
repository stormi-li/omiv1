package proxy

import (
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	register "github.com/stormi-li/omiv1/omiregister"
)

// WebSocket代理
type WebSocketProxy struct {
	Dialer websocket.Dialer
}

func NewWebSocketProxy(transport *http.Transport) *WebSocketProxy {
	return &WebSocketProxy{
		Dialer: websocket.Dialer{
			NetDialContext:  transport.DialContext,
			TLSClientConfig: transport.TLSClientConfig,
		},
	}
}

var upgrader = websocket.Upgrader{}

func (wp *WebSocketProxy) ServeWebSocket(w http.ResponseWriter, r *http.Request, targetURL *url.URL) *CapturedResponse {
	clientConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return &CapturedResponse{
			Error: err,
		}
	}
	defer clientConn.Close()

	if targetURL.Scheme == string(register.HTTPS) {
		targetURL.Scheme = "wss"
	} else {
		targetURL.Scheme = "ws"
	}

	targetConn, _, err := wp.Dialer.Dial(targetURL.String(), nil)
	if err != nil {
		return &CapturedResponse{
			Error:     err,
			TargetURL: targetURL,
		}
	}
	defer targetConn.Close()

	wp.proxy(clientConn, targetConn)

	return &CapturedResponse{
		TargetURL: targetURL,
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
