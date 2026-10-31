package response

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/IzomSoftware/GinWrapper/cidr"
	"github.com/IzomSoftware/GinWrapper/configuration"
	"github.com/IzomSoftware/GinWrapper/logger"
	"github.com/IzomSoftware/GinWrapper/storage/redis"
	"github.com/gin-gonic/gin"
)

var pathScript = redis.Script(`
	local c = redis.call("INCR", KEYS[1])
	if c == 1 then redis.call("EXPIRE", KEYS[1], ARGV[1]) end
	return c
`)

func Abort(c *gin.Context, status int, config configuration.ReverseProxyMode) {
	logger.Info("aborted", "ip", cidr.GetIP(c, config), "status", status)
	c.AbortWithStatus(status)
}

func AbortForbidden(c *gin.Context, config configuration.ReverseProxyMode) {
	Abort(c, http.StatusForbidden, config)
}

func AbortUnauthorized(c *gin.Context, config configuration.ReverseProxyMode) {
	Abort(c, http.StatusUnauthorized, config)
}

func AbortInternalError(c *gin.Context, config configuration.ReverseProxyMode) {
	Abort(c, http.StatusInternalServerError, config)
}

func BanIP(redis *redis.Storage, ip string, duration time.Duration) error {
	return redis.Set(fmt.Sprintf("ban:%s", ip), "1", duration)
}

func Ban(c *gin.Context, redis *redis.Storage, duration time.Duration, config configuration.ReverseProxyMode) {
	if duration != 0 && redis != nil {
		BanIP(redis, cidr.GetIP(c, config), duration)
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

// Renamed to invalidresponse because of potential usage on NoMethod
func InvalidResponse(c *gin.Context) {
	c.String(http.StatusNotFound, "404 Not Found")
}

func ProtectedInvalidResponse(engine *gin.Engine, redis *redis.Storage, threshold int64, window time.Duration, banDuration time.Duration, config configuration.ReverseProxyMode) gin.HandlerFunc {
	paths := RegisteredPaths(engine)
	return func(c *gin.Context) {
		ip := cidr.GetIP(c, config)
		key := "probe:" + ip
		i, err := redis.RunScript(pathScript, []string{key}, int64(window.Seconds())).Int64()
		if err != nil {
			logger.Error("probing counter failed", "ip", ip, "err", err)
			return
		}
		if i >= threshold {
			AbortForbidden(c, config)
			BanIP(redis, ip, banDuration)
			return
		}
		parts := strings.Split(strings.TrimPrefix(c.Request.URL.Path, "/"), "/")
		if parts[0] != "" {
			base := "/" + parts[0]
			for _, registered := range paths {
				if strings.HasPrefix(strings.ToLower(registered), strings.ToLower(base)) {
					AbortForbidden(c, config)
					BanIP(redis, ip, banDuration)
					return
				}
			}
		}
		InvalidResponse(c)
	}
}
