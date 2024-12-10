package main

import (
	omi "github.com/stormi-li/omiv1"
)

var RedisAddr = "localhost:6379"

func main() {
	options := &omi.Options{Addr: RedisAddr}

	omiClient := omi.NewClient(options)

	p := &Person{Name: "stormi-li"}

	response, err := omiClient.Call("server_demo", "/hello", p)
	if err == nil {
		response.Read(p)
	}
}

type Person struct {
	Name string
}
