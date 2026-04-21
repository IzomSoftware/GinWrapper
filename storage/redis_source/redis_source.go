package redis_source

import (
	"context"
	"crypto/tls"
	"fmt"
	"strconv"
	"time"

	"github.com/IzomSoftware/GinWrapper/configuration"
	"github.com/IzomSoftware/GinWrapper/logger"
	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var (
	redisCl            *redis.Client
	redisCtx           context.Context
	RedisNotConfigured = fmt.Errorf("No redis implemnetation configured")
)

type User struct {
	Rate     int64
	LastRate int64
}

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
func Init() error {
	config := configuration.ConfigHolder.DatabaseConfiguration

	if config.EmbeddedRedisConfiguration.Enabled {
		_, err := InitEmbeddedRedis()
		return err
	}

	if config.RedisConfiguration.Enabled {
		return InitRedis(config.RedisConfiguration)
	}

	return RedisNotConfigured
}

func SetValue(c context.Context, name string, arg any) error {
	if c == nil {
		c = redisCtx
	}
	return redisCl.Set(c, name, arg, 0).Err()
}

func GetValue(c context.Context, name string) (any, error) {
	if c == nil {
		c = redisCtx
	}
	return redisCl.Get(c, name).Result()
}

func HSetValue(c context.Context, name string, arg any) error {
	if c == nil {
		c = redisCtx
	}
	return redisCl.HSet(c, name, arg).Err()
}

func UpdateHashValue(c context.Context, name string, field string, arg any) error {
	if c == nil {
		c = redisCtx
	}
	return redisCl.HSet(c, name, field, arg).Err()
}

func HGetValue(c context.Context, name string, field string) (any, error) {
	if c == nil {
		c = redisCtx
	}
	return redisCl.HGet(c, name, field).Result()
}

func HGetAllValue(c context.Context, name string) (any, error) {
	if c == nil {
		c = redisCtx
	}

	return redisCl.HGetAll(c, name).Result()
}

func GetRateLimit(c context.Context, ip string) (int64, error) {
	rate, err := HGetValue(c, ip, "Rate")

	if err != nil {
		if err == redis.Nil {
			user := User{
				Rate:     1,
				LastRate: time.Now().UnixMilli(),
			}
			err = HSetValue(c, ip, user)
			rate = user.Rate
		}

		if err != nil {
			return 0, err
		}
	}

	rateStr, ok := rate.(string)
	val, err := strconv.ParseInt(rateStr, 10, 64)

	if !ok {
		if err != nil {
			logger.LogError(fmt.Sprintf("%s", err))
		}

		return 0, err
	}

	return val, nil
}

func GetLastRateLimit(c *gin.Context, ip string) (int64, error) {
	lastRate, err := HGetValue(c, ip, "LastRate")

	if err != nil {
		if err == redis.Nil {
			user := User{
				Rate:     1,
				LastRate: time.Now().UnixMilli(),
			}
			err = HSetValue(c, ip, user)
			lastRate = user.LastRate
		}

		if err != nil {
			return 0, err
		}
	}

	lastRateStr, ok := lastRate.(string)
	val, err := strconv.ParseInt(lastRateStr, 10, 64)

	if !ok {
		if err != nil {
			logger.LogError(fmt.Sprintf("%s", err))
		}

		return 0, err
	}

	return val, nil
}

func IncrementRateLimit(c context.Context, ip string) error {
	rate, err := GetRateLimit(c, ip)
	if err != nil {
		return err
	}

	newValue := rate + 1

	err = UpdateHashValue(c, ip, "LastRate", time.Now().UnixMilli())
	if err != nil {
		return err
	}

	return UpdateHashValue(c, ip, "Rate", newValue)
}
