package tools

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
}

func (r *Redis) NewClient() error {

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

func (r *Redis) SetSKUStock(skuId string, stock int64) error {
	err := r.Client.Set(context.Background(), "sku:"+skuId+":stock", stock, 0).Err()
	if err != nil {
		return fmt.Errorf("redis set sku stock failed: %#+v", err)
	}
	log.Println("库存已同步到 Redis:", skuId, stock)
	return nil
}

func (r *Redis) DeductStock(skuId string, quantity int64) (bool, error) {
	// Lua 脚本，保证扣减库存的原子性
	luaScript := `
		local stock = tonumber(redis.call("GET", KEYS[1]))
		if stock >= tonumber(ARGV[1]) then
			redis.call("DECRBY", KEYS[1], ARGV[1])
			return 1
		else
			return 0
		end
		`
	// 执行 Lua 脚本
	result, err := r.Client.Eval(context.Background(), luaScript, []string{"sku:" + skuId + ":stock"}, quantity).Result()
	if err != nil {
		return false, err
	}

	// 如果返回 1，表示库存扣减成功
	if result.(int64) == 1 {
		log.Println("扣减库存成功:", skuId, "数量:", quantity)
		return true, nil
	}

	// 否则表示库存不足
	log.Println("库存不足:", skuId)
	return false, nil
}

func (r *Redis) DeleteSkuStock(skuId string) error {
	err := r.Client.Del(context.Background(), "sku:"+skuId+":stock").Err()
	if err != nil {
		return fmt.Errorf("redis delete sku stock failed: %#+v", err)
	}
	log.Println("删除 Redis 中的库存:", skuId)
	return nil
}
