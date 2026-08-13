package middleware

import (
	"net/http"
	"time"

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

type RateLimitConfig struct {
	Rate   int64
	Window int64
}

func RateLimit(redis *redis.Storage, configuration RateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now().UnixMilli()

		result, err := redis.RunScript(script, []string{ip}, configuration.Rate, configuration.Window*1000, now).Int64()

		if err != nil {
			c.Next()
			return
		}

		if result == 1 {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}

		c.Next()
	}
}
