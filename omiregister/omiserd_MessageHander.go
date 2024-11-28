package register

import (
	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omiv1/omiconst"
	"github.com/stormi-li/omiv1/omipc"
)

type MessageHandler struct {
	ompcClient  *omipc.Client
	handleFuncs map[string]func(message string)
}

func newMessageHander(redisClient *redis.Client) *MessageHandler {
	return &MessageHandler{
		ompcClient:  omipc.NewClient(redisClient),
		handleFuncs: map[string]func(message string){},
	}
}

func (messageHandler *MessageHandler) AddHandleFunc(command string, handleFunc func(message string)) {
	messageHandler.handleFuncs[command] = handleFunc
}

func (messageHandler *MessageHandler) Handle(channel string) {
	messageHandler.ompcClient.Listen(channel, 0, func(message string) {
		command, message := splitMessage(message, omiconst.Namespace_separator)
		if handleFunc, ok := messageHandler.handleFuncs[command]; ok {
			handleFunc(message)
		}
	})
}
