package configuration

import (
	"os"

	"github.com/BurntSushi/toml"
)

type HttpsTlsConfiguration struct {
	Enable   bool   `toml:"enable"`
	CertFile string `toml:"cert_file"`
	KeyFile  string `toml:"key_file"`
}
type HTTPServer struct {
	Enabled          bool                  `toml:"enabled"`
	Address          string                `toml:"address"`
	Port             int                   `toml:"port"`
	TlsConfiguration HttpsTlsConfiguration `toml:"tls_configuration"`
}

type DatabaseConfiguration struct {
	Enabled            bool   `toml:"enabled"`
	SQLiteFileLocation string `toml:"sqlite_file_location"`
}

type JWTProtection struct {
	JWTSecret     string `toml:"jwt_secret"`
	JWTExpiration int    `toml:"jwt_expiration"`
}

type Protections struct {
	JWTProtection JWTProtection `toml:"jwt_protection"`
	APIUserAgent  string        `toml:"api_user_agent"`
}

type Configuration struct {
	Debug                 bool                  `toml:"debug"`
	HTTPServer            HTTPServer            `toml:"http_server"`
	DatabaseConfiguration DatabaseConfiguration `toml:"database"`
	Protections           Protections           `toml:"protections"`
}

var ConfigHolder Configuration
var DefaultConfig = Configuration{}

func SetupConfig(fileName string) error {
	ConfigHolder = DefaultConfig

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		file, err := os.Create(fileName)
		if err != nil {
			return err
		}

		defer file.Close()

		encoder := toml.NewEncoder(file)
		if err := encoder.Encode(ConfigHolder); err != nil {
			return err
		}

		return nil
	}

	if _, err := toml.DecodeFile(fileName, &ConfigHolder); err != nil {
		return err
	}

	return nil
}
