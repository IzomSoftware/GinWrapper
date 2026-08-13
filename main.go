package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/IzomSoftware/GinWrapper/authentication"
	"github.com/IzomSoftware/GinWrapper/configuration"
	"github.com/IzomSoftware/GinWrapper/logger"
	"github.com/IzomSoftware/GinWrapper/middleware"
	"github.com/IzomSoftware/GinWrapper/response"
	"github.com/IzomSoftware/GinWrapper/server"
	"github.com/IzomSoftware/GinWrapper/storage"
	"github.com/gin-gonic/gin"
)

const creationSchema = `
	CREATE TABLE IF NOT EXISTS Users (
		username TEXT PRIMARY KEY,
		hash TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS BannedIPs (
		ip TEXT PRIMARY KEY
	);
`

func main() {
	configuration, err := configuration.LoadConfiguration("config.toml")
	if err != nil {
		panic("Failed to initialize configuration")
	}

	logLevel := slog.LevelInfo
	if configuration.Debug {
		logLevel = slog.LevelDebug
	}
	logger.SetupLogger(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))

	storage, err := storage.New(configuration, creationSchema)
	if err != nil {
		panic("Failed to intiialize storage")
	}
	defer storage.Close()

	jwtSecret := configuration.Protections.JWTProtection.JWTSecret
	if jwtSecret == "" {
		jwtSecret, err = authentication.GenerateRandomSecret(32)
		if err != nil {
			panic("Failed to generate a random secret for JWT")
		}
	}

	jwtManager := authentication.NewJWTManager(
		jwtSecret,
		"GinWrapper",
		time.Duration(configuration.Protections.JWTProtection.JWTExpiration)*time.Second,
		24*time.Hour,
	)

	server := server.NewServer(configuration, storage, jwtManager)
	server.Use(gin.Recovery(), middleware.Logging())

	if storage.Redis != nil {
		server.Use(middleware.BanCheck(storage.Redis))
	}
	if configuration.Protections.APIUserAgent != "" {
		server.Use(middleware.UserAgent(configuration.Protections.APIUserAgent))
	}
	if configuration.Protections.UserPassAPI {
		server.Use(middleware.Authentication(jwtManager))
	}
	if storage.Redis != nil {
		server.Use(middleware.RateLimit(storage.Redis, middleware.RateLimitConfig{Rate: 30, Window: 60}))
	}

	server.RegisterRoute("POST", "/api/auth/register", func(c *gin.Context) {
		username, password := c.PostForm("username"), c.PostForm("password")
		hash, err := authentication.GenerateHash(password)

		if err != nil {
			response.AbortInternalError(c)
			return
		}

		err = storage.SQL.ExecuteUpdate("INSERT INTO Users (username, hash) VALUES (?, ?)", username, hash)
		if err != nil {
			response.Abort(c, http.StatusBadRequest)
			return
		}

		pair, err := jwtManager.GenerateJWTPair(username, username)
		if err != nil {
			response.AbortInternalError(c)
			return
		}

		c.JSON(http.StatusOK, pair)
	})

	server.RegisterRoute("POST", "/api/auth/login", func(c *gin.Context) {
		username, password := c.PostForm("username"), c.PostForm("password")
		var hash string
		err := storage.SQL.QueryRow("SELECT hash FROM Users WHERE username = ?", username).Scan(&hash)
		if err != nil || authentication.ValidateHash(hash, password) != nil {
			response.AbortUnauthorized(c)
			return
		}

		pair, err := jwtManager.GenerateJWTPair(username, username)
		if err != nil {
			response.AbortInternalError(c)
			return
		}

		c.JSON(http.StatusOK, pair)
	})

	server.RegisterRoute("POST", "/api/auth/refresh", func(c *gin.Context) {
		refreshToken := c.PostForm("refresh_token")
		pair, err := jwtManager.RefreshToken(refreshToken)

		if err != nil {
			response.AbortUnauthorized(c)
			return
		}

		c.JSON(http.StatusOK, pair)
	})

	protected := server.Engine.Group("/api/protected")
	protected.Use(middleware.Authentication(jwtManager))

	server.LoadTemplates(configuration.HTTPServer.TemplatesDir + "*")
	server.LoadStatics(configuration.HTTPServer.AssetsDir, "."+configuration.HTTPServer.AssetsDir)
	if err := server.ListenAndServe(); err != nil {
		panic(fmt.Sprintf("Failed to listen: %v", err))
	}
}
