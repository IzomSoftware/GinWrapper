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

var pathScript = redis.Script(`
	local c = redis.call("INCR", KEYS[1])
	if c == 1 then redis.call("EXPIRE", KEYS[1], ARGV[1]) end
	return c
`)

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
		c.Next()

		status := c.Writer.Status()
		if status != http.StatusNotFound && status != http.StatusForbidden {
			return
		}

		ip := cidr.GetIP(c, config)
		key := "probe:" + ip
		i, err := redis.RunScript(pathScript, []string{key}, int64(window.Seconds())).Int64()
		if err != nil {
			logger.Error("probing counter failed", "ip", ip, "err", err)
			return
		}

		if i >= threshold {
			BanIP(redis, ip, banDuration)
		}
	}
}

func RegisteredPaths(engine *gin.Engine) []string {
	routes := engine.Routes()
	result := make([]string, 0, len(routes))
	for _, r := range routes {
		result = append(result, r.Path)
	}
	return result
}

func IncompleteRouteCheck(engine *gin.Engine, config configuration.ReverseProxyMode) gin.HandlerFunc {
	return func(c *gin.Context) {
		paths := RegisteredPaths(engine)
		parts := strings.Split(strings.TrimPrefix(c.Request.URL.Path, "/"), "/")
		if parts[0] != "" {
			base := "/" + parts[0]
			for _, registered := range paths {
				if strings.HasPrefix(strings.ToLower(registered), strings.ToLower(base)) {
					response.AbortForbidden(c, config)
					return
				}
			}
		}
		response.NoRoute(c)
	}
}
