package sql

import (
	"database/sql"

	"github.com/IzomSoftware/GinWrapper/configuration"
)

type StorageImplementation interface {
	GetDBPool(config *configuration.SQLConfiguration) (*sql.DB, error)
}

type Storage struct {
	pool           *sql.DB
	CreationSchema string
}

func New(config *configuration.SQLConfiguration, impl StorageImplementation, creationSchema string) (*Storage, error) {
	pool, err := impl.GetDBPool(config)
	if err != nil {
		return nil, err
	}
	return &Storage{
		pool:           pool,
		CreationSchema: creationSchema,
	}, nil
}

func (S *Storage) SetupTables() error {
	_, err := S.pool.Exec(S.CreationSchema)
	return err
}

func (S *Storage) Ping() error {
	return S.pool.Ping()
}

func (S *Storage) Close() error {
	return S.pool.Close()
}

func (S *Storage) ExecuteUpdate(query string, args ...any) error {
	tx, err := S.pool.Begin()
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

func (S *Storage) QueryRow(query string, args ...any) *sql.Row {
	return S.pool.QueryRow(query, args...)
}

func (S *Storage) Query(query string, args ...any) (*sql.Rows, error) {
	return S.pool.Query(query, args...)
}
