package response

import (
	"net/http"
	"path"
	"strings"

	"github.com/IzomSoftware/GinWrapper/cidr"
	"github.com/IzomSoftware/GinWrapper/configuration"
	"github.com/IzomSoftware/GinWrapper/logger"
	"github.com/gin-gonic/gin"
)

func Abort(c *gin.Context, status int, config configuration.ReverseProxyMode) {
	c.AbortWithStatus(status)
	logger.Info("aborted", "ip", cidr.GetIP(c, config), "status", status)
}

func AbortForbidden(c *gin.Context, config configuration.ReverseProxyMode) {
	Abort(c, http.StatusForbidden, config)
}

func AbortUnauthorized(c *gin.Context, config configuration.ReverseProxyMode) {
	Abort(c, http.StatusUnauthorized, config)
}

func AbortInternalError(c *gin.Context, config configuration.ReverseProxyMode) {
	Abort(c, http.StatusInternalServerError, config)
}

func NoRoute(c *gin.Context) {
	c.String(http.StatusNotFound, "404 Not Found")
}

func NoRouteWithProtection(registeredPaths []string, aggressive bool, ban func(string), config configuration.ReverseProxyMode) gin.HandlerFunc {
	return func(c *gin.Context) {
		uriPath := c.Request.URL.Path
		cleanedUriPath := path.Clean(uriPath)
		ip := cidr.GetIP(c, config)

		if cleanedUriPath != uriPath {
			AbortForbidden(c, config)
			if aggressive && ban != nil {
				ban(ip)
			}
			return
		}

		parts := strings.Split(strings.TrimPrefix(cleanedUriPath, "/"), "/")
		if parts[0] != "" {
			base := "/" + parts[0]

			for _, registered := range registeredPaths {
				if strings.HasPrefix(registered, base) {
					AbortForbidden(c, config)

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
