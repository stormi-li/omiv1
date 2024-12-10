package main

import (
	"net/http"

	omi "github.com/stormi-li/omiv1"
	"github.com/stormi-li/omiv1/omihttp"
)

var RedisAddr = "localhost:6379"

func main() {
	options := &omi.Options{
		Addr:     RedisAddr,
		Password: "",
		DB:       0,
		CacheDir: "cache", //启用缓存，缓存路径为“cache”
	}

	omiClient := omi.NewClient(options)

	mux := omiClient.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request, rw *omihttp.ReadWriter) {
		omiClient.ServeDomainProxy(w, r) //域名代理，解析域名，并将请求转发
	})

	omiClient.RegisterAndServe("http-80代理", "localhost:80", mux) //注册并启动代理服务，代理80端口
}
