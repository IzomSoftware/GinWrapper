package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/IzomSoftware/GinWrapper/authentication"
	"github.com/IzomSoftware/GinWrapper/configuration"
	"github.com/IzomSoftware/GinWrapper/file"
	"github.com/IzomSoftware/GinWrapper/logger"
	"github.com/IzomSoftware/GinWrapper/response"
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

func (S *Server) RegisterAuthRoutes(identifier string, secret string, query string, config configuration.ReverseProxyMode) {
	S.Engine.POST("/api/auth/login", func(c *gin.Context) {
		key, val := c.PostForm(identifier), c.PostForm(secret)
		var hash string
		err := S.storage.SQL.QueryRow(
			query, key,
		).Scan(&hash)
		if err != nil || authentication.ValidateHash(hash, val) != nil {
			response.AbortUnauthorized(c, config)
			return
		}
		pair, err := S.jwtManager.GenerateJWTPair(key, key)
		if err != nil {
			response.AbortInternalError(c, config)
			return
		}
		S.jwtManager.SetAuthCookies(c, pair, S.configuration.Protections.JWTProtection.JWTExpiration)
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	S.Engine.POST("/api/auth/refresh", func(c *gin.Context) {
		refreshToken := c.PostForm("refresh_token")
		pair, err := S.jwtManager.RefreshToken(refreshToken)
		if err != nil {
			response.AbortUnauthorized(c, config)
			return
		}
		S.jwtManager.SetAuthCookies(c, pair, S.configuration.Protections.JWTProtection.JWTExpiration)
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	S.Engine.POST("/api/auth/logout", func(c *gin.Context) {
		S.jwtManager.ClearAuthCookies(c)
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
}
func (S *Server) RegisterFileUploadRoute(field string, path string, urlPrefix string, uploadDir string, allowedExt map[string]struct{}, maxSize int64, config configuration.ReverseProxyMode, middlewares ...gin.HandlerFunc) {
	handler := func(c *gin.Context) {
		uploadedFile, err := c.FormFile(field)
		if err != nil {
			response.Abort(c, http.StatusBadRequest, config)
			return
		}
		if uploadedFile.Size > maxSize {
			response.Abort(c, http.StatusRequestEntityTooLarge, config)
			return
		}
		extension := file.GetFileExtension(uploadedFile.Filename)
		if _, ok := allowedExt[extension]; !ok {
			response.Abort(c, http.StatusUnsupportedMediaType, config)
			return
		}
		_ = os.MkdirAll(uploadDir, 0755)
		fileName := authentication.GenerateUUID() + extension
		if err := c.SaveUploadedFile(uploadedFile, fmt.Sprintf("%s/%s", uploadDir, fileName)); err != nil {
			response.AbortInternalError(c, config)
            return
		}
		c.JSON(http.StatusOK, gin.H{"ok": true, "url": urlPrefix + "/" + fileName})
	}
	S.Engine.POST(path, append(middlewares, handler)...)
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

	if !S.configuration.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	if S.configuration.ReverseProxyMode.Enabled && !S.configuration.ReverseProxyMode.ReadHeaderExplicitly{
		var trustedProxies []string
		for _, ip := range S.configuration.ReverseProxyMode.IPList {
			trustedProxies = append(trustedProxies, ip)
		}
		S.Engine.SetTrustedProxies(trustedProxies)
	}

	addr := fmt.Sprintf("%s:%d", httpServerConfiguration.Address, httpServerConfiguration.Port)
	logger.Info("listening", "addr", addr)

	if httpServerConfiguration.TlsConfiguration.Enable {
		return S.Engine.RunTLS(addr, httpServerConfiguration.TlsConfiguration.CertFile, httpServerConfiguration.TlsConfiguration.KeyFile)
	}
	return S.Engine.Run(addr)
}
