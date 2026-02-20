package main

import (
	"github.com/IzomSoftware/GinWrapper/configuration"
	httpscore "github.com/IzomSoftware/GinWrapper/https/core"
	utils "github.com/IzomSoftware/GinWrapper/utils"
	"github.com/IzomSoftware/GinWrapper/logger"
	"github.com/sirupsen/logrus"
)

var (
	HttpsServer httpscore.HttpsServer
)

func main() {
	logger.SetupLogger("Test Website", logrus.DebugLevel)

	secret, _ := utils.GenerateJWTRandomSecret(32)
	// Adjust as needed
	configuration.DefaultConfig =
		configuration.Configuration{
			Debug: false,
			HTTPServer: configuration.HTTPServer{
				Enabled: true,
				Address: "0.0.0.0",
				Port:    2009,
				TlsConfiguration: configuration.HttpsTlsConfiguration{
					Enable:   false,
					CertFile: "cert.pem",
					KeyFile:  "key.pem",
				},
			},
			DatabaseConfiguration: configuration.DatabaseConfiguration{
				Enabled:            false,
				SQLiteConfiguration: configuration.SQLiteConfiguration{
					Enabled: true,
					DatabaseLocation: "db.sqlite",
				},
			},
			Protections: configuration.Protections{
				APIUserAgent: "Test Client 1.0/b (Software)",
				JWTProtection: configuration.JWTProtection{
					JWTSecret:     secret,
					JWTExpiration: 60,
				},
			},
		}
	// setup configuration
	configuration.SetupConfig("config.toml")

	// add responses
	// httpscore.Responses["index"] = httpscore.Response{
	// 	Fn: func(c *gin.Context) {
	// 		c.String(http.StatusOK, "test", nil)
	// 	},
	// 	Method:    "GET",
	// 	Addresses: []string{"/", "/index.html"},
	// }

	// first argument is templateDir and second one is assetsDir
	HttpsServer.ListenAndServe("./*", "./")
}
