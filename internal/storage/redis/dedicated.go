package redis

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/IzomSoftware/GinWrapper/internal/configuration"
	"github.com/redis/go-redis/v9"
)

type DedicatedRedisStorage struct{}

func (D *DedicatedRedisStorage) GetRedisOpts(config *configuration.RedisConfiguration) (*redis.Options, error) {
	dedicatedConfig := config.DedicatedRedisConfiguration
	redisOpts := &redis.Options{
		Addr:         fmt.Sprintf("%s:%x", dedicatedConfig.Hostname, dedicatedConfig.Port),
		Username:     dedicatedConfig.Username,
		Password:     dedicatedConfig.Password,
		DB:           dedicatedConfig.Database,
		PoolSize:     dedicatedConfig.PoolSize,
		MinIdleConns: dedicatedConfig.MinIdleConnections,
		MaxRetries:   dedicatedConfig.MaxRetries,
		PoolTimeout:  time.Duration(dedicatedConfig.PoolTimeout) * time.Second,
		DialTimeout:  time.Duration(dedicatedConfig.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(dedicatedConfig.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(dedicatedConfig.WriteTimeoutSec) * time.Second,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: dedicatedConfig.SkipTLSVerification,
		},
	}

	if dedicatedConfig.TLSEnabled {
		redisOpts.TLSConfig = &tls.Config{
			InsecureSkipVerify: dedicatedConfig.SkipTLSVerification,
		}
	}
	return redisOpts, nil
}
