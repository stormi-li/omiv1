package proxy

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/go-redis/redis/v8"
)

func jsonStrToMap(jsonStr string) map[string]string {
	var dataMap map[string]string
	json.Unmarshal([]byte(jsonStr), &dataMap)
	return dataMap
}

func getKeysByNamespace(redisClient *redis.Client, prefix string) []string {
	var keys []string
	cursor := uint64(0)
	for {
		res, newCursor, err := redisClient.Scan(context.Background(), cursor, prefix+"*", 0).Result()
		if err != nil {
			return nil
		}
		for _, key := range res {
			keyWithoutNamespace := key[len(prefix):]
			keys = append(keys, keyWithoutNamespace[1:])
		}
		cursor = newCursor
		if cursor == 0 {
			break
		}
	}
	return keys
}

func splitMessage(input, delimiter string) (string, string) {
	parts := strings.SplitN(input, delimiter, 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return parts[0], ""
}
