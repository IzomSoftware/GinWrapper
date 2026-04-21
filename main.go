package main

import (
	"fmt"
	"net/http"

	"github.com/IzomSoftware/GinWrapper/configuration"
	httpscore "github.com/IzomSoftware/GinWrapper/https/core"
	"github.com/IzomSoftware/GinWrapper/logger"
	"github.com/IzomSoftware/GinWrapper/storage/redis_source"
	"github.com/IzomSoftware/GinWrapper/storage/sql_source"
	"github.com/IzomSoftware/GinWrapper/utils/jwt_util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	HttpsServer httpscore.HttpsServer
)

func main() {
	logger.SetupLogger("Test Website", logrus.DebugLevel)

	secret, _ := jwt_util.GenerateJWTRandomSecret(32)

	// Adjust as needed
	configuration.DefaultConfig.Protections.JWTProtection.JWTSecret = secret

	config := configuration.DefaultConfig

	// setup configuration
	configuration.SetupConfig("config.toml")

	// Intiialize storage sources
	sql_source.Init()
	redis_source.Init()

	// add responses
	httpscore.Responses["home"] = &httpscore.Response{
		Handler: func(c *gin.Context) {
			c.String(http.StatusOK, "")
		},
		Type:      "GET",
		Addresses: []string{"/home"},
		Protections: httpscore.Protections{
			RateLimit: true,
		},
	}

	// first argument is templateDir and second one is assetsDir
	HttpsServer.ListenAndServe(fmt.Sprintf("%s*", config.HTTPServer.TemplatesDir), config.HTTPServer.AssetsDir)
}
