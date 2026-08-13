package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func UserAgent(useragent string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if useragent == "" {
			c.Next()
			return
		}

		headerAgent := c.GetHeader("User-Agent")

		if !strings.EqualFold(headerAgent, useragent) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}
