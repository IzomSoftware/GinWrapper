package https_core

import (
	"fmt"
	"github.com/IzomSoftware/GinWrapper/common/logger"
	"time"

	"github.com/gin-gonic/gin"
)

func LogInfo(s string) {
	logger.Logger.Info(fmt.Sprintf("[HTTPS] %s", s))
}

var lastLog time.Time
var logs []string

func Log() {
	if !(len(logs) > 0) {
		return
	}

	for _, log := range logs {
		time.Sleep(time.Second * 3)
		LogInfo(log)
	}
	logs = nil
}
func LogConnection(connection *gin.Context) {
	formatted := fmt.Sprintf("[%s] -> %s", connection.ClientIP(), connection.FullPath())
	logs = append(logs, formatted)

	if !(time.Since(lastLog).Seconds() >= 5) {
		go Log()
	}

	lastLog = time.Now()
}
