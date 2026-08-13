package response

import (
	"net/http"
	"strings"

	"github.com/IzomSoftware/GinWrapper/logger"
	"github.com/gin-gonic/gin"
)

func Abort(c *gin.Context, status int) {
	c.AbortWithStatus(status)
	logger.Info("aborted", "ip", c.ClientIP(), "status", status)
}

func AbortForbidden(c *gin.Context) {
	Abort(c, http.StatusForbidden)
}

func AbortUnauthorized(c *gin.Context) {
	Abort(c, http.StatusUnauthorized)
}

func AbortInternalError(c *gin.Context) {
	Abort(c, http.StatusInternalServerError)
}

func NoRoute(c *gin.Context) {
	c.String(http.StatusNotFound, "404 Not Found")
}

func NoRouteWithProtection(registeredPaths []string, aggressive bool, ban func(string)) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		ip := c.ClientIP()
		parts := strings.Split(strings.TrimPrefix(path, "/"), "/")

		if parts[0] != "" {
			base := "/" + parts[0]

			for _, registered := range registeredPaths {
				if strings.HasPrefix(registered, base) {
					AbortForbidden(c)

					if aggressive && ban != nil {
						ban(ip)
					}
					return
				}
			}
		}

		NoRoute(c)
	}
}
