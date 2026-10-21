package middleware

import (
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/IzomSoftware/GinWrapper/cidr"
	"github.com/IzomSoftware/GinWrapper/configuration"
	"github.com/IzomSoftware/GinWrapper/logger"
	"github.com/IzomSoftware/GinWrapper/response"
	"github.com/IzomSoftware/GinWrapper/storage/redis"
	"github.com/gin-gonic/gin"
)

func CleanCheck(config configuration.ReverseProxyMode) gin.HandlerFunc {
	return func(c *gin.Context) {
		uriPath := c.Request.URL.Path
		if path.Clean(uriPath) != uriPath {
			response.AbortForbidden(c, config)
			return
		}
		c.Next()
	}
}

func ProbingCheck(redis *redis.Storage, threshold int64, window, banDuration time.Duration, config configuration.ReverseProxyMode) gin.HandlerFunc {
    return func(c *gin.Context) {
        if err := redis.Ping(c.Request.Context()); err != nil {
            response.AbortInternalError(c, config)
            return
        }

        c.Next()

        status := c.Writer.Status()
        if status != http.StatusNotFound && status != http.StatusForbidden {
            return
        }

        ip := cidr.GetIP(c, config)
        key := "probe:" + ip
        i, err := redis.Incr(key)
        if err != nil {
            logger.Error("probing counter failed", "ip", ip, "err", err)
            return
        }

        if i == 1 {
            redis.Expire(key, window)
        } else if i >= threshold {
            BanIP(redis, ip, banDuration)
        }
    }
}

func IncompleteRouteCheck(registeredPaths []string, config configuration.ReverseProxyMode) gin.HandlerFunc {
	return func(c *gin.Context) {
		parts := strings.Split(strings.TrimPrefix(c.Request.URL.Path, "/"), "/")
		if parts[0] != "" {
			base := "/" + parts[0]
			for _, registered := range registeredPaths {
				if strings.HasPrefix(registered, base) {
					response.AbortForbidden(c, config)
					return
				}
			}
		}
		response.NoRoute(c)
	}
}
