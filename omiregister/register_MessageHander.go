package register

import (
	"github.com/go-redis/redis/v8"
)

type MessageHandler struct {
	ompcClient  *Omipc
	handleFuncs map[string]func(message string)
}

func newMessageHander(redisClient *redis.Client) *MessageHandler {
	return &MessageHandler{
		ompcClient:  NewOmipc(redisClient),
		handleFuncs: map[string]func(message string){},
	}
}

func (messageHandler *MessageHandler) AddHandleFunc(command string, handleFunc func(message string)) {
	messageHandler.handleFuncs[command] = handleFunc
}

func (messageHandler *MessageHandler) Handle(channel string) {
	messageHandler.ompcClient.Listen(channel, 0, func(message string) {
		command, message := splitMessage(message, Namespace_separator)
		if handleFunc, ok := messageHandler.handleFuncs[command]; ok {
			handleFunc(message)
		}
	})
}
