package main

import (
	"net/http"

	"github.com/go-redis/redis/v8"
	omi "github.com/stormi-li/omiv1"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	redisClient := redis.NewClient(&redis.Options{Addr: redisAddr, Password: password})

	web := omi.NewWeb("static", "/index.html", nil)
	web.GenerateTemplate()

	proxy := omi.NewProxy(redisClient)
	register := omi.NewRegister(redisClient)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if web.ServeWeb(w, r) {
			return
		}
		proxy.ServeProxy(w, r)
	})

	register.Register("web2", "118.25.196.166:8889", nil)
}
