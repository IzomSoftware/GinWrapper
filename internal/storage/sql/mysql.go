package sql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"

	"github.com/IzomSoftware/GinWrapper/internal/configuration"
)

type MYSQLStorage struct{}

func (M *MYSQLStorage) GetDBPool(config *configuration.SQLConfiguration) (*sql.DB, error) {
	mysqlConfig := config.MySQLConfiguration
	dataSourceName := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&tls=%s",
		mysqlConfig.Username,
		mysqlConfig.Password,
		mysqlConfig.Hostname,
		mysqlConfig.Port,
		mysqlConfig.Database,
		mysqlConfig.Charset,
		mysqlConfig.ParseTime,
		func() string {
			if mysqlConfig.TLSEnabled {
				if mysqlConfig.SkipTLSVerification {
					return "skip-verify"
				}
				return "true"
			}
			return "false"
		}(),
	)

	var err error

	dbPool, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	dbPool.SetMaxOpenConns(mysqlConfig.MaxOpenConnections)
	dbPool.SetMaxIdleConns(mysqlConfig.MaxIdleConnections)
	dbPool.SetConnMaxLifetime(time.Duration(mysqlConfig.ConnectionsMaxLifetime) * time.Second)

	return dbPool, nil
}
