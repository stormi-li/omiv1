package proxy

import (
	"math/rand/v2"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	register "github.com/stormi-li/omiv1/omiregister"
)

type Router struct {
	Discover        *register.Discover
	addressMap      map[string]map[string]map[string]string
	addressPool     map[string][]string
	mutex           sync.RWMutex
	RefreshInterval time.Duration
}

func NewRouter(redisClient *redis.Client) *Router {
	router := &Router{
		Discover:        register.NewDiscover(redisClient),
		addressMap:      map[string]map[string]map[string]string{},
		addressPool:     map[string][]string{},
		mutex:           sync.RWMutex{},
		RefreshInterval: 10 * time.Second,
	}
	go router.Refresh()
	return router
}

func (router *Router) Update() {
	addrs := router.Discover.GetAll()
	addrPool := map[string][]string{}
	for name, addrs := range addrs {
		for _, addr := range addrs {
			data := router.Discover.GetData(name, addr)
			weight, _ := strconv.Atoi(data["Weight"])
			for i := 0; i < weight; i++ {
				addrPool[name] = append(addrPool[name], addr)
			}
		}
	}
	addrMap := map[string]map[string]map[string]string{}
	for name, addrs := range addrs {
		if addrMap[name] == nil {
			addrMap[name] = map[string]map[string]string{}
		}
		for _, addr := range addrs {
			data := router.Discover.GetData(name, addr)
			addrMap[name][addr] = data
		}
	}
	router.mutex.Lock()
	router.addressMap = addrMap
	router.addressPool = addrPool
	router.mutex.Unlock()
}

func (router *Router) Refresh() {
	time.Sleep(100 * time.Millisecond)
	for {
		router.Update()
		time.Sleep(router.RefreshInterval)
	}
}

func (router *Router) GetAddressMap() map[string]map[string]map[string]string {
	router.Update()
	router.mutex.RLock()
	defer router.mutex.RUnlock()
	return router.addressMap
}

func (router *Router) GetNodeInfo(serverName, address string) map[string]string {
	return router.Discover.GetData(serverName, address)
}

func (router *Router) GetAddress(serverName string) string {
	router.mutex.RLock()
	defer router.mutex.RUnlock()
	if len(router.addressMap[serverName]) == 0 {
		return ""
	}
	return router.addressPool[serverName][rand.IntN(len(router.addressMap[serverName]))]
}

func (router *Router) Has(serverName string) bool {
	router.mutex.RLock()
	defer router.mutex.RUnlock()
	return len(router.addressMap[serverName]) != 0
}
