package infrastructure

import (
	"context"
	"crypto/tls"
	"strings"

	"github.com/iwariku/golang_twitter/utils"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	redisAddress := utils.GetEnvOrDefault("REDIS_ADDRESS", "redis_server:6379")

	// アドレスをベースに、AWS環境（TLS必須）かどうかの判定を行う
	var tlsConfig *tls.Config
	if strings.Contains(redisAddress, "amazonaws.com") {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:      redisAddress,
		Password:  "",
		DB:        0,
		TLSConfig: tlsConfig,
	})

	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		panic("Redisへの接続に失敗しました: " + err.Error())
	}

	return redisClient
}
