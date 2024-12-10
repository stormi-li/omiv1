package main

import (
	"fmt"
	"net/http"

	omi "github.com/stormi-li/omiv1"
	"github.com/stormi-li/omiv1/omihttp"
)

var RedisAddr = "localhost:6379"

func main() {
	options := &omi.Options{Addr: RedisAddr}

	omiClient := omi.NewClient(options)
	mux := omiClient.NewServeMux()

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request, rw *omihttp.ReadWriter) {
		fmt.Fprintf(w, "hello, send by http")

	})

	// mux.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request, rw *omihttp.ReadWriter) {
	// 	upgrader := websocket.Upgrader{}
	// 	c, _ := upgrader.Upgrade(w, r, nil)
	// 	c.WriteMessage(1, []byte("hello, send by websocket"))
	// 	time.Sleep(100 * time.Millisecond)
	// 	c.Close()
	// })

	omiClient.RegisterAndServe("hello_service", "localhost:9014", mux)
}
