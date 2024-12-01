package main

import (
	"net/http"

	omi "github.com/stormi-li/omiv1"
)

var RedisAddr = "localhost:6379"

func main() {
	options := &omi.Options{Addr: RedisAddr}

	proxy := omi.NewProxy(options)
	register := omi.NewRegister(options)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeProxy(w, r)
	})

	register.RegisterAndServe("http-80代理", "localhost:80", func(port string) {
		http.ListenAndServe(port, nil)
	})
}
