package main

import (
	"fmt"
	"net/http"

	"github.com/IzomSoftware/GinWrapper/internal/configuration"
	httpscore "github.com/IzomSoftware/GinWrapper/internal/https/core"
	"github.com/IzomSoftware/GinWrapper/internal/logger"
	"github.com/IzomSoftware/GinWrapper/internal/responses"
	"github.com/IzomSoftware/GinWrapper/internal/storage/redis_source"
	"github.com/IzomSoftware/GinWrapper/internal/storage/sql_source"
	"github.com/IzomSoftware/GinWrapper/internal/utils/jwt_util"
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
	responses.Responses["home"] = &responses.Response{
		Handler: func(c *gin.Context) {
			c.String(http.StatusOK, "")
		},
		Type:      "GET",
		Addresses: []string{"/home"},
		Protections: responses.Protections{
			RateLimit: responses.RateLimitProtection{
				Enabled: true,
				Rate:    5,
				Time:    1,
			},
		},
	}

	// first argument is templateDir and second one is assetsDir
	HttpsServer.ListenAndServe(fmt.Sprintf("%s*", config.HTTPServer.TemplatesDir), config.HTTPServer.AssetsDir)
}
