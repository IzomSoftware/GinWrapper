package main

import (
	"net/http"

	"github.com/IzomSoftware/GinWrapper/common/configuration"
	"github.com/IzomSoftware/GinWrapper/common/logger"
	httpscore "github.com/IzomSoftware/GinWrapper/https/core"
	"github.com/gin-gonic/gin"
)

var (
	HttpsServer httpscore.HttpsServer
)

func main() {
	logger.SetupLogger("Test Website")

	// Adjust as needed
	configuration.DefaultConfig =
		configuration.Holder{
			Debug: false,
			HTTPSServer: configuration.HTTPSServer{
				Enabled:      true,
				Address:      "0.0.0.0",
				Port:         2009,
				APIUserAgent: "Test Client 1.0/b (Software)",
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
			c.String(http.StatusOK, "test", nil)
		},
		Method:    "GET",
		Addresses: []string{"/", "/index.html"},
	}

	// first argument is templateDir and second one is assetsDir
	HttpsServer.ListenAndServe("./*", "./")
}
