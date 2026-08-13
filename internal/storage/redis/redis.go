package redis

import (
	"context"
	"time"

	"github.com/IzomSoftware/GinWrapper/internal/configuration"
	"github.com/redis/go-redis/v9"
)

type StorageImplementation interface {
	GetRedisOpts(config *configuration.RedisConfiguration) (*redis.Options, error)
}

type Storage struct {
	client *redis.Client
	ctx    context.Context
}

func New(cfg *configuration.RedisConfiguration, ctx context.Context, impl StorageImplementation) (*Storage, error) {
	opts, err := impl.GetRedisOpts(cfg)
	if err != nil {
		return nil, err
	}
	return &Storage{
		client: redis.NewClient(opts),
		ctx:    ctx,
	}, nil
}

func (S *Storage) Close() error {
	return S.client.Close()
}

func (S *Storage) Set(key string, value any, expiration time.Duration) error {
	return S.client.Set(S.ctx, key, value, expiration).Err()
}

func (S *Storage) Get(key string) (string, error) {
	return S.client.Get(S.ctx, key).Result()
}

func (S *Storage) HSet(key string, values map[string]any) error {
	return S.client.HSet(S.ctx, key, values).Err()
}

func (S *Storage) HUpdate(key string, field string, value any) error {
	return S.client.HSet(S.ctx, key, field, value).Err()
}

func (S *Storage) HGet(key string, field string) (string, error) {
	return S.client.HGet(S.ctx, key, field).Result()
}

func (S *Storage) HGetAll(key string) (map[string]string, error) {
	return S.client.HGetAll(S.ctx, key).Result()
}

func (S *Storage) Exists(key string) (bool, error) {
	count, err := S.client.Exists(S.ctx, key).Result()
	return count > 0, err
}

func (S *Storage) Del(keys ...string) error {
	return S.client.Del(S.ctx, keys...).Err()
}

func (S *Storage) Incr(key string) (int64, error) {
	return S.client.Incr(S.ctx, key).Result()
}

func (S *Storage) Expire(key string, time time.Duration) error {
	return S.client.Expire(S.ctx, key, time).Err()
}

func Script(script string) *redis.Script {
	return redis.NewScript(script)
}

func (S *Storage) RunScript(script *redis.Script, keys []string, args ...any) *redis.Cmd {
	return script.Run(S.ctx, S.client, keys, args)
}
