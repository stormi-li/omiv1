package register

import (
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omiv1/omiconst"
)

type RegisterHandler struct {
	ompcClient  *Omipc
	handleFuncs map[string]func() string
}

func newRegisterHandler(redisClient *redis.Client) *RegisterHandler {
	return &RegisterHandler{
		ompcClient:  NewOmipc(redisClient),
		handleFuncs: map[string]func() string{},
	}
}

func (registerHandler *RegisterHandler) AddHandleFunc(key string, handleFunc func() string) {
	registerHandler.handleFuncs[key] = handleFunc
}

func (registerHandler *RegisterHandler) Handle(register *Register) {
	for {
		for key, handleFunc := range registerHandler.handleFuncs {
			register.Info[key] = handleFunc()
		}
		jsonStrData := mapToJsonStr(register.Info)
		key := register.Prefix + register.ServerName + omiconst.Namespace_separator + register.Address
		register.RedisClient.Set(register.ctx, key, jsonStrData, RegisterInterval)
		time.Sleep(RegisterInterval / 2)
	}
}
