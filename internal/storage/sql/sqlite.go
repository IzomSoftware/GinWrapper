package sql

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	"github.com/IzomSoftware/GinWrapper/internal/configuration"
)

type SQLiteStorage struct{}

func (S *SQLiteStorage) GetDBPool(config *configuration.SQLConfiguration) (*sql.DB, error) {
	sqliteConfiguration := config.SQLiteConfiguration
	return sql.Open("sqlite3", sqliteConfiguration.DatabaseLocation)
}
