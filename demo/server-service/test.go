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

	http.HandleFunc("/http", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello, send by http")
	})

	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		c, _ := upgrader.Upgrade(w, r, nil)
		c.WriteMessage(1, []byte("hello, send by websocket"))
		time.Sleep(100 * time.Millisecond)
		c.Close()
	})

	register.RegisterAndServe("hello", "localhost:9015", func(port string) {
		http.ListenAndServe(port, nil)
	})
}
