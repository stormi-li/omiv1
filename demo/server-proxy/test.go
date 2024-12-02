package main

import (
	"fmt"
	"net/http"

	omi "github.com/stormi-li/omiv1"
)

var RedisAddr = "localhost:6379"

func main() {
	options := &omi.Options{Addr: RedisAddr}

	proxy := omi.NewProxy(options)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cw := proxy.ServeProxy(w, r)
		fmt.Println(cw.Error)
	})

	register := omi.NewRegister(options)
	register.RegisterAndServe("fsdf", "localhost:80", func(port string) {
		http.ListenAndServe(port, nil)
	})
}
