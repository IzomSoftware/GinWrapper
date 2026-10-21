package middleware

import (
	"net/http"
	"time"

	"github.com/IzomSoftware/GinWrapper/cidr"
	"github.com/IzomSoftware/GinWrapper/configuration"
	"github.com/IzomSoftware/GinWrapper/logger"
	"github.com/IzomSoftware/GinWrapper/response"
	"github.com/IzomSoftware/GinWrapper/storage/redis"
	"github.com/gin-gonic/gin"
)

var script = redis.Script(`
    local key = KEYS[1]
    local rate = tonumber(ARGV[1])
    local window = tonumber(ARGV[2])
    local now = tonumber(ARGV[3])

    local data = redis.call("HMGET", key, "rate", "last")
    local currentRate = tonumber(data[1])
    local lastTime = tonumber(data[2])

    if lastTime == nil or (now - lastTime) > window then
        currentRate = 1
        lastTime = now
    else
        currentRate = currentRate + 1
    end

    redis.call("HSET", key, "rate", currentRate, "last", lastTime)
    redis.call("EXPIRE", key, math.ceil(window / 1000) + 1)

    if currentRate > rate then
        return 1
    end

    return 0
`)

func RateLimitCheck(redis *redis.Storage, rateLimitConfig configuration.RateLimitProtection, banDuration time.Duration, reverseProxyConfig configuration.ReverseProxyMode) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := cidr.GetIP(c, reverseProxyConfig)
		now := time.Now().UnixMilli()

		result, err := redis.RunScript(script, []string{ip}, rateLimitConfig.Rate, rateLimitConfig.Window*1000, now).Int64()
		if err != nil {
			logger.Error("rate limit check failed", "ip", ip, "err", err)
			c.Next()
			return
		}

		if result == 1 {
			response.Abort(c, http.StatusTooManyRequests, reverseProxyConfig)
			Ban(c, redis, banDuration, reverseProxyConfig)
			return
		}

		c.Next()
	}
}
