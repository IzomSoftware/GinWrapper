package https_core

import (
	"fmt"

	"github.com/IzomSoftware/GinWrapper/configuration"
	"github.com/IzomSoftware/GinWrapper/logger"

	"github.com/gin-gonic/gin"
)

type HttpsServer struct {
	Router *gin.Engine
}

func middleware(context *gin.Context) {
	logger.LogConnection(context)

	connectionRequest := context.FullPath()

	if configuration.IsStorageConfigured() {
		for _, req := range Responses {
			for _, address := range req.Addresses {
				if connectionRequest == address {
					req.OnProtected(context)
				}
			}
		}
	}

	context.Next()
}

func (H *HttpsServer) ListenAndServe(templatesDir string, assetsDir string) {
	httpConfig := configuration.ConfigHolder.HTTPServer

	if !httpConfig.Enabled {
		return
	}

	gin.SetMode(gin.ReleaseMode)

	H.Router = gin.New()
	H.Router.Use(middleware)
	H.Router.LoadHTMLGlob(templatesDir)
	H.Router.Static(assetsDir, "."+assetsDir)

	// Register No Route
	H.Router.NoRoute(NoRoute)

	isAnyJWTAPIUsed := false

	// Registering the Paths and responses
	for name, req := range Responses {
		for _, address := range req.Addresses {
			if req.Protections.JWT {
				isAnyJWTAPIUsed = true
			}
			logger.Logger.Info(fmt.Sprintf("Registering Route -> %s - %s <-", name, address))
			H.Router.Handle(req.Type, address, req.Handler)
		}
	}

	// Register additional JWT Apis
	if isAnyJWTAPIUsed {
		ActivateJWTAPI()
	}

	addr := fmt.Sprintf("%s:%d", configuration.ConfigHolder.HTTPServer.Address, configuration.ConfigHolder.HTTPServer.Port)

	logger.Logger.Info(fmt.Sprintf("Listening on %s", addr))

	var err error

	if httpConfig.TlsConfiguration.Enable {
		err = H.Router.RunTLS(addr, httpConfig.TlsConfiguration.CertFile, httpConfig.TlsConfiguration.KeyFile)
	} else {
		err = H.Router.Run(addr)
	}

	if err != nil {
		logger.Logger.Error(err)
	}
}
