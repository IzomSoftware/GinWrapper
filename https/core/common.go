package https_core

import (
	"fmt"
	"net/http"

	"github.com/IzomSoftware/GinWrapper/logger"

	// "github.com/IzomSoftware/GinWrapper/redis"
	mysql "github.com/IzomSoftware/GinWrapper/sql"
	"github.com/gin-gonic/gin"
)


func AbortConnection(ip string, c *gin.Context) {
	logger.LogInfo(fmt.Sprintf("Connection %s is being aborted", ip))

	c.AbortWithStatus(http.StatusForbidden)
}

func BanConnection(ip string, c *gin.Context) {
	AbortConnection(ip, c)

	logger.LogInfo(fmt.Sprintf("Connection %s is being banned", ip))

	if err := mysql.BanIP(ip); err != nil {
		logger.LogInfo(fmt.Sprintf("Connection %s is already banned", ip))
	}
}

func (R *Response) IsAnyProtectionEnabled() bool {
	return R.Protections.UserAgent || R.Protections.JWT || R.Protections.RateLimit
}

