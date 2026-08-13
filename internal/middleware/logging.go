package middleware

import (
	"github.com/IzomSoftware/GinWrapper/internal/logger"
	"github.com/gin-gonic/gin"
)

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.LogConnection(c)
		c.Next()
	}
}
