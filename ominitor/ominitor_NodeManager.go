package monitor

import (
	"github.com/go-redis/redis/v8"
	proxy "github.com/stormi-li/omiv1/omiproxy"
	register "github.com/stormi-li/omiv1/omiregister"
)

type NodeManager struct {
	Router      *proxy.Router
	RedisClient *redis.Client
}

func NewManager(router *proxy.Router) *NodeManager {
	return &NodeManager{
		Router:      router,
		RedisClient: router.Discover.RedisClient,
	}
}

func (nodeManager *NodeManager) GetNodes() map[string]map[string]map[string]string {
	return nodeManager.Router.GetAddressMap()
}

func (nodeManager *NodeManager) GetNodeInfo(name, address string) map[string]string {
	return nodeManager.Router.GetNodeInfo(name, address)
}

func (nodeManager *NodeManager) SendCommand(name, address, command, message string) {
	register := register.NewRegister(nodeManager.RedisClient)
	if register == nil {
		return
	}
	register.SendMessage(name, address, command, message)
}
