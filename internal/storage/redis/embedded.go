package redis

import (
	"github.com/IzomSoftware/GinWrapper/internal/configuration"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

type EmbeddedRedisStorage struct{}

func (E *EmbeddedRedisStorage) GetRedisOpts(config *configuration.RedisConfiguration) (*redis.Options, error) {
	miniRedis, err := miniredis.Run()
	return &redis.Options{
		Addr:  miniRedis.Addr(),
	}, err
}