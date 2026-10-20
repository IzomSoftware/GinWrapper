package cidr

import (
	"github.com/IzomSoftware/GinWrapper/configuration"
	"github.com/gin-gonic/gin"
)

func GetIP(c *gin.Context, config configuration.ReverseProxyMode) string {
	var ip string
	if config.Enabled {
		if config.ReadHeaderExplicitly {
			if ip = c.GetHeader(config.ExplicitReadHeader); ip == "" {
				ip = c.ClientIP()
			}
		} else {
			ip = c.ClientIP()
		}
	} else {
		ip = c.Request.RemoteAddr
	}
	return ip
}
