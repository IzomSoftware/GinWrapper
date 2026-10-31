package response

import (
	"net/http"
	"strings"
	"time"

	"github.com/IzomSoftware/GinWrapper/cidr"
	"github.com/IzomSoftware/GinWrapper/configuration"
	"github.com/IzomSoftware/GinWrapper/logger"
	"github.com/IzomSoftware/GinWrapper/middleware"
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
func RegisteredPaths(engine *gin.Engine) []string {
	routes := engine.Routes()
	result := make([]string, 0, len(routes))
	for _, r := range routes {
		result = append(result, r.Path)
	}
	return result
}
func NoRoute(c *gin.Context) {
	c.String(http.StatusNotFound, "404 Not Found")
}
func ProtectedNoRoute(engine *gin.Engine, redis *redis.Storage, threshold int64, window, banDuration time.Duration, config configuration.ReverseProxyMode) gin.HandlerFunc {
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
			middleware.BanIP(redis, ip, banDuration)
		}
		parts := strings.Split(strings.TrimPrefix(c.Request.URL.Path, "/"), "/")
		if parts[0] != "" {
			base := "/" + parts[0]
			for _, registered := range paths {
				if strings.HasPrefix(strings.ToLower(registered), strings.ToLower(base)) {
					middleware.BanIP(redis, ip, banDuration)
					return
				}
			}
		}
		NoRoute(c)
	}
}
