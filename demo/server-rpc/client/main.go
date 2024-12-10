package main

import (
	"crypto/tls"
	"fmt"
	"net/http"

	omi "github.com/stormi-li/omiv1"
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

	p := Person{Name: "stormi-li"}

	response, err := omiClient.Post("rpc_hello", "/hello", &p)
	if err == nil {
		response.Read(&p)
		fmt.Println(p.Name)
	}
}

type Person struct {
	Name string
}
