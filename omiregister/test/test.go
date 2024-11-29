package main

import (
	"github.com/go-redis/redis/v8"
	register "github.com/stormi-li/omiv1/omiregister"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	redisClient := redis.NewClient(&redis.Options{Addr: redisAddr, Password: password})
	register.NewRegister(redisClient).Register("server1", "localhost:9999")
	register.NewRegister(redisClient).Register("server1", "localhost:9998")
	register.NewRegister(redisClient).Register("server1", "localhost:9997")
	register.NewRegister(redisClient).Register("server2", "localhost:9996")
	r := register.NewRegister(redisClient)
	r.AddMessageHandleFunc("SetCache", func(message string) {})
	r.Register("server2", "localhost:9995")
	register.NewRegister(redisClient).Register("server3", "localhost:9994")
	register.NewRegister(redisClient).Register("server4", "localhost:9993")
	select {}
}
