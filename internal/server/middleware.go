package server

import (
	"net/http"
	"strings"

	"github.com/IzomSoftware/GinWrapper/internal/authentication"
	"github.com/IzomSoftware/GinWrapper/internal/logger"
	"github.com/gin-gonic/gin"
)

func (S *Server) authentication(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" {
		S.AbortConnection(c, http.StatusForbidden)
		return
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		S.AbortConnection(c, http.StatusForbidden)
		return
	}

	claims, err := S.JWTManager.ValidateJWT(parts[1])
	if err != nil {
		S.AbortConnection(c, http.StatusForbidden)
		return
	}

	if err != nil {
		S.AbortConnection(c, http.StatusInternalServerError)
		return
	}
	// TODO: ban check

	c.Set("claims", claims)
	c.Set("uuid", claims.Uuid)
	c.Next()
}

func (S *Server) rateLimit(c *gin.Context) {
	c.Next()
}

func (S *Server) checkBan(c *gin.Context) {
	c.Next()
}

func (S *Server) context(c *gin.Context) {
	if val, ok := c.Get("claims"); ok {
		if claims, ok := val.(*authentication.JWTClaims); ok {
			_ = claims
		}
	}
	c.Next()
}

func (S *Server) log(c *gin.Context) {
	logger.LogConnection(c)
	c.Next()
}
//