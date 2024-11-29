package main

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	omi "github.com/stormi-li/omiv1"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

type User struct {
	ID    int
	Name  string
	Email string
}

func main() {
	proxy := omi.NewProxy(redis.NewClient(&redis.Options{Addr: redisAddr, Password: password}))
	resp, err := proxy.Post("hello_server", "/user", &User{ID: 1, Name: ""})
	if err != nil {
		fmt.Println(err)
		return
	}
	var user User
	resp.Read(&user)
	fmt.Println(user)
}
