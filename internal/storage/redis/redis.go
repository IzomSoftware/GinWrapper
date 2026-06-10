package redis

import (
	"context"

	"github.com/IzomSoftware/GinWrapper/internal/configuration"
	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	dbClient *redis.Client
	context  *context.Context
}

type RedisStorageImplementation interface {
	GetRedisOpts(config *configuration.RedisConfiguration) (*redis.Options, error)
}

func (R *RedisStorage) NewRedisStorage(config *configuration.RedisConfiguration, ctx *context.Context, implementation RedisStorageImplementation) (*RedisStorage, error) {
	opts, err := implementation.GetRedisOpts(config)
	return &RedisStorage{
		dbClient: redis.NewClient(opts),
		context: ctx,
	}, err
}

func (R *RedisStorage) SetValue(name string, arg any) error {
	return R.dbClient.Set(*R.context, name, arg, 0).Err()
}

func (R *RedisStorage) GetValue(name string) (any, error) {
	return R.dbClient.Get(*R.context, name).Result()
}

func (R *RedisStorage) HSetValue(name string, arg map[string]interface{}) error {
	return R.dbClient.HSet(*R.context, name, arg).Err()
}

func (R *RedisStorage) UpdateHashValue(name string, field string, arg any) error {
	return R.dbClient.HSet(*R.context, name, field, arg).Err()
}

func (R *RedisStorage) HGetValue(name string, field string) (any, error) {
	return R.dbClient.HGet(*R.context, name, field).Result()
}

func (R *RedisStorage) HGetAllValue(name string) (any, error) {
	return R.dbClient.HGetAll(*R.context, name).Result()
}