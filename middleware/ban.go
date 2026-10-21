package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/IzomSoftware/GinWrapper/cidr"
	"github.com/IzomSoftware/GinWrapper/configuration"
	"github.com/IzomSoftware/GinWrapper/logger"
	"github.com/IzomSoftware/GinWrapper/storage/redis"
	"github.com/gin-gonic/gin"
)

func BanCheck(redis *redis.Storage, config configuration.ReverseProxyMode) gin.HandlerFunc {
    return func(c *gin.Context) {
        ip := cidr.GetIP(c, config)
        banned, err := redis.Exists(fmt.Sprintf("ban:%s", ip))
        if err != nil {
            logger.Error("ban check failed", "ip", ip, "err", err)
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

func BanIP(redis *redis.Storage, ip string, duration time.Duration) error {
    return redis.Set(fmt.Sprintf("ban:%s", ip), "1", duration)
}