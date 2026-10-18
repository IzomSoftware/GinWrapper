package logger

import (
	"log/slog"

	"github.com/IzomSoftware/GinWrapper/cidr"
	"github.com/IzomSoftware/GinWrapper/configuration"
	"github.com/gin-gonic/gin"
)

func SetupLogger(handler slog.Handler) {
	slog.SetDefault(slog.New(handler))
}

func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}
func Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}
func Error(msg string, args ...any) {
	slog.Error(msg, args...)
}
func LogConnection(c *gin.Context, config configuration.ReverseProxyMode) {
	
	Info("connection", "ip", cidr.GetIP(c, config), "path", c.Request.URL.Path)
}
