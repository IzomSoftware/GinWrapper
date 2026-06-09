package sql

import (
	"database/sql"

	"github.com/IzomSoftware/GinWrapper/internal/configuration"
)

type SQLStorageImplementation interface {
	GetDBPool(config *configuration.SQLConfiguration) (*sql.DB, error)
}

type SQLStorage struct {
	dbPool            *sql.DB
	SQLCreationSchema string
}

func (S *SQLStorage) NewSQLStorage(config *configuration.SQLConfiguration, implementation SQLStorageImplementation, creationSchema string) (*SQLStorage, error) {
	dbPool, err := implementation.GetDBPool(config)
	if err != nil {
		return nil, err
	}
	return &SQLStorage{
		dbPool:            dbPool,
		SQLCreationSchema: creationSchema,
	}, nil
}

func (S *SQLStorage) SetupTables() error {
	_, err := S.dbPool.Exec(
		S.SQLCreationSchema,
	)
	return err
}

func (S *SQLStorage) Ping() error {
	return S.dbPool.Ping()
}

func (S *SQLStorage) ExecuteUpdate(query string, args ...any) error {
	tx, err := S.dbPool.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(query, args...)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (S *SQLStorage) Get(query string, args ...any) (any, error) {
	var result any

	err := S.dbPool.QueryRow(query, args...).Scan(&result)

	return result, err
}
