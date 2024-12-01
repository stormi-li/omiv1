package register

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

// Omipc 是 Redis 客户端的包装，用于封装与 Redis 的交互逻辑
type Omipc struct {
	redisClient *redis.Client   // Redis 客户端实例
	ctx         context.Context // 用于 Redis 操作的上下文
}

func NewOmipc(redisClient *redis.Client) *Omipc {
	return &Omipc{
		redisClient: redisClient,
		ctx:         context.Background(),
	}
}

// Notify 用于向指定频道发送消息
func (c *Omipc) Notify(channel, msg string) {
	c.redisClient.Publish(c.ctx, channel, msg)
}

// Listen 订阅 Redis 频道并处理接收到的消息
// 参数说明：
// - channel: 订阅的 Redis 频道
// - timeout: 超时时间，为 0 表示无超时
// - handFuncs: 可选的处理函数，用于处理收到的消息
// 返回值：如果收到消息并未超时，返回消息的内容；超时时返回空字符串。
func (c *Omipc) Listen(channel string, timeout time.Duration, handFuncs ...func(message string)) string {
	// 如果没有超时设置并且未提供处理函数，则 panic，避免死循环或逻辑错误
	if timeout == 0 && len(handFuncs) == 0 {
		panic("no handFunc provided")
	}

	// 订阅指定的频道
	sub := c.redisClient.Subscribe(c.ctx, channel)
	defer sub.Close() // 确保在函数退出时释放订阅资源

	// 获取 Redis 消息通道
	msgChan := sub.Channel()

	// 初始化超时通道
	var tickerC <-chan time.Time
	if timeout != 0 {
		ticker := time.NewTicker(timeout)
		defer ticker.Stop() // 确保定时器停止，避免资源泄漏
		tickerC = ticker.C
	}

	// 循环监听消息通道和超时通道
	for {
		select {
		case msg := <-msgChan:
			// 如果没有设置超时，则异步调用所有处理函数
			if timeout == 0 {
				go func(payload string) {
					// 使用 recover 捕获处理函数中的 panic，避免程序崩溃
					defer func() {
						if r := recover(); r != nil {
							log.Printf("Recovered from handler panic: %v", r)
						}
					}()
					// 遍历并调用所有处理函数
					for _, handFunc := range handFuncs {
						handFunc(payload)
					}
				}(msg.Payload)
			}
			// 如果设置了超时，则直接返回消息内容
			if timeout != 0 {
				return msg.Payload
			}
		case <-tickerC:
			// 超时通道触发时返回空字符串，表示超时
			return ""
		}
	}
}
