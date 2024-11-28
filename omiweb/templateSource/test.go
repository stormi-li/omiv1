package main

import (
	"embed"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

//go:embed static/*
var embeddedSource embed.FS
