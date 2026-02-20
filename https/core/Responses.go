package https_core

import (
	"net/http"

	"github.com/IzomSoftware/GinWrapper/configuration"
	mysql "github.com/IzomSoftware/GinWrapper/sql"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Fn                 gin.HandlerFunc
	Method             string
	Addresses          []string
	UserAgentProtected bool
	JWTProtected       bool
	مکBanOnFail        bool
}

var (
	Responses = map[string]Response{
		"not-found-screen": {
			Fn: func(c *gin.Context) {
				c.HTML(http.StatusNotFound, "not-found.html", nil)
			},
			Method: "GET",
		},
	}
)

func BanConnection(ip string, c *gin.Context) {
	c.AbortWithStatus(http.StatusForbidden)
	mysql.BanIP(ip)
}

func (R *Response) OnProtected(c *gin.Context) {
	if !R.JWTProtected || !R.UserAgentProtected {
		return
	}

	userAgent := c.Request.Header.Get("User-Agent")
	apiUserAgent := configuration.ConfigHolder.Protections.APIUserAgent
	ip := c.ClientIP()

	// Check for user agent
	if R.UserAgentProtected && userAgent == apiUserAgent {
		return
	}

	// TODO: add token check

	BanConnection(ip, c)
}