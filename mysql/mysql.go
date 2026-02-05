package mysql

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"github.com/SyNdicateFoundation/GinWrapper/common/configuration"
	"github.com/SyNdicateFoundation/GinWrapper/common/logger"
	_ "github.com/glebarez/go-sqlite"
	"time"
)

var dbPool *sql.DB

func InitDatabase() time.Duration {
	start := time.Now()

	var err error

	sqlLiteConfig := configuration.ConfigHolder.SQLLiteConfiguration

	if dbPool, err = sql.Open("sqlite", sqlLiteConfig.DatabaseFileLocation); err != nil {
		logger.Logger.Error(err)
	}

	if _, err = dbPool.Exec(`
		CREATE TABLE IF NOT EXISTS Users (
			username TEXT PRIMARY KEY,
			password TEXT NOT NULL,
			hwid TEXT,
			banned INTEGER NOT NULL CHECK (banned IN (0, 1))
		)`); err != nil {
		logger.Logger.Fatal("Failed to create Users table:", err)
	}

	if _, err = dbPool.Exec(`
		CREATE TABLE IF NOT EXISTS BlockedIPs (
			ip TEXT PRIMARY KEY
		)`); err != nil {
		logger.Logger.Fatal("Failed to create BlockedIPs table:", err)
	}

	return time.Since(start)
}

func ExecuteUpdate(query string, args ...any) {
	tx, err := dbPool.Begin()
	if err != nil {
		logger.Logger.Error(err)
	}

	defer func(tx *sql.Tx) {
		if err := tx.Rollback(); err != nil {
			logger.Logger.Error(err)
		}
	}(tx)

	stmt, err := tx.Prepare(query)
	if err != nil {
		logger.Logger.Error(err)
	}

	defer func(stmt *sql.Stmt) {
		if err := stmt.Close(); err != nil {
			logger.Logger.Error(err)
		}
	}(stmt)

	_, err = stmt.Exec(args...)
	if err != nil {
		logger.Logger.Error(err)
	}

	if err := tx.Commit(); err != nil {
		logger.Logger.Error(err)
	}
}

func GetData(query string, args ...any) any {
	var result any
	if err := dbPool.QueryRow(query, args...).Scan(&result); err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Logger.Error(err)
	}
	return result
}

func CreateUser(username, password string) {
	if CheckUserExists(username) {
		return
	}
	ExecuteUpdate("INSERT INTO Users (username, password, banned) VALUES (?, ?, 0)", username, SHAPass(password))
}

func BanUser(username string) {
	ExecuteUpdate("UPDATE Users SET banned = 1 WHERE username = ?", username)
}

func BanIp(ip string) {
	if CheckIPExists(ip) {
		return
	}
	ExecuteUpdate("INSERT INTO BlockedIPs (ip) VALUES (?)", ip)
}

func UnbanUser(username string) {
	ExecuteUpdate("UPDATE Users SET banned = 0 WHERE username = ?", username)
}

func UnbanIp(ip string) {
	ExecuteUpdate("DELETE FROM BlockedIPs WHERE ip = ?", ip)
}

func DeleteUser(username string) {
	ExecuteUpdate("DELETE FROM Users WHERE username = ?", username)
}

func IsBanned(username string) bool {
	return GetData("SELECT banned FROM Users WHERE username = ?", username).(int64) == 1
}

func CheckUserExists(username string) bool {
	return GetData("SELECT COUNT(*) FROM Users WHERE username = ?", username).(int64) > 0
}

func CheckIPExists(ip string) bool {
	return GetData("SELECT COUNT(*) FROM BlockedIPs WHERE ip = ?", ip).(int64) > 0
}

func CheckPassword(username, password string) bool {
	return SHAPass(password) == GetData("SELECT password FROM Users WHERE username = ?", username).(string)
}

func HWIDExists(username string) bool {
	result := GetData("SELECT hwid FROM Users WHERE username = ?", username)
	hwid, ok := result.(string)
	return ok && hwid != ""
}

func CheckHWID(username, hwid string) bool {
	result := GetData("SELECT hwid FROM Users WHERE username = ?", username)
	stored, ok := result.(string)
	return ok && stored == hwid
}

func AddHWID(username, hwid string) {
	if !HWIDExists(username) {
		ExecuteUpdate("UPDATE Users SET hwid = ? WHERE username = ?", hwid, username)
	}
}

func SHAPass(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash[:])
}
