package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	omi "github.com/stormi-li/omiv1"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {

	register := omi.NewRegister(&omi.Options{Addr: redisAddr, Password: password})
	http.HandleFunc("/http_hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello, send by http")
	})
	http.HandleFunc("/websocket_hello", func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		c, _ := upgrader.Upgrade(w, r, nil)
		c.WriteMessage(1, []byte("hello, send by websocket"))
	})
	register.RegisterAndServeTLS("hello_service", "stormili.site:8100", "../../../certs/stormili.crt", "../../../certs/stormili.key", nil)
}
