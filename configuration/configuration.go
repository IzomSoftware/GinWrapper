package configuration

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/IzomSoftware/GinWrapper/authentication"
)

type TlsConfiguration struct {
	Enable   bool   `toml:"enable"`
	CertFile string `toml:"cert_file"`
	KeyFile  string `toml:"key_file"`
}

type HTTPServer struct {
	Enabled          bool             `toml:"enabled"`
	Address          string           `toml:"address"`
	Port             int              `toml:"port"`
	TemplatesDir     string           `toml:"template_dir"`
	AssetsDir        string           `toml:"assets_dir"`
	TlsConfiguration TlsConfiguration `toml:"tls_configuration"`
}

type SQLiteConfiguration struct {
	Enabled          bool   `toml:"enabled"`
	DatabaseLocation string `toml:"database_location"`
}

type MySQLConfiguration struct {
	Enabled                bool   `toml:"enabled"`
	Hostname               string `toml:"hostname"`
	Port                   uint16 `toml:"port"`
	Username               string `toml:"username"`
	Password               string `toml:"password"`
	Database               string `toml:"database"`
	TLSEnabled             bool   `toml:"tls_enabled"`
	SkipTLSVerification    bool   `toml:"skip_tls_verification"`
	Charset                string `toml:"charset"`
	MaxOpenConnections     int    `toml:"max_open_connections"`
	MaxIdleConnections     int    `toml:"max_idle_connections"`
	ConnectionsMaxLifetime int    `toml:"connections_max_lifetime_seconds"`
	ParseTime              bool   `toml:"parse_time"`
}

type PostgreSQLConfiguration struct {
	Enabled                bool   `toml:"enabled"`
	Hostname               string `toml:"hostname"`
	Port                   uint16 `toml:"port"`
	Username               string `toml:"username"`
	Password               string `toml:"password"`
	Database               string `toml:"database"`
	SSLMode                string `toml:"ssl_mode"`
	MaxOpenConnections     int    `toml:"max_open_connections"`
	MaxIdleConnections     int    `toml:"max_idle_connections"`
	ConnectionsMaxLifetime int    `toml:"connections_max_lifetime_seconds"`
}

type SQLConfiguration struct {
	SQLiteConfiguration     SQLiteConfiguration     `toml:"sqlite_configuration"`
	MySQLConfiguration      MySQLConfiguration      `toml:"mysql_configuration"`
	PostgreSQLConfiguration PostgreSQLConfiguration `toml:"postgresql_configuration"`
}

type EmbeddedRedisConfiguration struct {
	Enabled bool `toml:"enabled"`
}

