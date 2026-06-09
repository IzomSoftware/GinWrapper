package logger

import (
	"fmt"
	"time"

	"container/list"

	"github.com/IzomSoftware/GinWrapper/internal/configuration"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger
var logs = list.New()

func Log() {
	for {
		time.Sleep(1 * time.Second)

		if logs.Len() < 1 {
			continue
		}

		front := logs.Front()

		val, _ := front.Value.(func())

		val()

		logs.Remove(front)
	}
}

func LogInfo(s string) {
	if !configuration.ConfigHolder.Debug {
		return
	}

	logs.PushBack(func() {
		Logger.Info(fmt.Sprintf("[HTTPS] %s", s))
	})
}

func LogError(s string) {
	if !configuration.ConfigHolder.Debug {
		return
	}

	logs.PushBack(func() {
		Logger.Error(fmt.Sprintf("[HTTPS] %s", s))
	})
}

func LogConnection(connection *gin.Context) {
	LogInfo(fmt.Sprintf("[%s] -> %s", connection.ClientIP(), connection.Request.URL.Path))
}

func SetupLogger(appName string, level logrus.Level) {
	Logger = logrus.New()

	Logger.SetFormatter(&customFormat{
		AppName: appName,
	})

	Logger.SetLevel(level)

	go Log()
}

type customFormat struct {
	AppName string
}

func (F *customFormat) Format(ent *logrus.Entry) ([]byte, error) {
	return fmt.Appendf(nil, "[%s] [%s] %s\n", F.AppName, ent.Level.String(), ent.Message), nil
}
