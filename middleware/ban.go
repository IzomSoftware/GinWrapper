package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/IzomSoftware/GinWrapper/cidr"
	"github.com/IzomSoftware/GinWrapper/configuration"
	"github.com/IzomSoftware/GinWrapper/logger"
	"github.com/IzomSoftware/GinWrapper/response"
	"github.com/IzomSoftware/GinWrapper/storage/redis"
	"github.com/gin-gonic/gin"
)

func BanOnFail() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.IsAborted() {
			
		}
	}
}

func BanCheck(redis *redis.Storage, config configuration.ReverseProxyMode) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := cidr.GetIP(c, config)
		banned, err := redis.Exists(fmt.Sprintf("ban:%s", ip))

		if err != nil {
			logger.Error("ban check failed for ip", ip, "err", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if banned {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}

func BanIP(redis *redis.Storage, ip string, time time.Duration) error {
	return redis.Set(fmt.Sprintf("ban:%s", ip), "1", time)
}

func Punish(c *gin.Context, protectionConfig configuration.Protections, reverseProxyConfig configuration.ReverseProxyMode, redis *redis.Storage, ip string, time time.Duration) error {
	if protectionConfig.Ban {
		return BanIP(redis, ip, time)
	} else {
		response.AbortForbidden(c, reverseProxyConfig)
	}
	return nil

}