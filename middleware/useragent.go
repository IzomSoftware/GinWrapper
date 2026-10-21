package middleware

import (
	"strings"

	"github.com/IzomSoftware/GinWrapper/configuration"
	"github.com/IzomSoftware/GinWrapper/response"
	"github.com/gin-gonic/gin"
)

func UserAgentCheck(useragent string, config configuration.ReverseProxyMode) gin.HandlerFunc {
	return func(c *gin.Context) {
		if useragent == "" {
			c.Next()
			return
		}

		headerAgent := c.GetHeader("User-Agent")

		if !strings.EqualFold(headerAgent, useragent) {
			response.AbortForbidden(c, config)
			return
		}

		c.Next()
	}
}
