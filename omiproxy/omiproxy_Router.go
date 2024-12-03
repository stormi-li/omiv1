package proxy

import (
	"math/rand"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	register "github.com/stormi-li/omiv1/omiregister"
)

type Router struct {
	Discover        *register.Discover
	addressMap      map[string]map[string]map[string]string
	addressPool     map[string][]string
	addressIndex    map[string]int
	mutex           sync.RWMutex
	RefreshInterval time.Duration
}

func NewRouter(redisClient *redis.Client) *Router {
	router := &Router{
		Discover:        register.NewDiscover(redisClient),
		addressMap:      map[string]map[string]map[string]string{},
		addressPool:     map[string][]string{},
		addressIndex:    map[string]int{},
		mutex:           sync.RWMutex{},
		RefreshInterval: 2 * time.Second,
	}
	router.Update()
	go router.Refresh()
	return router
}

func (router *Router) Update() {
	addrs := router.Discover.GetAll()
	addrPool := map[string][]string{}
	for name := range addrs {
		addresses := router.Discover.GetByWeightAndFilter(name)
		addrPool[name] = append(router.addressPool[name], addresses...)
	}
	addrMap := map[string]map[string]map[string]string{}
	for name, addrs := range addrs {
		if addrMap[name] == nil {
			addrMap[name] = map[string]map[string]string{}
		}
		for _, addr := range addrs {
			data := router.Discover.GetInfo(name, addr)
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
	return router.Discover.GetInfo(serverName, address)
}

func (router *Router) GetAddress(serverName string) string {
	router.mutex.RLock()
	defer router.mutex.RUnlock()
	if len(router.addressMap[serverName]) == 0 || len(router.addressPool[serverName]) == 0 {
		return ""
	}
	address := router.addressPool[serverName][rand.Intn(len(router.addressPool[serverName]))]
	return address
}

func (router *Router) Has(serverName string) bool {
	router.mutex.RLock()
	defer router.mutex.RUnlock()
	return len(router.addressPool[serverName]) != 0
}
