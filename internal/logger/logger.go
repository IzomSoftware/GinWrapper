package logger

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

func SetupLogger(format string, handler *slog.Handler) {
	slog.SetDefault(slog.New(*handler))
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

func LogConnection(connection *gin.Context) {
	Info("[%s] -> %s", connection.ClientIP(), connection.Request.URL.Path)
}