package main

import (
	"embed"
	"net/http"

	web "github.com/stormi-li/omiv1/omiweb"
)

//go:embed static/*
var embeddedSource embed.FS

func main() {
	omiweb := web.NewWeb(nil)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		omiweb.ServeWeb(w, r)
	})
	http.ListenAndServe(":8789", nil)
}
