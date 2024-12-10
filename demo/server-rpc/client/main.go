package main

import (
	"crypto/tls"
	"fmt"
	"net/http"

	omi "github.com/stormi-li/omiv1"
	"github.com/stormi-li/omiv1/demo/server-rpc/client/packages/person"
	"github.com/stormi-li/omiv1/omihttp/serialization"
)

var RedisAddr = "localhost:6379"

func main() {
	options := &omi.Options{Addr: RedisAddr}

	omiClient := omi.NewClient(options)

	omiClient.SetTransport(&http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	})

	p := person.Person{Name: "stormi-li"}

	response, err := omiClient.Post("rpc_hello", "/hello", &p, serialization.Protobuf)
	if err == nil {
		response.Read(&p, serialization.Protobuf)
		fmt.Println(p.Name)
	}
}
