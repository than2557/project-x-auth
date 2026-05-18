package database

import (
	"auth-service/internal/config"
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func ConnectRedis(
	cfg *config.Config,
) *redis.Client {

	addr := fmt.Sprintf(
		"%s:%s",
		cfg.RedisHost,
		cfg.RedisPort,
	)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisPassword,
		DB:       0,
	})

	_, err := client.Ping(Ctx).Result()

	if err != nil {
		log.Fatal("Redis connection failed:", err)
	}

	log.Println("Redis connected")

	return client
}
