package main

import (
	"net/http"

	omi "github.com/stormi-li/omiv1"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	// cache := omi.NewCache("cache", 1024*1024*1024)
	proxy := omi.NewProxy(&omi.Options{Addr: redisAddr, Password: password})
	register := omi.NewRegister(&omi.Options{Addr: redisAddr, Password: password})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// if cache.ServeCache(w, r) {
		// 	return
		// }
		proxy.ServeProxy(w, r)
		// cache.UpdateCache(cr)
	})

	register.RegisterTLS("代理http-8000", "stormili.site:8000", "../../../certs/stormili.crt", "../../../certs/stormili.key", nil)
}
