package main

import (
	"fmt"
	"net/http"

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
	c := omi.NewClient(&redis.Options{Addr: redisAddr, Password: password})
	proxy := c.NewProxy(&http.Transport{})
	resp, err := proxy.Post("hello_server", "/user", &User{ID: 1, Name: ""})
	if err != nil {
		fmt.Println(err)
		return
	}
	var user User
	resp.UnMarshal(&user)
	fmt.Println(user)
}
