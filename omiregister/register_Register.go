package register

import (
	"context"
	"log"
	"os"
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
	OmipcClient     *Omipc
	ctx             context.Context
	RegisterHandler *RegisterHandler
	MessageHandler  *MessageHandler
	StartTime       time.Time
	Port            string
	ReadWriter      omihttp.ReadWriter
	regestered      bool
}

func NewRegister(redisClient *redis.Client) *Register {
	register := &Register{
		RedisClient:     redisClient,
		Weight:          1,
		Info:            map[string]string{},
		Prefix:          Prefix,
		ctx:             context.Background(),
		OmipcClient:     NewOmipc(redisClient),
		RegisterHandler: newRegisterHandler(redisClient),
		MessageHandler:  newMessageHander(redisClient),
		StartTime:       time.Now(),
		ReadWriter:      *omihttp.NewReadWriter(),
	}

	register.AddRegisterHandleFunc("Weight", func() string {
		return strconv.Itoa(register.Weight)
	})
	register.AddRegisterHandleFunc("ProcessId", func() string {
		return strconv.Itoa(os.Getpid())
	})
	register.AddRegisterHandleFunc("Host", func() string {
		host, _ := os.Hostname()
		return host
	})
	register.AddRegisterHandleFunc("ServerType", func() string {
		return "server"
	})
	register.AddRegisterHandleFunc("StartTime", func() string {
		return register.StartTime.Format("2006-01-02 15:04:05")
	})
	register.AddRegisterHandleFunc("RunTime", func() string {
		return time.Since(register.StartTime).String()
	})
	register.AddRegisterHandleFunc("MessageHandlers", func() string {
		handlerNames := []string{}
		for name := range register.MessageHandler.handleFuncs {
			handlerNames = append(handlerNames, name)
		}
		return strings.Join(handlerNames, ", ")
	})

	register.AddMessageHandleFunc(Command_UpdateWeight, func(message string) {
		weight, err := strconv.Atoi(message)
		if err == nil {
			register.Weight = weight
		}
	})

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
	register.Port = ":" + strings.Split(address, ":")[1]

	log.Printf("%s server is registered on redis:%s with %s://%s", register.ServerName, register.RedisClient.Options().Addr, protocal, register.Address)
	register.AddRegisterHandleFunc("Protocal", func() string {
		return string(protocal)
	})

	go register.RegisterHandler.Handle(register)
	go register.MessageHandler.Handle(register.Channel)
}

func (register *Register) RegisterAndServe(serverName, address string, serveHandle func(port string)) {
	register.register(HTTP, serverName, address)
	serveHandle(register.Port)
}

func (register *Register) RegisterAndServeTLS(serverName, address string, serveHandle func(port string)) {
	register.register(HTTPS, serverName, address)
	serveHandle(register.Port)
}

func (register *Register) SendMessage(serverName, address, command, message string) {
	channel := Prefix + serverName + Namespace_separator + address
	register.OmipcClient.Notify(channel, command+Namespace_separator+message)
}
