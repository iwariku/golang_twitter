package infrastructure

import (
	"context"
	"github.com/iwariku/golang_twitter/utils"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	redisAddress := utils.GetEnvOrDefault("REDIS_ADDRESS", "redis_server:6379")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: "",
		DB:       0,
	})

	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		panic("Redisへの接続に失敗しました: " + err.Error())
	}

	return redisClient
}
