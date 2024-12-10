package main

import (
	"fmt"
	"net/http"

	omi "github.com/stormi-li/omiv1"
	"github.com/stormi-li/omiv1/omihttp"
)

var RedisAddr = "localhost:6379"

func main() {
	options := &omi.Options{Addr: RedisAddr}

	omiClient := omi.NewClient(options)
	mux := omiClient.NewServeMux()

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request, rw *omihttp.ReadWriter) {
		fmt.Fprintf(w, "hello, send by http")

	})
	mux.HandleFunc("/hello2", func(w http.ResponseWriter, r *http.Request, rw *omihttp.ReadWriter) {
		fmt.Fprintf(w, "hello2, send by http")

	})

	omiClient.RegisterAndServe("http_hello_service", "localhost:9014", mux)
}
