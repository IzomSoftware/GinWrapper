package main

import (
	"fmt"
	"net/http"

	"github.com/IzomSoftware/GinWrapper/configuration"
	httpscore "github.com/IzomSoftware/GinWrapper/https/core"
	"github.com/IzomSoftware/GinWrapper/logger"
	"github.com/IzomSoftware/GinWrapper/redis"
	"github.com/IzomSoftware/GinWrapper/sql"
	utils "github.com/IzomSoftware/GinWrapper/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	HttpsServer httpscore.HttpsServer
)

func main() {
	logger.SetupLogger("Test Website", logrus.DebugLevel)

	secret, _ := utils.GenerateJWTRandomSecret(32)

	// Adjust as needed
	configuration.DefaultConfig.Protections.JWTProtection.JWTSecret = secret

	config := configuration.DefaultConfig

	// setup configuration
	configuration.SetupConfig("config.toml")

	// Intiialize storage sources
	sql.Init()
	redis.Init()

	// add responses
	httpscore.Responses["home"] = httpscore.Response{
		Handler: func(c *gin.Context) {
			c.String(http.StatusOK, "")
		},
		Type:        "GET",
		Addresses:   []string{"/home", "/home/", "/kos/nago/kooni"},
		Protections: httpscore.Protections{},
	}

	// first argument is templateDir and second one is assetsDir
	HttpsServer.ListenAndServe(fmt.Sprintf("%s*", config.HTTPServer.TemplatesDir), config.HTTPServer.AssetsDir)
}
