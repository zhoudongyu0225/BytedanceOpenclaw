package store

import (
	"admin-backend/config"
	"context"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func InitRedis(cfg config.RedisConfig) error {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	// 测试连接
	if err := client.Ping(context.TODO()).Err(); err != nil {
		return err
	}
	RedisClient = client
	return nil
}

func CloseRedis() {
	if RedisClient != nil {
		RedisClient.Close()
	}
}
