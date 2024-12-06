package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	omi "github.com/stormi-li/omiv1"
	"github.com/stormi-li/omiv1/omirpc"
)

var RedisAddr = "localhost:6379"

func main() {
	options := &omi.Options{Addr: RedisAddr}

	omiClient := omi.NewClient(options)
	mux := omiClient.NewServeMux()

	mux.HandleFunc("/http", func(w http.ResponseWriter, r *http.Request, rw *omirpc.ReadWriter) {
		fmt.Fprintf(w, "hello, send by http")

	})

	mux.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request, rw *omirpc.ReadWriter) {
		upgrader := websocket.Upgrader{}
		c, _ := upgrader.Upgrade(w, r, nil)
		c.WriteMessage(1, []byte("hello, send by websocket"))
		time.Sleep(100 * time.Millisecond)
		c.Close()
	})

	omiClient.RegisterAndServe("hello", "localhost:9014", mux)
}
