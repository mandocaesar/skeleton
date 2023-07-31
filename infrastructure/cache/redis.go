package cache

import (
	"fmt"

	redis "github.com/go-redis/redis/v9"
)

type RedisConfig struct {
	Host string
	Port string
	DB   int
}

func CreateRedisConnection(config RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%v:%v", config.Host, config.Port),
		DB:   config.DB,
	})
}
