package middleware

import (
	"github.com/IzomSoftware/GinWrapper/configuration"
	"github.com/IzomSoftware/GinWrapper/logger"
	"github.com/gin-gonic/gin"
)

func Logging(config configuration.ReverseProxyMode) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.LogConnection(c, config)
		c.Next()
	}
}
