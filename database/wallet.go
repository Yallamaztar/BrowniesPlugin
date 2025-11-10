package database

import (
	"database/sql"
	"sync"
)

type Wallet struct {
	player  string
	xuid    string
	balance int64
	db      *sql.DB
	mu      sync.RWMutex
}

func NewWallet(player, xuid string, balance int64, db *sql.DB) *Wallet {
	schema := `CREATE TABLE IF NOT EXISTS wallets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		player TEXT NOT NULL,
		xuid TEXT NOT NULL UNIQUE,
		balance INTEGER NOT NULL DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(schema); err != nil {
		return nil
	}

	_, err := db.Exec("INSERT OR IGNORE INTO wallets (player, xuid, balance) VALUES (?, ?, ?)", player, xuid, balance)
	if err != nil {
		return nil
	}

	var actualBalance int64
	err = db.QueryRow("SELECT balance FROM wallets WHERE xuid = ?", xuid).Scan(&actualBalance)
	if err != nil {
		actualBalance = balance
	}

	return &Wallet{player: player, xuid: xuid, balance: actualBalance, db: db}
}

func (w *Wallet) Balance() int64 {
	w.mu.RLock()
	defer w.mu.RUnlock()

	var balance int64
	err := w.db.QueryRow("SELECT balance FROM wallets WHERE xuid = ?", w.xuid).Scan(&balance)
	if err == nil {
		w.balance = balance
	}

	return w.balance
}

func (w *Wallet) AddBalance(amount int64) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	_, err := w.db.Exec("UPDATE wallets SET balance = balance + ? WHERE xuid = ?", amount, w.xuid)
	if err != nil {
		return err
	}

	w.balance += amount
	return nil
}

func (w *Wallet) SubtractBalance(amount int64) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	_, err := w.db.Exec("UPDATE wallets SET balance = balance - ? WHERE xuid = ? AND balance >= ?", amount, w.xuid, amount)
	if err != nil {
		return err
	}

	w.balance -= amount
	return nil
}

func GetWallet(player, xuid string, db *sql.DB) *Wallet {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM wallets WHERE xuid = ?)", xuid).Scan(&exists)

	if err != nil || !exists {
		return NewWallet(player, xuid, 1000, db)
	}

	var balance int64
	err = db.QueryRow("SELECT balance FROM wallets WHERE xuid = ?", xuid).Scan(&balance)
	if err != nil {
		balance = 1000
	}

	return &Wallet{player: player, xuid: xuid, balance: balance, db: db}
}
