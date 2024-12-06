package main

import (
	omi "github.com/stormi-li/omiv1"
	monitor "github.com/stormi-li/omiv1/ominitor"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	omiClient := omi.NewClient(&omi.Options{Addr: redisAddr, Password: password})
	monitorMux := monitor.NewMonitorMux(omiClient)
	omiClient.RegisterAndServe("monitor", "118.25.196.166:8989", monitorMux)
}
