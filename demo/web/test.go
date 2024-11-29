package main

// import (
// 	"fmt"
// 	"net/http"

// 	"github.com/go-redis/redis/v8"
// 	"github.com/stormi-li/omi-v1"
// 	"github.com/stormi-li/omi-v1/omidto"
// )

// var redisAddr = "118.25.196.166:3934"
// var password = "12982397StrongPassw0rd"

// func main() {
// 	c := omi.NewClient(&redis.Options{Addr: redisAddr, Password: password})
// 	proxy := c.NewProxy(&http.Transport{})
// 	omiserver := c.NewServer("web_demo", "118.25.196.166:8999")
// 	omiserver.HandleFunc("/", func(w omidto.ResponseWriter, r *omidto.Request) {
// 		err := proxy.ServeHttp(w, r.Request)
// 		fmt.Println(err)
// 	})
// 	omiserver.HandleFunc("/", func(w omidto.ResponseWriter, r *omidto.Request) {

// 	})
// 	omiserver.Start(nil)
// }
