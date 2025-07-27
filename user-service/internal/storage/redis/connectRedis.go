package redis

import (
	"log"
	"user-service/config"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient(cfg config.Config) *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Host + ":" + cfg.Redis.Port,
	})

	log.Println("Connected to Redis")
	return client
}
