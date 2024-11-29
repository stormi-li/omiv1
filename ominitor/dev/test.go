package main

import (
	"github.com/go-redis/redis/v8"
	omi "github.com/stormi-li/omiv1"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	rc := redis.NewClient(&redis.Options{Addr: redisAddr, Password: password})
	m := omi.NewMonitor(rc)
	m.EmbedModel = false
	m.Start("118.25.196.166:8989")
}
