package register

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omiv1/omihttp"
)

const Command_UpdateWeight = "UpdateWeight"

var RegisterInterval = 2 * time.Second

var Address = ""

type Register struct {
	RedisClient     *redis.Client
	ServerName      string
	Address         string
	Weight          int
	Info            map[string]string
	Prefix          string
	Channel         string
	omipc           *Omipc
	ctx             context.Context
	RegisterHandler *RegisterHandler
	MessageHandler  *MessageHandler
	StartTime       time.Time
	Port            string
	regestered      bool
}

func NewRegister(redisClient *redis.Client) *Register {
	register := &Register{
		RedisClient:     redisClient,
		Weight:          1,
		Info:            map[string]string{},
		Prefix:          Prefix,
		ctx:             context.Background(),
		omipc:           NewOmipc(redisClient),
		RegisterHandler: newRegisterHandler(redisClient),
		MessageHandler:  newMessageHander(redisClient),
		StartTime:       time.Now(),
	}
	Init(register)
	return register
}

func (register *Register) AddRegisterHandleFunc(key string, handleFunc func() string) {
	register.RegisterHandler.AddHandleFunc(key, handleFunc)
}

func (register *Register) AddMessageHandleFunc(command string, handleFunc func(message string)) {
	register.MessageHandler.AddHandleFunc(command, handleFunc)
}

type Protocal string

var HTTP Protocal = "http"
var HTTPS Protocal = "https"

func (register *Register) register(protocal Protocal, serverName, address string) {
	if register.regestered {
		panic("该注册器已注册服务：" + register.ServerName)
	}
	register.regestered = true
	if strings.Contains(serverName, Namespace_separator) {
		panic("名字里不能包含字符\"" + Namespace_separator + "\"")
	}
	register.ServerName = serverName
	register.Address = address
	Address = address
	register.Channel = Prefix + serverName + Namespace_separator + address
	parts := strings.Split(address, ":")
	if len(parts) != 2 {
		panic("非法地址：" + address)
	}
	register.Port = ":" + parts[1]

	log.Printf("%s server is registered on redis:%s with %s://%s", register.ServerName, register.RedisClient.Options().Addr, protocal, register.Address)
	register.AddRegisterHandleFunc("Protocal", func() string {
		return string(protocal)
	})

	go register.RegisterHandler.Handle(register)
	go register.MessageHandler.Handle(register.Channel)
}

func (register *Register) RegisterAndServe(serverName, address string, handler http.Handler) {
	register.registerAndServe(serverName, address, "", "", handler)
}

func (register *Register) RegisterAndServeTLS(serverName, address string, certFile, keyFile string, handler http.Handler) {
	register.registerAndServe(serverName, address, certFile, keyFile, handler)
}

func (register *Register) registerAndServe(serverName, address string, certFile, keyFile string, handler http.Handler) {
	serveMux, ok := handler.(omihttp.ServeMux)
	if ok {
		register.registerMux(&serveMux)
	}
	serveMuxRef, ok := handler.(*omihttp.ServeMux)
	if ok {
		register.registerMux(serveMuxRef)
	}
	var err error
	if certFile == "" || keyFile == "" {
		register.register(HTTP, serverName, address)
		err = http.ListenAndServe(register.Port, handler)
	} else {
		register.register(HTTPS, serverName, address)
		err = http.ListenAndServeTLS(register.Port, certFile, keyFile, handler)
	}
	log.Fatalln(err)
}

func (register *Register) SendMessage(serverName, address, command, message string) {
	channel := Prefix + serverName + Namespace_separator + address
	register.omipc.Notify(channel, command+Namespace_separator+message)
}

func (register *Register) registerMux(mux *omihttp.ServeMux) {
	for pattern := range mux.RouterMap {
		register.AddRegisterHandleFunc("["+pattern+"]", func() string {
			if mux.RouterMap[pattern] {
				return "open"
			} else {
				return "closed"
			}
		})
		register.AddMessageHandleFunc("Switch["+pattern+"]", func(message string) {
			state, err := strconv.Atoi(message)
			if err == nil {
				if state == 1 {
					mux.RouterMap[pattern] = true
				}
				if state == 0 {
					mux.RouterMap[pattern] = false
				}
			}
		})
	}
}
