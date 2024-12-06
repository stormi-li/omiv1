package main

import (
	"fmt"
	"net/http"

	omi "github.com/stormi-li/omiv1"
	"github.com/stormi-li/omiv1/demo/server-rpc/server/packages/person"
	"github.com/stormi-li/omiv1/omihttp"
)

var RedisAddr = "localhost:6379"

func main() {
	options := &omi.Options{Addr: RedisAddr}

	c := omi.NewClient(options)
	mux := c.NewServeMux()
	mux.HandleFunc("/protobuf", func(w http.ResponseWriter, r *http.Request, rw *omihttp.ReadWriter) {
		p := person.Person{}
		err := rw.Read(&p)
		if err != nil {
			fmt.Println(err)
			return
		}
		p.Name = "hello " + p.Name
		rw.Write(&p)
	})

	mux.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request, rw *omihttp.ReadWriter) {
		p := Person{}
		err := rw.Read(&p)
		if err != nil {
			fmt.Println(err)
			return
		}
		p.Name = "hello " + p.Name
		rw.Write(&p)
	})

	c.RegisterAndServe("rpc", "localhost:9015", mux)
}

type Person struct {
	Name  string
	Age   int32
	Email string
}
