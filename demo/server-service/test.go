package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	omi "github.com/stormi-li/omiv1"
)

var RedisAddr = "localhost:6379"

func main() {
	options := &omi.Options{Addr: RedisAddr}

	register := omi.NewRegister(options)

	register.AddRegisterHandleFunc("Handlers", func() string {
		return "http,websocket"
	})
	openHttp := true
	register.AddMessageHandleFunc("SwitchHttpFunc", func(message string) {
		if message == "open" {
			openHttp = true
			register.AddRegisterHandleFunc("Handlers", func() string {
				return "http,websocket"
			})
		} else if message == "close" {
			openHttp = false
			register.AddRegisterHandleFunc("Handlers", func() string {
				return "websocket"
			})
		}
	})
	
	http.HandleFunc("/http", func(w http.ResponseWriter, r *http.Request) {
		if openHttp {
			fmt.Fprintf(w, "hello, send by http")
		} else {
			fmt.Fprintf(w, "http service is closed")
		}
	})

	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		c, _ := upgrader.Upgrade(w, r, nil)
		c.WriteMessage(1, []byte("hello, send by websocket"))
		time.Sleep(100 * time.Millisecond)
		c.Close()
	})

	register.RegisterAndServe("hello", "localhost:9014", func(port string) {
		http.ListenAndServe(port, nil)
	})
}
