package https_core

// import (
// 	"fmt"

// 	"github.com/IzomSoftware/GinWrapper/internal/configuration"
// 	"github.com/IzomSoftware/GinWrapper/internal/logger"
// 	"github.com/IzomSoftware/GinWrapper/internal/responses"

// 	"github.com/gin-gonic/gin"
// )

// // The HTTPS server struct. contains the gin.Engine which is the Router.
// type HttpsServer struct {
// 	Router *gin.Engine
// }

// // The middleware, which gets called when every single connection begins.
// func middleware(context *gin.Context) {
// 	logger.LogConnection(context)

// 	connectionRequest := context.FullPath()

// 	// We can't just provide every single protection without any
// 	// Storage configured.
// 	if configuration.IsStorageConfigured() {
// 		for _, req := range responses.Responses {
// 			for _, address := range req.Addresses {
// 				if connectionRequest == address {
// 					req.OnProtected(context)
// 				}
// 			}
// 		}
// 	}

// 	context.Next()
// }

// // Registers every single path
// func (H *HttpsServer) RegisterPaths() {
// 	for _, req := range responses.Responses {
// 		for _, address := range req.Addresses {
// 			H.Router.Handle(req.Type, address, req.Handler)
// 		}
// 	}
// }

// // Listens & serves the HTTP/HTTPS server
// func (H *HttpsServer) ListenAndServe(templatesDir string, assetsDir string) {
// 	httpConfig := configuration.ConfigHolder.HTTPServer

// 	if !httpConfig.Enabled {
// 		return
// 	}

// 	gin.SetMode(gin.ReleaseMode)

// 	H.Router = gin.New()
// 	H.Router.Use(middleware)
// 	H.Router.LoadHTMLGlob(templatesDir)
// 	H.Router.Static(assetsDir, "."+assetsDir)

// 	// Register the 404 screen
// 	H.Router.NoRoute(responses.NoRoute)

// 	isAnyUserPassAPIUsed := false
// 	isAnyJWTAPIUsed := false

// 	for name, req := range responses.Responses {
// 		for _, address := range req.Addresses {
// 			if req.Protections.JWT {
// 				isAnyUserPassAPIUsed = true
// 				isAnyJWTAPIUsed = true
// 			}
// 			logger.Logger.Info(fmt.Sprintf("Registering Route -> %s - %s <-", name, address))
// 		}
// 	}

// 	// Register additional UserPass APIs
// 	if isAnyUserPassAPIUsed || configuration.ConfigHolder.Protections.UserPassAPI {
// 		responses.ActivateUserPassAPI()
// 	}

// 	// Register additional JWT APIs
// 	if isAnyJWTAPIUsed {
// 		responses.ActivateJWTAPI()
// 	}

// 	// Finalize the path registeration (including APIs we provide)
// 	H.RegisterPaths()

// 	addr := fmt.Sprintf("%s:%d", configuration.ConfigHolder.HTTPServer.Address, configuration.ConfigHolder.HTTPServer.Port)

// 	logger.Logger.Info(fmt.Sprintf("Listening on %s", addr))

// 	var err error

// 	if httpConfig.TlsConfiguration.Enable {
// 		err = H.Router.RunTLS(addr, httpConfig.TlsConfiguration.CertFile, httpConfig.TlsConfiguration.KeyFile)
// 	} else {
// 		err = H.Router.Run(addr)
// 	}

// 	if err != nil {
// 		logger.Logger.Error(err)
// 	}
// }
