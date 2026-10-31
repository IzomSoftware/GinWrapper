package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/IzomSoftware/GinWrapper/authentication"
	"github.com/IzomSoftware/GinWrapper/configuration"
	"github.com/IzomSoftware/GinWrapper/file"
	"github.com/IzomSoftware/GinWrapper/logger"
	"github.com/IzomSoftware/GinWrapper/response"
	"github.com/IzomSoftware/GinWrapper/storage"
	"github.com/gin-gonic/gin"
)

type Server struct {
	configuration  *configuration.Config
	storage        *storage.Storage
	jwtManager     *authentication.JWTManager
	Engine         *gin.Engine
	invalidHandler *gin.HandlerFunc
}

func NewServer(configuration *configuration.Config, storage *storage.Storage, jwtManager *authentication.JWTManager, invalidHandler *gin.HandlerFunc) *Server {
	if !configuration.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(gin.Recovery())

	if configuration.ReverseProxyMode.Enabled && !configuration.ReverseProxyMode.ReadHeaderExplicitly {
		var trustedProxies []string
		for _, ip := range configuration.ReverseProxyMode.IPList {
			trustedProxies = append(trustedProxies, ip)
		}
		engine.SetTrustedProxies(trustedProxies)
	}

	return &Server{
		configuration:  configuration,
		storage:        storage,
		jwtManager:     jwtManager,
		Engine:         engine,
		invalidHandler: invalidHandler,
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
		if c.PostForm(identifier) == "" || c.PostForm(secret) == "" {
			response.Abort(c, http.StatusBadRequest, config)
			return
		}
		key, val := c.PostForm(identifier), c.PostForm(secret)
		var hash string
		err := S.storage.SQL.QueryRow(
			query, key,
		).Scan(&hash)
		if err != nil || authentication.ValidateHash(hash, val) != nil {
			response.AbortUnauthorized(c, config)
			return
		}
		pair, err := S.jwtManager.GenerateJWTPair(authentication.GenerateUUID(), key)
		if err != nil {
			response.AbortInternalError(c, config)
			return
		}
		S.jwtManager.SetAuthCookies(c, pair, S.configuration.Protections.JWTProtection.JWTExpiration)
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	S.Engine.POST("/api/auth/refresh", func(c *gin.Context) {
		refreshToken := S.jwtManager.GetRefreshCookie(c)
		if refreshToken == "" {
			refreshToken = c.PostForm("refresh_token")
		}
		pair, err := S.jwtManager.RefreshToken(refreshToken)
		if err != nil {
			response.AbortUnauthorized(c, config)
			return
		}
		S.jwtManager.SetAuthCookies(c, pair, S.configuration.Protections.JWTProtection.JWTExpiration)
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	S.Engine.POST("/api/auth/logout", func(c *gin.Context) {
		if claims, exists := c.Get("claims"); exists {
			if jc, ok := claims.(*authentication.JWTClaims); ok && jc.JTI != "" && S.storage != nil && S.storage.Redis != nil {
				_ = S.storage.Redis.Set("revoked:"+jc.JTI, "1", S.jwtManager.RefreshTokenExpiry())
			}
		}
		S.jwtManager.ClearAuthCookies(c)
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
}
func (S *Server) RegisterFileUploadRoute(field string, path string, urlPrefix string, uploadDir string, allowedExt map[string]struct{}, maxSize int64, config configuration.ReverseProxyMode, middlewares ...gin.HandlerFunc) {
	handler := func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
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
		f, err := uploadedFile.Open()
		if err != nil {
			response.AbortInternalError(c, config)
			return
		}
		mime, err := file.DetectMime(f)
		f.Close()
		if err != nil {
			response.AbortInternalError(c, config)
			return
		}
		if _, ok := allowedExt[mime]; !ok {
			response.Abort(c, http.StatusUnsupportedMediaType, config)
			return
		}
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			response.AbortInternalError(c, config)
			return
		}
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

func (S *Server) LoadStatics(protection bool, basePath string, path string) {
	if protection {
		S.Static(basePath, path, *S.invalidHandler)
	} else {
		S.Engine.Static(basePath, path)
	}
}

func (S *Server) SetInvalid(handler gin.HandlerFunc) {
	S.invalidHandler = &handler
	S.Engine.HandleMethodNotAllowed = true
	S.Engine.NoRoute(*S.invalidHandler)
	S.Engine.NoMethod(*S.invalidHandler)
}

func (s *Server) Static(urlPrefix, fsRoot string, invalidHandler gin.HandlerFunc) {
	fileServer := http.StripPrefix(urlPrefix, http.FileServer(http.Dir(fsRoot)))
	root := fmt.Sprintf("%s%s", fsRoot, string(os.PathSeparator))

	s.Engine.GET(urlPrefix+"/*filepath", func(c *gin.Context) {
		filePath := filepath.Join(fsRoot, filepath.FromSlash(c.Param("filepath")))

		if !strings.HasPrefix(fmt.Sprintf("%s%s", filePath, string(os.PathSeparator)), root) {
			invalidHandler(c)
			return
		}

		if info, err := os.Stat(filePath); err != nil || info.IsDir() {
			invalidHandler(c)
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Request)
	})
}

func (s *Server) StaticFile(urlPath, filePath string, invalidHandler gin.HandlerFunc) {
	s.Engine.GET(urlPath, func(c *gin.Context) {
		if _, err := os.Stat(filePath); err != nil {
			invalidHandler(c)
			return
		}
		http.ServeFile(c.Writer, c.Request, filePath)
	})
}

func (S *Server) ListenAndServe() error {
	httpServerConfiguration := S.configuration.HTTPServer
	if !httpServerConfiguration.Enabled {
		return nil
	}

	addr := fmt.Sprintf("%s:%d", httpServerConfiguration.Address, httpServerConfiguration.Port)
	logger.Info("listening", "addr", addr)

	if httpServerConfiguration.TlsConfiguration.Enable {
		return S.Engine.RunTLS(addr, httpServerConfiguration.TlsConfiguration.CertFile, httpServerConfiguration.TlsConfiguration.KeyFile)
	}
	return S.Engine.Run(addr)
}
