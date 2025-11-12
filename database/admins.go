package database

import (
	"context"
	"database/sql"
	"errors"
	"log"
)

func NewAdmins(db *sql.DB) error {
	schema := `CREATE TABLE IF NOT EXISTS admins (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        player TEXT NOT NULL,
        xuid TEXT NOT NULL UNIQUE,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`
	_, err := db.Exec(schema)
	return err
}

func NewAdminLimits(db *sql.DB) error {
	schema := `CREATE TABLE IF NOT EXISTS adminlimits (
		xuid TEXT PRIMARY KEY,
		total_given INTEGER NOT NULL DEFAULT 0,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := db.Exec(schema)
	return err
}

func AddAdmin(db *sql.DB, player, xuid string) error {
	_, err := db.Exec("INSERT OR REPLACE INTO admins (player, xuid) VALUES (?, ?)", player, xuid)
	if err == nil {
		log.Printf("Added admin: %s (%s)", player, xuid)
	}
	return err
}

func RemoveAdmin(db *sql.DB, xuid string) error {
	_, err := db.Exec("DELETE FROM admins WHERE xuid = ?", xuid)
	if err == nil {
		log.Printf("Removed admin with XUID: %s", xuid)
	}
	return err
}

func IsAdmin(db *sql.DB, xuid string) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM admins WHERE xuid = ?)", xuid).Scan(&exists)
	return exists, err
}

func ListAdmins(db *sql.DB) ([]struct{ Player, XUID string }, error) {
	rows, err := db.Query("SELECT player, xuid FROM admins ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []struct{ Player, XUID string }
	for rows.Next() {
		var p, x string
		if err := rows.Scan(&p, &x); err != nil {
			return nil, err
		}
		out = append(out, struct{ Player, XUID string }{Player: p, XUID: x})
	}
	return out, rows.Err()
}

func AdminGiveTotal(db *sql.DB, xuid string) (int64, error) {
	if db == nil {
		return 0, errors.New("database handle is nil")
	}

	var total int64
	err := db.QueryRow("SELECT total_given FROM admin_limits WHERE xuid = ?", xuid).Scan(&total)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	return total, err
}

func IncrementAdminGiveTotal(db *sql.DB, xuid string, amount int64) error {
	if db == nil {
		return errors.New("database handle is nil")
	}
	if amount <= 0 {
		return nil
	}

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	if _, err := tx.Exec(`
		INSERT INTO admin_limits (xuid, total_given)
		VALUES (?, ?)
		ON CONFLICT(xuid) DO UPDATE SET total_given = total_given + excluded.total_given,
		updated_at = CURRENT_TIMESTAMP`,
		xuid, amount,
	); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}
