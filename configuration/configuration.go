package configuration

import (
	"fmt"
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

type SQLiteConfiguration struct {
	Enabled          bool   `toml:"enabled"`
	DatabaseLocation string `toml:"database_location"`
}

type EmbeddedRedisConfiguration struct {
	Enabled bool `toml:"enabled"`
}

type MySQLConfiguration struct {
	Enabled             bool   `toml:"enabled"`
	Hostname            string `toml:"hostname"`
	Port                uint16 `toml:"port"`
	Username            string `toml:"username"`
	Password            string `toml:"password"`
	Database            string `toml:"database"`
	TLSEnabled          bool   `toml:"tls_enabled"`
	SkipTLSVerification bool   `toml:"skip_tls_verification"`
	// utf8mb4
	Charset                string `toml:"charset"`
	MaxOpenConnections     int    `toml:"max_open_connections"`
	MaxIdleConnections     int    `toml:"max_idle_connections"`
	ConnectionsMaxLifetime int    `toml:"connections_max_lifetime_seconds"`
	ParseTime              bool   `toml:"parse_time"`
}

type RedisConfiguration struct {
	Enabled             bool   `toml:"enabled"`
	Hostname            string `toml:"hostname"`
	Port                uint16 `toml:"port"`
	Username            string `toml:"username"`
	Password            string `toml:"password"`
	Database            int    `toml:"database"`
	PoolSize            int    `toml:"pool_size"`
	MinIdleConnections  int    `toml:"min_idle_connections"`
	MaxRetries          int    `toml:"max_retries"`
	PoolTimeout         int    `toml:"pool_timeout"`
	DialTimeout         int    `toml:"dial_timeout"`
	ReadTimeout         int    `toml:"read_timeout"`
	WriteTimeoutSec     int    `toml:"write_timeout_sec"`
	TLSEnabled          bool   `toml:"tls_enabled"`
	SkipTLSVerification bool   `toml:"skip_tls_verification"`
}

type DatabaseConfiguration struct {
	Enabled                    bool                       `toml:"enabled"`
	SQLiteConfiguration        SQLiteConfiguration        `toml:"sqlite_configuration"`
	EmbeddedRedisConfiguration EmbeddedRedisConfiguration `toml:"embedded_redis_configuration"`
	MySQLConfiguration         MySQLConfiguration         `toml:"mysql_configuration"`
	RedisConfiguration		RedisConfiguration	`toml:"redis_configuration"`
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

	if (ConfigHolder.DatabaseConfiguration.MySQLConfiguration.Enabled &&
	 ConfigHolder.DatabaseConfiguration.SQLiteConfiguration.Enabled) || 
	 (ConfigHolder.DatabaseConfiguration.RedisConfiguration.Enabled && 
		ConfigHolder.DatabaseConfiguration.EmbeddedRedisConfiguration.Enabled) {
		return fmt.Errorf("Can't enable multiple Redis/SQL based databases at once")
	}

	return nil
}
