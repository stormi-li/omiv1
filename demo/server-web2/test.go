package main

import (
	"net/http"

	omi "github.com/stormi-li/omiv1"
	"github.com/stormi-li/omiv1/omirpc"
)

var RedisAddr = "localhost:6379"

func main() {
	web := omi.NewWeb(nil)

	options := &omi.Options{Addr: RedisAddr}

	omiClient := omi.NewClient(options)
	mux := omiClient.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request, rw *omirpc.ReadWriter) {
		if web.ServeWeb(w, r) {
			return
		}
		omiClient.ServePathProxy(w, r)
	})
	omiClient.RegisterAndServeTLS("localhost", "localhost:8081", "server.crt", "server.key", mux)
}
