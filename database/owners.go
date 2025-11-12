package database

import (
	"database/sql"
	"log"
)

func NewOwners(db *sql.DB) error {
	schema := `CREATE TABLE IF NOT EXISTS owners (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		player TEXT NOT NULL,
		xuid TEXT NOT NULL UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := db.Exec(schema)
	return err
}

func AddOwner(db *sql.DB, player, xuid string) error {
	_, err := db.Exec("INSERT OR REPLACE INTO owners (player, xuid) VALUES (?, ?)", player, xuid)
	if err == nil {
		log.Printf("[Gambling] Added owner: %s (%s)", player, xuid)
	}
	return err
}

func RemoveOwner(db *sql.DB, xuid string) error {
	_, err := db.Exec("DELETE FROM owners WHERE xuid = ?", xuid)
	if err == nil {
		log.Printf("[Gambling] Removed owner with XUID: %s", xuid)
	}
	return err
}

func IsOwner(db *sql.DB, xuid string) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM owners WHERE xuid = ?)", xuid).Scan(&exists)
	return exists, err
}

func ListOwners(db *sql.DB) ([]struct{ Player, XUID string }, error) {
	rows, err := db.Query("SELECT player, xuid FROM owners ORDER BY created_at DESC")
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
