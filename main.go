package main

import (
	"net/http"

	"github.com/SyNdicateFoundation/GinWrapper/common/configuration"
	"github.com/SyNdicateFoundation/GinWrapper/common/logger"
	httpscore "github.com/SyNdicateFoundation/GinWrapper/https/core"
	"github.com/gin-gonic/gin"
)

var (
	HttpsServer httpscore.HttpsServer
)

func main() {
	logger.SetupLogger("SyNdicateWebsite")

	// Adjust as needed
	configuration.DefaultConfig =
		configuration.Holder{
			Debug: false,
			HTTPSServer: configuration.HTTPSServer{
				Enabled:      true,
				Address:      "0.0.0.0",
				Port:         2009,
				APIUserAgent: "LiteGuard Client 1.0/b (Software)",
				TlsConfiguration: configuration.HttpsTlsConfiguration{
					Enable:   false,
					CertFile: "cert.pem",
					KeyFile:  "key.pem",
				},
			},
			SQLLiteConfiguration: configuration.SQLLiteConfiguration{
				Enabled:              false,
				DatabaseFileLocation: "db.sqlite",
			},
		}
	// setup configuration
	configuration.SetupConfig("config.toml")

	// add responses
	httpscore.Responses["index"] = httpscore.Response{
		Fn: func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", nil)
		},
		Method:    "GET",
		Addresses: []string{"/", "/index.html"},
	}
	httpscore.Responses["projects"] = httpscore.Response{
		Fn: func(c *gin.Context) {
			c.HTML(http.StatusOK, "projects.html", nil)
		},
		Method:    "GET",
		Addresses: []string{"/projects", "/projects.html"},
	}
	httpscore.Responses["projects"] = httpscore.Response{
		Fn: func(c *gin.Context) {
			c.HTML(http.StatusOK, "members.html", nil)
		},
		Method:    "GET",
		Addresses: []string{"/members", "/members.html"},
	}
	httpscore.Responses["technologies"] = httpscore.Response{
		Fn: func(c *gin.Context) {
			c.HTML(http.StatusOK, "technologies.html", nil)
		},
		Method:    "GET",
		Addresses: []string{"/technologies", "/technologies.html"},
	}
	httpscore.Responses["colleagues"] = httpscore.Response{
		Fn: func(c *gin.Context) {
			c.HTML(http.StatusOK, "colleagues.html", nil)
		},
		Method:    "GET",
		Addresses: []string{"/colleagues", "/colleagues.html"},
	}

	// first argument is templateDir and second one is assetsDir
	HttpsServer.ListenAndServe("assets/templates/*", "/assets")
}
