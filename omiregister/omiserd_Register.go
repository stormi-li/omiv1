package register

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omiv1/omiconst"
	"github.com/stormi-li/omiv1/omihttp"
	"github.com/stormi-li/omiv1/omipc"
)

const Command_UpdateWeight = "UpdateWeight"

var RegisterInterval = 2 * time.Second

// Register 是服务注册和消息处理的核心结构
type Register struct {
	RedisClient     *redis.Client     // Redis 客户端实例
	ServerName      string            // 服务名
	Address         string            // 服务地址（包含主机和端口）
	Weight          int               // 服务权重
	Info            map[string]string // 服务的元数据，如权重、主机名等
	Prefix          string            // 命名空间前缀
	Channel         string            // Redis 发布/订阅使用的频道名
	OmipcClient     *omipc.Client     // omipc 客户端，用于异步通信
	ctx             context.Context   // 上下文，用于 Redis 操作
	RegisterHandler *RegisterHandler  // 注册处理器，管理服务注册逻辑
	MessageHandler  *MessageHandler   // 消息处理器，处理接收到的消息
	StartTime       time.Time
	Port            string
	ReadWriter      omihttp.ReadWriter
}

// NewRegister 创建一个新的 Register 实例
// 参数：
// - opts: Redis 连接配置
// - serverName: 服务名称
// - address: 服务地址（格式为 "host:port"）
// - prefix: 命名空间前缀
// 返回值：*Register
func NewRegister(redisClient *redis.Client, serverName, address string) *Register {
	if strings.Contains(serverName, omiconst.Namespace_separator) {
		panic("名字里不能包含字符冒号" + omiconst.Namespace_separator)
	}
	register := &Register{
		RedisClient:     redisClient, // 初始化 Redis 客户端
		ServerName:      serverName,
		Address:         address,
		Weight:          1,
		Info:            map[string]string{}, // 初始化空元数据
		Prefix:          omiconst.Prefix,
		ctx:             context.Background(),                                                  // 默认上下文
		OmipcClient:     omipc.NewClient(redisClient),                                          // 创建 omipc 客户端
		RegisterHandler: newRegisterHandler(redisClient),                                       // 创建服务注册处理器
		MessageHandler:  newMessageHander(redisClient),                                         // 创建消息处理器
		Channel:         omiconst.Prefix + serverName + omiconst.Namespace_separator + address, // 频道名称由前缀、服务名和地址拼接而成
		StartTime:       time.Now(),
		Port:            ":" + strings.Split(address, ":")[1],
		ReadWriter:      *omihttp.NewReadWriter(),
	}

	// 添加默认的注册逻辑处理函数
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

	// 添加消息权重修改回调函数
	register.AddMessageHandleFunc(Command_UpdateWeight, func(message string) {
		weight, err := strconv.Atoi(message)
		if err == nil {
			register.Weight = weight
		}
	})

	return register
}

// AddRegisterHandleFunc 添加额外的注册处理函数
func (register *Register) AddRegisterHandleFunc(key string, handleFunc func() string) {
	register.RegisterHandler.AddHandleFunc(key, handleFunc)
}

// AddMessageHandleFunc 添加额外的消息处理函数
func (register *Register) AddMessageHandleFunc(command string, handleFunc func(message string)) {
	register.MessageHandler.AddHandleFunc(command, handleFunc)
}

// RegisterAndServe 启动服务注册并运行服务
// 参数：
// - weight: 服务权重
// - serverFunc: 服务的启动函数，通常是一个 HTTP 或 TCP 服务器
type Protocal string

var HTTP Protocal = "http"
var HTTPS Protocal = "https"

func (register *Register) register(protocal Protocal) {
	log.Printf("%s is registered and starting at %s://%s", register.ServerName, protocal, register.Address)
	register.AddRegisterHandleFunc("Protocal", func() string {
		return string(protocal)
	})
	// 启动服务注册逻辑和消息处理逻辑
	go register.RegisterHandler.Handle(register)
	go register.MessageHandler.Handle(register.Channel)
}

func (register *Register) Register() {
	register.register(HTTP)
}

func (register *Register) RegisterTLS(protocal Protocal) {
	register.register(HTTPS)
}

// SendMessage 发送消息到指定频道
// 参数：
// - command: 消息命令
// - message: 消息内容
func (register *Register) SendMessage(command string, message string) {
	register.OmipcClient.Notify(register.Channel, command+omiconst.Namespace_separator+message)
}
