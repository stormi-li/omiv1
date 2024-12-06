package main

import (
	"net/http"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	mux := http.NewServeMux()
	http.ListenAndServe(":8000", mux)

}
