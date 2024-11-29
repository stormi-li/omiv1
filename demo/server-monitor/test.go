package main

import (
	omi "github.com/stormi-li/omiv1"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	m := omi.NewMonitor(&omi.Options{Addr: redisAddr, Password: password})
	m.Start("118.25.196.166:8989")
}
