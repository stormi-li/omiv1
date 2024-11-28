package server

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-redis/redis/v8"
	register "github.com/stormi-li/omiv1/omiregister"
	"github.com/vmihailenco/msgpack"
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

func (Server) HandleFunc(pattern string, handleFunc func(w http.ResponseWriter, r *http.Request, rw ReadWriter)) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		handleFunc(w, r, rw)
	})
}

var rw = ReadWriter{}

type ReadWriter struct{}

func (ReadWriter) Write(w http.ResponseWriter, v any) error {
	// 序列化为 MsgPack
	data, err := msgpack.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	// 写入响应
	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write response: %w", err)
	}

	return nil
}

func (ReadWriter) ReadRead(r *http.Request, v any) error {
	// 确保读取 Body 的内容
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}

	// 解码到目标对象
	if err := msgpack.Unmarshal(body, v); err != nil {
		return fmt.Errorf("failed to unmarshal body: %w", err)
	}

	return nil
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
