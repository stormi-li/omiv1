package main

import (
	"net/http"

	omi "github.com/stormi-li/omiv1"
)

var RedisAddr = "118.25.196.166:3934"
var Password = "12982397StrongPassw0rd"

func main() {
	options := &omi.Options{Addr: RedisAddr, Password: Password}
	
	proxy := omi.NewProxy(options)
	register := omi.NewRegister(options)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeProxy(w, r)
	})

	register.RegisterAndServe("代理http-80", "stormili.site:8000", func(port string) {
		http.ListenAndServe(port, nil)
	})
}
