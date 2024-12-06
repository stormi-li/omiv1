package main

import (
	"net/http"

	"github.com/go-redis/redis/v8"
	register "github.com/stormi-li/omiv1/omiregister"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	mux := http.NewServeMux()
	http.ListenAndServe(":8000", mux)

}

func Test() {
	redisClient := redis.NewClient(&redis.Options{Addr: redisAddr, Password: password})
	register.NewRegister(redisClient).RegisterAndServe("server1", "localhost:9999", func(port string) {})
	register.NewRegister(redisClient).RegisterAndServe("server1", "localhost:9998", func(port string) {})
	register.NewRegister(redisClient).RegisterAndServe("server1", "localhost:9997", func(port string) {})
	register.NewRegister(redisClient).RegisterAndServe("server2", "localhost:9996", func(port string) {})
	r := register.NewRegister(redisClient)
	r.AddMessageHandleFunc("SetCache", func(message string) {})
	r.RegisterAndServe("server2", "localhost:9995", func(port string) {})
	register.NewRegister(redisClient).RegisterAndServe("server3", "localhost:9994", func(port string) {})
	register.NewRegister(redisClient).RegisterAndServe("server4", "localhost:9993", func(port string) {})
	register.NewRegister(redisClient).RegisterAndServe("server4", "localhost:9293", func(port string) {})
	register.NewRegister(redisClient).RegisterAndServe("server4", "localhost:9193", func(port string) {})
	register.NewRegister(redisClient).RegisterAndServe("server5", "localhost:9393", func(port string) {})
	register.NewRegister(redisClient).RegisterAndServe("server5", "localhost:9493", func(port string) {})
	register.NewRegister(redisClient).RegisterAndServe("server5", "localhost:9593", func(port string) {})
	register.NewRegister(redisClient).RegisterAndServe("server6", "localhost:9693", func(port string) {})
	register.NewRegister(redisClient).RegisterAndServe("server7", "localhost:9793", func(port string) {})
	register.NewRegister(redisClient).RegisterAndServe("server8", "localhost:9893", func(port string) {})
	register.NewRegister(redisClient).RegisterAndServe("server9", "localhost:9093", func(port string) {})
	register.NewRegister(redisClient).RegisterAndServe("server9", "localhost:9903", func(port string) {})
	select {}
}
