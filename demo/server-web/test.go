package main

import (
	"net/http"

	omi "github.com/stormi-li/omiv1"
	"github.com/stormi-li/omiv1/omihttp"
)

var RedisAddr = "localhost:6379"

func main() {

	options := &omi.Options{Addr: RedisAddr}

	omiClient := omi.NewClient(options)
	web := omiClient.NewWebServer(nil)
	mux := omiClient.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request, rw *omihttp.ReadWriter) {
		if web.ServeWeb(w, r) {
			return
		}
		omiClient.ServePathProxy(w, r) //启动路径代理服务
	})

	omiClient.RegisterAndServe("localhost", "localhost:8080", mux)
}
