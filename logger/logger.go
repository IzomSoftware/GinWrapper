package logger

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger
var name string

func SetupLogger(appName string) {
	Logger = logrus.New()
	name = appName
	Logger.SetFormatter(&customFormat{})
	Logger.SetLevel(logrus.DebugLevel)
}

type customFormat struct{}

func (F *customFormat) Format(ent *logrus.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("[%s] [%s] %s\n", name, ent.Level.String(), ent.Message)), nil
}
