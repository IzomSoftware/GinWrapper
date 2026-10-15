package middleware

import (
	"net/http"

	"github.com/IzomSoftware/GinWrapper/response"
	"github.com/gin-gonic/gin"
)

func ValidatePOST(fields...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, field := range(fields) {
			if c.Query(field) == "" {
				response.Abort(c, http.StatusBadRequest)
				break
			}
		}
	}
}