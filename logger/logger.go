package logger

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func SetupLogger(appName string, level logrus.Level) {
	Logger = logrus.New()

	Logger.SetFormatter(&customFormat{
        AppName: appName,
    })
	Logger.SetLevel(level)
}

type customFormat struct {
    AppName string
}

func (F *customFormat) Format(ent *logrus.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("[%s] [%s] %s\n", F.AppName, ent.Level.String(), ent.Message)), nil
}
