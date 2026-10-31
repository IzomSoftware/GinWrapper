package middleware

import (
	"path"

	"github.com/IzomSoftware/GinWrapper/configuration"
	"github.com/IzomSoftware/GinWrapper/response"
	"github.com/gin-gonic/gin"
)

func CleanCheck(config configuration.ReverseProxyMode) gin.HandlerFunc {
	return func(c *gin.Context) {
		uriPath := c.Request.URL.Path
		if path.Clean(uriPath) != uriPath {
			response.AbortForbidden(c, config)
			return
		}
		c.Next()
	}
}
