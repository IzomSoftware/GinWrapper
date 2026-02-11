package main

import (
	"net/http"

	"github.com/IzomSoftware/GinWrapper/common/configuration"
	"github.com/IzomSoftware/GinWrapper/common/logger"
	utils "github.com/IzomSoftware/GinWrapper/https/utils"
	httpscore "github.com/IzomSoftware/GinWrapper/https/core"
	"github.com/gin-gonic/gin"
)

var (
	HttpsServer httpscore.HttpsServer
)
func main() {
	logger.SetupLogger("Test Website")

	secret, _ := utils.GenerateRandomSecret(32)
	// Adjust as needed
	configuration.DefaultConfig =
		configuration.Holder{
			Debug: false,
			HTTPSServer: configuration.HTTPSServer{
				Enabled:      true,
				Address:      "0.0.0.0",
				Port:         2009,
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
			Protections: configuration.Protections {
				APIUserAgent: "Test Client 1.0/b (Software)",
				Tokenizer: configuration.Tokenizer {
					TokenizerSecret: secret,
					TokenExpiration: 60,
				},
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
