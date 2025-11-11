package database

import (
	"database/sql"
	"errors"
	"time"
)

type ShopItem struct {
	ID          int64
	Name        string
	Price       int64
	Command     string
	Description string
	CreatedAt   time.Time
}

func EnsureShop(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS shop_items (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL UNIQUE,
        price INTEGER NOT NULL,
        command TEXT,
        description TEXT,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    )`)
	return err
}

func addShopItem(db *sql.DB, name string, price int64, command, description string) error {
	if price < 0 {
		return errors.New("price cannot be negative")
	}
	_, err := db.Exec(`INSERT INTO shop_items (name, price, command, description) VALUES (?,?,?,?)`, name, price, command, description)
	return err
}

func GetShopItem(db *sql.DB, name string) (*ShopItem, error) {
	it := &ShopItem{}
	row := db.QueryRow(`SELECT id, name, price, command, description, created_at FROM shop_items WHERE name = ?`, name)
	if err := row.Scan(&it.ID, &it.Name, &it.Price, &it.Command, &it.Description, &it.CreatedAt); err != nil {
		return nil, err
	}
	return it, nil
}

func ListShopItems(db *sql.DB) ([]ShopItem, error) {
	rows, err := db.Query(`SELECT id, name, price, command, description, created_at FROM shop_items ORDER BY price ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ShopItem
	for rows.Next() {
		var it ShopItem
		if err := rows.Scan(&it.ID, &it.Name, &it.Price, &it.Command, &it.Description, &it.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, it)
	}
	return items, nil
}

func SeedShop(db *sql.DB) error {
	var count int
	if err := db.QueryRow(`SELECT COUNT(1) FROM shop_items`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	baseline := []struct {
		name, cmd, desc string
		price           int64
	}{
		{"killplayer", "killplayer %PLAYER%", "Kill A Player", 1_500_000},
		{"slap", "slap %PLAYER%", "Slap A Player", 20_000},
		{"spectator", "spectator %PLAYER%", "Switch to spectator team", 750_000},
		{"takeweapons", "takeweapons %PLAYER%", "Remove all your weapons", 250_000},
		{"hide", "hide %PLAYER%", "Become invisible (toggle)", 5_000_000},
		{"freeze", "freeze %PLAYER%", "Freeze/unfreeze A Player", 600_000},
		// {"setspeed", "setspeed %PLAYER% 1.5", "Set A Players Movement Speed", 700_500},
		// {"teleport", "teleport %PLAYER% 100 100 100", "Teleport to preset coordinates", 9_000},
		// {"giveweapon", "giveweapon %PLAYER% ak47", "Receive a weapon (AK-47)", 12_000},
	}
	for _, b := range baseline {
		_ = addShopItem(db, b.name, b.price, b.cmd, b.desc)
	}
	return nil
}
