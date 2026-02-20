package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/IzomSoftware/GinWrapper/configuration"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

var redisCl *redis.Client
var redisCtx context.Context

/*
 * MiniRedis
 */
func InitEmbeddedRedis() (*miniredis.Miniredis, error) {
	miniRedis, err := miniredis.Run()
	if err != nil {
		return nil, err
	}

	redisCl = redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

	redisCtx = ctx

	return miniRedis, nil
}

func ShutdownEmbeddedRedis(miniRedis *miniredis.Miniredis) {
	miniRedis.Close()
}

/*
 * Redis
 */
func InitRedis(config configuration.RedisConfiguration) error {
	redisOpts := &redis.Options{
		Addr:         fmt.Sprintf("%s:%x", config.Hostname, config.Port),
		Username:     config.Username,
		Password:     config.Password,
		DB:           config.Database,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConnections,
		MaxRetries:   config.MaxRetries,
		PoolTimeout:  time.Duration(config.PoolTimeout) * time.Second,
		DialTimeout:  time.Duration(config.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(config.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.WriteTimeoutSec) * time.Second,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: config.SkipTLSVerification,
		},
	}

	if config.TLSEnabled {
		redisOpts.TLSConfig = &tls.Config{
			InsecureSkipVerify: config.SkipTLSVerification,
		}
	}

	redisCl = redis.NewClient(redisOpts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

	redisCtx = ctx
	
	return nil
}

/*
 * Common Methods
 */
func Init() {
	config := configuration.ConfigHolder.DatabaseConfiguration

	if config.EmbeddedRedisConfiguration.Enabled {
		InitEmbeddedRedis()
	}

	if config.RedisConfiguration.Enabled {
		InitRedis(config.RedisConfiguration)
	}
}

func SetRedisValue(name string, arg any) error {
	return redisCl.Set(redisCtx, name, arg, 0).Err()
}

func GetRedisValue(name string) (any, error) {
	return redisCl.Get(redisCtx, name).Result()
}