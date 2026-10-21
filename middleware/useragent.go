package middleware

import (
	"strings"

	"github.com/IzomSoftware/GinWrapper/configuration"
	"github.com/IzomSoftware/GinWrapper/response"
	"github.com/gin-gonic/gin"
)

func UserAgentCheck(userAgentConfig configuration.Protections, reverseProxyConfig configuration.ReverseProxyMode) gin.HandlerFunc {
	return func(c *gin.Context) {
		useragent := userAgentConfig.APIUserAgent
		if useragent == "" {
			c.Next()
			return
		}

		headerAgent := c.GetHeader("User-Agent")

		if !strings.EqualFold(headerAgent, useragent) {
			response.AbortForbidden(c, reverseProxyConfig)
			return
		}

		c.Next()
	}
}
