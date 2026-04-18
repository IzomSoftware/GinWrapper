package https_core

import (
	"fmt"
	"net/http"

	"github.com/IzomSoftware/GinWrapper/logger"
	"github.com/IzomSoftware/GinWrapper/storage/sql"
	"github.com/gin-gonic/gin"
)

/*
 * Aborts & logs the connection with the given http status code.
 */
func AbortConnection(ip string, c *gin.Context, status int) {
	c.AbortWithStatus(status)

	logger.LogInfo(fmt.Sprintf("Connection %s aborted with: %d", ip, status))
}

/*
 * Aborts & logs the suspicious connection with 403 forbidden http status code.
 */
func AbortSuspiciousConnection(ip string, c *gin.Context) {
	AbortConnection(ip, c, http.StatusForbidden)

	logger.LogInfo(fmt.Sprintf("Connection %s rejected", ip))
}

/*
 * Bans & logs the banned suspicious connection with 403 forbidden http status code.
 */
func BanConnection(ip string, c *gin.Context) {
	AbortSuspiciousConnection(ip, c)

	if err := sql.BanIP(ip); err != nil {
		logger.LogError(fmt.Sprintf("Connection %s is already banned", ip))
	}

	logger.LogInfo(fmt.Sprintf("Connection %s banned", ip))
}

/*
 * Returns true if any protection is enabled
 */
func (R *Response) IsAnyProtectionEnabled() bool {
	return R.Protections.BasicProtections || R.Protections.UserAgent || R.Protections.JWT || R.Protections.RateLimit
}
