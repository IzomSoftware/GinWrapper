package https_core

import (
	"net/http"
	"strings"

	"github.com/IzomSoftware/GinWrapper/configuration"
	"github.com/gin-gonic/gin"
)

type Protections struct {
	/*
	 * Contains protections as basic as path check.
	 * Example: if path is /api/auth or /api/ (and is it malicious?)
	 */
	BasicProtections bool
	UserAgent        bool
	JWT              bool
	RateLimit        bool
	Ban              bool
}

type Response struct {
	Handler     gin.HandlerFunc
	Type        string
	Addresses   []string
	Protections Protections
}

var (
	Responses    = map[string]Response{}
	NoRouteRoute = func(c *gin.Context) { c.String(http.StatusNotFound, "404 Not Found") }
)

func NoRoute(c *gin.Context) {
	if configuration.ConfigHolder.Protections.ProvideBasicProtections {
		path := c.Request.URL.Path
		basePath := strings.Split(path, "/")[1]

		for _, response := range Responses {
			for _, address := range response.Addresses {
				if strings.Contains(address, basePath) {
					ip := c.ClientIP()

					AbortConnection(ip, c)

					if configuration.ConfigHolder.Protections.AggressiveBasicProtections {
						BanConnection(ip, c)
					}

					return
				}
			}
		}
	}

	NoRouteRoute(c)
}

func ActivateJWTAPI() {
	Responses["JWTApiAuth"] = Response{
		Handler: func(c *gin.Context) {
			c.JSON(http.StatusOK, "")
		},
		Type:      "GET",
		Addresses: []string{"/api/auth"},
		Protections: Protections{
			BasicProtections: true,
			RateLimit:        true,
		},
	}

	Responses["JWTAPIValidate"] = Response{
		Handler: func(c *gin.Context) {
			c.JSON(http.StatusForbidden, "")
		},
		Type:      "GET",
		Addresses: []string{},
		Protections: Protections{
			BasicProtections: true,
			RateLimit:        true,
		},
	}
}

func (R *Response) OnProtected(c *gin.Context) {
	response := R

	if !response.IsAnyProtectionEnabled() {
		return
	}

	protections := response.Protections
	ip := c.ClientIP()
	header := c.Request.Header

	// Path check
	// if response.Protections.BasicProtections {

	// }

	// User agent check
	if userAgent, apiUserAgent := header.Get("User-Agent"),
		configuration.ConfigHolder.Protections.APIUserAgent;

	// The actual check, if user agent is valid, perform no further action
	protections.UserAgent && userAgent == apiUserAgent {
		return
	}

	if response.Protections.JWT {
		// Redirect to auth API
		// TODO: more logic here
	}

	// if response.RateLimit {
	// 	//
	// }

	if protections.Ban {
		BanConnection(ip, c)
	}
}
