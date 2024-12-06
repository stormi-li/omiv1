package omihttp

import (
	"net/http"
)

type ServeMux struct {
	*http.ServeMux
	RouterMap map[string]bool
}

func NewServeMux() *ServeMux {
	return &ServeMux{ServeMux: http.NewServeMux(), RouterMap: map[string]bool{}}
}

func (mux *ServeMux) Handle(pattern string, handler Handler) {
	mux.HandleFunc(pattern, handler.ServeHTTP)
}

func (mux *ServeMux) HandleFunc(pattern string, handler func(w http.ResponseWriter, r *http.Request, rw *ReadWriter)) {
	mux.RouterMap[pattern] = true
	mux.ServeMux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if mux.RouterMap[pattern] {
			handler(w, r, NewReadWriter(w, r))
		}
	})
}
