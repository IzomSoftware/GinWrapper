package redis

import (
	"context"
	"time"

	"github.com/IzomSoftware/GinWrapper/configuration"
	"github.com/redis/go-redis/v9"
)

type StorageImplementation interface {
	GetRedisOpts(config *configuration.RedisConfiguration) (*redis.Options, error)
	Close(client *redis.Client)
}

type Storage struct {
	client *redis.Client
	impl   StorageImplementation
}

func New(cfg *configuration.RedisConfiguration, impl StorageImplementation) (*Storage, error) {
	opts, err := impl.GetRedisOpts(cfg)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return &Storage{
		client: redis.NewClient(opts),
		impl:   impl,
	}, nil
}

func (S *Storage) Close() {
	S.impl.Close(S.client)
}

func (S *Storage) Set(key string, value any, expiration time.Duration) error {
	return S.client.Set(context.Background(), key, value, expiration).Err()
}

func (S *Storage) Get(key string) (string, error) {
	return S.client.Get(context.Background(), key).Result()
}

func (S *Storage) HSet(key string, values map[string]any) error {
	return S.client.HSet(context.Background(), key, values).Err()
}

func (S *Storage) HUpdate(key string, field string, value any) error {
	return S.client.HSet(context.Background(), key, field, value).Err()
}

func (S *Storage) HGet(key string, field string) (string, error) {
	return S.client.HGet(context.Background(), key, field).Result()
}

func (S *Storage) HGetAll(key string) (map[string]string, error) {
	return S.client.HGetAll(context.Background(), key).Result()
}

func (S *Storage) Exists(key string) (bool, error) {
	count, err := S.client.Exists(context.Background(), key).Result()
	return count > 0, err
}

func (S *Storage) Del(keys ...string) error {
	return S.client.Del(context.Background(), keys...).Err()
}

func (S *Storage) Incr(key string) (int64, error) {
	return S.client.Incr(context.Background(), key).Result()
}

func (S *Storage) Expire(key string, time time.Duration) error {
	return S.client.Expire(context.Background(), key, time).Err()
}

func (S *Storage) RunScript(script *redis.Script, keys []string, args ...any) *redis.Cmd {
	return script.Run(context.Background(), S.client, keys, args)
}

func Script(script string) *redis.Script {
	return redis.NewScript(script)
}

func (S *Storage) Ping(ctx context.Context) error {
	return S.client.Ping(ctx).Err()
}
