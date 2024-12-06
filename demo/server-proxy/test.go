package main

import (
	"crypto/tls"
	"net/http"

	omi "github.com/stormi-li/omiv1"
	"github.com/stormi-li/omiv1/omihttp"
)

var RedisAddr = "localhost:6379"

func main() {
	options := &omi.Options{Addr: RedisAddr, CacheDir: "cache"}

	omiClient := omi.NewClient(options)

	omiClient.Proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	mux := omiClient.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request, rw *omihttp.ReadWriter) {
		omiClient.ServeDomainProxy(w, r)
	})

	omiClient.RegisterAndServe("http-80代理", "localhost:80", mux)
}
