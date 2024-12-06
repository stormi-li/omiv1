package main

import (
	"net/http"

	omi "github.com/stormi-li/omiv1"
	"github.com/stormi-li/omiv1/omirpc"
)

func main() {
	http.NewServeMux()
	c := omi.NewClient(&omi.Options{
		Addr:     "localhost:6379",
		CacheDir: "cache",
	})
	mux := c.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request, rw *omirpc.ReadWriter) {})
	c.RegisterAndServe("test", "localhost:8899", mux)
}
