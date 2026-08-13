package storage

import (
	"context"
	"fmt"

	"github.com/IzomSoftware/GinWrapper/internal/configuration"
	"github.com/IzomSoftware/GinWrapper/internal/storage/redis"
	"github.com/IzomSoftware/GinWrapper/internal/storage/sql"
)

type Storage struct {
	SQL   *sql.Storage
	Redis *redis.Storage
}

var ErrNoStorageEnabled = fmt.Errorf("no storage backend enabled")

func New(config *configuration.Config, creationSchema string) (*Storage, error) {
	storage := &Storage{}
	ctx := context.Background()

	sql, err := initSQL(config, creationSchema)
	if err != nil {
		return nil, fmt.Errorf("sql init: %w", err)
	}
	storage.SQL = sql

	redis, err := initRedis(config, ctx)
	if err != nil {
		return nil, fmt.Errorf("redis init: %w", err)
	}
	storage.Redis = redis

	return storage, nil
}

func initSQL(config *configuration.Config, creationSchema string) (*sql.Storage, error) {
	databaseConfiguration := config.DatabaseConfiguration

	var storage *sql.Storage
	var err error

	if databaseConfiguration.MySQLConfiguration.Enabled {
		storage, err = sql.New(
			&configuration.SQLConfiguration{MySQLConfiguration: databaseConfiguration.MySQLConfiguration},
			&sql.MYSQLStorage{},
			creationSchema,
		)
	} else if databaseConfiguration.SQLiteConfiguration.Enabled {
		storage, err = sql.New(
			&configuration.SQLConfiguration{SQLiteConfiguration: databaseConfiguration.SQLiteConfiguration},
			&sql.SQLiteStorage{},
			creationSchema,
		)
	} else {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if err := storage.SetupTables(); err != nil {
		return nil, err
	}

	return storage, nil
}

func initRedis(config *configuration.Config, ctx context.Context) (*redis.Storage, error) {
	databaseConfiguration := config.DatabaseConfiguration

	if databaseConfiguration.DedicatedRedisConfiguration.Enabled {
		return redis.New(
			&configuration.RedisConfiguration{DedicatedRedisConfiguration: databaseConfiguration.DedicatedRedisConfiguration},
			ctx,
			&redis.DedicatedRedisStorage{},
		)
	}

	if databaseConfiguration.EmbeddedRedisConfiguration.Enabled {
		return redis.New(
			&configuration.RedisConfiguration{EmbeddedRedisConfiguration: databaseConfiguration.EmbeddedRedisConfiguration},
			ctx,
			&redis.EmbeddedRedisStorage{},
		)
	}

	return nil, nil
}

func (storage *Storage) Close() error {
	if storage.SQL != nil {
		if err := storage.SQL.Close(); err != nil {
			return err
		}
	}
	if storage.Redis != nil {
		if err := storage.Redis.Close(); err != nil {
			return err
		}
	}
	return nil
}
