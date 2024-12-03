package main

import (
	"fmt"
	"net/http"

	omi "github.com/stormi-li/omiv1"
	"github.com/stormi-li/omiv1/demo/server-rpc/server/packages/person"
)

var RedisAddr = "localhost:6379"

func main() {

	options := &omi.Options{Addr: RedisAddr}

	register := omi.NewRegister(options)

	http.HandleFunc("/protobuf", func(w http.ResponseWriter, r *http.Request) {
		p := person.Person{}
		err := omi.ReadWriter.Read(r, &p)
		if err != nil {
			fmt.Println(err)
			return
		}
		p.Name = "hello " + p.Name
		omi.ReadWriter.Write(w, &p)
	})
	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		p := Person{}
		err := omi.ReadWriter.Read(r, &p)
		if err != nil {
			fmt.Println(err)
			return
		}
		p.Name = "hello " + p.Name
		omi.ReadWriter.Write(w, &p)
	})

	register.RegisterAndServe("rpc", "localhost:9015", func(port string) {
		http.ListenAndServe(port, nil)
	})
}

type Person struct {
	Name  string
	Age   int32
	Email string
}
