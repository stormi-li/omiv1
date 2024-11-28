package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-redis/redis/v8"
	register "github.com/stormi-li/omiv1/omiregister"
)

type Server struct {
	ServerName string
	Address    string
	Register   *register.Register
	Port       string
	Weight     int
}

func NewServer(redisClient *redis.Client, serverName, address string) *Server {
	port := ":" + strings.Split(address, ":")[1]
	return &Server{
		ServerName: serverName,
		Address:    address,
		Register:   register.NewRegister(redisClient, serverName, address),
		Weight:     1,
		Port:       port,
	}
}

func (server *Server) Start(handler http.Handler) {
	server.RegisterServer()
	err := http.ListenAndServe(server.Port, handler)
	fmt.Println(err)
}

func (server *Server) StartTLS(certFile, keyFile string, handler http.Handler) {
	server.RegisterServerTLS()
	err := http.ListenAndServeTLS(server.Port, certFile, keyFile, handler)
	fmt.Println(err)
}

func (server *Server) RegisterServer() {
	server.Register.Register(register.HTTP)
}

func (server *Server) RegisterServerTLS() {
	server.Register.Register(register.HTTPS)
}
