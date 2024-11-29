package main

import (
	"net/http"

	web "github.com/stormi-li/omiv1/omiweb"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	omiweb := web.NewWeb("static", "/index.html", nil)
	// client := omi.NewClient(&redis.Options{Addr: redisAddr, Password: password})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		omiweb.ServeWeb(w, r)
	})
	// omihttp.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request, rw omihttp.ReadWriter) {
	// 	var user User
	// 	rw.Read(r, &user)
	// 	user.Name = "hello " + user.Name
	// 	rw.Write(w, &user)
	// })
	// omihttp.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request, rw omihttp.ReadWriter) {
	// 	var user = User{Name: "lili"}
	// 	rw.Write(w, &user)
	// })
	http.ListenAndServe(":8789", nil)
}

type User struct {
	Id    int
	ID    int
	Name  string
	Email string
}
