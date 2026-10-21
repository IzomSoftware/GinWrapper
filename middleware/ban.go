package middleware

import (
	"fmt"
	"time"

	"github.com/IzomSoftware/GinWrapper/cidr"
	"github.com/IzomSoftware/GinWrapper/configuration"
	"github.com/IzomSoftware/GinWrapper/logger"
	"github.com/IzomSoftware/GinWrapper/response"
	"github.com/IzomSoftware/GinWrapper/storage/redis"
	"github.com/gin-gonic/gin"
)

func BanCheck(redis *redis.Storage, config configuration.ReverseProxyMode) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := cidr.GetIP(c, config)
		banned, err := redis.Exists(fmt.Sprintf("ban:%s", ip))
		if err != nil {
			logger.Error("ban check failed", "ip", ip, "err", err)
			response.AbortInternalError(c, config)
			return
		}
		if banned {
			response.AbortForbidden(c, config)
			return
		}
		c.Next()
	}
}

func BanIP(redis *redis.Storage, ip string, duration time.Duration) error {
	return redis.Set(fmt.Sprintf("ban:%s", ip), "1", duration)
}

func Ban(c *gin.Context, redis *redis.Storage, duration time.Duration, config configuration.ReverseProxyMode) {
	if duration != 0 && redis != nil {
		BanIP(redis, cidr.GetIP(c, config), duration)
	}
}
