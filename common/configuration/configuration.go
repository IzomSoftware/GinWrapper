package configuration

import (
	"os"

	"github.com/IzomSoftware/GinWrapper/common/logger"

	"github.com/BurntSushi/toml"
)

// HTTPSServer -------------- HTTPS config holders --------------
type HttpsTlsConfiguration struct {
	Enable   bool   `toml:"enable"`
	CertFile string `toml:"cert_file"`
	KeyFile  string `toml:"key_file"`
}
type HTTPSServer struct {
	Enabled          bool                  `toml:"enabled"`
	Address          string                `toml:"address"`
	Port             int                   `toml:"port"`
	TlsConfiguration HttpsTlsConfiguration `toml:"tls_configuration"`
}

// SQLLiteConfiguration -------------- SQLLite config holders --------------
type SQLLiteConfiguration struct {
	Enabled              bool   `toml:"enabled"`
	DatabaseFileLocation string `toml:"file_location"`
}

type JWTProtection struct {
	JWTSecret     string `toml:"jwt_secret"`
	JWTExpiration int    `toml:"jwt_expiration"`
}

type Protections struct {
	JWTProtection JWTProtection `toml:"jwt_protection"`
	APIUserAgent  string        `toml:"api_user_agent"`
}

type Holder struct {
	Debug                bool                 `toml:"debug"`
	HTTPSServer          HTTPSServer          `toml:"https_server"`
	SQLLiteConfiguration SQLLiteConfiguration `toml:"database"`
	Protections          Protections          `toml:"protections"`
}

var ConfigHolder Holder
var DefaultConfig = Holder{}

func SetupConfig(fileName string) {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		file, err := os.Create(fileName)
		ConfigHolder = DefaultConfig
		if err != nil {
			logger.Logger.Error(err)
		}
		defer func(file *os.File) {

			if err := file.Close(); err != nil {
				logger.Logger.Error(err)
			}
		}(file)

		encoder := toml.NewEncoder(file)
		if err := encoder.Encode(ConfigHolder); err != nil {
			logger.Logger.Error(err)
		}
	}

	if _, err := toml.DecodeFile(fileName, &ConfigHolder); err != nil {
		logger.Logger.Error(err)
	}
}
