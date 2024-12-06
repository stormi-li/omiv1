package main

import (
	omi "github.com/stormi-li/omiv1"
	monitor "github.com/stormi-li/omiv1/ominitor"
)

var RedisAddr = "localhost:6379"

func main() {
	omiClient := omi.NewClient(&omi.Options{Addr: RedisAddr})
	monitorMux := monitor.NewMonitorMux(omiClient)
	omiClient.RegisterAndServe("monitor", "118.25.196.166:9013", monitorMux)
}
