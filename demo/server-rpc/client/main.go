package main

import (
	"fmt"

	omi "github.com/stormi-li/omiv1"
	"github.com/stormi-li/omiv1/demo/server-rpc/client/packages/person"
)

var RedisAddr = "localhost:6379"

func main() {
	options := &omi.Options{Addr: RedisAddr}

	proxy := omi.NewProxy(options)

	p := &person.Person{Name: "lili"}
	response, err := proxy.Post("rpc", "/protobuf", p)
	if err != nil {
		fmt.Println("请求错误:", err)
	} else {
		// 读取返回数据
		response.Read(p)
		fmt.Println(p.Name)
	}

	pr := &Person{Name: "hushuang"}
	// 发送请求
	response, err = proxy.Post("rpc", "/json", pr)
	if err != nil {
		fmt.Println("请求错误:", err)
	} else {
		// 读取返回数据
		response.Read(pr)
		fmt.Println(pr.Name)
	}
}

type Person struct {
	Name string
}
