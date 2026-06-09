package responses

import (
	"fmt"

	"github.com/IzomSoftware/GinWrapper/internal/logger"
	"github.com/gin-gonic/gin"
)


var (
	RedisSourceError = fmt.Errorf("Unable to access Redis source value")
)

/*
 * Aborts & logs the connection with the given http status code.
 */
func AbortConnection(ip string, c *gin.Context, status int) {
	c.AbortWithStatus(status)

	logger.Info("Connection %s aborted with: %d", ip, status)
}

/*
 * Aborts & logs the suspicious connection with 403 forbidden http status code.
 */
func AbortSuspiciousConnection(ip string, c *gin.Context) {
	// AbortConnection(ip, c, http.StatusForbidden)

	logger.Info(fmt.Sprintf("Connection %s rejected", ip))
}

/*
 * Bans & logs the banned suspicious connection with 403 forbidden http status code.
 */
func BanConnection(ip string, c *gin.Context) {
	AbortSuspiciousConnection(ip, c)

	// if err := sql_source.BanIP(ip); err != nil {
	// 	logger.Error(fmt.Sprintf("Connection %s is already banned", ip))
	// }

	logger.Info(fmt.Sprintf("Connection %s banned", ip))
}
