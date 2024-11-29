package main

import (
	"net/http"

	omi "github.com/stormi-li/omiv1"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	web := omi.NewWeb("static", "/index.html", nil)
	web.GenerateTemplate()

	proxy := omi.NewProxy(&omi.Options{Addr: redisAddr, Password: password})
	register := omi.NewRegister(&omi.Options{Addr: redisAddr, Password: password})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if web.ServeWeb(w, r) {
			return
		}
		proxy.ServeProxy(w, r)
	})

	register.RegisterAndServe("web2", "118.25.196.166:8889", nil)
}