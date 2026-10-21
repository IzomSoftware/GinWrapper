package response

import (
	"net/http"

	"github.com/IzomSoftware/GinWrapper/cidr"
	"github.com/IzomSoftware/GinWrapper/configuration"
	"github.com/IzomSoftware/GinWrapper/logger"
	"github.com/gin-gonic/gin"
)

func Abort(c *gin.Context, status int, config configuration.ReverseProxyMode) {
	logger.Info("aborted", "ip", cidr.GetIP(c, config), "status", status)
	c.AbortWithStatus(status)
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
