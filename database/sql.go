package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/IzomSoftware/GinWrapper/configuration"
	_ "github.com/glebarez/go-sqlite"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var dbPool *sql.DB
var unexpectedTypeError = fmt.Errorf("Unexpected type for value")

/*
 * SQLite
 */
func InitSQLite(config configuration.SQLiteConfiguration) error {
	var err error

	dbPool, err = sql.Open("sqlite", config.DatabaseLocation)

	if err != nil {
		return err
	}

	err = SetupTables()
	if err != nil {
		return err
	}

	dbPool.SetMaxOpenConns(1)
	dbPool.SetMaxIdleConns(1)
	dbPool.SetConnMaxLifetime(0)

	return PingDatabase()
}

/*
 * MySQL
 */
func InitMySQL(config configuration.MySQLConfiguration) error {
	dataSourceName := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&tls=%s",
		config.Username,
		config.Password,
		config.Hostname,
		config.Port,
		config.Database,
		config.Charset,
		config.ParseTime,
		func() string {
			if config.TLSEnabled {
				if config.SkipTLSVerification {
					return "skip-verify"
				}
				return "true"
			}
			return "false"
		}(),
	)

	var err error

	dbPool, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		return err
	}

	dbPool.SetMaxOpenConns(config.MaxOpenConnections)
	dbPool.SetMaxIdleConns(config.MaxIdleConnections)
	dbPool.SetConnMaxLifetime(time.Duration(config.ConnectionsMaxLifetime) * time.Second)

	return PingDatabase()
}

/*
 * Common Methods
 */
func Init() error {
	databaseConfig := configuration.ConfigHolder.DatabaseConfiguration

	if databaseConfig.SQLiteConfiguration.Enabled {
		return InitSQLite(databaseConfig.SQLiteConfiguration)
	}

	if databaseConfig.MySQLConfiguration.Enabled {
		return InitMySQL(databaseConfig.MySQLConfiguration)
	}

	panic("Unreachable state")
}

func ExecuteUpdate(query string, args ...any) error {
	tx, err := dbPool.Begin()
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

func GetData(query string, args ...any) (any, error) {
	var result any

	err := dbPool.QueryRow(query, args...).Scan(&result)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return result, nil
}

func PingDatabase() error {
	return dbPool.Ping()
}
func SetupTables() error {
	_, err := dbPool.Exec(`
		CREATE TABLE IF NOT EXISTS Users (
			id UUID PRIMARY KEY,
            username TEXT UNIQUE NOT NULL,
            password TEXT NOT NULL,
            jwt_version INTEGER NOT NULL DEFAULT 0,
			banned INTEGER NOT NULL CHECK (banned IN (0, 1))
		)`)
	if err != nil {
		return err
	}

	_, err = dbPool.Exec(`
		CREATE TABLE IF NOT EXISTS BlockedIPs (
			ip TEXT PRIMARY KEY
		)`)
	if err != nil {
		return err
	}

	_, err = dbPool.Exec(`
		CREATE TABLE IF NOT EXISTS BlockedHWIDs (
			hwid TEXT PRIMARY KEY
		)`)
	if err != nil {
		return err
	}

	_, err = dbPool.Exec(`
		CREATE TABLE IF NOT EXISTS RefreshJWTs (
			id UUID PRIMARY KEY,
            user_id UUID REFERENCES Users(id) ON DELETE CASCADE,
            token_hash TEXT NOT NULL,
            expires_at TIMESTAMP NOT NULL,
            created_at TIMESTAMP NOT NULL DEFAULT NOW(),
            revoked_at TIMESTAMP
		)`)
	if err != nil {
		return err
	}

	return nil
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CheckUsernameExists(username string) (bool, error) {
	result, err := GetData("SELECT COUNT(*) FROM Users WHERE username = ?", username)
	if err != nil {
		return false, err
	}

	val, ok := result.(int64)
	if !ok {
		return false, unexpectedTypeError
	}

	return val > 0, nil
}

func CheckUserExists(uuid string) (bool, error) {
	result, err := GetData("SELECT COUNT(*) FROM Users WHERE id = ?", uuid)
	if err != nil {
		return false, err
	}

	val, ok := result.(int64)
	if !ok {
		return false, unexpectedTypeError
	}

	return val > 0, nil
}

func CheckIsUserBanned(uuid string) (bool, error) {
	result, err := GetData("SELECT COUNT(*) FROM Users WHERE id = ?", uuid)
	if err != nil {
		return false, err
	}

	val, ok := result.(int64)
	if !ok {
		return false, unexpectedTypeError
	}

	return val == 1, nil
}

func CheckIsBannedHWID(hwid string) (bool, error) {
	result, err := GetData("SELECT COUNT(*) FROM BlockedHWIDs WHERE hwid = ?", hwid)
	if err != nil {
		return false, err
	}

	val, ok := result.(int64)
	if !ok {
		return false, unexpectedTypeError
	}

	return val > 0, nil
}

func CheckIsBannedIP(ip string) (bool, error) {
	result, err := GetData("SELECT COUNT(*) FROM BlockedIPs WHERE ip = ?", ip)
	if err != nil {
		return false, err
	}

	val, ok := result.(int64)
	if !ok {
		return false, fmt.Errorf("Unexpected type for value")
	}

	return val > 0, nil
}

func CreateUser(username, password string) error {
	userExists, err := CheckUsernameExists(username)
	if err != nil {
		return err
	}
	if userExists {
		return fmt.Errorf("Username already exists")
	}

	uuid := uuid.NewString()
	uuidExists, err := CheckUserExists(uuid)
	if uuidExists {
		return fmt.Errorf("User already exists")
	}

	hash, err := HashPassword(password)
	if err != nil {
		return err
	}

	return ExecuteUpdate("INSERT INTO Users (id, username, password, jwt_version, banned) VALUES (?, ?, ?, 0, 0)", uuid, username, hash)
}

func BanUser(uuid string) error {
	result, err := CheckIsUserBanned(uuid)
	if err != nil {
		return err
	}

	if result {
		return fmt.Errorf("User is already banned")
	}
	return ExecuteUpdate("UPDATE Users SET banned = 1 WHERE id = ?", uuid)
}

func BanIP(ip string) error {
	result, err := CheckIsBannedIP(ip)
	if err != nil {
		return err
	}

	if result {
		return fmt.Errorf("IP is already banned")
	}
	return ExecuteUpdate("INSERT INTO BlockedIPs (ip) VALUES (?)", ip)
}

func BanHWID(hwid string) error {
	result, err := CheckIsBannedHWID(hwid)
	if err != nil {
		return err
	}

	if result {
		return fmt.Errorf("HWID is already banned")
	}
	return ExecuteUpdate("INSERT INTO BlockedHWIDs (hwid) VALUES (?)", hwid)
}

func UnbanUser(uuid string) error {
	return ExecuteUpdate("UPDATE Users SET banned = 0 WHERE id = ?", uuid)
}

func UnbanIP(ip string) error {
	return ExecuteUpdate("DELETE FROM BlockedIPs WHERE ip = ?", ip)
}

func UnbanHWID(hwid string) error {
	return ExecuteUpdate("DELETE FROM BlockedHWIDs WHERE hwid = ?", hwid)
}