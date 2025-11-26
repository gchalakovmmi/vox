package redis

import (
	"github.com/redis/go-redis/v9"
	"vox/internal/config"
)

func Init(cfg *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:		cfg.GetRedisAddress(),
		Password:	cfg.GetRedisPassword(),
		DB:		0,
	})
}
