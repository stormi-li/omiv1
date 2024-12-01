package main

import (
	omi "github.com/stormi-li/omiv1"
)

var RedisAddr = "localhost:6379"

func main() {
	options := &omi.Options{Addr: RedisAddr}

	m := omi.NewMonitor(options)

	m.Start("localhost:9013")
}
