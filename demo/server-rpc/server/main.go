package main

import (
	"net/http"

	omi "github.com/stormi-li/omiv1"
	"github.com/stormi-li/omiv1/demo/server-rpc/server/packages/person"
	"github.com/stormi-li/omiv1/omihttp"
	"github.com/stormi-li/omiv1/omihttp/serialization"
)

var RedisAddr = "localhost:6379"

func main() {
	options := &omi.Options{Addr: RedisAddr}

	c := omi.NewClient(options)
	mux := c.NewServeMux()

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request, rw *omihttp.ReadWriter) {
		p := person.Person{}
		err := rw.Read(&p, serialization.Protobuf)
		if err == nil {
			p.Name = "hello " + p.Name
			rw.Write(&p, serialization.Protobuf)
		}
	})
	c.RegisterAndServeTLS("rpc_hello", "localhost:9015", "server.crt", "server.key", mux)
}
