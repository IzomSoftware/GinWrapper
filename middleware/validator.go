package middleware

import (
	"net/http"

	"github.com/IzomSoftware/GinWrapper/configuration"
	"github.com/IzomSoftware/GinWrapper/response"
	"github.com/gin-gonic/gin"
)

func POSTFormCheck(config configuration.ReverseProxyMode, fields ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, field := range fields {
			if c.PostForm(field) == "" {
				response.Abort(c, http.StatusBadRequest, config)
				return
			}
		}
		c.Next()
	}
}
