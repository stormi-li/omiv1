package main

import (
	"net/http"

	omi "github.com/stormi-li/omiv1"
)

func main() {
	web := omi.NewWeb(nil)
	web.GenerateTemplate()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		web.ServeWeb(w, r)
	})

	http.ListenAndServe(":5500", nil)
}
