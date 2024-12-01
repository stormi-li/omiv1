package main

import (
	"net/http"

	omi "github.com/stormi-li/omiv1"
)

var RedisAddr = "localhost:6379"

func main() {
	web := omi.NewWeb(nil)
	web.GenerateTemplate()

	options := &omi.Options{Addr: RedisAddr}

	proxy := omi.NewProxy(options)
	register := omi.NewRegister(options)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if web.ServeWeb(w, r) {
			return
		}
		proxy.ServeProxy(w, r)
	})

	register.RegisterAndServe("localhost", "localhost:9014", func(address string) {
		http.ListenAndServe(address, nil)
	})
}
