package main

import (
	"net/http"

	omi "github.com/stormi-li/omiv1"
	"github.com/stormi-li/omiv1/omihttp"
)

var RedisAddr = "localhost:6379"

func main() {
	options := &omi.Options{Addr: RedisAddr}

	c := omi.NewClient(options)
	mux := c.NewServeMux()

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request, rw *omihttp.ReadWriter) {
		p := Person{}
		err := rw.Read(&p)
		if err == nil {
			p.Name = "hello " + p.Name
			rw.Write(&p)
		}
	})
	c.RegisterAndServeTLS("rpc_hello", "localhost:9015", "server.crt", "server.key", mux)
}

type Person struct {
	Name string
}
