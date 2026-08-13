package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/IzomSoftware/GinWrapper/internal/logger"
	"github.com/IzomSoftware/GinWrapper/internal/storage/redis"
	"github.com/gin-gonic/gin"
)

func BanCheck(redis *redis.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		banned, err := redis.Exists(fmt.Sprintf("ban:%s", ip))

		if err != nil {
			logger.Error("ban check failed for ip", ip, "err", err)
			c.Next()
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
