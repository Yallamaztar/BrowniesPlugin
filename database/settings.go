package database

import (
	"database/sql"

	"github.com/Yallamaztar/BrowniesPlugin/helpers"
)

func EnsureSettings(db *sql.DB) error {
	schema := `CREATE TABLE IF NOT EXISTS settings (
		key TEXT PRIMARY KEY,
		value INTEGER NOT NULL
	);`
	_, err := db.Exec(schema)
	return err
}

func SetMaxBet(db *sql.DB, amount int64) error {
	if db == nil {
		return sql.ErrConnDone
	}
	_, err := db.Exec(
		"INSERT INTO settings (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value=excluded.value",
		"max_bet", amount,
	)

	helpers.GamblingMaxBetWebhook(amount)
	return err
}

func GetMaxBet(db *sql.DB) int64 {
	if db == nil {
		return 0
	}
	var v int64
	err := db.QueryRow("SELECT value FROM settings WHERE key = ?", "max_bet").Scan(&v)
	if err != nil {
		return 0
	}
	return v
}

func EnableGambling(db *sql.DB, enable bool) error {
	if db == nil {
		return sql.ErrConnDone
	}
	val := 0
	if enable {
		val = 1
	}

	_, err := db.Exec(
		"INSERT INTO settings (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value=excluded.value",
		"gambling_enabled", val,
	)

	helpers.GamblingEnabledWebhook(enable)
	return err
}
func IsGamblingEnabled(db *sql.DB) bool {
	if db == nil {
		return false
	}

	var v int64
	err := db.QueryRow("SELECT value FROM settings WHERE key = ?", "gambling_enabled").Scan(&v)
	if err != nil {
		return false
	}
	return v == 1
}
