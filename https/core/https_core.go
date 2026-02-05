package https_core

import (
	"fmt"

	"github.com/IzomSoftware/GinWrapper/common/configuration"
	"github.com/IzomSoftware/GinWrapper/common/logger"

	"github.com/gin-gonic/gin"
)

type HttpsServer struct {
	Router *gin.Engine
}

func middleware(context *gin.Context) {
	LogConnection(context)

	connectionRequest := context.Request.URL.Path

	if !configuration.ConfigHolder.SQLLiteConfiguration.Enabled {
		for _, req := range Responses {
			for _, address := range req.Addresses {
				if req.Protected && connectionRequest == address {
					req.OnProtected(context)
				}
			}
		}
	}

	context.Next()
}

func (H *HttpsServer) ListenAndServe(templatesDir string, assetsDir string) {
	httpConfig := configuration.ConfigHolder.HTTPSServer

	if !httpConfig.Enabled {
		return
	}

	gin.SetMode(gin.ReleaseMode)

	H.Router = gin.New()
	H.Router.Use(middleware)
	H.Router.LoadHTMLGlob(templatesDir)
	H.Router.Static(assetsDir, "."+assetsDir)

	//Handling 404 error
	if unknw, ok := Responses["not-found-screen"]; ok {
		H.Router.NoRoute(unknw.Fn)
		delete(Responses, "not-found-screen")
	}

	//Registering the Paths and responses
	for name, req := range Responses {
		for _, address := range req.Addresses {
			logger.Logger.Debug(fmt.Sprintf("Registering Route -> %s - %s <-", name, address))
			H.Router.Handle(req.Method, address, req.Fn)
		}
	}

	addr := fmt.Sprintf("%s:%d", configuration.ConfigHolder.HTTPSServer.Address, configuration.ConfigHolder.HTTPSServer.Port)

	LogInfo(fmt.Sprintf("Listening on %s", addr))

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
