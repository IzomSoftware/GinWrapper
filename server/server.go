package server

import (
	"fmt"

	"github.com/IzomSoftware/GinWrapper/authentication"
	"github.com/IzomSoftware/GinWrapper/configuration"
	"github.com/IzomSoftware/GinWrapper/logger"
	"github.com/IzomSoftware/GinWrapper/storage"
	"github.com/gin-gonic/gin"
)

type Server struct {
	configuration *configuration.Config
	storage       *storage.Storage
	jwtManager    *authentication.JWTManager
	Engine        *gin.Engine
}

func NewServer(configuration *configuration.Config, storage *storage.Storage, jwtManager *authentication.JWTManager) *Server {
	return &Server{
		configuration: configuration,
		storage:       storage,
		jwtManager:    jwtManager,
		Engine:        gin.New(),
	}
}

func (S *Server) Use(handlerfuncs ...gin.HandlerFunc) {
	S.Engine.Use(handlerfuncs...)
}

func (S *Server) RegisterRoutes(routes map[string]map[string]gin.HandlerFunc) {
	for method, paths := range routes {
		for path, handler := range paths {
			S.Engine.Handle(method, path, handler)
		}
	}
}

func (S *Server) RegisterRoute(method string, path string, handler gin.HandlerFunc) {
	S.Engine.Handle(method, path, handler)
}

func (S *Server) LoadTemplates(path string) {
	S.Engine.LoadHTMLGlob(path)
}

func (S *Server) LoadStatics(basePath string, path string) {
	S.Engine.Static(basePath, path)
}

func (S *Server) SetNoRoute(handler gin.HandlerFunc) {
	S.Engine.NoRoute(handler)
}

func (S *Server) ListenAndServe() error {
	httpServerConfiguration := S.configuration.HTTPServer
	if !httpServerConfiguration.Enabled {
		return nil
	}

	gin.SetMode(gin.ReleaseMode)

	addr := fmt.Sprintf("%s:%d", httpServerConfiguration.Address, httpServerConfiguration.Port)
	logger.Info("listening", "addr", addr)

	if httpServerConfiguration.TlsConfiguration.Enable {
		return S.Engine.RunTLS(addr, httpServerConfiguration.TlsConfiguration.CertFile, httpServerConfiguration.TlsConfiguration.KeyFile)
	}
	return S.Engine.Run(addr)
}
