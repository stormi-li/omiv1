package main

import (
	"net/http"

	omi "github.com/stormi-li/omiv1"
)

var RedisAddr = "localhost:6379"

func main() {
	web := omi.NewWeb(nil)

	options := &omi.Options{Addr: RedisAddr}

	register := omi.NewRegister(options)
	proxy := omi.NewProxy(options)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if web.ServeWeb(w, r) {
			return
		}
		proxy.ServePathProxy(w, r)
	})
	http.NewServeMux()
	register.RegisterAndServeTLS("localhost", "localhost:8081", func(port string) {
		http.ListenAndServeTLS(port, "server.crt", "server.key", nil)
	})
}