type DedicatedRedisConfiguration struct {
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

type RedisConfiguration struct {
	EmbeddedRedisConfiguration  EmbeddedRedisConfiguration  `toml:"embedded_redis_configuration"`
	DedicatedRedisConfiguration DedicatedRedisConfiguration `toml:"dedicated_redis_configuration"`
}

type DatabaseConfiguration struct {
	SQLiteConfiguration         SQLiteConfiguration         `toml:"sqlite_configuration"`
	MySQLConfiguration          MySQLConfiguration          `toml:"mysql_configuration"`
	PostgreSQLConfiguration     PostgreSQLConfiguration     `toml:"postgresql_configuration"`
	EmbeddedRedisConfiguration  EmbeddedRedisConfiguration  `toml:"embedded_redis_configuration"`
	DedicatedRedisConfiguration DedicatedRedisConfiguration `toml:"dedicated_redis_configuration"`
}

type RateLimitProtection struct {
	Enabled bool `toml:"enabled"`
	Rate    int  `toml:"rate"`
	Window  int  `toml:"window"`
}
type JWTProtection struct {
	JWTSecret     string `toml:"jwt_secret"`
	JWTExpiration int    `toml:"jwt_expiration"`
}

type OrderingProtection struct {
	Enabled bool                `toml:"enabled"`
	Orders  map[string][]string `toml:"orders"`
}

type Protections struct {
	APIUserAgent        string              `toml:"api_user_agent_protection"`
	RateLimitProtection RateLimitProtection `toml:"rate_limit_protection"`
	JWTProtection       JWTProtection       `toml:"jwt_protection"`
	OrderingProtection  OrderingProtection  `toml:"ordering_protection"`
}

type Config struct {
	Debug                 bool                  `toml:"debug"`
	HTTPServer            HTTPServer            `toml:"http_server"`
	DatabaseConfiguration DatabaseConfiguration `toml:"database"`
	Protections           Protections           `toml:"protections"`
}

var Default = Config{
	Debug: true,
	HTTPServer: HTTPServer{
		Enabled:      true,
		Address:      "0.0.0.0",
		Port:         2009,
		TemplatesDir: "./assets/templates/",
		AssetsDir:    "./assets/",
		TlsConfiguration: TlsConfiguration{
			Enable:   false,
			CertFile: "cert.pem",
			KeyFile:  "key.pem",
		},
	},
	DatabaseConfiguration: DatabaseConfiguration{
		SQLiteConfiguration: SQLiteConfiguration{
			Enabled:          true,
			DatabaseLocation: "db.sqlite",
		},
		EmbeddedRedisConfiguration: EmbeddedRedisConfiguration{
			Enabled: true,
		},
		MySQLConfiguration: MySQLConfiguration{
			Enabled:                false,
			Hostname:               "127.0.0.1",
			Port:                   3306,
			Username:               "root",
			Password:               "",
			Database:               "GinWrapper",
			TLSEnabled:             true,
			SkipTLSVerification:    true,
			Charset:                "utf8mb4",
			MaxOpenConnections:     151,
			MaxIdleConnections:     10,
			ConnectionsMaxLifetime: 3600,
			ParseTime:              true,
		},
		PostgreSQLConfiguration: PostgreSQLConfiguration{
			Enabled:                false,
			Hostname:               "127.0.0.1",
			Port:                   5432,
			Username:               "postgres",
			Password:               "",
			Database:               "GinWrapper",
			SSLMode:                "disable",
			MaxOpenConnections:     25,
			MaxIdleConnections:     5,
			ConnectionsMaxLifetime: 3600,
		},
		DedicatedRedisConfiguration: DedicatedRedisConfiguration{
			Enabled:             false,
			Hostname:            "127.0.0.1",
			Port:                6379,
			Username:            "root",
			Password:            "",
			Database:            0,
			PoolSize:            20,
			MaxRetries:          5,
			PoolTimeout:         1,
			DialTimeout:         1,
			ReadTimeout:         2,
			WriteTimeoutSec:     3,
			TLSEnabled:          true,
			SkipTLSVerification: true,
		},
	},
	Protections: Protections{
		APIUserAgent: "",
		RateLimitProtection: RateLimitProtection{
			Enabled: true,
			Rate:    30,
			Window:  60,
		},
		OrderingProtection: OrderingProtection{
			Enabled: false,
			Orders: map[string][]string{
				"/auth":      {"/", "/home"},
				"/dashboard": {"/auth"},
			},
		},
		JWTProtection: JWTProtection{
			JWTSecret:     "",
			JWTExpiration: 60,
		},
	},
}

var ErrMultipleStorageSources = fmt.Errorf("cannot enable multiple Redis/SQL databases at once")

func (c *Config) IsStorageConfigured() bool {
	return (c.DatabaseConfiguration.DedicatedRedisConfiguration.Enabled || c.DatabaseConfiguration.EmbeddedRedisConfiguration.Enabled) && (c.DatabaseConfiguration.MySQLConfiguration.Enabled || c.DatabaseConfiguration.SQLiteConfiguration.Enabled || c.DatabaseConfiguration.PostgreSQLConfiguration.Enabled)
}

func LoadConfiguration(fileName string) (*Config, error) {
	configuration := Default

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		file, err := os.Create(fileName)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		encoder := toml.NewEncoder(file)
		secret, err := authentication.GenerateRandomSecret(32)
		if err != nil {
			return nil, err
		}
		configuration.Protections.JWTProtection.JWTSecret = secret
		if err := encoder.Encode(&configuration); err != nil {
			return nil, err
		}
		return &configuration, nil
	}

	if _, err := toml.DecodeFile(fileName, &configuration); err != nil {
		return nil, err
	}

	databaseConfiguration := configuration.DatabaseConfiguration
	if (databaseConfiguration.SQLiteConfiguration.Enabled && databaseConfiguration.MySQLConfiguration.Enabled) || (databaseConfiguration.DedicatedRedisConfiguration.Enabled && databaseConfiguration.EmbeddedRedisConfiguration.Enabled) {
		return nil, ErrMultipleStorageSources
	}

	return &configuration, nil
}
