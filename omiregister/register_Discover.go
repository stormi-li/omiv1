package register

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type Discover struct {
	RedisClient *redis.Client
	Prefix      string
	ctx         context.Context
	Filter      map[string]map[string]func(value string) bool
}

func NewDiscover(redisClient *redis.Client) *Discover {
	return &Discover{
		RedisClient: redisClient,
		Prefix:      Prefix,
		ctx:         context.Background(),
		Filter:      map[string]map[string]func(value string) bool{},
	}
}

func (discover *Discover) Close() {
	discover.RedisClient.Close()
}

func (discover *Discover) Get(serverName string) []string {
	return getKeysByNamespace(discover.RedisClient, discover.Prefix+serverName)
}
func (discover *Discover) AddFilter(serverName, key string, handle func(value string) bool) {
	if discover.Filter[serverName] == nil {
		discover.Filter[serverName] = map[string]func(value string) bool{}
	}
	discover.Filter[serverName][key] = handle
}

func (discover *Discover) GetByWeightAndFilter(serverName string) []string {
	addresses := discover.Get(serverName)
	var addressPool []string
	for _, address := range addresses {
		info := discover.GetInfo(serverName, address)

		skip := false
		for key, val := range info {
			if discover.Filter[serverName][key] != nil && discover.Filter[serverName][key](val) {
				skip = true
				break
			}
		}

		if skip {
			continue
		}

		weight, err := strconv.Atoi(info["Weight"])
		if err != nil {
			continue
		}
		for i := 0; i < weight; i++ {
			addressPool = append(addressPool, address)
		}
	}

	return addressPool
}

func (discover *Discover) GetInfo(serverName string, address string) map[string]string {
	key := discover.Prefix + serverName + Namespace_separator + address
	dataStr, err := discover.RedisClient.Get(discover.ctx, key).Result()
	if err != nil {
		return map[string]string{}
	}

	data := jsonStrToMap(dataStr)
	return data
}

func (discover *Discover) IsAlive(serverName string, address string) bool {
	data := discover.GetInfo(serverName, address)
	if len(data) == 0 {
		return false
	}
	if data["weight"] == "0" {
		return false
	}
	return true
}

func (discover *Discover) GetAll() map[string][]string {
	keys := getKeysByNamespace(discover.RedisClient, discover.Prefix[:len(discover.Prefix)-1])
	result := map[string][]string{}

	for _, key := range keys {
		name, address := splitMessage(key, Namespace_separator)

		if _, exists := result[name]; !exists {
			result[name] = []string{}
		}

		result[name] = append(result[name], address)
	}

	return result
}
