package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/IzomSoftware/GinWrapper/authentication"
	"github.com/IzomSoftware/GinWrapper/configuration"
	"github.com/IzomSoftware/GinWrapper/response"
	"github.com/IzomSoftware/GinWrapper/storage/redis"
	"github.com/gin-gonic/gin"
)

func Authentication(jwtManager *authentication.JWTManager, redis *redis.Storage, banDuration time.Duration, config configuration.ReverseProxyMode) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			response.AbortUnauthorized(c, config)
			Ban(c, redis, banDuration, config)
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.AbortUnauthorized(c, config)
			Ban(c, redis, banDuration, config)
			return
		}

		claims, err := jwtManager.ValidateJWT(parts[1], "access")
		if err != nil {
			response.AbortUnauthorized(c, config)
			Ban(c, redis, banDuration, config)
			return
		}

		c.Set("claims", claims)
		c.Set("uuid", claims.Uuid)
		c.Set("username", claims.Username)
		c.Next()
	}
}

func AuthenticationWithCookie(jwtManager *authentication.JWTManager, redis *redis.Storage, banDuration time.Duration, config configuration.ReverseProxyMode) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := ""
		if authHeader := c.GetHeader("Authorization"); authHeader != "" {
			if scheme, rest, ok := strings.Cut(authHeader, " "); ok && strings.EqualFold(scheme, "Bearer") {
				token = strings.TrimSpace(rest)
			}
		}
		if token == "" {
			if cookie, err := c.Cookie(jwtManager.CookieName); err == nil && cookie != "" {
				token = cookie
			}
		}
		if token == "" {
			response.AbortUnauthorized(c, config)
			Ban(c, redis, banDuration, config)
			return
		}

		claims, err := jwtManager.ValidateJWT(token, "access")
		if err != nil {
			response.AbortUnauthorized(c, config)
			Ban(c, redis, banDuration, config)
			return
		}

		c.Set("claims", claims)
		c.Set("uuid", claims.Uuid)
		c.Set("username", claims.Username)
		c.Next()
	}
}

func CSRFCheck(redis *redis.Storage, banDuration time.Duration, config configuration.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		default:
			c.Next()
			return
		}

		origin := c.GetHeader("Origin")
		if origin == "" {
			origin = c.GetHeader("Referer")
		}

		allowedOrigins := config.Protections.AllowedOrigins
		if len(allowedOrigins) == 0 {
			c.Next()
			return
		}

		allowed := false

		for _, a := range allowedOrigins {
			if origin == a {
				allowed = true
				break
			}
		}
		if !allowed {
			reverseProxyMode := config.ReverseProxyMode
			response.AbortForbidden(c, reverseProxyMode)

			Ban(c, redis, banDuration, reverseProxyMode)
			return
		}

		c.Next()
	}
}


func IsRevoked(redis *redis.Storage, jti string) bool {
    exists, err := redis.Exists("revoked:" + jti)
    if err != nil {
        return false
    }
    return exists
}