package main

import (
	"github.com/go-redis/redis/v8"
	register "github.com/stormi-li/omiv1/omiregister"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	redisClient := redis.NewClient(&redis.Options{Addr: redisAddr, Password: password})
	register.NewRegister(redisClient, "server1", "localhost:9999").Register()
	register.NewRegister(redisClient, "server1", "localhost:9998").Register()
	register.NewRegister(redisClient, "server1", "localhost:9997").Register()
	register.NewRegister(redisClient, "server2", "localhost:9996").Register()
	r := register.NewRegister(redisClient, "server2", "localhost:9995")
	r.AddMessageHandleFunc("SetCache", func(message string) {})
	r.Register()
	register.NewRegister(redisClient, "server3", "localhost:9994").Register()
	register.NewRegister(redisClient, "server4", "localhost:9993").Register()
	register.NewRegister(redisClient, "server4", "localhost:9992").Register()
	register.NewRegister(redisClient, "server5", "localhost:9991").Register()
	register.NewRegister(redisClient, "server6", "localhost:9919").Register()
	register.NewRegister(redisClient, "server7", "localhost:9929").Register()
	register.NewRegister(redisClient, "server7", "localhost:9939").Register()
	register.NewRegister(redisClient, "server8", "localhost:9949").Register()
	register.NewRegister(redisClient, "server8", "localhost:9959").Register()
	select {}
}
