package redis

import (
	"github.com/IzomSoftware/GinWrapper/configuration"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

type EmbeddedRedisStorage struct {
	MiniRedis *miniredis.Miniredis
}

func (E *EmbeddedRedisStorage) GetRedisOpts(config *configuration.RedisConfiguration) (*redis.Options, error) {
	miniRedis, err := miniredis.Run()
	if err != nil {
		return nil, err
	}
	E.MiniRedis = miniRedis
	return &redis.Options{
		Addr: miniRedis.Addr(),
	}, nil
}
