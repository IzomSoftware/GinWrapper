package middleware

import (
	"net/http"
	"strings"

	"github.com/IzomSoftware/GinWrapper/authentication"
	"github.com/gin-gonic/gin"
)

func Authentication(jwtManager *authentication.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, err := jwtManager.ValidateJWT(parts[1], "access")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("claims", claims)
		c.Set("uuid", claims.Uuid)
		c.Set("username", claims.Username)
		c.Set("token_type", claims.TokenType)
		c.Next()
	}
}

func AuthenticationWithCookie(jwtManager *authentication.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := ""
		if authHeader := c.GetHeader("Authorization"); len(authHeader) > 7 && strings.EqualFold(authHeader[:7], "Bearer ") {
			token = strings.TrimSpace(authHeader[7:])
		}
		if token == "" {
			if cookie, err := c.Cookie(jwtManager.CookieName); err == nil && cookie != "" {
				token = cookie
			}
		}
		if token == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, err := jwtManager.ValidateJWT(token, "access")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("claims", claims)
		c.Set("uuid", claims.Uuid)
		c.Set("username", claims.Username)
		c.Set("token_type", claims.TokenType)
		c.Next()
	}
}
