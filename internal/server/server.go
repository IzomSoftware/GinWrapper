package server

import (
	"net/http"

	"github.com/IzomSoftware/GinWrapper/internal/authentication"
	"github.com/IzomSoftware/GinWrapper/internal/logger"
	"github.com/gin-gonic/gin"
)

type Server struct {
	JWTManager authentication.JWTManager
}

/*
 * Aborts & logs the connection with the given http status code.
 */
func (S *Server) AbortConnection(c *gin.Context, status int) {
	c.AbortWithStatus(status)

	logger.Info("Connection %s aborted with: %d", c.ClientIP(), status)
}

/*
 * Aborts & logs the suspicious connection with 403 forbidden http status code.
 */
func (S *Server) AbortSuspiciousConnection(c *gin.Context) {
	S.AbortConnection(c, http.StatusForbidden)

	logger.Info("Connection %s rejected", c.ClientIP())
}

/*
 * Bans & logs the banned suspicious connection with 403 forbidden http status code.
 */
func (S *Server) BanConnection(ip string, c *gin.Context) {
	S.AbortSuspiciousConnection(c)

	// if err := sql_source.BanIP(ip); err != nil {
	// 	logger.Error(fmt.Sprintf("Connection %s is already banned", ip))
	// }

	logger.Info("Connection %s banned", c.ClientIP())
}
