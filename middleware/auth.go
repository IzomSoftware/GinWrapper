package middleware

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/IzomSoftware/GinWrapper/authentication"
	"github.com/IzomSoftware/GinWrapper/configuration"
	"github.com/IzomSoftware/GinWrapper/response"
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

func AuthenticationWithCookie(jwtManager *authentication.JWTManager, config configuration.ReverseProxyMode) gin.HandlerFunc {
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
			response.AbortUnauthorized(c, config)
			return
		}

		claims, err := jwtManager.ValidateJWT(token, "access")
		if err != nil {
			response.AbortUnauthorized(c, config)
			return
		}

		c.Set("claims", claims)
		c.Set("uuid", claims.Uuid)
		c.Set("username", claims.Username)
		c.Set("token_type", claims.TokenType)
		c.Next()
	}
}

func CSRFCheck(config configuration.ReverseProxyMode) gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		default:
			c.Next()
			return
		}

		origin := c.GetHeader("Origin")
		if origin == "" {
			if ref := c.GetHeader("Referer"); ref != "" {
				origin = ref
			} else {
				c.Next()
				return
			}
		}

		if link, err := url.Parse(origin); err != nil || link.Host != c.Request.Host {
			response.AbortForbidden(c, config)
			return
		}
		c.Next()
	}
}
