package main

import (
	"encoding/json"
	"fmt"

	"github.com/stormi-li/omiv1/demo/server-rpc/client/packages/person"
	rpc "github.com/stormi-li/omiv1/omirpc"
)

func main() {
	p := &person.Person{Name: "lili", Age: 24, Email: "23df@fs.com"}
	data, _ := rpc.ProtobufMarshal(p)
	fmt.Println(len(data))
	p.Email = "ff"
	data, _ = rpc.ProtobufMarshal(p)
	fmt.Println(len(data))
	p.Email = "f"
	data, _ = rpc.ProtobufMarshal(p)
	fmt.Println(len(data))
	p.Email = ""
	p.Age = 0
	data, _ = rpc.ProtobufMarshal(p)
	fmt.Println(len(data))
	pr := Person{Name: "lili"}
	data, _ = json.Marshal(pr)
	fmt.Println(len(data))
}

type Person struct {
	Name  string
	Age   int32
	Email string
}
