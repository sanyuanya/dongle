package tools

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
}

func (r Redis) NewClient() error {

	// 创建 Redis 客户端
	r.Client = redis.NewClient(&redis.Options{
		Addr:     "218.11.1.36:6379", // Redis 服务器地址
		DB:       0,                  // 使用的数据库
		Password: "iKf7wwxtx3Hgu2LhVt3enL9Z3EGjsyVD0tKFxXXMTMd5Qeqv2MjktRkLP8XKwBjG7qynrWdiXMYsrDdoPFJr3qJqJCa0WC",
	})

	if r.Client == nil {
		return fmt.Errorf("redis client is nil")
	}

	_, err := r.Client.Ping(context.Background()).Result()
	if err != nil {
		return fmt.Errorf("redis ping failed: %#+v", err)
	}

	return nil
}
