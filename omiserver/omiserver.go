package server

import (
	"fmt"
	"net/http"

	register "github.com/stormi-li/omiv1/omiregister"
)

type Server struct {
	ReadWriter ReadWriter
}

func NewServer() *Server {
	return &Server{
		ReadWriter: ReadWriter{MarshalFunc: MarshalFunc, UnMarshalFunc: UnMarshalFunc},
	}
}

func (server *Server) Start(r *register.Register, handler http.Handler) {
	r.Register(register.HTTP)
	err := http.ListenAndServe(r.Port, handler)
	fmt.Println(err)
}

func (server *Server) StartTLS(r *register.Register, certFile, keyFile string, handler http.Handler) {
	r.Register(register.HTTPS)
	err := http.ListenAndServeTLS(r.Port, certFile, keyFile, handler)
	fmt.Println(err)
}
