package https_core

import (
	"fmt"
	"net/http"

	"github.com/IzomSoftware/GinWrapper/common/configuration"
	"github.com/IzomSoftware/GinWrapper/mysql"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Fn        gin.HandlerFunc
	Method    string
	Addresses []string
	Protected bool
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

func (R *Response) OnProtected(c *gin.Context) {
	if !R.Protected {
		return
	}

	userAgent := c.Request.Header.Get("User-Agent")
	apiUserAgent := configuration.ConfigHolder.HTTPSServer.APIUserAgent
	ip := c.ClientIP()

	fmt.Println(userAgent, " ", apiUserAgent)

	if userAgent != apiUserAgent {
		c.AbortWithStatus(http.StatusForbidden)
		if userAgent != apiUserAgent {
			mysql.BanIp(ip)
		}
		return
	}
}
