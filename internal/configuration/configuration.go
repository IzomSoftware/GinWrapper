package configuration

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

/*
 * the TLS configuration struct
 */
type HttpsTlsConfiguration struct {
	Enable   bool   `toml:"enable"`
	CertFile string `toml:"cert_file"`
	KeyFile  string `toml:"key_file"`
}

/*
 * the HTTP configuration struct
 */
type HTTPServer struct {
	Enabled          bool                  `toml:"enabled"`
	Address          string                `toml:"address"`
	Port             int                   `toml:"port"`
	TemplatesDir     string                `toml:"template_dir"`
	AssetsDir        string                `toml:"assets_dir"`
	TlsConfiguration HttpsTlsConfiguration `toml:"tls_configuration"`
}

/*
 * the SQLite configuration struct
 */
type SQLiteConfiguration struct {
	Enabled          bool   `toml:"enabled"`
	DatabaseLocation string `toml:"database_location"`
}

/*
 * the MYSQL configuration struct
 */
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

type SQLConfiguration struct {
	SQLiteConfiguration SQLiteConfiguration `toml:"sqlite_configuration"`
	MySQLConfiguration  MySQLConfiguration  `toml:"mysql_configuration"`
}


type RedisConfiguration struct {
	EmbeddedRedisConfiguration EmbeddedRedisConfiguration `toml:"embedded_redis_configuration"`
	DedicatedRedisConfiguration DedicatedRedisConfiguration `toml:"dedicated_redis_configuration"`
}

/*
 * the EmbeddedRedis configuration struct
 */
type EmbeddedRedisConfiguration struct {
	Enabled bool `toml:"enabled"`
}

/*
 * the Redis configuration struct
 */
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

/*
 * the whole storage configuration struct
 */
type DatabaseConfiguration struct {
	SQLiteConfiguration        SQLiteConfiguration        `toml:"sqlite_configuration"`
	EmbeddedRedisConfiguration EmbeddedRedisConfiguration `toml:"embedded_redis_configuration"`
	MySQLConfiguration         MySQLConfiguration         `toml:"mysql_configuration"`
	DedicatedRedisConfiguration         DedicatedRedisConfiguration         `toml:"redis_configuration"`
}

/*
 * the JWT related configurations struct
 */
type JWTProtection struct {
	JWTSecret     string `toml:"jwt_secret"`
	JWTExpiration int    `toml:"jwt_expiration"`
}

/*
 * the BasicProtections configuration struct
 */
type BasicProtections struct {
	Provide    bool `toml:"provide"`
	Aggressive bool `toml:"aggressive"`
}

/*
 * the OrderingProtection configuration struct
 */
type OrderingProtection struct {
	Enabled bool                `toml:"enabled"`
	Orders  map[string][]string `toml:"orders"`
}

/*
 * the whole protections configuration struct
 */
type Protections struct {
	UserPassAPI        bool               `toml:"user_pass_api"`
	APIUserAgent       string             `toml:"api_user_agent"`
	BasicProtections   BasicProtections   `toml:"basic_protections"`
	JWTProtection      JWTProtection      `toml:"jwt_protection"`
	OrderingProtection OrderingProtection `toml:"ordering_protection"`
}

/*
 * the whole configuration itself
 */
type Configuration struct {
	Debug                 bool                  `toml:"debug"`
	HTTPServer            HTTPServer            `toml:"http_server"`
	DatabaseConfiguration DatabaseConfiguration `toml:"database"`
	Protections           Protections           `toml:"protections"`
}

// The configuration holder (inside ram of course)
var ConfigHolder *Configuration

// The default configuration we provide to help developers
var DefaultConfig = Configuration{
	Debug: true,
	HTTPServer: HTTPServer{
		Enabled:      true,
		Address:      "0.0.0.0",
		Port:         2009,
		TemplatesDir: "./assets/templates/",
		AssetsDir:    "./assets/",
		TlsConfiguration: HttpsTlsConfiguration{
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
		UserPassAPI: true,
		BasicProtections: BasicProtections{
			Provide:    true,
			Aggressive: true,
		},
		APIUserAgent: "Test Client 1.0/b (Software)",
		OrderingProtection: OrderingProtection{
			Enabled: true,
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

// Errors
var cannotInitializeMultipleStorageSourcesAtOnce = fmt.Errorf("Can't enable multiple Redis/SQL based databases at once")

// Returns true if storage is configured correctly
func IsStorageConfigured() bool {
	databaseConfig := ConfigHolder.DatabaseConfiguration
	isRedisEnabled := databaseConfig.DedicatedRedisConfiguration.Enabled || databaseConfig.EmbeddedRedisConfiguration.Enabled
	isSqlEnabled := databaseConfig.MySQLConfiguration.Enabled || databaseConfig.SQLiteConfiguration.Enabled

	return isRedisEnabled && isSqlEnabled
}

// Setups the config based on the Default config (which developers should configure)
func SetupConfig(fileName string) error {
	ConfigHolder = &DefaultConfig

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

	databaseConfig := ConfigHolder.DatabaseConfiguration
	if (databaseConfig.MySQLConfiguration.Enabled &&
		databaseConfig.SQLiteConfiguration.Enabled) ||
		(databaseConfig.DedicatedRedisConfiguration.Enabled &&
			databaseConfig.EmbeddedRedisConfiguration.Enabled) {
		return cannotInitializeMultipleStorageSourcesAtOnce
	}

	return nil
}
